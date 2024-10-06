import { useEffect, useState } from 'react';
import { ArrowDownward, ArrowUpward } from '@mui/icons-material';
import {
  Typography,
  Paper,
  Box,
  CircularProgress,
  TableContainer,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Tabs,
  Tab,
} from '@mui/material';
import { useAuth } from '../AuthContext';
import { Attendance, getSessionAttendances, markUsersAsPresent } from '../client/http';
import { toYYslashMMslashDDspaceHHcolonMMcolonSS } from '../common/date';
import useHandleError from '../common/error';
import AddAttendance from './AddAttendance';

type OrderBy = 'asc' | 'desc';

type OrderKeys = 'userExternalName' | 'userGeneration' | 'userJoinedAt';

type TabTypes = 'attendance' | 'addAttendance';

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

  const [nameOrder, setNameOrder] = useState<OrderBy>('asc');
  const [generationOrder, setGenerationOrder] = useState<OrderBy>('asc');
  const [joinedAtOrder, setJoinedAtOrder] = useState<OrderBy>('asc');
  const [orderBy, setOrderBy] = useState<OrderKeys>('userExternalName');
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
        console.error(error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchAttendances();
  }, [sessionId, authenticated, handleError]);

  const applyAttendances = async (userIds: string[]) => {
    await markUsersAsPresent(sessionId, userIds);
    const newAttendances = await getSessionAttendances(sessionId);
    setAttendances(newAttendances);
  };

  const onSortChange = (newOrderBy: OrderKeys) => {
    switch (newOrderBy) {
      case 'userExternalName':
        setNameOrder(oppositeOrder(nameOrder));
        setOrderBy(newOrderBy);
        break;
      case 'userGeneration':
        setGenerationOrder(oppositeOrder(generationOrder));
        setOrderBy(newOrderBy);
        break;
      case 'userJoinedAt':
        setJoinedAtOrder(oppositeOrder(joinedAtOrder));
        setOrderBy(newOrderBy);
        break;
      default:
        break;
    }
  };

  const oppositeOrder = (order: OrderBy) => (order === 'asc' ? 'desc' : 'asc');

  const sortedAttendances = attendances.slice().sort((a, b) => {
    switch (orderBy) {
      case 'userExternalName':
        return (a.userExternalName < b.userExternalName ? -1 : 1) * (nameOrder === 'asc' ? 1 : -1);
      case 'userGeneration':
        return (a.userGeneration < b.userGeneration ? -1 : 1) * (generationOrder === 'asc' ? 1 : -1);
      case 'userJoinedAt':
        return (a.userJoinedAt < b.userJoinedAt ? -1 : 1) * (joinedAtOrder === 'asc' ? 1 : -1);
      default:
        return 0;
    }
  });

  if (isLoading) {
    return (
      <Paper sx={{ p: 2 }} elevation={4}>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress />
        </Box>
      </Paper>
    );
  }

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

  return (
    <Paper sx={{ p: 2 }} elevation={4}>
      <Box display="flex" flexDirection="column">
        <Typography variant="h6">출석</Typography>
        <Tabs value={tab} onChange={(_, newTab) => setTab(newTab)}>
          <Tab label="출석 현황" value="attendance" />
          {/* TODO(#177): Hide it if the user is not admin. */}
          <Tab label="출석 추가" value="addAttendance" />
        </Tabs>

        {tab === 'attendance' && (
          <TableContainer sx={{ overflowY: 'auto', maxHeight: 400 }}>
            <Table>
              <TableHead sx={{ position: 'sticky', top: 0, backgroundColor: 'background.paper' }}>
                <TableRow>
                  <TableCell align="center" sx={{ width: '30%' }} onClick={() => onSortChange('userExternalName')}>
                    <Box display="flex" alignItems="center" gap={1}>
                      이름
                      <OrderArrows
                        active={orderBy === 'userExternalName'}
                        order={nameOrder}
                        onClick={() => onSortChange('userExternalName')}
                      />
                    </Box>
                  </TableCell>
                  <TableCell align="center" sx={{ width: '30%' }} onClick={() => onSortChange('userGeneration')}>
                    <Box display="flex" alignItems="center" gap={1}>
                      기수
                      <OrderArrows
                        active={orderBy === 'userGeneration'}
                        order={generationOrder}
                        onClick={() => onSortChange('userGeneration')}
                      />
                    </Box>
                  </TableCell>
                  <TableCell align="center" sx={{ width: '40%' }} onClick={() => onSortChange('userJoinedAt')}>
                    <Box display="flex" alignItems="center" gap={1}>
                      제출 시간
                      <OrderArrows
                        active={orderBy === 'userJoinedAt'}
                        order={joinedAtOrder}
                        onClick={() => onSortChange('userJoinedAt')}
                      />
                    </Box>
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {sortedAttendances.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={3}>출석 제출 목록이 없습니다.</TableCell>
                  </TableRow>
                ) : (
                  sortedAttendances.map((attendance) => (
                    <TableRow key={attendance.id}>
                      <TableCell align="center" sx={{ width: '30%' }}>
                        {attendance.userExternalName}
                      </TableCell>
                      <TableCell align="center" sx={{ width: '30%' }}>
                        {attendance.userGeneration}
                      </TableCell>
                      <TableCell align="center" sx={{ width: '40%' }}>
                        {toYYslashMMslashDDspaceHHcolonMMcolonSS(attendance.userJoinedAt)}
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}

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

const OrderArrows = ({ active, order, onClick }: { active: boolean; order: OrderBy; onClick: () => void }) => (
  <Box display="flex" alignItems="center" onClick={onClick}>
    {order === 'asc' ? (
      <ArrowUpward color={active ? 'primary' : 'action'} sx={{ width: 16, height: 16, p: 0 }} />
    ) : (
      <ArrowDownward color={active ? 'primary' : 'action'} sx={{ width: 16, height: 16, p: 0 }} />
    )}
  </Box>
);

export default SessionAttendanceTable;
