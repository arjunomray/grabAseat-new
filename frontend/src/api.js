// Base URL for API requests.
//
// In development the Vite dev server proxies /api to the backend, so the
// default relative path works without configuration. Set VITE_API_URL to point
// at a different backend (see .env.example).
const API_BASE = import.meta.env.VITE_API_URL ?? '/api/v1';

export async function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers);
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const token = localStorage.getItem('token');
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  return fetch(`${API_BASE}${path}`, { ...options, headers });
}
