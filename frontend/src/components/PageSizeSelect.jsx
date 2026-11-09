import { PAGE_SIZES, usePrefs } from '../context/PrefsContext';

// "Mostrar N registros": va junto al buscador de cada tabla y vale para todas.
export default function PageSizeSelect() {
  const { pageSize, setPageSize } = usePrefs();
  return (
    <label className="pagesize">
      <span>Mostrar</span>
      <div className="input input--select">
        <select value={pageSize} onChange={(e) => setPageSize(Number(e.target.value))} aria-label="Registros por página">
          {PAGE_SIZES.map((n) => <option key={n} value={n}>{n}</option>)}
        </select>
      </div>
      <span>registros</span>
    </label>
  );
}
