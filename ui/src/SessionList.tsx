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
  Typography,
  TablePagination,
  Button,
  Box,
  Modal,
  LinearProgress,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import SessionCreate from './SessionCreate';
import { Session, listSessions } from './client/http';
import toYYYY년MM월DD일HH시MM분 from './common/date';

const SessionList = () => {
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const pageSize = isMobile ? 8 : 10;
  const [sessions, setSessions] = useState<Session[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(0);
  const [isEnd, setIsEnd] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);

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
        console.error(e);
      } finally {
        setIsLoading(false);
      }
    },
    [pageSize],
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
      <Typography variant="h4" sx={{ mb: 5 }}>
        Sessions
      </Typography>
      <Box display="flex" flexDirection="column" gap={2}>
        <Box display="flex" justifyContent="flex-end">
          <Button
            variant="outlined"
            onClick={() => {
              setIsModalOpen(true);
            }}
          >
            New Session
          </Button>
        </Box>
        <TableContainer component={Paper}>
          {/* TODO(#31): Implements the common table UI with a loader. */}
          <Box sx={{ width: '100%', height: '4px', mb: 2 }}>{isLoading ? <LinearProgress /> : null}</Box>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Starts At</TableCell>
                <TableCell>Score</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {sessions.map((session) => (
                <TableRow key={session.id} onClick={() => handleRowClick(session)} style={{ cursor: 'pointer' }}>
                  <TableCell>{session.name}</TableCell>
                  <TableCell>{toYYYY년MM월DD일HH시MM분(session.startsAt)}</TableCell>
                  <TableCell>{session.score}</TableCell>
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
      <Modal
        open={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        aria-labelledby="modal-modal-title"
        aria-describedby="modal-modal-description"
      >
        <Box sx={ModalStyle}>
          <SessionCreate />
        </Box>
      </Modal>
    </Container>
  );
};

const ModalStyle = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: 'background.paper',
  boxShadow: 24,
  p: 4,
};

export default SessionList;
