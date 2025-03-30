import { useEffect, useState } from 'react';
import {
  Box,
  Button,
  CircularProgress,
  Container,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Paper,
  Typography,
} from '@mui/material';
import { DataGrid, GridColDef, GridPaginationModel } from '@mui/x-data-grid';
import { adminLateApplyAttendance, adminListSessions } from '../client/http/admin';
import { AdminSession } from '../client/http/data';
import { toYYslashMMslashDDspaceHHcolonMMwithDay } from '../common/date';
import useHandleError from '../common/error';
import { useSnackbar } from '../contexts/snackbar';
import AddAttendance from './SessionDetail/AddAttendance';

const Exception = () => {
  const { handleError } = useHandleError();
  const { showInfo } = useSnackbar();

  const [sessions, setSessions] = useState<AdminSession[]>([]);
  const [isLoadingSessions, setIsLoadingSessions] = useState(false);
  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(null);
  const [showUserSelection, setShowUserSelection] = useState(false);
  const [totalCount, setTotalCount] = useState(0);
  const [isEnd, setIsEnd] = useState(false);
  const [paginationModel, setPaginationModel] = useState<GridPaginationModel>({
    pageSize: 20,
    page: 0,
  });

  useEffect(() => {
    const fetchSessions = async () => {
      try {
        setIsLoadingSessions(true);
        const offset = paginationModel.page * paginationModel.pageSize;
        const fetchedSessions = await adminListSessions(offset, paginationModel.pageSize);
        setSessions(fetchedSessions.sessions);
        setTotalCount(fetchedSessions.totalCount);
        setIsEnd(fetchedSessions.isEnd);
      } catch (error) {
        handleError({
          error,
          messageAuth: 'Session list retrieval requires authentication',
          messageInternal: 'Failed to retrieve session list. Contact the dev.',
        });
      } finally {
        setIsLoadingSessions(false);
      }
    };

    fetchSessions();
  }, [handleError, paginationModel.page, paginationModel.pageSize]);

  const handlePaginationModelChange = (newPaginationModel: GridPaginationModel) => {
    if (isEnd && newPaginationModel.page > paginationModel.page) {
      return;
    }
    setPaginationModel(newPaginationModel);
  };

  const handleSessionSelect = (sessionId: string) => {
    setSelectedSessionId(sessionId);
    setShowUserSelection(true);
  };

  const handleCloseUserSelection = () => {
    setShowUserSelection(false);
  };

  const applyExceptionalAttendance = async (userIds: string[]) => {
    if (!selectedSessionId) return;

    try {
      await adminLateApplyAttendance(selectedSessionId, userIds);
      setShowUserSelection(false);
      showInfo('Exceptional attendance applied successfully');
    } catch (error) {
      handleError({
        error,
        messageAuth: 'Exceptional attendance application is restricted to admins.',
        messageInternal: 'Failed to apply exceptional attendance. Contact the dev.',
      });
    }
  };

  const columns: GridColDef[] = [
    { field: 'name', headerName: 'Name', width: 200 },
    {
      field: 'startsAt',
      headerName: 'Starts At',
      width: 200,
      renderCell: (params) => toYYslashMMslashDDspaceHHcolonMMwithDay(params.row.startsAt),
    },
    {
      field: 'actions',
      headerName: 'Actions',
      width: 200,
      renderCell: (params) => (
        <Button variant="outlined" size="small" onClick={() => handleSessionSelect(params.row.id)}>
          Apply Attendance
        </Button>
      ),
    },
  ];

  return (
    <Container>
      <Paper sx={{ p: 2, mt: 2 }} elevation={4}>
        <Typography variant="h6" gutterBottom>
          Exceptional Attendance
        </Typography>
        <Typography variant="body2" paragraph>
          This page allows admins to apply exceptional attendances for sessions when normal attendance recording was not
          possible.
          <br />
          <br />
          <b>Note:</b> Any attendance applied here will be recorded.
        </Typography>

        {isLoadingSessions ? (
          <Box display="flex" justifyContent="center" p={3}>
            <CircularProgress />
          </Box>
        ) : (
          <Box sx={{ height: 'calc(100vh - 240px)', width: '100%' }}>
            <DataGrid
              rows={sessions}
              columns={columns}
              rowCount={totalCount}
              loading={isLoadingSessions}
              pageSizeOptions={[10, 25, 50]}
              paginationModel={paginationModel}
              paginationMode="server"
              onPaginationModelChange={handlePaginationModelChange}
              disableRowSelectionOnClick
              initialState={{
                sorting: {
                  sortModel: [{ field: 'startsAt', sort: 'desc' }],
                },
              }}
            />
          </Box>
        )}
      </Paper>

      {/* User Selection Dialog */}
      <Dialog open={showUserSelection} onClose={handleCloseUserSelection} fullWidth maxWidth="md">
        <DialogTitle>
          Apply Exceptional Attendance
          {selectedSessionId && (
            <Typography variant="subtitle1">
              Session: {sessions.find((s) => s.id === selectedSessionId)?.name}
            </Typography>
          )}
        </DialogTitle>
        <DialogContent>
          <DialogContentText paragraph>
            Select users who should receive exceptional attendance for this session.
          </DialogContentText>
          <Box sx={{ height: '60vh' }}>
            {selectedSessionId && <AddAttendance applyAttendances={(userIds) => applyExceptionalAttendance(userIds)} />}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseUserSelection}>Cancel</Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default Exception;
