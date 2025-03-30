import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowForwardIos, CalendarTodayOutlined } from '@mui/icons-material';
import {
  Container,
  Box,
  LinearProgress,
  useTheme,
  useMediaQuery,
  Stack,
  Typography,
  Grid,
  Paper,
  TablePagination,
} from '@mui/material';
import dayjs from 'dayjs';
import { useHeader } from '../../Layout';
import { Session } from '../../client/http/data';
import { listSessions } from '../../client/http/default';
import { toYYYY년MM월DD일HH시MM분 } from '../../common/date';
import useHandleError from '../../common/error';
import { useAdminMode } from '../../mode';

const UserSessionList = () => {
  const navigate = useNavigate();
  const { adminMode } = useAdminMode();
  const { handleError } = useHandleError();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  useHeader({ newTitle: 'Sessions' });

  const pageSize = isMobile ? 6 : 18;
  const [sessions, setSessions] = useState<Session[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(0);
  const [isEnd, setIsEnd] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  // TODO(#209): Split user/admin UIs and routes all together.
  if (adminMode) {
    navigate('/admin/sessions');
  }

  const fetchSessions = useCallback(
    async (page: number) => {
      try {
        setIsLoading(true);
        const offset = page * pageSize;
        const listSessionsResponse = await listSessions(offset, pageSize);
        setSessions(listSessionsResponse.sessions);
        setIsEnd(listSessionsResponse.isEnd);
        setTotalCount(listSessionsResponse.totalCount);
        setCurrentPage(page);
      } catch (e) {
        handleError({
          error: e,
          messageAuth: '세션 목록을 불러오는 중 오류가 발생했습니다.',
          messageInternal: '세션 목록을 불러오는 중 오류가 발생했습니다.',
        });
      } finally {
        setIsLoading(false);
      }
    },
    [pageSize, handleError],
  );

  useEffect(() => {
    fetchSessions(0);
  }, [fetchSessions]);

  const handleChangePage = async (_: unknown, newPage: number) => {
    if (isEnd && newPage > currentPage) {
      return;
    }
    await fetchSessions(newPage);
  };

  const handleRowClick = (session: Session) => {
    navigate(`/sessions/${session.id}`);
  };

  return (
    <Container>
      <Box display="flex" flexDirection="column" gap={1}>
        <Box
          sx={{
            position: 'sticky',
            top: 0,
            zIndex: 1,
            backgroundColor: theme.palette.background.paper,
          }}
        >
          <TablePagination
            rowsPerPageOptions={[]}
            component="div"
            rowsPerPage={pageSize}
            page={currentPage}
            onPageChange={handleChangePage}
            count={totalCount}
            slotProps={{
              actions: {
                previousButton: { disabled: currentPage === 0 },
                nextButton: { disabled: isEnd },
              },
            }}
          />
        </Box>

        <Box sx={{ width: '100%' }}>
          {isLoading ? (
            <LinearProgress sx={{ height: 4, padding: 0, margin: 0 }} />
          ) : (
            <Box sx={{ height: 4, padding: 0, margin: 0 }} />
          )}
        </Box>

        <Grid container spacing={1} sx={{ overflowY: 'auto', maxHeight: 'calc(100vh - 200px)' }}>
          {sessions.map((session) => (
            <Grid item key={session.id} xs={12} sm={6} md={4} sx={{ padding: 1 }}>
              <SessionCard session={session} onClick={() => handleRowClick(session)} />
            </Grid>
          ))}
        </Grid>
      </Box>
    </Container>
  );
};

const SessionCard = ({ session, onClick }: { session: Session; onClick: () => void }) => {
  const now = dayjs();
  const startsAt = dayjs(session.startsAt);
  const isPast = startsAt < now;

  return (
    <Paper sx={{ p: 1.5, width: '100%', maxWidth: 400, boxSizing: 'border-box' }} elevation={4} onClick={onClick}>
      <Stack spacing={1}>
        <Stack direction="row" justifyContent="space-between">
          <Typography variant="body1" sx={{ fontSize: '0.9rem' }}>
            {session.name}
          </Typography>
          <ArrowForwardIos color="action" sx={{ width: 18, height: 18 }} />
        </Stack>

        <Box display="flex" alignItems="center">
          <CalendarTodayOutlined sx={{ mr: 1 }} color={isPast ? 'primary' : 'action'} />
          <Typography variant="body2" sx={{ fontSize: '0.8rem' }}>
            {toYYYY년MM월DD일HH시MM분(session.startsAt)}
          </Typography>
        </Box>
      </Stack>
    </Paper>
  );
};

export default UserSessionList;
