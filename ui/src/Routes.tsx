import { Navigate, Outlet, Route, Routes, useLocation } from 'react-router-dom';
import { Layout } from './Layout';
import SignIn from './SignIn';
import AdminLayout from './admin/Layout';
import { useAuth } from './auth';
import { Role } from './auth/role';
import HalfYearAttendances from './user/Attendance';
import MyPage from './user/MyPage';
import { AdminSessionDetail, UserSessionDetail } from './user/SessionDetail';
import { AdminSessionList, UserSessionList } from './user/SessionList';

const AppRoutes = () => (
  <Routes>
    <Route path="/">
      <Route element={<Layout />}>
        <Route index element={<UserSessionList />} />
        <Route path="/sessions" element={<UserSessionList />} />
        <Route path="/sessions/:id" element={<UserSessionDetail />} />
        <Route element={<AuthRoute />}>
          <Route path="/me" element={<MyPage />} />
          <Route path="/attendances" element={<HalfYearAttendances />} />
        </Route>
      </Route>
      <Route path="/signin" element={<SignIn />} />
    </Route>

    <Route path="/admin" element={<AdminRoute />}>
      <Route element={<AdminLayout />}>
        <Route index element={<AdminSessionList />} />
        <Route path="me" element={<MyPage />} />
        <Route path="sessions" element={<AdminSessionList />} />
        <Route path="sessions/:id" element={<AdminSessionDetail />} />
      </Route>
    </Route>
  </Routes>
);

const AuthRoute = () => {
  const location = useLocation();
  const { authenticated } = useAuth();
  return authenticated ? <Outlet /> : <Navigate to="/signin" state={{ from: location.pathname }} />;
};

const AdminRoute = () => {
  const { role } = useAuth();
  return role === Role.ADMIN ? <Outlet /> : <Navigate to="/" />;
};

export default AppRoutes;
