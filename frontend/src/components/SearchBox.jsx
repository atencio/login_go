import Icon from './Icon';

export default function SearchBox({ value, onChange, placeholder }) {
  return (
    <div className="input input--search">
      <Icon name="search" size={16} />
      <input value={value} onChange={(e) => onChange(e.target.value)} placeholder={placeholder} aria-label={placeholder} />
    </div>
  );
}
