import {
  AfterViewInit,
  ChangeDetectionStrategy, ChangeDetectorRef,
  Component,
} from '@angular/core';
import {latLng, LatLngBounds, tileLayer} from 'leaflet';

@Component({
  selector: 'track-tile-viewer',
  templateUrl: './track-tile-viewer.component.html',
  styleUrls: ['./track-tile-viewer.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: {class: 'overflow-hidden'}
})
export class TrackTileViewerComponent {
  options = {
    layers: [
      tileLayer('http://127.0.0.1:9551/{z}/{x}/{y}.png', {maxZoom: 18, minZoom: 18, attribution: '...',})
    ],
    zoom: 5,
    center: latLng(-0.00069, 0)
  };

  bounding = new LatLngBounds({lat:-0.00069*2, lng: -181}, {lat: 0, lng: 181})
}
