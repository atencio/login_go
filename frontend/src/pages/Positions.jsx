import { useState } from 'react';
import Icon from '../components/Icon';
import Modal from '../components/Modal';
import Field from '../components/Field';
import PageHead from '../components/PageHead';
import SearchBox from '../components/SearchBox';
import { useAuth } from '../context/AuthContext';
import { useDebounced, useResource } from '../hooks/useResource';
import { useCrud } from '../hooks/useCrud';

const EMPTY = { name: '', description: '', active: true };
const T = {
  created: 'Cargo creado', updated: 'Cargo actualizado', deleted: 'Cargo eliminado', saveError: 'No se pudo guardar el cargo',
  createdText: 'Ya puedes asignarlo a funcionarios.', updatedText: 'Los cambios se guardaron.',
};

export default function Positions() {
  const { can } = useAuth();
  const [q, setQ] = useState('');
  const dq = useDebounced(q);
  const { rows, loading, reload } = useResource('/positions', { q: dq });
  const { saving, save, remove } = useCrud('/positions', reload, T);
  const [form, setForm] = useState(null);
  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.type === 'checkbox' ? e.target.checked : e.target.value }));

  const submit = async (e) => {
    e.preventDefault();
    const { id, ...body } = form;
    if (await save(id, body)) setForm(null);
  };

  return (
    <>
      <PageHead title="Cargos" subtitle="Puestos de trabajo a los que pertenecen los funcionarios.">
        {can('positions.find') && <SearchBox value={q} onChange={setQ} placeholder="Buscar cargo" />}
        <button className="btn btn--outline btn--icon" onClick={reload} aria-label="Actualizar"><Icon name="refresh" size={16} /></button>
        {can('positions.create') && (
          <button className="btn btn--primary" onClick={() => setForm({ ...EMPTY })}><Icon name="plus" size={16} />Nuevo cargo</button>
        )}
      </PageHead>

      <div className="card card--table">
        <table className="table">
          <thead><tr><th>Cargo</th><th>Funcionarios</th><th>Estado</th><th className="col-actions" /></tr></thead>
          <tbody>
            {loading && <tr><td colSpan="4" className="muted pad">Cargando…</td></tr>}
            {!loading && rows.length === 0 && <tr><td colSpan="4" className="muted pad">No hay cargos para mostrar.</td></tr>}
            {rows.map((p) => (
              <tr key={p.id}>
                <td><strong>{p.name}</strong><small className="sub">{p.description}</small></td>
                <td>{p.employee_count}</td>
                <td><span className={`status ${p.active ? 'status--on' : 'status--off'}`}>{p.active ? 'Activo' : 'Inactivo'}</span></td>
                <td className="col-actions">
                  {can('positions.update') && (
                    <button className="icon-btn" aria-label={`Editar ${p.name}`}
                            onClick={() => setForm({ id: p.id, name: p.name, description: p.description, active: p.active })}><Icon name="edit" size={17} /></button>
                  )}
                  {can('positions.delete') && (
                    <button className="icon-btn icon-btn--danger" aria-label={`Eliminar ${p.name}`}
                            onClick={() => remove(p, { title: `¿Eliminar el cargo ${p.name}?` })}><Icon name="trash" size={17} /></button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={!!form} onClose={() => setForm(null)} title={form?.id ? 'Editar cargo' : 'Nuevo cargo'}>
        {form && (
          <form onSubmit={submit}>
            <Field label="Nombre del cargo"><input required autoFocus maxLength={80} value={form.name} onChange={set('name')} /></Field>
            <Field label="Descripción" area><textarea rows={3} value={form.description} onChange={set('description')} /></Field>
            <label className="check"><input type="checkbox" checked={form.active} onChange={set('active')} />Cargo activo</label>
            <div className="modal__foot">
              <button type="button" className="btn btn--outline" onClick={() => setForm(null)}>Cancelar</button>
              <button className="btn btn--primary" disabled={saving}>{saving ? 'Guardando…' : 'Guardar'}</button>
            </div>
          </form>
        )}
      </Modal>
    </>
  );
}
