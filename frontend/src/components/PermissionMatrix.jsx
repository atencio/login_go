import Check from './Check';
import { ACTIONS, groupCatalog } from '../permissions';
import { matrixSections } from '../navigation';

// Permisos en paneles por grupo (los mismos del menú), en dos columnas. Cada panel tiene su tabla
// módulo × acción y enlaces "Todos | Ninguno". Las celdas que no aplican llevan un punto gris.
export default function PermissionMatrix({ catalog, value, onToggle, onToggleMany, readOnly = false }) {
  const sections = matrixSections(groupCatalog(catalog));
  return (
    <div className="matrix-grid">
      {sections.map((section) => {
        const codes = section.rows.flatMap((r) => Object.values(r.actions));
        return (
          <section key={section.title} className="perm-group">
            <h4 className="perm-group__title">{section.title}</h4>
            <table className="perm-table">
              <thead>
                <tr>
                  <th className="perm-table__name">Módulo</th>
                  {ACTIONS.map((a) => <th key={a.key}>{a.label}</th>)}
                </tr>
              </thead>
              <tbody>
                {section.rows.map((row) => (
                  <tr key={row.module}>
                    <td className="perm-table__name">{row.label}</td>
                    {ACTIONS.map((a) => {
                      const code = row.actions[a.key];
                      return (
                        <td key={a.key}>
                          {code
                            ? <Check checked={value.includes(code)} disabled={readOnly}
                                     label={`${a.label} ${row.label}`} onChange={() => onToggle(code)} />
                            : <span className="perm-table__na" title="No aplica" />}
                        </td>
                      );
                    })}
                  </tr>
                ))}
              </tbody>
            </table>
            {!readOnly && (
              <div className="perm-group__all">
                <button type="button" className="link" onClick={() => onToggleMany(codes, true)}>Todos</button>
                <span aria-hidden="true">|</span>
                <button type="button" className="link" onClick={() => onToggleMany(codes, false)}>Ninguno</button>
              </div>
            )}
          </section>
        );
      })}
    </div>
  );
}
