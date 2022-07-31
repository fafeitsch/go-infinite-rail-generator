export function fetchConfig(versionElement: HTMLElement) {
  const url = import.meta.env.VITE_BACKEND
  fetch(url + '/config')
    .then((response) => response.json())
    .then(config => {
      const buildTime = config.buildTime ? ` (${config.buildTime})` : ''
      versionElement.innerHTML = `${config.version}${buildTime}`;
    })
}
