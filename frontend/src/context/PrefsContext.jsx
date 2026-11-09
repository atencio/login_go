import { createContext, useContext, useEffect, useMemo, useState } from 'react';

export const PAGE_SIZES = [10, 25, 50, 100];
export const SCALES = [100, 90, 80]; // tamaño de la interfaz en %

const PrefsContext = createContext(null);
export const usePrefs = () => useContext(PrefsContext);

const read = (key, allowed, fallback) => {
  const v = Number(localStorage.getItem(key));
  return allowed.includes(v) ? v : fallback;
};

// Preferencias de la interfaz que se recuerdan en el navegador: registros por página y tamaño de fuente.
export function PrefsProvider({ children }) {
  const [pageSize, setPageSize] = useState(() => read('pageSize', PAGE_SIZES, 10));
  const [scale, setScale] = useState(() => read('uiScale', SCALES, 100));

  // Todo el CSS está en rem, así que cambiar el tamaño base escala texto, espacios e iconos a la vez.
  useEffect(() => {
    document.documentElement.style.fontSize = `${(16 * scale) / 100}px`;
    localStorage.setItem('uiScale', scale);
  }, [scale]);
  useEffect(() => { localStorage.setItem('pageSize', pageSize); }, [pageSize]);

  const value = useMemo(() => ({ pageSize, setPageSize, scale, setScale }), [pageSize, scale]);
  return <PrefsContext.Provider value={value}>{children}</PrefsContext.Provider>;
}
