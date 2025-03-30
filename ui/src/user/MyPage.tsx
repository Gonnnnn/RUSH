import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { RunCircleOutlined, StarBorderRounded } from '@mui/icons-material';
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
import { Attendance, User } from '../client/http/data';
import { getUser, getUserAttendances, getUserAuth } from '../client/http/default';
import { toYYslashMMslashDDspaceHHcolonMM } from '../common/date';
import { useHeader } from '../contexts/header';

const MyPage = () => {
  useHeader({ newTitle: 'Me' });

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const [user, setUser] = useState<User>();
  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const init = async () => {
      try {
        setIsLoading(true);
        const { user_id: userId } = await getUserAuth();
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

  return (
    <Container>
      <Paper sx={{ p: 2, mb: 3, display: 'flex', flexDirection: 'column', justifyContent: 'center', gap: 2 }}>
        <Typography variant="h6">
          {user.externalName} / {user.generation}기
        </Typography>
        <Box display="flex" alignItems="center">
          <StarBorderRounded sx={{ mr: 1 }} color="primary" />
          <Typography variant="body1">
            출석 총점: {attendances.reduce((acc, cur) => acc + cur.sessionScore, 0)}점
          </Typography>
        </Box>
        <Box display="flex" alignItems="center">
          <RunCircleOutlined sx={{ mr: 1 }} color="primary" />
          <Typography variant="body1">세션 참여: {attendances.length}회</Typography>
        </Box>
      </Paper>

      <TableContainer component={Paper}>
        <Box sx={{ width: '100%', height: '4px', mb: 2 }}>{isLoading ? <LinearProgress /> : null}</Box>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell align="center" sx={{ width: '40%' }}>
                참여 세션 ({attendances.length})
              </TableCell>
              <TableCell align="center" sx={{ width: '40%' }}>
                출석 시간
              </TableCell>
              <TableCell align="center" sx={{ width: '20%' }}>
                점수
              </TableCell>
            </TableRow>
          </TableHead>
          {attendances.length > 0 ? (
            <TableBody>
              {attendances.map((attendance) => (
                <TableRow
                  key={attendance.id}
                  onClick={() => navigate(`/sessions/${attendance.sessionId}`, { state: { from: pathname } })}
                  style={{ cursor: 'pointer' }}
                >
                  <TableCell align="center">{attendance.sessionName}</TableCell>
                  <TableCell align="center">{toYYslashMMslashDDspaceHHcolonMM(attendance.userJoinedAt)}</TableCell>
                  <TableCell align="center">{attendance.sessionScore}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          ) : (
            <TableBody>
              <TableRow>
                <TableCell colSpan={3} align="center">
                  출석 데이터가 없습니다. 출석하세요!
                </TableCell>
              </TableRow>
            </TableBody>
          )}
        </Table>
      </TableContainer>
    </Container>
  );
};

export default MyPage;
