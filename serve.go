package plugin

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-plugin"
	"github.com/libdns/libdns"
)

type LibDNSProvider interface {
	libdns.RecordGetter
	libdns.RecordSetter
	libdns.RecordAppender
	libdns.RecordDeleter
}

func Serve(provider LibDNSProvider) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]plugin.Plugin{
			"provider": &ProviderPlugin{Impl: &configurableProvider{provider: provider}},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type configurableProvider struct {
	provider LibDNSProvider
}

func (c *configurableProvider) Configure(_ context.Context, message json.RawMessage) error {
	return json.Unmarshal(message, &c.provider)
}

func (c *configurableProvider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	return c.provider.GetRecords(ctx, zone)
}

func (c *configurableProvider) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return c.provider.SetRecords(ctx, zone, recs)
}

func (c *configurableProvider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return c.provider.AppendRecords(ctx, zone, recs)
}

func (c *configurableProvider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	return c.provider.DeleteRecords(ctx, zone, recs)
}
