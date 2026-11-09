// Páginas a mostrar: primera, última y la actual con sus vecinas; los huecos se marcan con "…".
function pageList(current, pages) {
  const nums = [...new Set([1, pages, current - 1, current, current + 1])]
    .filter((n) => n >= 1 && n <= pages)
    .sort((a, b) => a - b);
  const out = [];
  nums.forEach((n, i) => {
    if (i && n - nums[i - 1] > 1) out.push(`gap-${n}`);
    out.push(n);
  });
  return out;
}

export default function Pager({ page, pages, total, pageSize, onPage }) {
  const from = total === 0 ? 0 : (page - 1) * pageSize + 1;
  const to = Math.min(total, page * pageSize);
  return (
    <div className="pager">
      <span>{total === 0 ? 'Sin registros' : `Mostrando ${from}–${to} de ${total} registros`}</span>
      <nav className="pager__nav" aria-label="Paginación">
        <button className="pager__btn" disabled={page <= 1} onClick={() => onPage(page - 1)} aria-label="Página anterior">‹</button>
        {pageList(page, pages).map((n) => (typeof n === 'string'
          ? <span key={n} className="pager__gap">…</span>
          : <button key={n} className={n === page ? 'pager__btn pager__btn--on' : 'pager__btn'}
                    aria-current={n === page ? 'page' : undefined} onClick={() => onPage(n)}>{n}</button>))}
        <button className="pager__btn" disabled={page >= pages} onClick={() => onPage(page + 1)} aria-label="Página siguiente">›</button>
      </nav>
    </div>
  );
}
