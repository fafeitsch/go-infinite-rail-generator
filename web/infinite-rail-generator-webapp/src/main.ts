import './style.scss'
import 'leaflet/dist/leaflet.css';
import {fetchConfig} from './config';
import {setupMap} from './map';

fetchConfig(document.querySelector<HTMLElement>('#version')!)
setupMap(document.querySelector<HTMLElement>('#map')!)
