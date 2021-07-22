import {Component, Injectable} from '@angular/core';
import {ComponentStore} from '@ngrx/component-store';
import {ConfigService} from '../config.service';
import {pluck, tap} from 'rxjs/operators';
import {Observable} from 'rxjs';

interface State {
  seed: string;
  version: string;
  buildTime: string;
}

@Injectable()
export class TrackTileViewerStore extends ComponentStore<State> {
  readonly seed$ = super.select(state => state.seed);
  readonly version$ = super.select(state => state.version);
  readonly buildTime$ = super.select(state => state.buildTime);

  readonly setSeed$ = super.updater((state, seed: string) => ({...state, seed}))

  readonly loadConfig = super.effect((trigger$: Observable<void>) => this.service.fetchConfig()
    .pipe(tap(config => super.patchState(config))))

  constructor(private readonly service: ConfigService) {
    super({seed: '', buildTime: '', version: ''});
  }
}
