import { useState } from 'react';
import Icon from '../components/Icon';
import Modal from '../components/Modal';
import Field from '../components/Field';
import PageHead from '../components/PageHead';
import SearchBox from '../components/SearchBox';
import StatusPill from '../components/StatusPill';
import HistoryModal from '../components/HistoryModal';
import { useAuth } from '../context/AuthContext';
import { useDebounced, useResource } from '../hooks/useResource';
import { useCrud } from '../hooks/useCrud';
import { STATUS } from '../lib/format';

const CATEGORIES = ['Computador portátil', 'Computador de escritorio', 'Monitor', 'Teclado', 'Mouse', 'Teléfono',
  'Impresora', 'Escáner', 'Proyector', 'Tablet', 'Celular', 'Silla', 'Escritorio', 'Otro'];
const EMPTY = { name: '', category: '', brand: '', model: '', serial_number: '', purchase_date: '', cost: '', notes: '', quantity: 1 };

const T = {
  created: 'Equipos registrados', updated: 'Equipo actualizado', deleted: 'Equipo eliminado', saveError: 'No se pudo guardar el equipo',
  createdText: (list) => list.length === 1
    ? `Se creó ${list[0].code} y quedó disponible para asignar.`
    : `Se crearon ${list.length} unidades (${list[0].code} a ${list[list.length - 1].code}), cada una con su propio código.`,
  updatedText: 'Los cambios se guardaron.',
};

