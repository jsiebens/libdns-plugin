package plugin

import (
	"context"
	"encoding/json"
	"github.com/libdns/libdns"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/jsiebens/libdns-plugin/internal/proto"
)

const ProviderPluginName = "provider"

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "LIBDNS_PLUGIN",
	MagicCookieValue: "k29lqez4d5dzfh629t2it4fv8n804blr38nj9495uj3wqrgoyzalmg2l7jnqmz6e",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	ProviderPluginName: &ProviderPlugin{},
}

type Provider interface {
	libdns.RecordGetter
	libdns.RecordSetter
	libdns.RecordAppender
	libdns.RecordDeleter

	Configure(context.Context, json.RawMessage) error
}

type ProviderPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Provider
}

func (p *ProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterProviderServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *ProviderPlugin) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewProviderClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &ProviderPlugin{}
var _ Provider = &GRPCClient{}
