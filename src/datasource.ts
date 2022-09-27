import { DataQueryRequest, DataQueryResponse, DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { Observable } from 'rxjs';
import { GRPCServerDataSourceOptions, GRPCServerQuery } from './types';

export class DataSource extends DataSourceWithBackend<GRPCServerQuery, GRPCServerDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GRPCServerDataSourceOptions>) {
    super(instanceSettings);
  }
  
  query(request: DataQueryRequest<GRPCServerQuery>): Observable<DataQueryResponse> {
    request.targets = request.targets.map((target) => ({...target, datasource:{type: this.type, uid: this.uid}, datasourceId: undefined}))
    console.log(this.type, this.uid, request.targets);
    return super.query(request);
  }
}
