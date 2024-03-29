import {latLng, LatLngBounds, map, tileLayer, Map} from 'leaflet'
import store from './store';
import {tap} from 'rxjs';

export function setupMap(element: HTMLElement) {
  let leafletMap: Map | undefined
  store.effect(store.get.seed$.pipe(tap(seed => {
    leafletMap?.remove()
    const url = import.meta.env.VITE_BACKEND + 'tiles?tile={x}&vertical={y}&seed=' + seed
    const bounding = new LatLngBounds(
      {lat: -0.00069 * 2, lng: -181},
      {lat: 0, lng: 181}
    );
    leafletMap = map(element)
      .setMaxZoom(18)
      .setMinZoom(18)
      .setView(latLng(-0.00069, 0), 18)
      .setMaxBounds(bounding)
    tileLayer(url, {
      attribution: 'go-infinite-rail-generator'
    }).addTo(leafletMap);
  })))
}
