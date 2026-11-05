import { Navigate, Route, Routes } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import DashboardLayout from './layouts/DashboardLayout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Users from './pages/Users';
import Reports from './pages/Reports';
import Roles from './pages/Roles';

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route element={<ProtectedRoute />}>
        <Route element={<DashboardLayout />}>
          <Route index element={<Dashboard />} />
          <Route element={<ProtectedRoute permission="users.list" />}>
            <Route path="usuarios" element={<Users />} />
          </Route>
          <Route element={<ProtectedRoute permission="roles.list" />}>
            <Route path="perfiles" element={<Roles />} />
          </Route>
          <Route element={<ProtectedRoute permission="reports.list" />}>
            <Route path="reportes" element={<Reports />} />
          </Route>
        </Route>
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
