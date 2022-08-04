import store from '../store';
import {tap} from 'rxjs';

const template = document.createElement('template')

template.innerHTML = `
<style>
    .group {
        position: relative;
    }

    .group > label {
        padding: 0 0.2em;
        position: absolute;
        top: -0.5em;
        left: 0.5em;
        background-color: white;
    }

    .group > input {
        padding: 0.8em;
        border-width: 1px;
        border-style: solid;
        border-color: var(--primary-color)
    }
</style>
<div class="group">
  <label>Seed</label>
  <input class="themed-border" id="input" type="text" value="my own input"/>
</div>
`

class SeedInput extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({mode: 'open'})
    this.shadowRoot!.appendChild(template.content.cloneNode(true))
    const inputElement = this.shadowRoot!.querySelector('input')!;
    inputElement.addEventListener('input', () => {
      store.set.seed(inputElement.value)
    })
    store.effect(store.get.seed$.pipe(tap(seed => inputElement.value = seed)))
  }
}

window.customElements.define('seed-input', SeedInput)
