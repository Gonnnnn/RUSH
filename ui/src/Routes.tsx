import { Navigate, Outlet, Route, Routes, useLocation } from 'react-router-dom';
import HalfYearAttendances from './Attendance';
import { Layout } from './Layout';
import MyPage from './MyPage';
import { AdminSessionDetail, UserSessionDetail } from './SessionDetail';
import { AdminSessionList, UserSessionList } from './SessionList';
import SignIn from './SignIn';
import { useAuth } from './auth';
import { Role } from './auth/role';

const AppRoutes = () => (
  <Routes>
    <Route element={<Layout />}>
      <Route index element={<UserSessionList />} />
      <Route path="/sessions" element={<UserSessionList />} />
      <Route path="/sessions/:id" element={<UserSessionDetail />} />
      <Route element={<AuthRoute />}>
        <Route path="/me" element={<MyPage />} />
        <Route path="/attendances" element={<HalfYearAttendances />} />
      </Route>
      <Route path="/signin" element={<SignIn />} />
      <Route path="/admin" element={<AdminRoute />}>
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
