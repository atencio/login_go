import { STATUS } from '../lib/format';

export default function StatusPill({ status }) {
  const s = STATUS[status] || { label: status, tone: 'muted' };
  return <span className={`pill pill--${s.tone}`}>{s.label}</span>;
}
