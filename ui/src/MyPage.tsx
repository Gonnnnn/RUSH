import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Typography,
  Paper,
  Box,
  Button,
  CircularProgress,
  LinearProgress,
  Table,
  TableHead,
  TableContainer,
  TableRow,
  TableCell,
  TableBody,
} from '@mui/material';
import { Attendance, User, getUser, getUserAttendances, getUserId } from './client/http';
import toYYYY년MM월DD일HH시MM분 from './common/date';

const MyPage = () => {
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
        <Typography variant="h4" sx={{ mb: 3 }}>
          My Page
        </Typography>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  // TODO(#42): Fetch user attendance and show it here.
  return (
    <Container>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
          gap: 1,
          mb: 3,
        }}
      >
        <Typography variant="h4">My Page</Typography>
        <Button variant="outlined" onClick={() => navigate('/')} sx={{ alignSelf: 'flex-start' }}>
          Back
        </Button>
      </Box>
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="h6">Details</Typography>
        <Typography>Name: {user.name}</Typography>
        <Typography>Generation: {user.generation}</Typography>
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
          <TableBody>
            {attendances.map((attendance) => (
              <TableRow
                key={attendance.id}
                onClick={() => navigate(`/session/${attendance.sessionId}`)}
                style={{ cursor: 'pointer' }}
              >
                <TableCell>{attendance.sessionName}</TableCell>
                <TableCell>{toYYYY년MM월DD일HH시MM분(attendance.joinedAt)}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
};

export default MyPage;
