import { Navigate, Route, Routes } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import DashboardLayout from './layouts/DashboardLayout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Users from './pages/Users';
import Reports from './pages/Reports';
import Roles from './pages/Roles';
import Positions from './pages/Positions';
import Employees from './pages/Employees';
import Equipment from './pages/Equipment';
import Assignments from './pages/Assignments';
import Categories from './pages/Categories';
import Brands from './pages/Brands';

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
          <Route element={<ProtectedRoute permission="assignments.list" />}>
            <Route path="asignaciones" element={<Assignments />} />
          </Route>
          <Route element={<ProtectedRoute permission="equipment.list" />}>
            <Route path="equipos" element={<Equipment />} />
          </Route>
          <Route element={<ProtectedRoute permission="employees.list" />}>
            <Route path="funcionarios" element={<Employees />} />
          </Route>
          <Route element={<ProtectedRoute permission="positions.list" />}>
            <Route path="cargos" element={<Positions />} />
          </Route>
          <Route element={<ProtectedRoute permission="categories.list" />}>
            <Route path="categorias" element={<Categories />} />
          </Route>
          <Route element={<ProtectedRoute permission="brands.list" />}>
            <Route path="marcas" element={<Brands />} />
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
