package plugin

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-plugin"
	"github.com/jsiebens/libdns-plugin/internal/proto"
	"github.com/libdns/libdns"
	"google.golang.org/protobuf/types/known/durationpb"
)

type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.ProviderClient
}

func (g *GRPCClient) Configure(ctx context.Context, message json.RawMessage) error {
	_, err := g.client.Configure(ctx, &proto.ConfigureRequest{Value: message})
	return err
}

func (g *GRPCClient) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	resp, err := g.client.GetRecords(ctx, &proto.GetRecordsRequest{Zone: zone})
	if err != nil {
		return nil, err
	}
	return fromProto(resp.Records), nil
}

func (g *GRPCClient) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	resp, err := g.client.SetRecords(ctx, &proto.RecordsRequest{Zone: zone, Records: toProto(recs)})
	if err != nil {
		return nil, err
	}
	return fromProto(resp.Records), nil
}

func (g *GRPCClient) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	resp, err := g.client.AppendRecords(ctx, &proto.RecordsRequest{Zone: zone, Records: toProto(recs)})
	if err != nil {
		return nil, err
	}
	return fromProto(resp.Records), nil
}

func (g *GRPCClient) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	resp, err := g.client.DeleteRecords(ctx, &proto.RecordsRequest{Zone: zone, Records: toProto(recs)})
	if err != nil {
		return nil, err
	}
	return fromProto(resp.Records), nil
}

type GRPCServer struct {
	Impl   Provider
	broker *plugin.GRPCBroker
}

func (g *GRPCServer) Configure(ctx context.Context, request *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	err := g.Impl.Configure(ctx, request.Value)
	if err != nil {
		return nil, err
	}
	return &proto.ConfigureResponse{}, nil
}

func (g *GRPCServer) GetRecords(ctx context.Context, request *proto.GetRecordsRequest) (*proto.RecordsResponse, error) {
	records, err := g.Impl.GetRecords(ctx, request.Zone)
	if err != nil {
		return nil, err
	}
	return &proto.RecordsResponse{Records: toProto(records)}, nil
}

func (g *GRPCServer) SetRecords(ctx context.Context, request *proto.RecordsRequest) (*proto.RecordsResponse, error) {
	records, err := g.Impl.SetRecords(ctx, request.Zone, fromProto(request.Records))
	if err != nil {
		return nil, err
	}
	return &proto.RecordsResponse{Records: toProto(records)}, nil
}

func (g *GRPCServer) AppendRecords(ctx context.Context, request *proto.RecordsRequest) (*proto.RecordsResponse, error) {
	records, err := g.Impl.AppendRecords(ctx, request.Zone, fromProto(request.Records))
	if err != nil {
		return nil, err
	}
	return &proto.RecordsResponse{Records: toProto(records)}, nil
}

func (g *GRPCServer) DeleteRecords(ctx context.Context, request *proto.RecordsRequest) (*proto.RecordsResponse, error) {
	records, err := g.Impl.DeleteRecords(ctx, request.Zone, fromProto(request.Records))
	if err != nil {
		return nil, err
	}
	return &proto.RecordsResponse{Records: toProto(records)}, nil
}

func fromProto(records []*proto.Record) []libdns.Record {
	result := make([]libdns.Record, len(records))
	for i, record := range records {
		result[i] = libdns.RR{
			Name: record.Name,
			Type: record.Type,
			Data: record.Data,
			TTL:  record.Ttl.AsDuration(),
		}
	}
	return result
}

func toProto(records []libdns.Record) []*proto.Record {
	result := make([]*proto.Record, len(records))
	for i, record := range records {
		result[i] = &proto.Record{
			Name: record.RR().Name,
			Type: record.RR().Type,
			Data: record.RR().Data,
			Ttl:  durationpb.New(record.RR().TTL),
		}
	}
	return result
}
