const hercules = window.hercules || {};
const apiHost = hercules.apiHost || 'http://127.0.0.1:8080';

export function fetch(endpoint) {
  return window.fetch(`${apiHost}${endpoint}`).then(r => {
    return r.json().catch(() => {
      if (r.ok) {
        return Promise.reject("Can't parse server response");
      }

      return Promise.reject(`Network error: ${r.status} ${r.statusText}`);
    });
  });
}
