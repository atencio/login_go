const BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

let onUnauthorized = () => {};
export const setUnauthorizedHandler = (fn) => { onUnauthorized = fn; };

// "Mantener sesión" guarda el token en localStorage; si no, solo vive mientras la pestaña esté abierta.
export const tokenStore = {
  get: () => localStorage.getItem('token') || sessionStorage.getItem('token'),
  set: (t, remember) => {
    tokenStore.clear();
    (remember ? localStorage : sessionStorage).setItem('token', t);
  },
  clear: () => { localStorage.removeItem('token'); sessionStorage.removeItem('token'); },
};

export async function api(path, { method = 'GET', body } = {}) {
  const token = tokenStore.get();
  let res;
  try {
    res = await fetch(`${BASE}${path}`, {
      method,
      headers: { 'Content-Type': 'application/json', ...(token && { Authorization: `Bearer ${token}` }) },
      body: body ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new Error('No se pudo conectar con el servidor. Revisa tu conexión e inténtalo de nuevo.');
  }
  if (res.status === 204) return null;
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    if (res.status === 401 && token) onUnauthorized();
    throw new Error(data.error || 'Ocurrió un error inesperado');
  }
  return data;
}
