import React from 'react';
import {
  Container,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from '@mui/material';

// Mock Data
interface User {
  name: string;
  gen: number;
}

interface Session {
  id: string;
  name: string;
  date: string;
}

interface Attendance {
  id: string;
  name: string;
  createdAt: Date;
  startDate: Date;
  endDate: Date;
  users: User[];
  sessions: Session[];
  attendanceRecords: { [userId: string]: { [sessionId: string]: 'v' | '' } };
}

const AttendanceDetail: React.FC = () => {
  const detail: Attendance = makeMockAttendanceRecords(mockUsers, mockSessions);

  return (
    <Container>
      <Typography variant="h4" sx={{ mb: 3 }}>
        Attendance Detail
      </Typography>
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="h6">Details</Typography>
        <Typography>ID: {detail.id}</Typography>
        <Typography>Name: {detail.name}</Typography>
        <Typography>Created At: {detail.createdAt.toDateString()}</Typography>
        <Typography>Start Date: {detail.startDate.toDateString()}</Typography>
        <Typography>End Date: {detail.endDate.toDateString()}</Typography>
      </Paper>
      <TableContainer component={Paper}>
        <Table stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell>User Name</TableCell>
              <TableCell>User Gen</TableCell>
              {detail.sessions.map((session) => (
                <TableCell key={session.id}>
                  {session.name}
                  <br />
                  {session.date}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {detail.users.map((user, userIndex) => (
              <TableRow key={userIndex}>
                <TableCell>{user.name}</TableCell>
                <TableCell>{user.gen}</TableCell>
                {detail.sessions.map((session) => (
                  <TableCell key={session.id}>{detail.attendanceRecords[user.name][session.id]}</TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Container>
  );
};

const makeMockAttendanceRecords = (users: User[], sessions: Session[]) => {
  const attendanceRecords: { [userId: string]: { [sessionId: string]: 'v' | '' } } = {};
  users.forEach((user) => {
    attendanceRecords[user.name] = {};
    sessions.forEach((session) => {
      attendanceRecords[user.name][session.id] = Math.random() < 0.5 ? 'v' : '';
    });
  });

  return {
    id: '1',
    name: 'Rush 2023 상반기',
    createdAt: new Date('2023-02-01'),
    startDate: new Date('2023-03-01'),
    endDate: new Date('2024-08-14'),
    users,
    sessions,
    attendanceRecords,
  };
};

const mockUsers = [
  { name: 'User1', gen: 1 },
  { name: 'User2', gen: 2 },
  { name: 'User3', gen: 3 },
  { name: 'User4', gen: 4 },
  { name: 'User5', gen: 5 },
  { name: 'User6', gen: 6 },
  { name: 'User7', gen: 7 },
  { name: 'User8', gen: 8 },
  { name: 'User9', gen: 9 },
  { name: 'User10', gen: 10 },
  { name: 'User11', gen: 11 },
  { name: 'User12', gen: 12 },
  { name: 'User13', gen: 13 },
  { name: 'User14', gen: 14 },
  { name: 'User15', gen: 15 },
];
const mockSessions = [
  { id: '1', name: '한강공원', date: '2023-03-01' },
  { id: '2', name: '러쉬마라톤', date: '2023-04-01' },
  { id: '3', name: '어쩌고', date: '2023-05-01' },
  { id: '4', name: '저쩌고', date: '2023-06-01' },
  { id: '5', name: '대충가짜세션1', date: '2023-07-01' },
  { id: '6', name: '대충가짜세션2', date: '2023-08-01' },
  { id: '7', name: '대충가짜세션3', date: '2023-09-01' },
  { id: '8', name: '대충가짜세션4', date: '2023-10-01' },
  { id: '9', name: '대충가짜세션5', date: '2023-11-01' },
  { id: '10', name: '대충가짜세션6', date: '2023-12-01' },
  { id: '11', name: '대충가짜세션7', date: '2024-01-01' },
  { id: '12', name: '대충가짜세션8', date: '2024-02-01' },
  { id: '13', name: '대충가짜세션9', date: '2024-03-01' },
  { id: '14', name: '대충가짜세션10', date: '2024-04-01' },
  { id: '15', name: '대충가짜세션11', date: '2024-05-01' },
];

export default AttendanceDetail;
