import { createContext, useCallback, useContext, useEffect, useRef, useState } from 'react';
import Icon from '../components/Icon';

const AlertContext = createContext(null);
export const useAlert = () => useContext(AlertContext);

const ICONS = { success: 'check', error: 'x', warning: 'alert', info: 'info' };

export function AlertProvider({ children }) {
  const [alert, setAlert] = useState(null);
  const resolver = useRef(null);

  const show = useCallback((opts) => new Promise((resolve) => {
    resolver.current = resolve;
    setAlert(opts);
  }), []);

  const close = (result) => { setAlert(null); resolver.current?.(result); };

  const api = {
    success: (title, text) => show({ type: 'success', title, text }),
    error: (title, text) => show({ type: 'error', title, text }),
    info: (title, text) => show({ type: 'info', title, text }),
    confirm: ({ title, text, confirmText = 'Confirmar', cancelText = 'Cancelar', danger = false }) =>
      show({ type: danger ? 'warning' : 'info', title, text, confirmText, cancelText, confirm: true, danger }),
  };

  useEffect(() => {
    if (!alert) return;
    const onKey = (e) => e.key === 'Escape' && close(false);
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [alert]);

  return (
    <AlertContext.Provider value={api}>
      {children}
      {alert && (
        <div className="overlay overlay--alert" onMouseDown={() => close(false)}>
          <div className={`alert alert--${alert.type}`} role="alertdialog" aria-modal="true"
               onMouseDown={(e) => e.stopPropagation()}>
            <div className="alert__icon"><Icon name={ICONS[alert.type]} size={34} /></div>
            <h3 className="alert__title">{alert.title}</h3>
            {alert.text && <p className="alert__text">{alert.text}</p>}
            <div className="alert__actions">
              {alert.confirm && (
                <button className="btn btn--outline" onClick={() => close(false)}>{alert.cancelText}</button>
              )}
              <button autoFocus className={`btn ${alert.danger ? 'btn--danger' : 'btn--primary'}`}
                      onClick={() => close(true)}>
                {alert.confirm ? alert.confirmText : 'Entendido'}
              </button>
            </div>
          </div>
        </div>
      )}
    </AlertContext.Provider>
  );
}
