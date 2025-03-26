import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import MainLayout from './layouts/MainLayout';
import Login from './pages/Login';
import Services from './pages/Services';
import Stats from './pages/Stats';
import Users from './pages/Users';
import Settings from './pages/Settings';
import AuthGuard from './components/AuthGuard';
import PermissionGuard from './components/PermissionGuard';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/" element={
          <AuthGuard>
            <MainLayout />
          </AuthGuard>
        }>
          <Route index element={<Navigate to="/services" replace />} />
          <Route path="services" element={<Services />} />
          <Route path="stats" element={<Stats />} />
          <Route path="settings" element={
            <PermissionGuard permission="manage_system">
              <Settings />
            </PermissionGuard>
          } />
          <Route path="users" element={
            <PermissionGuard permission="manage_users">
              <Users />
            </PermissionGuard>
          } />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
