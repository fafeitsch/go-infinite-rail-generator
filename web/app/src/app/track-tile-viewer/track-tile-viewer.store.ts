import {Component, Injectable} from '@angular/core';
import {ComponentStore} from '@ngrx/component-store';
import {ConfigService} from '../config.service';
import {pluck, tap} from 'rxjs/operators';
import {Observable} from 'rxjs';

interface State {
  defaultSeed: string;
  version: string;
  buildTime: string;
}

@Injectable()
export class TrackTileViewerStore extends ComponentStore<State> {
  readonly seed$ = super.select(state => state.defaultSeed);
  readonly version$ = super.select(state => state.version);
  readonly buildTime$ = super.select(state => state.buildTime);

  readonly setSeed$ = super.updater((state, seed: string) => ({...state, defaultSeed: seed}))

  readonly loadConfig = super.effect((trigger$: Observable<void>) => this.service.fetchConfig()
    .pipe(tap(config => super.patchState(config))))

  constructor(private readonly service: ConfigService) {
    super({defaultSeed: '', buildTime: '', version: ''});
  }
}