export default function Equipment() {
  const { can } = useAuth();
  const [q, setQ] = useState('');
  const [status, setStatus] = useState('');
  const dq = useDebounced(q);
  const canFind = can('equipment.find');
  const { rows, loading, reload } = useResource('/equipment', { q: dq, status });
  const { saving, save, remove } = useCrud('/equipment', reload, T);
  const [form, setForm] = useState(null);
  const [history, setHistory] = useState(null);
  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.value }));

  const submit = async (e) => {
    e.preventDefault();
    const { id, quantity, ...rest } = form;
    const body = { ...rest, cost: rest.cost === '' ? null : Number(rest.cost) };
    if (!id) body.quantity = Number(quantity); // la cantidad solo existe al crear
    if (await save(id, body)) setForm(null);
  };
  const edit = (r) => setForm({
    id: r.id, code: r.code, name: r.name, category: r.category, brand: r.brand, model: r.model, serial_number: r.serial_number,
    status: r.status, purchase_date: r.purchase_date, cost: r.cost ?? '', notes: r.notes,
  });
  const multiple = !form?.id && Number(form?.quantity) > 1;

  return (
    <>
      <PageHead title="Equipos de oficina" subtitle="Inventario de activos: cada unidad tiene su código, su estado y quién la tiene hoy.">
        {canFind && <SearchBox value={q} onChange={setQ} placeholder="Buscar por código, equipo, serie o responsable" />}
        {canFind && (
          <div className="input input--select">
            <select value={status} onChange={(e) => setStatus(e.target.value)} aria-label="Filtrar por estado">
              <option value="">Todos los estados</option>
              {Object.entries(STATUS).map(([k, v]) => <option key={k} value={k}>{v.label}</option>)}
            </select>
          </div>
        )}
        <button className="btn btn--outline btn--icon" onClick={reload} aria-label="Actualizar"><Icon name="refresh" size={16} /></button>
        {can('equipment.create') && (
          <button className="btn btn--primary" onClick={() => setForm({ ...EMPTY })}><Icon name="plus" size={16} />Nuevo equipo</button>
        )}
      </PageHead>

      <div className="card card--table">
        <table className="table">
          <thead><tr><th>Código</th><th>Equipo</th><th>Serie</th><th>Estado</th><th>Asignado a</th><th className="col-actions" /></tr></thead>
          <tbody>
            {loading && <tr><td colSpan="6" className="muted pad">Cargando…</td></tr>}
            {!loading && rows.length === 0 && <tr><td colSpan="6" className="muted pad">No hay equipos para mostrar.</td></tr>}
            {rows.map((r) => (
              <tr key={r.id}>
                <td><code className="code">{r.code}</code></td>
                <td>
                  <strong>{r.name}</strong>
                  <small className="sub">{[r.category, r.brand, r.model].filter(Boolean).join(' · ')}</small>
                </td>
                <td>{r.serial_number || '—'}</td>
                <td><StatusPill status={r.status} /></td>
                <td>{r.holder || '—'}</td>
                <td className="col-actions">
                  {can('assignments.list') && (
                    <button className="icon-btn" aria-label={`Historial de ${r.code}`} onClick={() => setHistory(r)}><Icon name="clock" size={17} /></button>
                  )}
                  {can('equipment.update') && (
                    <button className="icon-btn" aria-label={`Editar ${r.code}`} onClick={() => edit(r)}><Icon name="edit" size={17} /></button>
                  )}
                  {can('equipment.delete') && (
                    <button className="icon-btn icon-btn--danger" aria-label={`Eliminar ${r.code}`}
                            onClick={() => remove(r, { title: `¿Eliminar ${r.code}?`, text: `${r.name}. Esta acción no se puede deshacer.` })}><Icon name="trash" size={17} /></button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Modal open={!!form} onClose={() => setForm(null)} width={680}
             title={form?.id ? `Editar ${form.code}` : 'Nuevo equipo'}>
        {form && (
          <form onSubmit={submit}>
            <div className="form-grid">
              <Field label="Nombre del equipo" className="span-2">
                <input required autoFocus value={form.name} onChange={set('name')} placeholder="Ej. Portátil Latitude 5440" />
              </Field>
              <Field label="Categoría">
                <input list="categories" value={form.category} onChange={set('category')} placeholder="Escribe o elige" />
                <datalist id="categories">{CATEGORIES.map((c) => <option key={c} value={c} />)}</datalist>
              </Field>
              <Field label="Marca"><input value={form.brand} onChange={set('brand')} /></Field>
              <Field label="Modelo"><input value={form.model} onChange={set('model')} /></Field>
              <Field label="Número de serie" hint={multiple ? 'Se completa después en cada unidad.' : undefined}>
                <input value={multiple ? '' : form.serial_number} disabled={multiple} onChange={set('serial_number')} />
              </Field>
              <Field label="Fecha de compra"><input type="date" value={form.purchase_date} onChange={set('purchase_date')} /></Field>
              <Field label="Costo unitario"><input type="number" min="0" step="0.01" value={form.cost} onChange={set('cost')} /></Field>
              {form.id && (
                <Field label="Estado" hint={form.status === 'assigned' ? 'Para cambiarlo, devuélvelo desde Asignaciones.' : undefined}>
                  <select value={form.status} onChange={set('status')} disabled={form.status === 'assigned'}>
                    {form.status === 'assigned' && <option value="assigned">Asignado</option>}
                    <option value="available">Disponible</option>
                    <option value="maintenance">En mantenimiento</option>
                    <option value="retired">Dado de baja</option>
                  </select>
                </Field>
              )}
              {!form.id && (
                <Field label="Cantidad" hint="Se crea una unidad por cada una, con código propio.">
                  <input type="number" min="1" max="200" required value={form.quantity} onChange={set('quantity')} />
                </Field>
              )}
              <Field label="Notas" area className="span-2"><textarea rows={2} value={form.notes} onChange={set('notes')} /></Field>
            </div>
            <div className="modal__foot">
              <button type="button" className="btn btn--outline" onClick={() => setForm(null)}>Cancelar</button>
              <button className="btn btn--primary" disabled={saving}>
                {saving ? 'Guardando…' : form.id ? 'Guardar cambios' : `Registrar ${Number(form.quantity) > 1 ? `${form.quantity} equipos` : 'equipo'}`}
              </button>
            </div>
          </form>
        )}
      </Modal>

      {history && <HistoryModal title={`Historial de ${history.code} · ${history.name}`} query={{ equipment_id: history.id }} onClose={() => setHistory(null)} />}
    </>
  );
}
