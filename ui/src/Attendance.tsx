import { useState, useEffect } from 'react';
import { Box, Paper, Stack, Typography } from '@mui/material';
import { DataGrid, GridColDef, GridRowsProp } from '@mui/x-data-grid';
import { useHeader } from './Layout';
import { getHalfYearAttendances } from './client/http';

type Attendance = {
  id: string;
  sessionId: string;
  sessionName: string;
  sessionScore: number;
  sessionStartedAt: Date;
  userId: string;
  userName: string;
  userJoinedAt: Date;
  createdAt: Date;
};

type User = {
  id: string;
  name: string;
  generation: number;
};

type Session = {
  id: string;
  name: string;
  startedAt: Date;
};

const HalfYearAttendances = () => {
  useHeader({ newTitle: 'Attendance' });

  const [attendanceData, setAttendanceData] = useState<Attendance[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [sessions, setSessions] = useState<Session[]>([]);
  const [rows, setRows] = useState<GridRowsProp>([]);
  const [columns, setColumns] = useState<GridColDef[]>([]);

  useEffect(() => {
    const init = async () => {
      try {
        const response = await getHalfYearAttendances();
        setAttendanceData(
          response.attendances.map((attendance) => ({
            id: attendance.id,
            sessionId: attendance.sessionId,
            sessionName: attendance.sessionName,
            sessionScore: attendance.sessionScore,
            sessionStartedAt: attendance.sessionStartedAt,
            userId: attendance.userId,
            userName: attendance.userName,
            userJoinedAt: attendance.userJoinedAt,
            createdAt: attendance.createdAt,
          })),
        );
        setSessions(
          response.sessions.map((session) => ({
            id: session.id,
            name: session.name,
            startedAt: session.startedAt,
          })),
        );
        setUsers(
          response.users.map((user) => ({
            id: user.id,
            name: user.name,
            generation: user.generation,
          })),
        );
      } catch (error) {
        // eslint-disable-next-line no-console
        console.error(error);
      }
    };
    init();
  }, []);

  useEffect(() => {
    // Transform the data
    const userMap = new Map<string, User>();
    users.forEach((user) => userMap.set(user.id, user));

    const sessionMap = new Map<string, Session>();
    sessions.forEach((session) => sessionMap.set(session.id, session));

    const transformedRows = users.map((user) => {
      const row: any = {
        id: user.id,
        name: user.name,
        generation: user.generation,
      };
      attendanceData.forEach((attendance) => {
        if (attendance.userId === user.id) {
          row[attendance.sessionId] = attendance.sessionScore > 0 ? attendance.sessionScore : '';
        }
      });
      return row;
    });

    const transformedColumns: GridColDef[] = [
      { field: 'name', headerName: '이름', width: 150, sortable: true },
      { field: 'generation', headerName: '기수', width: 150, sortable: true },
      ...sessions.map((session) => ({
        field: session.id,
        headerName: `${session.name} (${new Date(session.startedAt).toLocaleDateString()})`,
        width: 200,
        sortable: true,
      })),
    ];

    setRows(transformedRows);
    setColumns(transformedColumns);
  }, [attendanceData, users, sessions]);

  return (
    <Stack spacing={2}>
      <Paper sx={{ p: 2 }}>
        <Typography variant="body2">Admin only page to check user attendance.</Typography>
        <Typography variant="body2">Click 이름 to search a specific user.</Typography>
        <Typography variant="body2">You can also sort the rows by 이름 and 유저 기수.</Typography>
        <Typography variant="body2">If necessary, select only certain columns to view.</Typography>
      </Paper>
      <Box style={{ height: '100%', width: '100%', overflow: 'auto' }}>
        {/* paginate by 20 */}
        <DataGrid
          rows={rows}
          columns={columns}
          initialState={{
            pagination: {
              paginationModel: { pageSize: 10, page: 0 },
            },
          }}
        />
      </Box>
    </Stack>
  );
};

export default HalfYearAttendances;
