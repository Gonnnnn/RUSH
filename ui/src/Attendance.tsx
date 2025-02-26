import { useState, useEffect } from 'react';
import { Box, Button, Stack } from '@mui/material';
import {
  DataGrid,
  GridColDef,
  GridRowsProp,
  GridToolbarColumnsButton,
  GridToolbarContainer,
  GridToolbarQuickFilter,
} from '@mui/x-data-grid';
import * as XLSX from 'xlsx';
import { useHeader } from './Layout';
import { getHalfYearAttendances } from './client/http/default';
import { toYYslashMMslashDDspaceHHcolonMM } from './common/date';

type Attendance = {
  id: string;
  sessionId: string;
  sessionName: string;
  sessionScore: number;
  sessionStartedAt: Date;
  userId: string;
  userExternalName: string;
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

type Row = User & {
  [sessionId: string]: string | number;
  totalScore: number;
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
            userExternalName: attendance.userExternalName,
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
    const userMap = new Map<string, User>();
    users.forEach((user) => userMap.set(user.id, user));

    const sessionMap = new Map<string, Session>();
    sessions.forEach((session) => sessionMap.set(session.id, session));

    const transformedRows: Row[] = users.map((user) => {
      const row: Row = {
        id: user.id,
        name: user.name,
        generation: user.generation,
        totalScore: 0,
      };
      attendanceData.forEach((attendance) => {
        if (attendance.userId === user.id) {
          const score = attendance.sessionScore > 0 ? attendance.sessionScore : 0;
          row[attendance.sessionId] = score;
          row.totalScore += score;
        }
      });
      return row;
    });

    const transformedColumns: GridColDef[] = [
      { field: 'name', headerName: '이름', width: 150, sortable: true },
      { field: 'generation', headerName: '기수', width: 100, sortable: true },
      { field: 'totalScore', headerName: '총점', width: 100, sortable: true },
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
      <Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Button variant="outlined" onClick={() => exportToExcel(attendanceData, users, sessions)}>
          Download
        </Button>
      </Box>
      <DataGrid
        rows={rows}
        columns={columns}
        slots={{
          toolbar: CustomToolbar,
        }}
        disableColumnFilter
        initialState={{
          pagination: {
            paginationModel: { pageSize: 20, page: 0 },
          },
        }}
      />
    </Stack>
  );
};

const CustomToolbar = () => (
  <GridToolbarContainer sx={{ pt: 2, pr: 2, pl: 2 }}>
    <GridToolbarQuickFilter />
    <Box sx={{ flexGrow: 1 }} />
    <GridToolbarColumnsButton slotProps={{ button: { children: <div>hi</div> } }} />
  </GridToolbarContainer>
);

const exportToExcel = (attendanceData: Attendance[], users: User[], sessions: Session[]) => {
  const data = toExcelFormat(attendanceData, users, sessions);
  const ws = XLSX.utils.aoa_to_sheet(data);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, '출석');

  const now = new Date();
  const filename = `출석_${toYYslashMMslashDDspaceHHcolonMM(now)}.xlsx`;
  XLSX.writeFile(wb, filename);
};

const toExcelFormat = (attendanceData: Attendance[], users: User[], sessions: Session[]) => {
  const headers = [
    ['이름', '기수', '총점', ...sessions.map((session) => session.name)],
    ['', '', '', ...sessions.map((session) => toYYslashMMslashDDspaceHHcolonMM(session.startedAt))],
  ];

  const rows = users.map((user) => {
    const row = [user.name, user.generation];
    const totalScore = sessions.reduce((acc, session) => {
      const attendance = attendanceData.find((a) => a.userId === user.id && a.sessionId === session.id);
      return acc + (attendance && attendance.sessionScore > 0 ? attendance.sessionScore : 0);
    }, 0);
    const sessionScores = sessions.map((session) => {
      const attendance = attendanceData.find((a) => a.userId === user.id && a.sessionId === session.id);
      const score = attendance && attendance.sessionScore > 0 ? attendance.sessionScore : 0;
      return score || '';
    });
    row.push(totalScore, ...sessionScores);
    return row;
  });

  return [...headers, ...rows];
};

export default HalfYearAttendances;
