import React, { useEffect, useState } from 'react';
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
} from '@mui/material';
import { Attendance, getAttendances } from './client/http';

const AttendanceList: React.FC = () => {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [isLoading, setIsLoading] = useState(false);
  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    const init = async () => {
      try {
        setIsLoading(true);
        const attendances = await getAttendances();
        setAttendances(attendances);
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
          Attendance List
        </Typography>
        <Typography>Loading...</Typography>
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

  const handleRowClick = (attendance: Attendance) => {
    navigate(`/attendances/${attendance.id}`);
  };

  const paginatedAttendances = attendances.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);

  return (
    <Container>
      <Typography variant="h4" sx={{ mb: 5 }}>
        Attendance List
      </Typography>
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Created At</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedAttendances.map((attendance, index) => (
              <TableRow key={index} onClick={() => handleRowClick(attendance)} style={{ cursor: 'pointer' }}>
                <TableCell>{attendance.name}</TableCell>
                <TableCell>{attendance.createdAt.toISOString()}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={[10, 20, 30]}
          component="div"
          count={attendances.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </TableContainer>
    </Container>
  );
};

export default AttendanceList;
