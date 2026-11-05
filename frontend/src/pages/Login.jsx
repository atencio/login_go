import { useState } from 'react';
import { Navigate, useNavigate } from 'react-router-dom';
import Icon from '../components/Icon';
import { useAuth } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';

export default function Login() {
  const { profile, login } = useAuth();
  const alert = useAlert();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [show, setShow] = useState(false);
  const [remember, setRemember] = useState(true);
  const [busy, setBusy] = useState(false);

  if (profile) return <Navigate to="/" replace />;

  const submit = async (e) => {
    e.preventDefault();
    setBusy(true);
    try {
      await login(email, password, remember);
      navigate('/', { replace: true });
    } catch (err) {
      await alert.error('No pudimos iniciar sesión', err.message);
    } finally {
      setBusy(false);
    }
  };

  const forgot = (e) => {
    e.preventDefault();
    alert.info('¿Olvidaste tu contraseña?', 'Pide a un administrador que la restablezca desde Usuarios.');
  };

  return (
    <div className="login">
      <section className="login__main">
        <div className="brand"><span className="logo">N</span><span className="brand-name">Nexo</span></div>

        <form className="login__form" onSubmit={submit}>
          <h1>Iniciar sesión</h1>
          <p className="muted">Utiliza el correo electrónico del trabajo para acceder al panel de control.</p>

          <label className="field">
            <span>Dirección de correo electrónico</span>
            <div className="input">
              <input type="email" required autoComplete="username" placeholder="tú@empresa.com"
                     value={email} onChange={(e) => setEmail(e.target.value)} />
            </div>
          </label>

          <label className="field">
            <span className="field__row">Contraseña <a href="#" onClick={forgot}>¿Olvidó?</a></span>
            <div className="input">
              <input type={show ? 'text' : 'password'} required autoComplete="current-password"
                     placeholder="Introduce tu contraseña" value={password} onChange={(e) => setPassword(e.target.value)} />
              <button type="button" className="input__btn" onClick={() => setShow((s) => !s)}
                      aria-label={show ? 'Ocultar contraseña' : 'Mostrar contraseña'}>
                <Icon name={show ? 'eyeOff' : 'eye'} size={18} />
              </button>
            </div>
          </label>

          <label className="check">
            <input type="checkbox" checked={remember} onChange={(e) => setRemember(e.target.checked)} />
            Mantener mi sesión iniciada durante 30 días.
          </label>

          <button className="btn btn--primary btn--block" disabled={busy}>{busy ? 'Ingresando…' : 'Iniciar sesión'}</button>
        </form>
      </section>

      <section className="login__aside" aria-hidden="false">
        <p>Todo lo que la capa de administración necesita, y nada que no necesite.</p>
        <div className="login__preview" aria-hidden="true"><i /><i /><i /></div>
      </section>
    </div>
  );
}
