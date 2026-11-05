import Icon from './Icon';

export default function Toggle({ checked, onChange, disabled, label }) {
  return (
    <button type="button" role="switch" aria-checked={checked} aria-label={label} disabled={disabled}
            className={`toggle ${checked ? 'toggle--on' : ''}`} onClick={onChange}>
      <span className="toggle__knob"><Icon name={checked ? 'check' : 'x'} size={11} /></span>
    </button>
  );
}
