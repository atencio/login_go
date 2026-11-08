import { useState } from 'react';
import Icon from '../components/Icon';
import Modal from '../components/Modal';
import Field from '../components/Field';
import PageHead from '../components/PageHead';
import SearchBox from '../components/SearchBox';
import HistoryModal from '../components/HistoryModal';
import { initials } from '../components/ProfileModal';
import { useAuth } from '../context/AuthContext';
import { useDebounced, useResource } from '../hooks/useResource';
import { useCrud } from '../hooks/useCrud';

const EMPTY = { first_name: '', last_name: '', email: '', document: '', phone: '', department: '', position_id: '', hire_date: '', active: true };
const T = {
  created: 'Funcionario creado', updated: 'Funcionario actualizado', deleted: 'Funcionario eliminado',
  saveError: 'No se pudo guardar el funcionario', createdText: 'Ya puedes asignarle equipos.', updatedText: 'Los cambios se guardaron.',
};

export default function Employees() {
  const { can } = useAuth();
  const [q, setQ] = useState('');
  const [positionId, setPositionId] = useState('');
  const dq = useDebounced(q);
  const canFind = can('employees.find');
  const canOptions = can('employees.create') || can('employees.update') || can('positions.list');
  const { rows, loading, reload } = useResource('/employees', { q: dq, position_id: positionId });
  const { rows: positions } = useResource('/positions/options', {}, { enabled: canOptions });
  const { saving, save, remove } = useCrud('/employees', reload, T);
  const [form, setForm] = useState(null);
  const [history, setHistory] = useState(null);
  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.type === 'checkbox' ? e.target.checked : e.target.value }));

  const submit = async (e) => {
    e.preventDefault();
    const { id, ...body } = form;
    if (await save(id, { ...body, position_id: Number(body.position_id) })) setForm(null);
  };
  const edit = (r) => setForm({
    id: r.id, first_name: r.first_name, last_name: r.last_name, email: r.email, document: r.document, phone: r.phone,
    department: r.department, position_id: r.position_id, hire_date: r.hire_date, active: r.active,
  });
  const selectable = positions.filter((p) => p.active || p.id === Number(form?.position_id));

  return (
    <>
      <PageHead title="Funcionarios" subtitle="Personas a las que se les asignan los equipos de oficina.">
        {canFind && <SearchBox value={q} onChange={setQ} placeholder="Buscar por nombre, correo o documento" />}
        {canFind && canOptions && (
          <div className="input input--select">
            <select value={positionId} onChange={(e) => setPositionId(e.target.value)} aria-label="Filtrar por cargo">
              <option value="">Todos los cargos</option>
              {positions.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
            </select>
          </div>
        )}
        <button className="btn btn--outline btn--icon" onClick={reload} aria-label="Actualizar"><Icon name="refresh" size={16} /></button>
        {can('employees.create') && (
          <button className="btn btn--primary" onClick={() => setForm({ ...EMPTY, position_id: selectable[0]?.id ?? '' })}>
            <Icon name="plus" size={16} />Nuevo funcionario
          </button>
        )}
      </PageHead>

      <div className="card card--table">
        <table className="table">
          <thead><tr><th>Funcionario</th><th>Cargo</th><th>Departamento</th><th>Equipos</th><th>Estado</th><th className="col-actions" /></tr></thead>
          <tbody>
            {loading && <tr><td colSpan="6" className="muted pad">Cargando…</td></tr>}
            {!loading && rows.length === 0 && <tr><td colSpan="6" className="muted pad">No hay funcionarios para mostrar.</td></tr>}
            {rows.map((r) => (
              <tr key={r.id}>
                <td>
                  <div className="cell-user">
                    <span className="avatar">{initials(r.full_name)}</span>
                    <div><strong>{r.full_name}</strong><small>{r.email}</small></div>
                  </div>
                </td>
                <td><span className="badge">{r.position}</span></td>
                <td>{r.department || '—'}</td>
                <td>{r.active_assignments}</td>
                <td><span className={`status ${r.active ? 'status--on' : 'status--off'}`}>{r.active ? 'Activo' : 'Inactivo'}</span></td>
                <td className="col-actions">
                  {can('assignments.list') && (
                    <button className="icon-btn" aria-label={`Historial de ${r.full_name}`} onClick={() => setHistory(r)}><Icon name="clock" size={17} /></button>
                  )}
                  {can('employees.update') && (
                    <button className="icon-btn" aria-label={`Editar ${r.full_name}`} onClick={() => edit(r)}><Icon name="edit" size={17} /></button>
                  )}
                  {can('employees.delete') && (
                    <button className="icon-btn icon-btn--danger" aria-label={`Eliminar ${r.full_name}`}
                            onClick={() => remove(r, { title: `¿Eliminar a ${r.full_name}?` })}><Icon name="trash" size={17} /></button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={!!form} onClose={() => setForm(null)} title={form?.id ? 'Editar funcionario' : 'Nuevo funcionario'} width={640}>
        {form && (
          <form onSubmit={submit}>
            <div className="form-grid">
              <Field label="Nombre"><input required autoFocus value={form.first_name} onChange={set('first_name')} /></Field>
              <Field label="Apellido"><input required value={form.last_name} onChange={set('last_name')} /></Field>
              <Field label="Correo" className="span-2"><input type="email" required value={form.email} onChange={set('email')} /></Field>
              <Field label="Documento de identidad"><input value={form.document} onChange={set('document')} /></Field>
              <Field label="Teléfono"><input value={form.phone} onChange={set('phone')} /></Field>
              <Field label="Cargo">
                <select required value={form.position_id} onChange={set('position_id')}>
                  <option value="" disabled>Selecciona un cargo</option>
                  {selectable.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                </select>
              </Field>
              <Field label="Departamento / área"><input value={form.department} onChange={set('department')} /></Field>
              <Field label="Fecha de ingreso"><input type="date" value={form.hire_date} onChange={set('hire_date')} /></Field>
            </div>
            <label className="check"><input type="checkbox" checked={form.active} onChange={set('active')} />Funcionario activo</label>
            <div className="modal__foot">
              <button type="button" className="btn btn--outline" onClick={() => setForm(null)}>Cancelar</button>
              <button className="btn btn--primary" disabled={saving}>{saving ? 'Guardando…' : 'Guardar'}</button>
            </div>
          </form>
        )}
      </Modal>

      {history && <HistoryModal title={`Historial de ${history.full_name}`} query={{ employee_id: history.id }} onClose={() => setHistory(null)} />}
    </>
  );
}
