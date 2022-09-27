package plugin

import (
	"context"
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/genproto/pluginv2"
	"github.com/toddtreece/grafana-grpc-server-example-datsource/pkg/plugin/settings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	_ backend.QueryDataHandler      = (*GRPCServerQueryDatasource)(nil)
	_ backend.CheckHealthHandler    = (*GRPCServerQueryDatasource)(nil)
	_ instancemgmt.InstanceDisposer = (*GRPCServerQueryDatasource)(nil)
)

func NewGRPCServerQueryDatasource(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	grpcSettings := settings.Load(s)
	conn, err := grpc.Dial(grpcSettings.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pluginv2.NewDataClient(conn)
	return &GRPCServerQueryDatasource{client: client, settings: grpcSettings}, nil
}

type GRPCServerQueryDatasource struct {
	client   pluginv2.DataClient
	settings *settings.Settings
}

func (d *GRPCServerQueryDatasource) Dispose() {}

func (d *GRPCServerQueryDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	req, err := d.setDataSourceFromTarget(ctx, req)
	protoReq := backend.ToProto().QueryDataRequest(req)
	ctx = d.addAuthToContext(ctx)
	res, err := d.client.QueryData(ctx, protoReq)
	if err != nil {
		return nil, err
	}
	return backend.FromProto().QueryDataResponse(res)
}

func (d *GRPCServerQueryDatasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if d.settings.URL == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "No URL configured",
		}, nil
	}

	if d.settings.Token == "" {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "No Token configured",
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

func (d *GRPCServerQueryDatasource) addAuthToContext(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + d.settings.Token,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

func (d *GRPCServerQueryDatasource) setDataSourceFromTarget(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataRequest, error) {
	for i, q := range req.Queries {
		parsed := make(map[string]interface{})
		err := json.Unmarshal(q.JSON, &parsed)
		if err != nil {
			return nil, err
		}

		parsed["datasource"] = parsed["target_datasource"]
		delete(parsed, "target_datasource")
		delete(parsed, "datasourceId")
		log.DefaultLogger.Info("parsed query", "parsed", parsed)
		raw, err := json.Marshal(&parsed)
		if err != nil {
			return nil, err
		}

		q.JSON = raw
		req.Queries[i] = q
	}
	return req, nil
}
