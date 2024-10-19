import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Container,
  TablePagination,
  Box,
  LinearProgress,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import { useHeader } from '../Layout';
import { Session } from '../client/http/data';
import { listSessions } from '../client/http/default';
import { toYYslashMMslashDDspaceHHcolonMMwithDay } from '../common/date';
import useHandleError from '../common/error';

const UserSessionList = () => {
  const navigate = useNavigate();
  const { handleError } = useHandleError();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  useHeader({ newTitle: 'Sessions' });

  const pageSize = isMobile ? 8 : 10;
  const [sessions, setSessions] = useState<Session[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(0);
  const [isEnd, setIsEnd] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchSessions = useCallback(
    async (page: number) => {
      try {
        setIsLoading(true);
        const offset = page * pageSize;
        // TODO(#200): Fetch sessions that don't have admin-specific fields. Split clients.
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
    fetchSessions(newPage);
  };

  const handleRowClick = (session: Session) => {
    navigate(`/sessions/${session.id}`);
  };

  return (
    <Container>
      <Box display="flex" flexDirection="column" gap={2}>
        <TableContainer component={Paper}>
          {/* TODO(#31): Implements the common table UI with a loader. */}
          <Box sx={{ width: '100%', height: '4px', mb: 2 }}>{isLoading ? <LinearProgress /> : null}</Box>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell align="center" sx={{ width: '50%' }}>
                  시작 시간
                </TableCell>
                <TableCell align="center" sx={{ width: '50%' }}>
                  이름
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {sessions.map((session) => (
                <TableRow key={session.id} onClick={() => handleRowClick(session)} style={{ cursor: 'pointer' }}>
                  <TableCell align="center">{toYYslashMMslashDDspaceHHcolonMMwithDay(session.startsAt)}</TableCell>
                  <TableCell align="center">{session.name}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
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
        </TableContainer>
      </Box>
    </Container>
  );
};

export default UserSessionList;
