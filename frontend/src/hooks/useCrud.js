import { useState } from 'react';
import { api } from '../api/client';
import { useAlert } from '../context/AlertContext';

// Guardar (crear/editar) y eliminar con modales de confirmación y resultado.
// t: { created, updated, deleted, saveError, createdText?, updatedText? } (los textos pueden ser funciones de la respuesta)
export function useCrud(path, reload, t) {
  const alert = useAlert();
  const [saving, setSaving] = useState(false);

  const save = async (id, body) => {
    setSaving(true);
    try {
      const res = await api(id ? `${path}/${id}` : path, { method: id ? 'PUT' : 'POST', body });
      await reload();
      const text = id ? t.updatedText : t.createdText;
      alert.success(id ? t.updated : t.created, typeof text === 'function' ? text(res) : text);
      return true;
    } catch (e) {
      alert.error(t.saveError, e.message);
      return false;
    } finally {
      setSaving(false);
    }
  };

  const remove = async (row, { title, text }) => {
    const ok = await alert.confirm({
      title, text: text || 'Esta acción no se puede deshacer.', confirmText: 'Eliminar', danger: true,
    });
    if (!ok) return;
    try {
      await api(`${path}/${row.id}`, { method: 'DELETE' });
      await reload();
      alert.success(t.deleted);
    } catch (e) {
      alert.error('No se pudo eliminar', e.message);
    }
  };

  return { saving, save, remove };
}
