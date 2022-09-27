import { DataQuery, DataSourceJsonData, DataSourceRef } from '@grafana/data';

export interface GRPCServerQuery extends DataQuery {
  target_datasource: DataSourceRef;
}

export interface GRPCServerDataSourceOptions extends DataSourceJsonData {
  url?: string;
}

export interface GRPCServerSecureData {
  token?: string;
}
