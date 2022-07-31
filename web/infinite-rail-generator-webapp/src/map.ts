import {latLng, LatLngBounds, map, tileLayer} from 'leaflet'

export function setupMap(element: HTMLElement) {
  const url = import.meta.env.VITE_BACKEND + 'tiles?tile={x}&vertical={y}'
  const bounding = new LatLngBounds(
    { lat: -0.00069 * 2, lng: -181 },
    { lat: 0, lng: 181 }
  );
  const leafletMap = map(element)
    .setMaxZoom(18)
    .setMinZoom(18)
    .setView(latLng(-0.00069, 0), 18)
    .setMaxBounds(bounding)
  tileLayer(url, {
    attribution: 'go-infinite-rail-generator'
  }).addTo(leafletMap);
}
