import { useState } from 'react';
import TableCard from '../components/TableCard';
import PageSizeSelect from '../components/PageSizeSelect';
import Icon from '../components/Icon';
import Modal from '../components/Modal';
import Field from '../components/Field';
import PageHead from '../components/PageHead';
import SearchBox from '../components/SearchBox';
import { useAuth } from '../context/AuthContext';
import { useDebounced, useResource } from '../hooks/useResource';
import { useCrud } from '../hooks/useCrud';

// CRUD de un catálogo simple (categorías, marcas). cfg define textos, ruta de la API y prefijo de permisos.
export default function CatalogPage({ cfg }) {
  const { can } = useAuth();
  const [q, setQ] = useState('');
  const dq = useDebounced(q);
  const { rows, loading, reload } = useResource(cfg.path, { q: dq });
  const { saving, save, remove } = useCrud(cfg.path, reload, cfg.texts);
  const [form, setForm] = useState(null);

  const submit = async (e) => {
    e.preventDefault();
    const { id, ...body } = form;
    if (await save(id, body)) setForm(null);
  };

  return (
    <div className="fill">
      <PageHead title={cfg.title} subtitle={cfg.subtitle}>
        <PageSizeSelect />
        {can(`${cfg.perm}.find`) && <SearchBox value={q} onChange={setQ} placeholder={`Buscar ${cfg.singular}`} />}
        <button className="btn btn--outline btn--icon" onClick={reload} aria-label="Actualizar"><Icon name="refresh" size={16} /></button>
        {can(`${cfg.perm}.create`) && (
          <button className="btn btn--primary" onClick={() => setForm({ name: '', active: true })}>
            <Icon name="plus" size={16} />Nueva {cfg.singular}
          </button>
        )}
      </PageHead>

      <TableCard rows={rows} loading={loading} cols={4} empty="No hay registros para mostrar."
        head={<tr><th>Nombre</th><th>Equipos</th><th>Estado</th><th className="col-actions" /></tr>}
        renderRow={(r) => (
          <tr key={r.id}>
                <td><strong>{r.name}</strong></td>
                <td>{r.equipment_count}</td>
                <td><span className={`status ${r.active ? 'status--on' : 'status--off'}`}>{r.active ? 'Activa' : 'Inactiva'}</span></td>
                <td className="col-actions">
                  {can(`${cfg.perm}.update`) && (
                    <button className="icon-btn" aria-label={`Editar ${r.name}`}
                            onClick={() => setForm({ id: r.id, name: r.name, active: r.active })}><Icon name="edit" size={17} /></button>
                  )}
                  {can(`${cfg.perm}.delete`) && (
                    <button className="icon-btn icon-btn--danger" aria-label={`Eliminar ${r.name}`}
                            onClick={() => remove(r, { title: `¿Eliminar la ${cfg.singular} ${r.name}?` })}><Icon name="trash" size={17} /></button>
                  )}
                </td>
              </tr>
        )} />

      <Modal open={!!form} onClose={() => setForm(null)} width={440}
             title={form?.id ? `Editar ${cfg.singular}` : `Nueva ${cfg.singular}`}>
        {form && (
          <form onSubmit={submit}>
            <Field label="Nombre">
              <input required autoFocus maxLength={80} value={form.name} placeholder={cfg.example}
                     onChange={(e) => setForm({ ...form, name: e.target.value })} />
            </Field>
            <label className="check">
              <input type="checkbox" checked={form.active} onChange={(e) => setForm({ ...form, active: e.target.checked })} />
              {cfg.activeLabel}
            </label>
            <div className="modal__foot">
              <button type="button" className="btn btn--outline" onClick={() => setForm(null)}>Cancelar</button>
              <button className="btn btn--primary" disabled={saving}>{saving ? 'Guardando…' : 'Guardar'}</button>
            </div>
          </form>
        )}
      </Modal>
    </div>
  );
}
