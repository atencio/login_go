import { useEffect, useState } from 'react';
import { api } from '../api/client';
import Icon from '../components/Icon';
import Modal from '../components/Modal';
import Field from '../components/Field';
import PageHead from '../components/PageHead';
import SearchBox from '../components/SearchBox';
import HistoryModal from '../components/HistoryModal';
import { useAuth } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';
import { useDebounced, useResource } from '../hooks/useResource';
import { fmtDateTime } from '../lib/format';

const equipmentLabel = (e) => `${e.code} · ${e.name}${e.brand || e.model ? ` — ${[e.brand, e.model].filter(Boolean).join(' ')}` : ''}`;

export default function Assignments() {
  const { can } = useAuth();
  const alert = useAlert();
  const [tab, setTab] = useState('active'); // active | history
  const [state, setState] = useState('');
  const [q, setQ] = useState('');
  const dq = useDebounced(q);
  const canFind = can('assignments.find');
  const { rows, loading, reload } = useResource('/assignments', { state: tab === 'active' ? 'active' : state, q: dq });
  const [modal, setModal] = useState(null); // { type: 'assign' | 'update' | 'return' | 'history', row? }

  const done = async (title, text) => {
    setModal(null);
    await reload();
    alert.success(title, text);
  };

  return (
    <>
      <PageHead title="Asignaciones" subtitle="Entrega, actualiza y devuelve equipos; cada movimiento queda en el historial.">
        {canFind && <SearchBox value={q} onChange={setQ} placeholder="Buscar por equipo, código o funcionario" />}
        {tab === 'history' && (
          <div className="input input--select">
            <select value={state} onChange={(e) => setState(e.target.value)} aria-label="Filtrar por estado">
              <option value="">Todos los movimientos</option>
              <option value="active">En uso</option>
              <option value="returned">Devueltos</option>
            </select>
          </div>
        )}
        <button className="btn btn--outline btn--icon" onClick={reload} aria-label="Actualizar"><Icon name="refresh" size={16} /></button>
        {can('assignments.create') && (
          <button className="btn btn--primary" onClick={() => setModal({ type: 'assign' })}><Icon name="plus" size={16} />Nueva asignación</button>
        )}
      </PageHead>

      <div className="tabs" role="tablist">
        <button role="tab" aria-selected={tab === 'active'} className={tab === 'active' ? 'tab tab--on' : 'tab'} onClick={() => setTab('active')}>Asignaciones activas</button>
        <button role="tab" aria-selected={tab === 'history'} className={tab === 'history' ? 'tab tab--on' : 'tab'} onClick={() => setTab('history')}>Historial</button>
      </div>

      <div className="card card--table">
        {tab === 'active' ? (
          <table className="table">
            <thead><tr><th>Equipo</th><th>Funcionario</th><th>Entregado</th><th>Notas</th><th className="col-actions col-actions--wide" /></tr></thead>
            <tbody>
              {loading && <tr><td colSpan="5" className="muted pad">Cargando…</td></tr>}
              {!loading && rows.length === 0 && <tr><td colSpan="5" className="muted pad">No hay equipos asignados en este momento.</td></tr>}
              {rows.map((a) => (
                <tr key={a.id}>
                  <td><code className="code">{a.equipment_code}</code><div>{a.equipment_name}</div></td>
                  <td>{a.employee}</td>
                  <td>{fmtDateTime(a.assigned_at)}<small className="sub">{a.assigned_by && `por ${a.assigned_by}`}</small></td>
                  <td className="notes">{a.delivery_notes || '—'}</td>
                  <td className="col-actions col-actions--wide">
                    <button className="icon-btn" aria-label={`Historial de ${a.equipment_code}`} onClick={() => setModal({ type: 'history', row: a })}><Icon name="clock" size={17} /></button>
                    {can('assignments.update') && (
                      <>
                        <button className="btn btn--outline btn--sm" onClick={() => setModal({ type: 'update', row: a })}>Actualizar</button>
                        <button className="btn btn--outline btn--sm" onClick={() => setModal({ type: 'return', row: a })}>Devolver</button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <table className="table">
            <thead><tr><th>Equipo</th><th>Funcionario</th><th>Entrega</th><th>Devolución</th><th>Notas</th></tr></thead>
            <tbody>
              {loading && <tr><td colSpan="5" className="muted pad">Cargando…</td></tr>}
              {!loading && rows.length === 0 && <tr><td colSpan="5" className="muted pad">No hay movimientos para mostrar.</td></tr>}
              {rows.map((a) => (
                <tr key={a.id}>
                  <td><code className="code">{a.equipment_code}</code><div>{a.equipment_name}</div></td>
                  <td>{a.employee}</td>
                  <td>{fmtDateTime(a.assigned_at)}<small className="sub">{a.assigned_by && `por ${a.assigned_by}`}</small></td>
                  <td>
                    {a.active
                      ? <span className="pill pill--info">En uso</span>
                      : <>{fmtDateTime(a.returned_at)}<small className="sub">{a.received_by && `por ${a.received_by}`}</small></>}
                  </td>
                  <td className="notes">{[a.delivery_notes, a.return_notes].filter(Boolean).join(' · ') || '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {modal?.type === 'assign' && <AssignModal onClose={() => setModal(null)} onSaved={done} />}
      {modal?.type === 'update' && <UpdateModal row={modal.row} onClose={() => setModal(null)} onSaved={done} />}
      {modal?.type === 'return' && <ReturnModal row={modal.row} onClose={() => setModal(null)} onSaved={done} />}
      {modal?.type === 'history' && (
        <HistoryModal title={`Historial de ${modal.row.equipment_code} · ${modal.row.equipment_name}`}
                      query={{ equipment_id: modal.row.equipment_id }} onClose={() => setModal(null)} />
      )}
    </>
  );
}

function AssignModal({ onClose, onSaved }) {
  const alert = useAlert();
  const [equipment, setEquipment] = useState(null);
  const [employees, setEmployees] = useState(null);
  const [form, setForm] = useState({ equipment_id: '', employee_id: '', notes: '' });
  const [saving, setSaving] = useState(false);
  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.value }));

  useEffect(() => {
    Promise.all([api('/equipment/options?status=available'), api('/employees/options')])
      .then(([eq, emp]) => { setEquipment(eq); setEmployees(emp); })
      .catch(async (e) => { onClose(); await alert.error('No se pudieron cargar las opciones', e.message); });
  }, []); // eslint-disable-line

  const submit = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      const res = await api('/assignments', {
        method: 'POST',
        body: { equipment_id: Number(form.equipment_id), employee_id: Number(form.employee_id), notes: form.notes },
      });
      onSaved('Equipo asignado', `${res.equipment_code} quedó a cargo de ${res.employee}.`);
    } catch (err) {
      setSaving(false);
      alert.error('No se pudo asignar', err.message);
    }
  };

  const ready = equipment && employees;
  const none = ready && equipment.length === 0;
  return (
    <Modal open onClose={onClose} title="Nueva asignación" width={560}>
      {!ready ? <p className="muted">Cargando…</p> : (
        <form onSubmit={submit}>
          {none && <p className="note">No hay equipos disponibles. Registra equipos nuevos o devuelve alguno desde las asignaciones activas.</p>}
          <Field label="Equipo disponible">
            <select required value={form.equipment_id} onChange={set('equipment_id')} disabled={none}>
              <option value="" disabled>Selecciona un equipo</option>
              {equipment.map((e) => <option key={e.id} value={e.id}>{equipmentLabel(e)}</option>)}
            </select>
          </Field>
          <Field label="Funcionario">
            <select required value={form.employee_id} onChange={set('employee_id')}>
              <option value="" disabled>Selecciona un funcionario</option>
              {employees.map((p) => <option key={p.id} value={p.id}>{p.name} · {p.position}</option>)}
            </select>
          </Field>
          <Field label="Notas de entrega (opcional)" area>
            <textarea rows={3} value={form.notes} onChange={set('notes')} placeholder="Estado, accesorios entregados…" />
          </Field>
          <div className="modal__foot">
            <button type="button" className="btn btn--outline" onClick={onClose}>Cancelar</button>
            <button className="btn btn--primary" disabled={saving || none}>{saving ? 'Asignando…' : 'Asignar equipo'}</button>
          </div>
        </form>
      )}
    </Modal>
  );
}

function UpdateModal({ row, onClose, onSaved }) {
  const alert = useAlert();
  const [employees, setEmployees] = useState(null);
  const [form, setForm] = useState({ employee_id: row.employee_id, notes: row.delivery_notes });
  const [saving, setSaving] = useState(false);
  const transfer = Number(form.employee_id) !== row.employee_id;

  useEffect(() => {
    api('/employees/options').then(setEmployees)
      .catch(async (e) => { onClose(); await alert.error('No se pudieron cargar los funcionarios', e.message); });
  }, []); // eslint-disable-line

  const submit = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      const res = await api(`/assignments/${row.id}`, { method: 'PUT', body: { employee_id: Number(form.employee_id), notes: form.notes } });
      onSaved(transfer ? 'Equipo transferido' : 'Asignación actualizada',
        transfer ? `${res.equipment_code} pasó de ${row.employee} a ${res.employee}. Ambos movimientos quedaron en el historial.` : 'Los cambios se guardaron.');
    } catch (err) {
      setSaving(false);
      alert.error('No se pudo actualizar', err.message);
    }
  };

  return (
    <Modal open onClose={onClose} title="Actualizar asignación" width={560}>
      {!employees ? <p className="muted">Cargando…</p> : (
        <form onSubmit={submit}>
          <p className="summary"><code className="code">{row.equipment_code}</code> {row.equipment_name}</p>
          <Field label="Funcionario responsable" hint={transfer ? 'Se registrará la devolución al funcionario actual y una nueva entrega.' : undefined}>
            <select required value={form.employee_id} onChange={(e) => setForm({ ...form, employee_id: e.target.value })}>
              {employees.map((p) => <option key={p.id} value={p.id}>{p.name} · {p.position}</option>)}
              {!employees.some((p) => p.id === row.employee_id) && <option value={row.employee_id}>{row.employee} (inactivo)</option>}
            </select>
          </Field>
          <Field label="Notas de entrega" area>
            <textarea rows={3} value={form.notes} onChange={(e) => setForm({ ...form, notes: e.target.value })} />
          </Field>
          <div className="modal__foot">
            <button type="button" className="btn btn--outline" onClick={onClose}>Cancelar</button>
            <button className="btn btn--primary" disabled={saving}>{saving ? 'Guardando…' : transfer ? 'Transferir equipo' : 'Guardar cambios'}</button>
          </div>
        </form>
      )}
    </Modal>
  );
}

function ReturnModal({ row, onClose, onSaved }) {
  const alert = useAlert();
  const [form, setForm] = useState({ status: 'available', notes: '' });
  const [saving, setSaving] = useState(false);

  const submit = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api(`/assignments/${row.id}/return`, { method: 'POST', body: form });
      onSaved('Equipo devuelto', `${row.equipment_code} fue recibido de ${row.employee}.`);
    } catch (err) {
      setSaving(false);
      alert.error('No se pudo registrar la devolución', err.message);
    }
  };

  return (
    <Modal open onClose={onClose} title="Devolver equipo" width={520}>
      <form onSubmit={submit}>
        <p className="summary"><code className="code">{row.equipment_code}</code> {row.equipment_name} — {row.employee}</p>
        <Field label="Estado del equipo al recibirlo">
          <select value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })}>
            <option value="available">Disponible para asignar</option>
            <option value="maintenance">Enviar a mantenimiento</option>
            <option value="retired">Dar de baja</option>
          </select>
        </Field>
        <Field label="Notas de devolución (opcional)" area>
          <textarea rows={3} value={form.notes} onChange={(e) => setForm({ ...form, notes: e.target.value })} placeholder="Daños, faltantes, observaciones…" />
        </Field>
        <div className="modal__foot">
          <button type="button" className="btn btn--outline" onClick={onClose}>Cancelar</button>
          <button className="btn btn--primary" disabled={saving}>{saving ? 'Registrando…' : 'Registrar devolución'}</button>
        </div>
      </form>
    </Modal>
  );
}
