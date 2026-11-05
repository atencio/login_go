import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function ProtectedRoute({ permission }) {
  const { profile, loading, can } = useAuth();
  if (loading) return <div className="splash">Cargando…</div>;
  if (!profile) return <Navigate to="/login" replace />;
  if (permission && !can(permission)) return <Navigate to="/" replace />;
  return <Outlet />;
}
