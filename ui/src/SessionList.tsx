import { useEffect, useState } from 'react';
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
  CircularProgress,
} from '@mui/material';
import SessionCreate from './SessionCreate';
import { Session, getSessions } from './client/http';
import toYYYY년MM월DD일HH시MM분 from './common/date';

const SessionList: React.FC = () => {
  const navigate = useNavigate();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [isLoading, setIsLoading] = useState(false);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    const init = async () => {
      try {
        setIsLoading(true);
        setSessions(await getSessions());
      } catch (e) {
        console.error(e);
      } finally {
        setIsLoading(false);
      }
    };

    init();
  }, []);

  if (isLoading) {
    return (
      <Container>
        <Typography variant="h4" sx={{ mb: 5 }}>
          Sessions
        </Typography>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  const handleChangePage = (_: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleRowClick = (session: Session) => {
    navigate(`/sessions/${session.id}`);
  };

  const paginatedSessions = sessions.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);

  return (
    <>
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
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Name</TableCell>
                  <TableCell>Starts At</TableCell>
                  <TableCell>Score</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {paginatedSessions.map((session) => (
                  <TableRow key={session.id} onClick={() => handleRowClick(session)} style={{ cursor: 'pointer' }}>
                    <TableCell>{session.name}</TableCell>
                    <TableCell>{toYYYY년MM월DD일HH시MM분(session.startsAt)}</TableCell>
                    <TableCell>{session.score}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
            <TablePagination
              rowsPerPageOptions={[10, 20, 30]}
              component="div"
              count={sessions.length}
              rowsPerPage={rowsPerPage}
              page={page}
              onPageChange={handleChangePage}
              onRowsPerPageChange={handleChangeRowsPerPage}
            />
          </TableContainer>
        </Box>
      </Container>
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
    </>
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
