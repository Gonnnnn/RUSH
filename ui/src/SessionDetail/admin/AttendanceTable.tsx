import { useEffect, useState } from 'react';
import { Typography, Paper, Box, CircularProgress, Tabs, Tab } from '@mui/material';
import { useAuth } from '../../auth';
import { Attendance, getSessionAttendances, markUsersAsPresent } from '../../client/http';
import useHandleError from '../../common/error';
import UserAttendance from '../common/UserAttendance';
import AddAttendance from './AddAttendance';

type TabTypes = 'attendance' | 'addAttendance';

/**
 * The attendance table for the session. It handles the attendance data view and the add attendance action.
 * It fetches the attendance and updates it by itself, not by the parent component on purpose
 * to avoid unnecessary props passing, and also to not affect the parent component render when it fails.
 */
const SessionAttendanceTable = ({ sessionId, reloadSession }: { sessionId: string; reloadSession: () => void }) => {
  const { authenticated } = useAuth();
  const { handleError } = useHandleError();

  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const [tab, setTab] = useState<TabTypes>('attendance');

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

  const applyAttendances = async (userIds: string[]) => {
    await markUsersAsPresent(sessionId, userIds);
    const newAttendances = await getSessionAttendances(sessionId);
    reloadSession();
    setAttendances(newAttendances);
  };

  if (!authenticated) {
    return (
      <Paper sx={{ p: 2 }} elevation={4}>
        <Box display="flex" flexDirection="column">
          <Typography variant="h6">출석</Typography>
          <Tabs value={tab} onChange={(_, newTab) => setTab(newTab)} sx={{ mb: 2 }}>
            <Tab label="출석 현황" value="attendance" disabled />
            {/* TODO(#177): Hide it if the user is not admin. */}
            <Tab label="출석 추가" value="addAttendance" disabled />
          </Tabs>
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
        <Tabs value={tab} onChange={(_, newTab) => setTab(newTab)}>
          <Tab label="출석 현황" value="attendance" />
          {/* TODO(#177): Hide it if the user is not admin. */}
          <Tab label="출석 추가" value="addAttendance" />
        </Tabs>

        {tab === 'attendance' && <UserAttendance attendances={attendances} isLoading={isLoading} />}

        {tab === 'addAttendance' && (
          <Box>
            <AddAttendance
              applyAttendances={async (userIds) => {
                await applyAttendances(userIds);
                setTab('attendance');
              }}
            />
          </Box>
        )}
      </Box>
    </Paper>
  );
};

export default SessionAttendanceTable;
