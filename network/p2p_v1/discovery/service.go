package discovery

import (
	"context"
	"github.com/bloxapp/ssv/network/p2p_v1/discovery/mdns"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"go.uber.org/zap"
)

type DiscType string

const (
	Mdns   DiscType = "mdns"
	DiscV5 DiscType = "discv5"
)

const (
	udp4 = "udp4"
	udp6 = "udp6"
	tcp  = "tcp"
)

// CheckPeerLimit enables listener to check peers limit
type CheckPeerLimit func() bool

// Connect is the function interface for connecting to some peer
type Connect func(*peer.AddrInfo) error

// PeerEvent is the data passed to peer handler
type PeerEvent struct {
	AddrInfo peer.AddrInfo
	Node     *enode.Node
}

// HandleNewPeer is the function interface for handling new peer
type HandleNewPeer func(e PeerEvent)

// Options represents the options passed to create a service
type Options struct {
	Logger *zap.Logger

	Type       DiscType
	Host       host.Host
	Connect    Connect
	CheckPeerLimit CheckPeerLimit
	DiscV5Opts *DiscV5Options
}

// Service is the interface for discovery
type Service interface {
	discovery.Discovery

	Bootstrap(handler HandleNewPeer) error
}

// NewService creates new discovery.Service
func NewService(ctx context.Context, opts Options) (Service, error) {
	if opts.Type == Mdns {
		return mdns.NewMdnsDiscovery(ctx, opts.Logger, opts.Host)
	}
	return newDiscV5Service(ctx, &opts)
}
