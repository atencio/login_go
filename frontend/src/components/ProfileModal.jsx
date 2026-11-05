import { useEffect, useState } from 'react';
import Modal from './Modal';
import PermissionMatrix from './PermissionMatrix';
import { useAuth } from '../context/AuthContext';
import { loadCatalog } from '../api/catalog';

export const initials = (name = '') =>
  name.split(' ').filter(Boolean).slice(0, 2).map((p) => p[0].toUpperCase()).join('');

export default function ProfileModal({ open, onClose }) {
  const { user, profile } = useAuth();
  const [catalog, setCatalog] = useState(null);

  useEffect(() => {
    if (open) loadCatalog().then(setCatalog).catch(() => setCatalog([]));
  }, [open]);

  if (!user) return null;
  return (
    <Modal open={open} onClose={onClose} title="Mi perfil" width={780}>
      <div className="profile__head">
        <div className="avatar avatar--lg">{initials(user.name)}</div>
        <div>
          <div className="profile__name">{user.name}</div>
          <div className="profile__mail">{user.email}</div>
          <span className="badge">Perfil: {user.role}</span>
        </div>
      </div>

      <h4 className="profile__section">Permisos de tu perfil</h4>
      {catalog
        ? <PermissionMatrix catalog={catalog} value={profile.permissions} readOnly onToggle={() => {}} />
        : <p className="muted">Cargando permisos…</p>}
    </Modal>
  );
}
