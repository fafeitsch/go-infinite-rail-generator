import {
  AfterViewInit,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
} from '@angular/core';
import { latLng, LatLngBounds, tileLayer } from 'leaflet';
import { ConfigService } from '../config.service';
import { pluck } from 'rxjs/operators';

@Component({
  selector: 'track-tile-viewer',
  templateUrl: './track-tile-viewer.component.html',
  styleUrls: ['./track-tile-viewer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: { class: 'overflow-hidden' },
})
export class TrackTileViewerComponent {
  options = {
    layers: [
      tileLayer('http://127.0.0.1:9551/tiles?hectometer={x}&vertical={y}', {
        maxZoom: 18,
        minZoom: 18,
        attribution: 'go-infinite-rail-generator',
      }),
    ],
    zoomControl: false,
    zoom: 18,
    center: latLng(-0.00069, 0),
  };
  defaultSeed$ = this.configService.fetchConfig().pipe(pluck('defaultSeed'));

  bounding = new LatLngBounds(
    { lat: -0.00069 * 2, lng: -181 },
    { lat: 0, lng: 181 }
  );

  constructor(private readonly configService: ConfigService) {}
}
