import { api } from './client';

let cache = null;

// Catálogo de permisos (módulo + acción). Se pide una sola vez por sesión de la página.
export function loadCatalog() {
  if (!cache) cache = api('/permissions').catch((e) => { cache = null; throw e; });
  return cache;
}
