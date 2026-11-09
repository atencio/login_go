import { useEffect, useState } from 'react';
import Pager from './Pager';
import { usePrefs } from '../context/PrefsContext';

// Tabla fija: el encabezado y el paginador no se mueven, solo se desplazan las filas.
// head: <tr> del encabezado · renderRow(fila) devuelve el <tr key=…> de cada registro.
export default function TableCard({ head, rows, loading, cols, empty = 'No hay registros para mostrar.', renderRow, compact = false }) {
  const { pageSize } = usePrefs();
  const [page, setPage] = useState(1);
  useEffect(() => { setPage(1); }, [pageSize]);

  const pages = Math.max(1, Math.ceil(rows.length / pageSize));
  const current = Math.min(page, pages);
  const slice = rows.slice((current - 1) * pageSize, current * pageSize);

  return (
    <div className="card card--grow">
      <div className="table-scroll">
        <table className={compact ? 'table table--compact' : 'table'}>
          <thead>{head}</thead>
          <tbody>
            {loading && <tr><td colSpan={cols} className="muted pad">Cargando…</td></tr>}
            {!loading && rows.length === 0 && <tr><td colSpan={cols} className="muted pad">{empty}</td></tr>}
            {!loading && slice.map(renderRow)}
          </tbody>
        </table>
      </div>
      <Pager page={current} pages={pages} total={rows.length} pageSize={pageSize} onPage={setPage} />
    </div>
  );
}
