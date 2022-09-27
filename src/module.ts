import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { GRPCServerQuery, GRPCServerDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, GRPCServerQuery, GRPCServerDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
