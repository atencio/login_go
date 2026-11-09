import { useCallback, useEffect, useState } from 'react';
import TableCard from '../components/TableCard';
import PageSizeSelect from '../components/PageSizeSelect';
import PageHead from '../components/PageHead';
import { api } from '../api/client';
import { loadCatalog } from '../api/catalog';
import Icon from '../components/Icon';
import Modal from '../components/Modal';
import PermissionMatrix from '../components/PermissionMatrix';
import { useAuth } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';

export default function Roles() {
  const { can } = useAuth();
  const alert = useAlert();
  const [roles, setRoles] = useState([]);
  const [catalog, setCatalog] = useState([]);
  const [loading, setLoading] = useState(true);
  const [editor, setEditor] = useState(null); // { id?, name, permissions, readOnly }
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    try {
      const [r, c] = await Promise.all([api('/roles'), loadCatalog()]);
      setRoles(r); setCatalog(c);
    } catch (e) {
      alert.error('No se pudo cargar la lista', e.message);
    } finally {
      setLoading(false);
    }
  }, []); // eslint-disable-line

  useEffect(() => { load(); }, [load]);

  const openCreate = () => setEditor({ name: '', permissions: [], readOnly: false });
  const openEdit = (r) => setEditor({ id: r.id, name: r.name, permissions: [...r.permissions], readOnly: false });
  const openView = async (r) => {
    try {
      const detail = await api(`/roles/${r.id}`); // requiere el permiso "Encontrar"
      setEditor({ id: detail.id, name: detail.name, permissions: detail.permissions, readOnly: true });
    } catch (e) {
      alert.error('No se pudo abrir el perfil', e.message);
    }
  };

  const toggle = (code) => setEditor((ed) => ({
    ...ed,
    permissions: ed.permissions.includes(code) ? ed.permissions.filter((c) => c !== code) : [...ed.permissions, code],
  }));
  const setMany = (codes, on) => setEditor((ed) => ({
    ...ed,
    permissions: on ? [...new Set([...ed.permissions, ...codes])] : ed.permissions.filter((c) => !codes.includes(c)),
  }));
  const setAll = (on) => setEditor((ed) => ({ ...ed, permissions: on ? catalog.map((p) => p.code) : [] }));

  const save = async (e) => {
    e.preventDefault();
    setSaving(true);
    const { id, name, permissions } = editor;
    try {
      await api(id ? `/roles/${id}` : '/roles', { method: id ? 'PUT' : 'POST', body: { name, permissions } });
      setEditor(null);
      await load();
      alert.success(id ? 'Perfil actualizado' : 'Perfil creado', id ? 'Los permisos nuevos se aplican de inmediato.' : `El perfil ${name} ya se puede asignar a usuarios.`);
    } catch (err) {
      alert.error('No se pudo guardar el perfil', err.message);
    } finally {
      setSaving(false);
    }
  };

  const remove = async (r) => {
    const ok = await alert.confirm({
      title: `¿Eliminar el perfil ${r.name}?`, text: 'Esta acción no se puede deshacer.',
      confirmText: 'Eliminar perfil', danger: true,
    });
    if (!ok) return;
    try {
      await api(`/roles/${r.id}`, { method: 'DELETE' });
      await load();
      alert.success('Perfil eliminado');
    } catch (err) {
      alert.error('No se pudo eliminar', err.message);
    }
  };

  return (
    <div className="fill">
      <PageHead title="Perfiles" subtitle="Define qué puede hacer cada perfil en cada módulo.">
        <PageSizeSelect />

          <button className="btn btn--outline btn--icon" onClick={load} aria-label="Actualizar"><Icon name="refresh" size={16} /></button>
          {can('roles.create') && (
            <button className="btn btn--primary" onClick={openCreate}><Icon name="plus" size={16} />Nuevo perfil</button>
          )}
      </PageHead>

      <TableCard rows={roles} loading={loading} cols={4} empty="Aún no hay perfiles."
        head={<tr><th>Perfil</th><th>Permisos</th><th>Usuarios</th><th className="col-actions" /></tr>}
        renderRow={(r) => (
          <tr key={r.id}>
                <td>
                  <strong>{r.name}</strong>
                  {r.locked && <span className="tag">Protegido</span>}
                </td>
                <td>{r.permissions.length} de {catalog.length}</td>
                <td>{r.user_count}</td>
                <td className="col-actions">
                  {can('roles.find') && (
                    <button className="icon-btn" onClick={() => openView(r)} aria-label={`Ver permisos de ${r.name}`}><Icon name="eye" size={17} /></button>
                  )}
                  {can('roles.update') && !r.locked && (
                    <button className="icon-btn" onClick={() => openEdit(r)} aria-label={`Editar ${r.name}`}><Icon name="edit" size={17} /></button>
                  )}
                  {can('roles.delete') && !r.locked && (
                    <button className="icon-btn icon-btn--danger" onClick={() => remove(r)} aria-label={`Eliminar ${r.name}`}><Icon name="trash" size={17} /></button>
                  )}
                </td>
              </tr>
        )} />

      <Modal open={!!editor} onClose={() => setEditor(null)} width={1180}
             title={editor?.readOnly ? `Perfil ${editor.name}` : editor?.id ? 'Editar perfil' : 'Nuevo perfil'}>
        {editor && (
          <form onSubmit={save}>
            {!editor.readOnly && (
              <label className="field"><span>Nombre del perfil</span>
                <div className="input"><input required autoFocus maxLength={60} value={editor.name}
                       placeholder="Ej. Reclutador" onChange={(e) => setEditor({ ...editor, name: e.target.value })} /></div>
              </label>
            )}
            <PermissionMatrix catalog={catalog} value={editor.permissions} readOnly={editor.readOnly} onToggle={toggle} onToggleMany={setMany} />
            <div className="modal__foot modal__foot--split modal__foot--sticky">
              {!editor.readOnly ? (
                <div className="link-group">
                  <button type="button" className="link" onClick={() => setAll(true)}>Marcar todo</button>
                  <button type="button" className="link" onClick={() => setAll(false)}>Quitar todo</button>
                </div>
              ) : <span />}
              <div className="btn-group">
                <button type="button" className="btn btn--outline" onClick={() => setEditor(null)}>{editor.readOnly ? 'Cerrar' : 'Cancelar'}</button>
                {!editor.readOnly && <button className="btn btn--primary" disabled={saving}>{saving ? 'Guardando…' : 'Guardar perfil'}</button>}
              </div>
            </div>
          </form>
        )}
      </Modal>
    </div>
  );
}
