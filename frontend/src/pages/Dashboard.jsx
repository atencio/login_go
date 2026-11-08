import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { loadCatalog } from '../api/catalog';
import Icon from '../components/Icon';
import { ACTIONS, groupCatalog } from '../permissions';
import { useAuth } from '../context/AuthContext';

export default function Dashboard() {
  const { user, profile, can } = useAuth();
  const [summary, setSummary] = useState(null);
  const [inv, setInv] = useState(null);
  const [catalog, setCatalog] = useState([]);

  useEffect(() => {
    if (can('reports.list')) {
      api('/reports/users-summary').then(setSummary).catch(() => {});
      api('/reports/inventory-summary').then(setInv).catch(() => {});
    }
    loadCatalog().then(setCatalog).catch(() => {});
  }, [can]);

  const maxRole = summary ? Math.max(1, ...summary.by_role.map((r) => r.total)) : 1;
  const maxCat = inv ? Math.max(1, ...inv.by_category.map((c) => c.total)) : 1;
  const modules = groupCatalog(catalog).map((m) => ({
    ...m,
    granted: ACTIONS.filter((a) => m.actions[a.key] && profile.permissions.includes(m.actions[a.key])),
  }));

  return (
    <>
      <div className="page-head">
        <div>
          <h1>Dashboard</h1>
          <p className="muted">Hola, {user.name.split(' ')[0]}. Este es el resumen de tu espacio de trabajo.</p>
        </div>
      </div>

      {inv && (
        <div className="stats">
          <Stat icon="monitor" label="Equipos" value={inv.total_equipment} />
          <Stat icon="check" label="Disponibles" value={inv.available} />
          <Stat icon="swap" label="Asignados" value={inv.assigned} />
          <Stat icon="alert" label="En mantenimiento" value={inv.maintenance} />
        </div>
      )}
      {summary && inv && (
        <div className="stats stats--spaced">
          <Stat icon="idcard" label="Funcionarios activos" value={inv.employees} />
          <Stat icon="users" label="Usuarios del sistema" value={summary.total_users} />
          <Stat icon="x" label="Dados de baja" value={inv.retired} />
          <Stat icon="shield" label="Perfiles en uso" value={summary.by_role.filter((r) => r.total > 0).length} />
        </div>
      )}

      <div className="grid-2">
        {inv && inv.by_category.length > 0 && (
          <section className="card">
            <h2 className="card__title">Equipos por categoría</h2>
            <div className="card__body">
              {inv.by_category.map((c) => (
                <div key={c.name} className="bar-row">
                  <span>{c.name}</span>
                  <div className="bar"><div style={{ width: `${(c.total / maxCat) * 100}%` }} /></div>
                  <strong>{c.total}</strong>
                </div>
              ))}
            </div>
          </section>
        )}
        {summary && (
          <section className="card">
            <h2 className="card__title">Usuarios por perfil</h2>
            <div className="card__body">
              {summary.by_role.map((r) => (
                <div key={r.role} className="bar-row">
                  <span>{r.role}</span>
                  <div className="bar"><div style={{ width: `${(r.total / maxRole) * 100}%` }} /></div>
                  <strong>{r.total}</strong>
                </div>
              ))}
            </div>
          </section>
        )}
        <section className="card">
          <h2 className="card__title">Tus accesos <span className="tag">{user.role}</span></h2>
          <ul className="access">
            {modules.map((m) => (
              <li key={m.module}>
                <strong>{m.label}</strong>
                <span className="muted">{m.granted.length ? m.granted.map((a) => a.label).join(', ') : 'Sin acceso'}</span>
              </li>
            ))}
          </ul>
        </section>
      </div>
    </>
  );
}

function Stat({ icon, label, value }) {
  return (
    <div className="stat">
      <div className="stat__icon"><Icon name={icon} size={22} /></div>
      <div><span>{label}</span><strong>{value}</strong></div>
    </div>
  );
}
