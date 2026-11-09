export default function Check({ checked, onChange, disabled, label }) {
  return <input type="checkbox" className="cbx" checked={checked} disabled={disabled} onChange={onChange} aria-label={label} />;
}
