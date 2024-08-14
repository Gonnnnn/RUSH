import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Typography,
  Paper,
  Box,
  CircularProgress,
  LinearProgress,
  Table,
  TableHead,
  TableContainer,
  TableRow,
  TableCell,
  TableBody,
} from '@mui/material';
import { useHeader } from './Layout';
import { Attendance, User, getUser, getUserAttendances, getUserId } from './client/http';
import { toYYYY년MM월DD일HH시MM분 } from './common/date';

const MyPage = () => {
  useHeader({ newTitle: 'Me' });

  const navigate = useNavigate();
  const [user, setUser] = useState<User>();
  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const init = async () => {
      try {
        setIsLoading(true);
        // TODO(#42): Fetch user data directly.
        const userId = await getUserId();
        setUser(await getUser(userId));
        setAttendances(await getUserAttendances(userId));
      } catch (error) {
        console.error(error);
        navigate('/');
      } finally {
        setIsLoading(false);
      }
    };

    init();
  }, [navigate]);

  if (isLoading || !user) {
    return (
      <Container>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  // TODO(#42): Fetch user attendance and show it here.
  return (
    <Container>
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="body1">
          {user.name} / {user.generation}기
        </Typography>
      </Paper>

      <TableContainer component={Paper}>
        <Box sx={{ width: '100%', height: '4px', mb: 2 }}>{isLoading ? <LinearProgress /> : null}</Box>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Attendance</TableCell>
              <TableCell>Joined at</TableCell>
            </TableRow>
          </TableHead>
          {attendances.length > 0 ? (
            <TableBody>
              {attendances.map((attendance) => (
                <TableRow
                  key={attendance.id}
                  onClick={() => navigate(`/sessions/${attendance.sessionId}`)}
                  style={{ cursor: 'pointer' }}
                >
                  <TableCell>{attendance.sessionName}</TableCell>
                  <TableCell>{toYYYY년MM월DD일HH시MM분(attendance.userJoinedAt)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          ) : (
            <TableBody>
              <TableRow>
                <TableCell colSpan={2}>출석 데이터가 없습니다. 출석하세요!</TableCell>
              </TableRow>
            </TableBody>
          )}
        </Table>
      </TableContainer>
    </Container>
  );
};

export default MyPage;
