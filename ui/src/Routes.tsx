import { Navigate, Outlet, Route, Routes, useLocation } from 'react-router-dom';
import { useAuth } from './AuthContext';
import Layout from './Layout';
import MyPage from './MyPage';
import SessionDetail from './SessionDetail';
import SessionList from './SessionList';
import SignIn from './SignIn';
import UserList from './UserList';

const AppRoutes = () => (
  <Routes>
    <Route path="/" element={<Layout />}>
      <Route path="/" element={<SessionList />} />
      <Route path="/sessions" element={<SessionList />} />
      <Route path="/sessions/:id" element={<SessionDetail />} />
      <Route path="/users" element={<UserList />} />
      <Route path="/" element={<AuthRoute />}>
        <Route path="/me" element={<MyPage />} />
      </Route>
      <Route path="signin" element={<SignIn />} />
    </Route>
  </Routes>
);

const AuthRoute = () => {
  const location = useLocation();
  const { authenticated } = useAuth();
  return authenticated ? <Outlet /> : <Navigate to="/signin" state={{ from: location.pathname }} />;
};

export default AppRoutes;
