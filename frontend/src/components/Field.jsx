export default function Field({ label, hint, area = false, className = '', children }) {
  return (
    <label className={`field ${className}`}>
      <span>{label}</span>
      <div className={`input ${area ? 'input--area' : ''}`}>{children}</div>
      {hint && <small className="hint">{hint}</small>}
    </label>
  );
}
