import { useCallback, useEffect, useState } from 'react';
import { api } from '../api/client';
import { useAlert } from '../context/AlertContext';

export function useDebounced(value, ms = 300) {
  const [v, setV] = useState(value);
  useEffect(() => {
    const t = setTimeout(() => setV(value), ms);
    return () => clearTimeout(t);
  }, [value, ms]);
  return v;
}

// Carga una lista de la API. Los parámetros vacíos no se envían; si cambian, vuelve a cargar.
export function useResource(path, params = {}, { enabled = true } = {}) {
  const alert = useAlert();
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);

  const qs = new URLSearchParams(
    Object.entries(params).filter(([, v]) => v !== '' && v != null && v !== 0),
  ).toString();
  const url = qs ? `${path}?${qs}` : path;

  const reload = useCallback(async () => {
    try {
      setRows(await api(url));
    } catch (e) {
      alert.error('No se pudo cargar la información', e.message);
    } finally {
      setLoading(false);
    }
  }, [url]); // eslint-disable-line

  useEffect(() => { if (enabled) reload(); else setLoading(false); }, [reload, enabled]);

  return { rows, loading, reload };
}
