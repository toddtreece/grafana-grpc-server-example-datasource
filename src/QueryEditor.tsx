import React from 'react';
import { QueryEditorProps, DataSourceInstanceSettings, DataSourceRef } from '@grafana/data';
import { DataSource } from './datasource';
import { GRPCServerDataSourceOptions, GRPCServerQuery } from './types';
import { DataSourcePicker } from '@grafana/runtime';

type Props = QueryEditorProps<DataSource, GRPCServerQuery, GRPCServerDataSourceOptions>;

export function QueryEditor (props: Props) {  
  const onDataSourceChange = (ref: DataSourceRef) => {
    const { onChange, query } = props;
    onChange({ ...query, datasource: ref, target_datasource: ref });
  };

    const { datasource} = props.query;

    return (
      <DataSourcePicker
        filter={ds => ds.meta.id != "grafana-grpc-server-example-datasource"}
        placeholder="Select a data source"
        onChange={(newSettings: DataSourceInstanceSettings) => {
          onDataSourceChange({ type: newSettings.type, uid: newSettings.uid });
        }}
        noDefault={true}
        current={datasource?.type != "grafana-grpc-server-example-datasource" ? datasource : undefined}
      />
    );
}
