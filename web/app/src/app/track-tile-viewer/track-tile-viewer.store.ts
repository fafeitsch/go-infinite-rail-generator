import {Component, Injectable} from '@angular/core';
import {ComponentStore} from '@ngrx/component-store';
import {ConfigService} from '../config.service';
import {pluck, tap} from 'rxjs/operators';
import {Observable} from 'rxjs';

interface State {
  seed: string;
}

@Injectable()
export class TrackTileViewerStore extends ComponentStore<State> {
  readonly seed$ = super.select(state => state.seed);
  readonly setSeed$ = super.updater((state, seed: string) => ({...state, seed}))

  readonly loadDefaultSeed$ = super.effect((trigger$: Observable<void>) => this.service.fetchConfig()
    .pipe(pluck('defaultSeed'), tap(seed => super.patchState({seed}))))

  constructor(private readonly service: ConfigService) {
    super({seed: ''});
  }
}
