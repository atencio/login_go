import Toggle from './Toggle';
import { ACTIONS, groupCatalog } from '../permissions';

// Matriz módulo × acción. Las celdas sin permiso en el catálogo se muestran con un guion.
export default function PermissionMatrix({ catalog, value, onToggle, readOnly = false }) {
  const rows = groupCatalog(catalog);
  return (
    <div className="matrix-wrap">
      <table className="matrix">
        <thead>
          <tr>
            <th className="matrix__num">#</th>
            <th>Tipo</th>
            {ACTIONS.map((a) => <th key={a.key} className="matrix__col">{a.label}</th>)}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={row.module}>
              <td className="matrix__num">{i + 1}</td>
              <td className="matrix__type">{row.label}</td>
              {ACTIONS.map((a) => {
                const code = row.actions[a.key];
                return (
                  <td key={a.key} className="matrix__col">
                    {code
                      ? <Toggle checked={value.includes(code)} disabled={readOnly}
                                label={`${a.label} ${row.label}`} onChange={() => onToggle(code)} />
                      : <span className="matrix__na" title="No aplica">—</span>}
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
