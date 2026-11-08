import { useEffect, useMemo, useRef, useState } from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import Icon from '../components/Icon';
import ProfileModal, { initials } from '../components/ProfileModal';
import { useAuth } from '../context/AuthContext';
import { useAlert } from '../context/AlertContext';

const GROUPS = [
  { label: null, items: [
    { to: '/', label: 'Dashboard', icon: 'dashboard', end: true },
    { to: '/reportes', label: 'Reportes', icon: 'chart', permission: 'reports.list' },
  ] },
  { label: 'Inventario de activos', items: [
    { to: '/asignaciones', label: 'Asignaciones', icon: 'swap', permission: 'assignments.list' },
    { to: '/equipos', label: 'Equipos de oficina', icon: 'monitor', permission: 'equipment.list' },
    { to: '/funcionarios', label: 'Funcionarios', icon: 'idcard', permission: 'employees.list' },
    { to: '/cargos', label: 'Cargos', icon: 'briefcase', permission: 'positions.list' },
  ] },
  { label: 'Administración', items: [
    { to: '/usuarios', label: 'Usuarios', icon: 'users', permission: 'users.list' },
    { to: '/perfiles', label: 'Perfiles', icon: 'shield', permission: 'roles.list' },
  ] },
];

export default function DashboardLayout() {
  const { user, can, logout } = useAuth();
  const alert = useAlert();
  const navigate = useNavigate();
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem('sidebar') === 'collapsed');
  const [profileOpen, setProfileOpen] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const [q, setQ] = useState('');
  const [searchFocus, setSearchFocus] = useState(false);
  const searchRef = useRef(null);
  const menuRef = useRef(null);

  const groups = useMemo(() => GROUPS
    .map((g) => ({ ...g, items: g.items.filter((i) => !i.permission || can(i.permission)) }))
    .filter((g) => g.items.length), [can]);
  const allItems = groups.flatMap((g) => g.items);
  const results = allItems.filter((i) => i.label.toLowerCase().includes(q.trim().toLowerCase()));

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
          {groups.map((g, i) => (
            <div key={i} className="nav-group">
              {g.label && <div className="nav-group__label">{g.label}</div>}
              {g.items.map((n) => (
                <NavLink key={n.to} to={n.to} end={n.end} title={collapsed ? n.label : undefined}
                         className={({ isActive }) => `nav-link ${isActive ? 'nav-link--active' : ''}`}>
                  <Icon name={n.icon} size={18} />
                  <span className="nav-link__label">{n.label}</span>
                </NavLink>
              ))}
            </div>
          ))}
        </nav>
      </aside>

      <main className="content"><Outlet /></main>
      <ProfileModal open={profileOpen} onClose={() => setProfileOpen(false)} />
    </div>
  );
}
