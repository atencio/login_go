import { useCallback, useEffect, useRef, useState } from 'react';
import TableCard from '../components/TableCard';
import PageSizeSelect from '../components/PageSizeSelect';
import PageHead from '../components/PageHead';
import { api } from '../api/client';
import Icon from '../components/Icon';
import Modal from '../components/Modal';
import { useAuth } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';
import { initials } from '../components/ProfileModal';

const EMPTY = { name: '', email: '', password: '', role_id: '', active: true };

export default function Users() {
  const { can, user: me } = useAuth();
  const alert = useAlert();
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [form, setForm] = useState(null); // null = cerrado; { id?, ...campos }
  const [saving, setSaving] = useState(false);

  const [query, setQuery] = useState('');
  const queryRef = useRef('');

  const load = useCallback(async () => {
    try {
      const q = queryRef.current.trim();
      const [u, r] = await Promise.all([
        api(q ? `/users/search?q=${encodeURIComponent(q)}` : '/users'),
        api('/roles/options'),
      ]);
      setUsers(u); setRoles(r);
    } catch (e) {
      alert.error('No se pudo cargar la lista', e.message);
    } finally {
      setLoading(false);
    }
  }, []); // eslint-disable-line

  // Buscar (permiso "Encontrar"): espera 300 ms tras la última tecla.
  useEffect(() => {
    queryRef.current = query;
    const t = setTimeout(load, 300);
    return () => clearTimeout(t);
  }, [query]); // eslint-disable-line

  const openCreate = () => setForm({ ...EMPTY, role_id: roles[0]?.id ?? '' });
  const openEdit = (u) => setForm({ id: u.id, name: u.name, email: u.email, password: '', role_id: u.role_id, active: u.active });
  const set = (k) => (e) => setForm((f) => ({ ...f, [k]: e.target.type === 'checkbox' ? e.target.checked : e.target.value }));

  const save = async (e) => {
    e.preventDefault();
    setSaving(true);
    const { id, ...body } = form;
    body.role_id = Number(body.role_id);
    try {
      await api(id ? `/users/${id}` : '/users', { method: id ? 'PUT' : 'POST', body });
      setForm(null);
      await load();
      alert.success(id ? 'Usuario actualizado' : 'Usuario creado', id ? 'Los cambios se guardaron.' : `${body.name} ya puede iniciar sesión.`);
    } catch (err) {
      alert.error('No se pudo guardar', err.message);
    } finally {
      setSaving(false);
    }
  };

  const remove = async (u) => {
    const ok = await alert.confirm({
      title: `¿Eliminar a ${u.name}?`, text: 'Esta acción no se puede deshacer.',
      confirmText: 'Eliminar usuario', danger: true,
    });
    if (!ok) return;
    try {
      await api(`/users/${u.id}`, { method: 'DELETE' });
      await load();
      alert.success('Usuario eliminado');
    } catch (err) {
      alert.error('No se pudo eliminar', err.message);
    }
  };

  return (
    <div className="fill">
      <PageHead title="Usuarios" subtitle="Crea cuentas y asigna el perfil que define sus permisos.">
        <PageSizeSelect />

          {can('users.find') && (
            <div className="input input--search">
              <Icon name="search" size={16} />
              <input value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Buscar por nombre, correo o perfil" aria-label="Buscar usuarios" />
            </div>
          )}
          <button className="btn btn--outline btn--icon" onClick={load} aria-label="Actualizar"><Icon name="refresh" size={16} /></button>
          {can('users.create') && (
            <button className="btn btn--primary" onClick={openCreate}><Icon name="plus" size={16} />Nuevo usuario</button>
          )}
      </PageHead>

      <TableCard rows={users} loading={loading} cols={4} empty="No hay usuarios para mostrar."
        head={<tr><th>Usuario</th><th>Perfil</th><th>Estado</th><th className="col-actions" /></tr>}
        renderRow={(u) => (
          <tr key={u.id}>
                <td>
                  <div className="cell-user">
                    <span className="avatar">{initials(u.name)}</span>
                    <div><strong>{u.name}</strong><small>{u.email}</small></div>
                  </div>
                </td>
                <td><span className="badge">{u.role}</span></td>
                <td><span className={`status ${u.active ? 'status--on' : 'status--off'}`}>{u.active ? 'Activo' : 'Inactivo'}</span></td>
                <td className="col-actions">
                  {can('users.update') && (
                    <button className="icon-btn" onClick={() => openEdit(u)} aria-label={`Editar ${u.name}`}><Icon name="edit" size={18} /></button>
                  )}
                  {can('users.delete') && u.id !== me.id && (
                    <button className="icon-btn icon-btn--danger" onClick={() => remove(u)} aria-label={`Eliminar ${u.name}`}><Icon name="trash" size={18} /></button>
                  )}
                </td>
              </tr>
        )} />

      <Modal open={!!form} onClose={() => setForm(null)} title={form?.id ? 'Editar usuario' : 'Nuevo usuario'}>
        {form && (
          <form onSubmit={save} className="form">
            <label className="field"><span>Nombre completo</span>
              <div className="input"><input required value={form.name} onChange={set('name')} /></div></label>
            <label className="field"><span>Correo</span>
              <div className="input"><input type="email" required value={form.email} onChange={set('email')} /></div></label>
            <label className="field"><span>{form.id ? 'Nueva contraseña (opcional)' : 'Contraseña'}</span>
              <div className="input"><input type="password" minLength={8} required={!form.id} autoComplete="new-password"
                     placeholder="Mínimo 8 caracteres" value={form.password} onChange={set('password')} /></div></label>
            <label className="field"><span>Perfil</span>
              <div className="input"><select required value={form.role_id} onChange={set('role_id')}>
                {roles.map((r) => <option key={r.id} value={r.id}>{r.name}</option>)}
              </select></div></label>
            <label className="check"><input type="checkbox" checked={form.active} onChange={set('active')} />Usuario activo</label>
            <div className="modal__foot">
              <button type="button" className="btn btn--outline" onClick={() => setForm(null)}>Cancelar</button>
              <button className="btn btn--primary" disabled={saving}>{saving ? 'Guardando…' : 'Guardar cambios'}</button>
            </div>
          </form>
        )}
      </Modal>
    </div>
  );
}
