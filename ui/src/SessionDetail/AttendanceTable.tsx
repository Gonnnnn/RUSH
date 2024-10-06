import { useState } from 'react';
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
  TableSortLabel,
} from '@mui/material';
import { Attendance } from '../client/http';
import { toYYslashMMslashDDspaceHHcolonMMcolonSS } from '../common/date';

type OrderBy = 'asc' | 'desc';

type OrderKeys = 'userExternalName' | 'userGeneration' | 'userJoinedAt';

const AttendanceTable = ({ isLoading, attendances }: { isLoading: boolean; attendances: Attendance[] }) => {
  const [order, setOrder] = useState<OrderBy>('asc');
  const [orderBy, setOrderBy] = useState<OrderKeys>('userJoinedAt');

  const handleSort = (property: OrderKeys) => {
    const isAsc = orderBy === property && order === 'asc';
    setOrder(isAsc ? 'desc' : 'asc');
    setOrderBy(property);
  };

  const sortedAttendances = attendances.slice().sort((a, b) => {
    switch (orderBy) {
      case 'userExternalName':
        return (a.userExternalName < b.userExternalName ? -1 : 1) * (order === 'asc' ? 1 : -1);
      case 'userGeneration':
        return (a.userGeneration < b.userGeneration ? -1 : 1) * (order === 'asc' ? 1 : -1);
      case 'userJoinedAt':
        return (a.userJoinedAt < b.userJoinedAt ? -1 : 1) * (order === 'asc' ? 1 : -1);
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
    // make it scrollable
    <Paper sx={{ p: 2 }} elevation={4}>
      <Typography variant="h6">출석 제출 목록</Typography>
      <TableContainer sx={{ overflowY: 'auto', maxHeight: 400 }}>
        <Table>
          {/* make the header sticky */}
          <TableHead sx={{ position: 'sticky', top: 0, backgroundColor: 'background.paper' }}>
            <TableRow>
              <TableCell align="center" sx={{ width: '30%' }}>
                <TableSortLabel
                  active={orderBy === 'userExternalName'}
                  direction={orderBy === 'userExternalName' ? order : 'asc'}
                  onClick={() => handleSort('userExternalName')}
                >
                  이름
                </TableSortLabel>
              </TableCell>
              <TableCell align="center" sx={{ width: '30%' }}>
                <TableSortLabel
                  active={orderBy === 'userGeneration'}
                  direction={orderBy === 'userGeneration' ? order : 'asc'}
                  onClick={() => handleSort('userGeneration')}
                >
                  기수
                </TableSortLabel>
              </TableCell>
              <TableCell align="center" sx={{ width: '40%' }}>
                <TableSortLabel
                  active={orderBy === 'userJoinedAt'}
                  direction={orderBy === 'userJoinedAt' ? order : 'asc'}
                  onClick={() => handleSort('userJoinedAt')}
                >
                  제출 시간
                </TableSortLabel>
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

export default AttendanceTable;
