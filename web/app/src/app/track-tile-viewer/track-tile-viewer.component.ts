import {
  AfterViewInit,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
} from '@angular/core';
import { latLng, LatLngBounds, tileLayer } from 'leaflet';
import { ConfigService } from '../config.service';
import {map, pluck, tap} from 'rxjs/operators';
import {TrackTileViewerStore} from './track-tile-viewer.store';
import {FormControl} from '@angular/forms';

@Component({
  selector: 'track-tile-viewer',
  templateUrl: './track-tile-viewer.component.html',
  styleUrls: ['./track-tile-viewer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: { class: 'overflow-hidden' },
  providers: [TrackTileViewerStore]
})
export class TrackTileViewerComponent {
  options = {
    layers: [],
    zoomControl: false,
    zoom: 18,
    center: latLng(-0.00069, 0),
  }

  layers$  = this.store.seed$.pipe(map(seed => [
     tileLayer(`http://127.0.0.1:9551/tiles?hectometer={x}&vertical={y}&seed=${seed}`, {
        maxZoom: 18,
        minZoom: 18,
        attribution: 'go-infinite-rail-generator',
      }),
      ]))

  defaultSeed$ = this.store.seed$;

  bounding = new LatLngBounds(
    { lat: -0.00069 * 2, lng: -181 },
    { lat: 0, lng: 181 }
  );

  seedControl = new FormControl('')

  constructor(private readonly configService: ConfigService, private readonly store: TrackTileViewerStore) {
    this.store.setSeed$(this.seedControl.valueChanges)
  }
}
