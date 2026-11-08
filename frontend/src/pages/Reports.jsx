import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { useAlert } from '../context/AlertContext';
import { STATUS } from '../lib/format';

function Bars({ items, max }) {
  return items.map((i) => (
    <div key={i.label} className="bar-row">
      <span>{i.label}</span>
      <div className="bar"><div style={{ width: `${(i.total / max) * 100}%` }} /></div>
      <strong>{i.total}</strong>
    </div>
  ));
}

export default function Reports() {
  const alert = useAlert();
  const [users, setUsers] = useState(null);
  const [inv, setInv] = useState(null);

  useEffect(() => {
    Promise.all([api('/reports/users-summary'), api('/reports/inventory-summary')])
      .then(([u, i]) => { setUsers(u); setInv(i); })
      .catch((e) => alert.error('No se pudo cargar el reporte', e.message));
  }, []); // eslint-disable-line

  const roles = users ? users.by_role.map((r) => ({ label: r.role, total: r.total })) : [];
  const states = inv ? ['available', 'assigned', 'maintenance', 'retired'].map((k) => ({ label: STATUS[k].label, total: inv[k] })) : [];
  const cats = inv ? inv.by_category.map((c) => ({ label: c.name, total: c.total })) : [];
  const max = (list) => Math.max(1, ...list.map((i) => i.total));

  return (
    <>
      <div className="page-head">
        <div>
          <h1>Reportes</h1>
          <p className="muted">Estado del inventario de activos y de los usuarios del sistema.</p>
        </div>
      </div>
      <div className="grid-2 grid-2--top">
        <section className="card">
          <h2 className="card__title">Equipos por estado</h2>
          <div className="card__body">{inv ? <Bars items={states} max={max(states)} /> : <p className="muted">Cargando…</p>}</div>
        </section>
        <section className="card">
          <h2 className="card__title">Equipos por categoría</h2>
          <div className="card__body">
            {!inv ? <p className="muted">Cargando…</p> : cats.length ? <Bars items={cats} max={max(cats)} /> : <p className="muted">Aún no hay equipos registrados.</p>}
          </div>
        </section>
        <section className="card">
          <h2 className="card__title">Usuarios por perfil</h2>
          <div className="card__body">{users ? <Bars items={roles} max={max(roles)} /> : <p className="muted">Cargando…</p>}</div>
        </section>
      </div>
    </>
  );
}
