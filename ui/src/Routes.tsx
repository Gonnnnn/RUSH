import { Navigate, Outlet, Route, Routes, useLocation } from 'react-router-dom';
import HalfYearAttendances from './Attendance';
import { useAuth } from './AuthContext';
import { Layout } from './Layout';
import MyPage from './MyPage';
import SessionDetail from './SessionDetail';
import SessionList from './SessionList';
import SignIn from './SignIn';
import UserList from './UserList';

const AppRoutes = () => (
  <Routes>
    <Route element={<Layout />}>
      <Route index element={<SessionList />} />
      <Route path="/sessions" element={<SessionList />} />
      <Route path="/sessions/:id" element={<SessionDetail />} />
      <Route path="/users" element={<UserList />} />
      <Route element={<AuthRoute />}>
        <Route path="/me" element={<MyPage />} />
        <Route path="/attendances" element={<HalfYearAttendances />} />
      </Route>
      <Route path="/signin" element={<SignIn />} />
    </Route>
  </Routes>
);

const AuthRoute = () => {
  const location = useLocation();
  const { authenticated } = useAuth();
  return authenticated ? <Outlet /> : <Navigate to="/signin" state={{ from: location.pathname }} />;
};

export default AppRoutes;
