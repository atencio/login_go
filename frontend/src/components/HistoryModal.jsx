import Modal from './Modal';
import { useResource } from '../hooks/useResource';
import { fmtDateTime } from '../lib/format';

// Historial de movimientos (entregas y devoluciones). query: { equipment_id } o { employee_id }
export default function HistoryModal({ title, query, onClose }) {
  const { rows, loading } = useResource('/assignments', query);
  return (
    <Modal open onClose={onClose} title={title} width={940}>
      <div className="table-scroll">
        <table className="table table--compact">
          <thead>
            <tr><th>Equipo</th><th>Funcionario</th><th>Entrega</th><th>Devolución</th><th>Notas</th></tr>
          </thead>
          <tbody>
            {loading && <tr><td colSpan="5" className="muted pad">Cargando…</td></tr>}
            {!loading && rows.length === 0 && <tr><td colSpan="5" className="muted pad">Aún no hay movimientos.</td></tr>}
            {rows.map((a) => (
              <tr key={a.id}>
                <td><code className="code">{a.equipment_code}</code><div>{a.equipment_name}</div></td>
                <td>{a.employee}</td>
                <td>{fmtDateTime(a.assigned_at)}<small className="sub">{a.assigned_by && `por ${a.assigned_by}`}</small></td>
                <td>
                  {a.active
                    ? <span className="pill pill--info">En uso</span>
                    : <>{fmtDateTime(a.returned_at)}<small className="sub">{a.received_by && `por ${a.received_by}`}</small></>}
                </td>
                <td className="notes">{[a.delivery_notes, a.return_notes].filter(Boolean).join(' · ') || '—'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="modal__foot"><button className="btn btn--outline" onClick={onClose}>Cerrar</button></div>
    </Modal>
  );
}
