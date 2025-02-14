package connections

import (
	"context"
	"time"

	"github.com/bloxapp/ssv/logging/fields"
	"github.com/bloxapp/ssv/network/peers"
	"github.com/bloxapp/ssv/network/records"
	"github.com/bloxapp/ssv/network/streams"
	"github.com/bloxapp/ssv/operator/storage"
	libp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/identify"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// errHandshakeInProcess is thrown when and handshake process for that peer is already running
var errHandshakeInProcess = errors.New("handshake already in process")

// errPeerWasFiltered is thrown when a peer is filtered during handshake
var errPeerWasFiltered = errors.New("peer was filtered during handshake")

// errUnknownUserAgent is thrown when a peer has an unknown user agent
var errUnknownUserAgent = errors.New("user agent is unknown")

// errConsumingMessage is thrown when we сan't consume(parse) message: data is broken or incoming msg is from node with different Permissioned mode
// example: the Node in NON-Permissoned mode receives SignedNodeInfo; the Node in Permissoned mode receives NodeInfo
var errConsumingMessage = errors.New("error consuming message")

// errPeerPruned is thrown when remote peer is pruned
var errPeerPruned = errors.New("peer is pruned")

// HandshakeFilter can be used to filter nodes once we handshaked with them
type HandshakeFilter func(senderID peer.ID, sni records.AnyNodeInfo) error

// SubnetsProvider returns the subnets of or node
type SubnetsProvider func() records.Subnets

// Handshaker is the interface for handshaking with peers.
// it uses node info protocol to exchange information with other nodes and decide whether we want to connect.
//
// NOTE: due to compatibility with v0,
// we accept nodes with user agent as a fallback when the new protocol is not supported.
type Handshaker interface {
	Handshake(logger *zap.Logger, conn libp2pnetwork.Conn) error
	Handler(logger *zap.Logger) libp2pnetwork.StreamHandler
}

type handshaker struct {
	ctx context.Context

	Permissioned func() bool
	filters      func() []HandshakeFilter

	streams     streams.StreamController
	nodeInfoIdx peers.NodeInfoIndex
	states      peers.NodeStates
	connIdx     peers.ConnectionIndex
	subnetsIdx  peers.SubnetsIndex
	ids         identify.IDService
	net         libp2pnetwork.Network
	nodeStorage storage.Storage

	subnetsProvider SubnetsProvider
}

// HandshakerCfg is the configuration for creating an handshaker instance
type HandshakerCfg struct {
	Network         libp2pnetwork.Network
	Streams         streams.StreamController
	NodeInfoIdx     peers.NodeInfoIndex
	States          peers.NodeStates
	ConnIdx         peers.ConnectionIndex
	SubnetsIdx      peers.SubnetsIndex
	IDService       identify.IDService
	NodeStorage     storage.Storage
	SubnetsProvider SubnetsProvider
	Permissioned    func() bool
}

// NewHandshaker creates a new instance of handshaker
func NewHandshaker(ctx context.Context, cfg *HandshakerCfg, filters func() []HandshakeFilter) Handshaker {
	h := &handshaker{
		ctx:             ctx,
		streams:         cfg.Streams,
		nodeInfoIdx:     cfg.NodeInfoIdx,
		connIdx:         cfg.ConnIdx,
		subnetsIdx:      cfg.SubnetsIdx,
		ids:             cfg.IDService,
		filters:         filters,
		states:          cfg.States,
		subnetsProvider: cfg.SubnetsProvider,
		net:             cfg.Network,
		nodeStorage:     cfg.NodeStorage,
		Permissioned:    cfg.Permissioned,
	}
	return h
}

// Handler returns the handshake handler
func (h *handshaker) Handler(logger *zap.Logger) libp2pnetwork.StreamHandler {
	return func(stream libp2pnetwork.Stream) {
		// start by marking the peer as pending
		pid := stream.Conn().RemotePeer()
		pidStr := pid.String()

		req, res, done, err := h.streams.HandleStream(logger, stream)
		defer done()
		if err != nil {
			return
		}

		logger := logger.With(zap.String("otherPeer", pidStr))
		permissioned := h.Permissioned()
		var ani records.AnyNodeInfo

		if permissioned {

			sni := &records.SignedNodeInfo{}
			err = sni.Consume(req)
			if err != nil {
				logger.Warn("could not consume node info request", zap.Error(err))
				return
			}

			ani = sni
		} else {

			ni := &records.NodeInfo{}
			err = ni.Consume(req)
			if err != nil {
				logger.Warn("could not consume node info request", zap.Error(err))
				return
			}

			ani = ni
		}
		// process the node info in a new goroutine so we won't block the stream
		go func() {
			err := h.processIncomingNodeInfo(logger, pid, ani)
			if err != nil {
				if errors.Is(err, errPeerWasFiltered) {
					logger.Debug("peer was filtered", zap.Error(err))
					return
				}
				logger.Warn("could not process node info", zap.Error(err))
			}
		}()

		privateKey, found, err := h.nodeStorage.GetPrivateKey()
		if !found {
			logger.Warn("could not get private key", zap.Error(err))
			return
		}

		self, err := h.nodeInfoIdx.SelfSealed(h.net.LocalPeer(), pid, permissioned, privateKey)
		if err != nil {
			logger.Warn("could not seal self node info", zap.Error(err))
			return
		}

		if err := res(self); err != nil {
			logger.Warn("could not send self node info", zap.Error(err))
			return
		}
	}
}

func (h *handshaker) processIncomingNodeInfo(logger *zap.Logger, sender peer.ID, ani records.AnyNodeInfo) error {
	h.updateNodeSubnets(logger, sender, ani.GetNodeInfo())
	if err := h.applyFilters(sender, ani); err != nil {
		return err
	}

	if _, err := h.nodeInfoIdx.AddNodeInfo(logger, sender, ani.GetNodeInfo()); err != nil {
		return err
	}
	return nil
}

// preHandshake makes sure that we didn't reach peers limit and have exchanged framework information (libp2p)
// with the peer on the other side of the connection.
// it should enable us to know the supported protocols of peers we connect to
func (h *handshaker) preHandshake(conn libp2pnetwork.Conn) error {
	ctx, cancel := context.WithTimeout(h.ctx, time.Second*15)
	defer cancel()
	select {
	case <-ctx.Done():
		return errors.New("identity protocol (libp2p) timeout")
	case <-h.ids.IdentifyWait(conn):
	}
	return nil
}

// Handshake initiates handshake with the given conn
func (h *handshaker) Handshake(logger *zap.Logger, conn libp2pnetwork.Conn) error {
	pid := conn.RemotePeer()
	// check if the peer is known before we continue
	ni, err := h.getNodeInfo(pid)
	if err != nil || ni != nil {
		return err
	}
	if err := h.preHandshake(conn); err != nil {
		return errors.Wrap(err, "could not perform pre-handshake")
	}

	var ani records.AnyNodeInfo

	ani, err = h.nodeInfoFromStream(logger, conn)
	if err != nil {
		return err
	}

	logger = logger.With(zap.String("otherPeer", pid.String()), zap.Any("info", ani))

	err = h.processIncomingNodeInfo(logger, pid, ani)
	if err != nil {
		logger.Debug("could not process node info", zap.Error(err))
		return err
	}

	return nil
}

func (h *handshaker) getNodeInfo(pid peer.ID) (*records.NodeInfo, error) {
	ni, err := h.nodeInfoIdx.GetNodeInfo(pid)
	if err != nil && err != peers.ErrNotFound {
		return nil, errors.Wrap(err, "could not read node info")
	}
	if ni != nil {
		switch h.states.State(pid) {
		case peers.StateIndexing:
			return nil, errHandshakeInProcess
		case peers.StatePruned:
			return nil, errors.Wrap(errPeerPruned, pid.String())
		case peers.StateReady:
			return ni, nil
		default: // unknown > continue the flow
		}
	}
	return nil, nil
}

// updateNodeSubnets tries to update the subnets of the given peer
func (h *handshaker) updateNodeSubnets(logger *zap.Logger, pid peer.ID, ni *records.NodeInfo) {
	if ni.Metadata != nil {
		subnets, err := records.Subnets{}.FromString(ni.Metadata.Subnets)
		if err == nil && len(subnets) > 0 {
			updated := h.subnetsIdx.UpdatePeerSubnets(pid, subnets)
			if updated {
				logger.Debug("[handshake] peer subnets were updated", fields.PeerID(pid),
					zap.String("subnets", subnets.String()))
			}
		}
	}
}

func (h *handshaker) nodeInfoFromStream(logger *zap.Logger, conn libp2pnetwork.Conn) (records.AnyNodeInfo, error) {
	res, err := h.net.Peerstore().FirstSupportedProtocol(conn.RemotePeer(), peers.NodeInfoProtocol)
	if err != nil {
		return nil, errors.Wrapf(err, "could not check supported protocols of peer %s",
			conn.RemotePeer().String())
	}

	permissioned := h.Permissioned()

	privateKey, found, err := h.nodeStorage.GetPrivateKey()
	if !found {
		return nil, err
	}
	data, err := h.nodeInfoIdx.SelfSealed(h.net.LocalPeer(), conn.RemotePeer(), permissioned, privateKey)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, errors.Errorf("peer [%s] doesn't supports handshake protocol", conn.RemotePeer().String())
	}
	resBytes, err := h.streams.Request(logger, conn.RemotePeer(), peers.NodeInfoProtocol, data)
	if err != nil {
		return nil, err
	}

	var ani records.AnyNodeInfo
	if permissioned {
		sni := &records.SignedNodeInfo{}
		err = sni.Consume(resBytes)
		ani = sni
	} else {
		ni := &records.NodeInfo{}
		err = ni.Consume(resBytes)
		ani = ni
	}

	if err != nil {
		return nil, errors.Wrap(errConsumingMessage, err.Error())
	}
	return ani, nil
}

func (h *handshaker) applyFilters(sender peer.ID, ani records.AnyNodeInfo) error {
	fltrs := h.filters()
	for i := range fltrs {
		err := fltrs[i](sender, ani)
		if err != nil {
			return errors.Wrap(errPeerWasFiltered, err.Error())
		}
	}

	return nil
}
