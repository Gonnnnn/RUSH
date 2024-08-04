import { Route, Routes } from 'react-router-dom';
import Layout from './Layout';
import MyPage from './MyPage';
import SessionDetail from './SessionDetail';
import SessionList from './SessionList';
import UserList from './UserList';

const AppRoutes = () => (
  <Routes>
    <Route path="/" element={<Layout />}>
      <Route path="/" element={<SessionList />} />
      <Route path="/sessions" element={<SessionList />} />
      <Route path="/sessions/:id" element={<SessionDetail />} />
      <Route path="/users" element={<UserList />} />
      <Route path="/me" element={<MyPage />} />
    </Route>
  </Routes>
);

export default AppRoutes;
