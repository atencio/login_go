// Única fuente de la navegación: la usan el menú lateral y la matriz de permisos de los perfiles,
// así que ambos muestran los mismos grupos, el mismo orden y los mismos nombres.
// El módulo de permisos de cada ítem es el prefijo de su permiso ("users.list" → "users").
export const NAV_GROUPS = [
  { label: null, title: 'General', items: [
    { to: '/', label: 'Dashboard', icon: 'dashboard', end: true },
    { to: '/reportes', label: 'Reportes', icon: 'chart', permission: 'reports.list' },
  ] },
  { label: 'Inventario de activos', items: [
    { to: '/asignaciones', label: 'Asignaciones', icon: 'swap', permission: 'assignments.list' },
    { to: '/equipos', label: 'Equipos de oficina', icon: 'monitor', permission: 'equipment.list' },
    { to: '/funcionarios', label: 'Funcionarios', icon: 'idcard', permission: 'employees.list' },
    { to: '/cargos', label: 'Cargos', icon: 'briefcase', permission: 'positions.list' },
    { to: '/categorias', label: 'Categorías', icon: 'tag', permission: 'categories.list' },
    { to: '/marcas', label: 'Marcas', icon: 'award', permission: 'brands.list' },
  ] },
  { label: 'Administración', items: [
    { to: '/usuarios', label: 'Usuarios', icon: 'users', permission: 'users.list' },
    { to: '/perfiles', label: 'Perfiles', icon: 'shield', permission: 'roles.list' },
  ] },
];

const moduleOf = (item) => item.permission?.split('.')[0];

// Secciones de la matriz: [{ title, rows: [{ module, label, actions }] }] con el orden y los nombres del menú.
// Los módulos del catálogo que aún no estén en el menú caen en "Otros".
export function matrixSections(modules) {
  const byModule = new Map(modules.map((m) => [m.module, m]));
  const used = new Set();
  const sections = [];
  NAV_GROUPS.forEach((g) => {
    const rows = g.items
      .map((item) => {
        const m = byModule.get(moduleOf(item));
        if (!m) return null;
        used.add(m.module);
        return { ...m, label: item.label };
      })
      .filter(Boolean);
    if (rows.length) sections.push({ title: g.label || g.title, rows });
  });
  const rest = modules.filter((m) => !used.has(m.module));
  if (rest.length) sections.push({ title: 'Otros', rows: rest });
  return sections;
}
