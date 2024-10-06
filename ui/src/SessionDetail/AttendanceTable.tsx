import { useState } from 'react';
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
} from '@mui/material';
import { Attendance } from '../client/http';
import { toYYslashMMslashDDspaceHHcolonMMcolonSS } from '../common/date';

type OrderBy = 'asc' | 'desc';

type OrderKeys = 'userExternalName' | 'userGeneration' | 'userJoinedAt';

const AttendanceTable = ({ isLoading, attendances }: { isLoading: boolean; attendances: Attendance[] }) => {
  const [nameOrder, setNameOrder] = useState<OrderBy>('asc');
  const [generationOrder, setGenerationOrder] = useState<OrderBy>('asc');
  const [joinedAtOrder, setJoinedAtOrder] = useState<OrderBy>('asc');

  const [orderBy, setOrderBy] = useState<OrderKeys>('userExternalName');

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

  return (
    <Paper sx={{ p: 2 }} elevation={4}>
      <Typography variant="h6">출석 제출 목록</Typography>
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

export default AttendanceTable;
