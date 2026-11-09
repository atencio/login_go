import { useEffect, useMemo, useRef, useState } from 'react';
import { NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom';
import Icon from '../components/Icon';
import ProfileModal, { initials } from '../components/ProfileModal';
import { useAuth } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';
import { SCALES, usePrefs } from '../context/PrefsContext';
import { NAV_GROUPS } from '../navigation';

export default function DashboardLayout() {
  const { user, can, logout } = useAuth();
  const alert = useAlert();
  const { scale, setScale } = usePrefs();
  const navigate = useNavigate();
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem('sidebar') === 'collapsed');
  const [profileOpen, setProfileOpen] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const { pathname } = useLocation();
  const [openGroup, setOpenGroup] = useState(undefined); // un solo grupo desplegado a la vez; null = el usuario recogió todos
  const [q, setQ] = useState('');
  const [searchFocus, setSearchFocus] = useState(false);
  const searchRef = useRef(null);
  const menuRef = useRef(null);

  const groups = useMemo(() => NAV_GROUPS
    .map((g) => ({ ...g, items: g.items.filter((i) => !i.permission || can(i.permission)) }))
    .filter((g) => g.items.length), [can]);
  const isActive = (to) => (to === '/' ? pathname === '/' : pathname === to || pathname.startsWith(`${to}/`));
  const activeGroup = groups.find((g) => g.label && g.items.some((i) => isActive(i.to)))?.label;
  const firstGroup = groups.find((g) => g.label)?.label ?? null;
  const toggleGroup = (label) => setOpenGroup((cur) => (cur === label ? null : label));
  const allItems = groups.flatMap((g) => g.items);
  const results = allItems.filter((i) => i.label.toLowerCase().includes(q.trim().toLowerCase()));

  // Al entrar a una página (o navegar con el buscador) se despliega su grupo y se recogen los demás.
  useEffect(() => {
    setOpenGroup((cur) => activeGroup ?? (cur === undefined ? firstGroup : cur));
  }, [activeGroup, firstGroup]);

  useEffect(() => {
    const onKey = (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') { e.preventDefault(); searchRef.current?.focus(); }
    };
    const onClick = (e) => { if (menuRef.current && !menuRef.current.contains(e.target)) setMenuOpen(false); };
    window.addEventListener('keydown', onKey);
    window.addEventListener('mousedown', onClick);
    return () => { window.removeEventListener('keydown', onKey); window.removeEventListener('mousedown', onClick); };
  }, []);

  const toggleSidebar = () => setCollapsed((c) => {
    localStorage.setItem('sidebar', c ? 'open' : 'collapsed');
    return !c;
  });

  const goTo = (to) => { navigate(to); setQ(''); searchRef.current?.blur(); };

  const toggleFullscreen = () => {
    if (document.fullscreenElement) document.exitFullscreen();
    else document.documentElement.requestFullscreen?.();
  };

  const onLogout = async () => {
    setMenuOpen(false);
    const ok = await alert.confirm({
      title: '¿Cerrar sesión?', text: 'Tendrás que iniciar sesión de nuevo para continuar.',
      confirmText: 'Cerrar sesión', danger: true,
    });
    if (ok) logout();
  };

  return (
    <div className={`app ${collapsed ? 'app--collapsed' : ''}`}>
      <header className="header">
        <div className="header__brand">
          <span className="logo">N</span>
          <span className="brand-name">Nexo</span>
        </div>
        <button className="icon-btn" onClick={toggleSidebar} aria-label={collapsed ? 'Expandir menú' : 'Contraer menú'}>
          <Icon name="menu" size={18} />
        </button>

        <div className="gsearch">
          <input ref={searchRef} value={q} placeholder="Buscar página… (Ctrl+K)" aria-label="Buscar página"
                 onChange={(e) => setQ(e.target.value)} onFocus={() => setSearchFocus(true)}
                 onBlur={() => setTimeout(() => setSearchFocus(false), 120)}
                 onKeyDown={(e) => { if (e.key === 'Enter' && results[0]) goTo(results[0].to); }} />
          <Icon name="search" size={18} />
          {searchFocus && q.trim() && (
            <ul className="gsearch__results">
              {results.length === 0 && <li className="muted">Sin resultados</li>}
              {results.map((i) => (
                <li key={i.to}><button onMouseDown={() => goTo(i.to)}><Icon name={i.icon} size={16} />{i.label}</button></li>
              ))}
            </ul>
          )}
        </div>

        <div className="header__right">
          <div className="scale" role="group" aria-label="Tamaño de la fuente" title="Tamaño de la fuente">
            <span className="scale__icon" aria-hidden="true">Aa</span>
            {SCALES.map((s) => (
              <button key={s} className={s === scale ? 'scale__btn scale__btn--on' : 'scale__btn'}
                      aria-pressed={s === scale} onClick={() => setScale(s)}>{s}%</button>
            ))}
          </div>
          <button className="icon-btn" onClick={toggleFullscreen} aria-label="Pantalla completa"><Icon name="expand" size={18} /></button>
          <div className="usermenu" ref={menuRef}>
            <button className="usermenu__btn" onClick={() => setMenuOpen((o) => !o)} aria-expanded={menuOpen}>
              <span className="avatar avatar--sm">{initials(user.name)}</span>
              <strong className="usermenu__name">{user.name}</strong>
              <Icon name="chevronDown" size={16} />
            </button>
            {menuOpen && (
              <div className="usermenu__panel">
                <div className="usermenu__who"><strong>{user.name}</strong><small>{user.role}</small></div>
                <button onClick={() => { setMenuOpen(false); setProfileOpen(true); }}><Icon name="user" size={16} />Mi perfil</button>
                <button onClick={onLogout}><Icon name="logout" size={16} />Cerrar sesión</button>
              </div>
            )}
          </div>
        </div>
      </header>

      <aside className="sidebar">
        <nav>
          {groups.map((g, i) => {
            const open = !g.label || openGroup === g.label; // el grupo sin título (Dashboard, Reportes) siempre se ve
            return (
              <div key={i} className={open ? 'nav-group nav-group--open' : 'nav-group'}>
                {g.label && (
                  <button className="nav-group__label" aria-expanded={open} aria-controls={`nav-group-${i}`}
                          onClick={() => toggleGroup(g.label)}>
                    <span>{g.label}</span>
                    <Icon name="chevronDown" size={14} />
                  </button>
                )}
                <div className="nav-group__items" id={`nav-group-${i}`}>
                  <div>
                    {g.items.map((n) => (
                      <NavLink key={n.to} to={n.to} end={n.end} title={collapsed ? n.label : undefined}
                               className={({ isActive: on }) => `nav-link ${on ? 'nav-link--active' : ''}`}>
                        <Icon name={n.icon} size={18} />
                        <span className="nav-link__label">{n.label}</span>
                      </NavLink>
                    ))}
                  </div>
                </div>
              </div>
            );
          })}
        </nav>
      </aside>

      <main className="content"><Outlet /></main>
      <ProfileModal open={profileOpen} onClose={() => setProfileOpen(false)} />
    </div>
  );
}
