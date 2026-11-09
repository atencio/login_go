import { useEffect } from 'react';
import Icon from './Icon';

export default function Modal({ open, title, onClose, children, width = 520 }) {
  useEffect(() => {
    if (!open) return;
    const onKey = (e) => e.key === 'Escape' && onClose();
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [open, onClose]);

  if (!open) return null;
  return (
    <div className="overlay" onMouseDown={onClose}>
      <div className={width >= 1000 ? 'modal modal--xl' : 'modal'} style={{ maxWidth: `${width / 16}rem` }} role="dialog" aria-modal="true"
           aria-label={title} onMouseDown={(e) => e.stopPropagation()}>
        <header className="modal__head">
          <h2>{title}</h2>
          <button className="icon-btn" onClick={onClose} aria-label="Cerrar"><Icon name="x" size={18} /></button>
        </header>
        <div className="modal__body">{children}</div>
      </div>
    </div>
  );
}
