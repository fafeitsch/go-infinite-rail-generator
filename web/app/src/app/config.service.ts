import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';

export interface ConfigDto {
  defaultSeed: string;
  buildTime: string;
  version: string;
}

@Injectable({ providedIn: 'root' })
export class ConfigService {
  constructor(private readonly http: HttpClient) {}

  fetchConfig() {
    return this.http.get<ConfigDto>(`${environment.url}/config`);
  }
}
