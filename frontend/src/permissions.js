// Columnas de la matriz de permisos. El catálogo de módulos viene del backend (GET /permissions),
// así que un módulo nuevo aparece solo en cuanto se agrega a la tabla permissions.
export const ACTIONS = [
  { key: 'create', label: 'Crear' },
  { key: 'list', label: 'Listar' },
  { key: 'find', label: 'Encontrar' },
  { key: 'update', label: 'Editar' },
  { key: 'delete', label: 'Eliminar' },
];

// [{ module, label, actions: { create: 'users.create', ... } }]
export function groupCatalog(catalog = []) {
  const map = new Map();
  catalog.forEach((p) => {
    if (!map.has(p.module)) map.set(p.module, { module: p.module, label: p.module_label, actions: {} });
    map.get(p.module).actions[p.action] = p.code;
  });
  return [...map.values()];
}
