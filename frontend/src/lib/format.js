export const fmtDateTime = (iso) =>
  iso ? new Intl.DateTimeFormat('es', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(iso)) : '—';

// Las fechas sin hora ("2025-03-01") se leen como locales para no correrse un día por la zona horaria.
export const fmtDate = (d) =>
  d ? new Intl.DateTimeFormat('es', { dateStyle: 'medium' }).format(new Date(`${d}T00:00:00`)) : '—';

export const fmtMoney = (n) =>
  n == null ? '—' : new Intl.NumberFormat('es', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(n);

export const STATUS = {
  available: { label: 'Disponible', tone: 'ok' },
  assigned: { label: 'Asignado', tone: 'info' },
  maintenance: { label: 'En mantenimiento', tone: 'warn' },
  retired: { label: 'Dado de baja', tone: 'muted' },
};
