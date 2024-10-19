import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Tooltip } from 'react-tooltip';
import { CheckCircleOutlineRounded } from '@mui/icons-material';
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
  Button,
  Box,
  Modal,
  LinearProgress,
  useTheme,
  useMediaQuery,
  Typography,
} from '@mui/material';
import { useHeader } from '../Layout';
import SessionCreate from '../SessionCreate';
import { Session, listSessions } from '../client/http';
import { toYYslashMMslashDDspaceHHcolonMMwithDay } from '../common/date';

const AdminSessionList = () => {
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  useHeader({ newTitle: 'Sessions' });

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
    navigate(`/admin/sessions/${session.id}`);
  };

  return (
    <Container>
      <Box display="flex" flexDirection="column" gap={2}>
        <Box display="flex" justifyContent="flex-end">
          <Button
            variant="outlined"
            onClick={() => {
              setIsModalOpen(true);
            }}
          >
            New
          </Button>
        </Box>
        <Box display="flex" alignItems="center" gap={1}>
          <Typography variant="body2">출석 sync 상태</Typography>
          <Box display="flex" alignItems="center">
            <CheckCircleOutlineRounded color="primary" data-tooltip-id="attendance-status-applied-tooltip" />
            <CheckCircleOutlineRounded color="warning" data-tooltip-id="attendance-status-ignored-tooltip" />
            <CheckCircleOutlineRounded color="disabled" data-tooltip-id="attendance-status-not-applied-tooltip" />
            <Tooltip
              id="attendance-status-applied-tooltip"
              place="top"
              content="출석 반영 완료"
              openEvents={{ click: true, mouseover: true }}
            />
            <Tooltip
              id="attendance-status-ignored-tooltip"
              place="top"
              content="출석 반영이 시도됐으나 무시된 상태"
              openEvents={{ click: true, mouseover: true }}
            />
            <Tooltip
              id="attendance-status-not-applied-tooltip"
              place="top"
              content="세션 출석이 마감되지 않아, 아직 반영되지 않은 상태"
              openEvents={{ click: true, mouseover: true }}
            />
          </Box>
        </Box>
        <TableContainer component={Paper}>
          {/* TODO(#31): Implements the common table UI with a loader. */}
          <Box sx={{ width: '100%', height: '4px', mb: 2 }}>{isLoading ? <LinearProgress /> : null}</Box>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell align="center" sx={{ width: '35%' }}>
                  시작 시간
                </TableCell>
                <TableCell align="center" sx={{ width: '45%' }}>
                  이름
                </TableCell>
                <TableCell align="center" sx={{ width: '20%' }}>
                  Sync
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {sessions.map((session) => (
                <TableRow key={session.id} onClick={() => handleRowClick(session)} style={{ cursor: 'pointer' }}>
                  <TableCell align="center">{toYYslashMMslashDDspaceHHcolonMMwithDay(session.startsAt)}</TableCell>
                  <TableCell align="center">{session.name}</TableCell>
                  <TableCell align="center">
                    {(() => {
                      switch (session.attendanceStatus) {
                        case 'applied':
                          return <CheckCircleOutlineRounded color="primary" />;
                        case 'ignored':
                          return <CheckCircleOutlineRounded color="warning" />;
                        case 'not_applied_yet':
                          return <CheckCircleOutlineRounded color="disabled" />;
                        default:
                          return 'UNKNOWN - Contact dev';
                      }
                    })()}
                  </TableCell>
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

export default AdminSessionList;
