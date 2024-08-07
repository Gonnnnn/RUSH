import { useCallback, useEffect, useState } from 'react';
import { CheckCircleOutlineRounded } from '@mui/icons-material';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Container,
  TablePagination,
  Box,
  LinearProgress,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import { useHeader } from './Layout';
import { User, listUsers } from './client/http';

const UserList = () => {
  useHeader({ newTitle: 'Users' });
  const theme = useTheme();
  // TODO(#31): Centralize the isMobile logic.
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const pageSize = isMobile ? 8 : 10;
  const [users, setUsers] = useState<User[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [currentPage, setCurrentPage] = useState(0);
  const [isEnd, setIsEnd] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const fetchUsers = useCallback(
    async (page: number) => {
      try {
        setIsLoading(true);
        const offset = page * pageSize;
        const listUsersResponse = await listUsers(offset, pageSize);
        setUsers(listUsersResponse.users);
        setIsEnd(listUsersResponse.isEnd);
        setTotalCount(listUsersResponse.totalCount);
        setCurrentPage(page);
      } catch (e) {
        console.error(e);
      } finally {
        setIsLoading(false);
      }
    },
    [pageSize],
  );

  useEffect(() => {
    fetchUsers(0);
  }, [fetchUsers]);

  const handleChangePage = async (_: unknown, newPage: number) => {
    if (isEnd && newPage > currentPage) {
      return;
    }
    fetchUsers(newPage);
  };

  return (
    <Container>
      <TableContainer component={Paper}>
        <Box sx={{ width: '100%', height: '4px', mb: 2 }}>{isLoading ? <LinearProgress /> : null}</Box>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell align="center">이름</TableCell>
              <TableCell align="center">기수</TableCell>
              <TableCell align="center">활동 중</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {users.map((user) => (
              <TableRow key={user.id}>
                <TableCell align="center">{user.name}</TableCell>
                <TableCell align="center">{user.generation}</TableCell>
                <TableCell align="center">
                  <CheckCircleOutlineRounded color={user.isActive ? 'primary' : 'disabled'} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={[]}
          component="div"
          rowsPerPage={pageSize}
          page={currentPage}
          onPageChange={handleChangePage}
          count={totalCount}
          slotProps={{
            actions: {
              previousButton: { disabled: currentPage === 0 },
              nextButton: { disabled: isEnd },
            },
          }}
        />
      </TableContainer>
    </Container>
  );
};

export default UserList;
