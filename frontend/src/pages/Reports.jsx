import { useEffect, useState } from 'react';
import { api } from '../api/client';
import { useAlert } from '../context/AlertContext';

export default function Reports() {
  const alert = useAlert();
  const [summary, setSummary] = useState(null);

  useEffect(() => {
    api('/reports/users-summary').then(setSummary).catch((e) => alert.error('No se pudo cargar el reporte', e.message));
  }, []); // eslint-disable-line

  const max = summary ? Math.max(1, ...summary.by_role.map((r) => r.total)) : 1;

  return (
    <>
      <div className="page-head">
        <div>
          <h1>Reportes</h1>
          <p className="muted">Usuarios agrupados por perfil.</p>
        </div>
      </div>
      <div className="card">
        {!summary ? <p className="muted">Cargando…</p> : summary.by_role.map((r) => (
          <div key={r.role} className="bar-row">
            <span>{r.role}</span>
            <div className="bar"><div style={{ width: `${(r.total / max) * 100}%` }} /></div>
            <strong>{r.total}</strong>
          </div>
        ))}
      </div>
    </>
  );
}
