// Base URL for API requests.
//
// In development the Vite dev server proxies /api to the backend, so the
// default relative path works without configuration. Set VITE_API_URL to point
// at a different backend (see .env.example).
const API_BASE = import.meta.env.VITE_API_URL ?? '/api/v1';

// Endpoints that must not trigger the global 401 redirect (a 401 there means
// "wrong credentials", not "your session expired").
const AUTH_PATHS = ['/login', '/register'];

function decodeToken(token) {
  const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
  const padded = base64 + '='.repeat((4 - (base64.length % 4)) % 4);
  return JSON.parse(atob(padded));
}

// isAuthenticated reports whether a non-expired token is stored.
export function isAuthenticated() {
  const token = localStorage.getItem('token');
  if (!token) return false;

  try {
    const { exp } = decodeToken(token);
    return !exp || exp * 1000 > Date.now();
  } catch {
    return false;
  }
}

function clearSessionAndRedirect() {
  localStorage.removeItem('token');
  localStorage.removeItem('userId');
  if (window.location.pathname !== '/login') {
    window.location.assign('/login');
  }
}

export async function apiFetch(path, options = {}) {
  const headers = new Headers(options.headers);
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const token = localStorage.getItem('token');
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const response = await fetch(`${API_BASE}${path}`, { ...options, headers });

  if (response.status === 401 && !AUTH_PATHS.some((p) => path.startsWith(p))) {
    clearSessionAndRedirect();
  }

  return response;
}
