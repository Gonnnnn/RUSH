import { useEffect, useState } from 'react';
import { Typography, Paper, Box, CircularProgress } from '@mui/material';
import { useAuth } from '../../auth';
import { Attendance } from '../../client/http/data';
import { getSessionAttendances } from '../../client/http/default';
import useHandleError from '../../common/error';
import UserAttendance from '../common/UserAttendance';

/**
 * The attendance table for the session. It handles the attendance data view and the add attendance action.
 * It fetches the attendance and updates it by itself, not by the parent component on purpose
 * to avoid unnecessary props passing, and also to not affect the parent component render when it fails.
 */
const SessionAttendanceTable = ({ sessionId }: { sessionId: string }) => {
  const { authenticated } = useAuth();
  const { handleError } = useHandleError();

  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const fetchAttendances = async () => {
      if (!authenticated) {
        return;
      }
      try {
        setIsLoading(true);
        const fetchedAttendances = await getSessionAttendances(sessionId);
        setAttendances(fetchedAttendances);
      } catch (error) {
        handleError({
          error,
          messageAuth: 'Requires login.',
          messageInternal: 'Failed to load attendance list. Contact the dev.',
        });
      } finally {
        setIsLoading(false);
      }
    };

    fetchAttendances();
  }, [sessionId, authenticated, handleError]);

  if (!authenticated) {
    return (
      <Paper sx={{ p: 2 }} elevation={4}>
        <Box display="flex" flexDirection="column">
          <Typography variant="h6">출석</Typography>
          <Typography variant="body1">로그인이 필요한 서비스입니다.</Typography>
        </Box>
      </Paper>
    );
  }

  if (isLoading) {
    return (
      <Paper sx={{ p: 2 }} elevation={4}>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress />
        </Box>
      </Paper>
    );
  }

  return (
    <Paper sx={{ p: 2 }} elevation={4}>
      <Box display="flex" flexDirection="column">
        <Typography variant="h6">출석</Typography>
        <UserAttendance attendances={attendances} isLoading={isLoading} />
      </Box>
    </Paper>
  );
};

export default SessionAttendanceTable;
