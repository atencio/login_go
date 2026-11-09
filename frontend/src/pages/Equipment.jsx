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

const MAX_BATCH = 200;
const EMPTY = {
  name: '', category_id: '', brand_id: '', model: '', purchase_date: '', cost: '', notes: '',
  quantity: 1, serials: [''],
};

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
  const canOptions = can('equipment.create') || can('equipment.update');
  const { rows, loading, reload } = useResource('/equipment', { q: dq, status });
  const { rows: categories } = useResource('/categories/options', {}, { enabled: canOptions });
  const { rows: brands } = useResource('/brands/options', {}, { enabled: canOptions });
  const { saving, save, remove } = useCrud('/equipment', reload, T);
  const [form, setForm] = useState(null);
  const [history, setHistory] = useState(null);
  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.value }));

  // Cantidad efectiva y una serie por cada unidad (la serie es única por producto)
  const n = Math.min(MAX_BATCH, Math.max(1, Math.floor(Number(form?.quantity)) || 1));
  const setSerial = (i, v) => setForm((f) => {
    const s = [...f.serials];
    while (s.length <= i) s.push('');
    s[i] = v;
    return { ...f, serials: s };
  });

  const submit = async (e) => {
    e.preventDefault();
    const { id, quantity, serials, ...rest } = form;
    const body = {
      ...rest,
      category_id: Number(rest.category_id),
      brand_id: rest.brand_id === '' ? 0 : Number(rest.brand_id),
      cost: rest.cost === '' ? null : Number(rest.cost),
    };
    if (!id) {
      body.quantity = n; // la cantidad solo existe al crear
      body.serials = Array.from({ length: n }, (_, i) => (serials[i] || '').trim());
    }
    if (await save(id, body)) setForm(null);
  };

  const edit = (r) => setForm({
    id: r.id, code: r.code, name: r.name, category_id: r.category_id || '', brand_id: r.brand_id || '', model: r.model,
    serial_number: r.serial_number, status: r.status, purchase_date: r.purchase_date, cost: r.cost ?? '', notes: r.notes,
  });

  const categoryOptions = categories.filter((c) => c.active || c.id === Number(form?.category_id));
  const brandOptions = brands.filter((b) => b.active || b.id === Number(form?.brand_id));

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
                <select required value={form.category_id} onChange={set('category_id')}>
                  <option value="" disabled>Selecciona una categoría</option>
                  {categoryOptions.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
                </select>
              </Field>
              <Field label="Marca">
                <select value={form.brand_id} onChange={set('brand_id')}>
                  <option value="">Sin marca</option>
                  {brandOptions.map((b) => <option key={b.id} value={b.id}>{b.name}</option>)}
                </select>
              </Field>
              <Field label="Modelo"><input value={form.model} onChange={set('model')} /></Field>
              <Field label="Fecha de compra"><input type="date" value={form.purchase_date} onChange={set('purchase_date')} /></Field>
              <Field label="Costo unitario"><input type="number" min="0" step="0.01" value={form.cost} onChange={set('cost')} /></Field>

              {form.id ? (
                <>
                  <Field label="Número de serie"><input value={form.serial_number} onChange={set('serial_number')} /></Field>
                  <Field label="Estado" hint={form.status === 'assigned' ? 'Para cambiarlo, devuélvelo desde Asignaciones.' : undefined}>
                    <select value={form.status} onChange={set('status')} disabled={form.status === 'assigned'}>
                      {form.status === 'assigned' && <option value="assigned">Asignado</option>}
                      <option value="available">Disponible</option>
                      <option value="maintenance">En mantenimiento</option>
                      <option value="retired">Dado de baja</option>
                    </select>
                  </Field>
                </>
              ) : (
                <>
                  <Field label="Cantidad" hint="Se crea una unidad por cada una, con código propio.">
                    <input type="number" min="1" max={MAX_BATCH} required value={form.quantity} onChange={set('quantity')} />
                  </Field>
                  <fieldset className="serials span-2">
                    <legend>{n === 1 ? 'Número de serie' : `Números de serie (${n} unidades)`}</legend>
                    <p className="hint">La serie es única por equipo. Puedes dejar vacías las que aún no tengas y completarlas después.</p>
                    <div className={n === 1 ? 'serials__grid serials__grid--one' : 'serials__grid'}>
                      {Array.from({ length: n }, (_, i) => (
                        <div key={i} className="input">
                          {n > 1 && <span className="serials__n">{i + 1}</span>}
                          <input aria-label={`Serie de la unidad ${i + 1}`} placeholder="Número de serie"
                                 value={form.serials[i] || ''} onChange={(e) => setSerial(i, e.target.value)} />
                        </div>
                      ))}
                    </div>
                  </fieldset>
                </>
              )}
              <Field label="Notas" area className="span-2"><textarea rows={2} value={form.notes} onChange={set('notes')} /></Field>
            </div>
            <div className="modal__foot">
              <button type="button" className="btn btn--outline" onClick={() => setForm(null)}>Cancelar</button>
              <button className="btn btn--primary" disabled={saving}>
                {saving ? 'Guardando…' : form.id ? 'Guardar cambios' : n > 1 ? `Registrar ${n} equipos` : 'Registrar equipo'}
              </button>
            </div>
          </form>
        )}
      </Modal>

      {history && <HistoryModal title={`Historial de ${history.code} · ${history.name}`} query={{ equipment_id: history.id }} onClose={() => setHistory(null)} />}
    </>
  );
}
