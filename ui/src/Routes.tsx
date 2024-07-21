import { Route, Routes } from 'react-router-dom';
import AttendanceDetail from './AttendanceDetail';
import AttendanceList from './AttendanceList';
import Home from './Home';
import Layout from './Layout';
import SessionDetail from './SessionDetail';
import SessionList from './SessionList';
import UserList from './UserList';

const AppRoutes = () => (
  <Routes>
    <Route path="/" element={<Layout />}>
      <Route path="/" element={<Home />} />
      <Route path="/sessions" element={<SessionList />} />
      <Route path="/sessions/:id" element={<SessionDetail />} />
      <Route path="/users" element={<UserList />} />
      <Route path="/attendances" element={<AttendanceList />} />
      <Route path="/attendances/:id" element={<AttendanceDetail />} />
    </Route>
  </Routes>
);

export default AppRoutes;
