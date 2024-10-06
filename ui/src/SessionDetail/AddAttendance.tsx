import { useEffect, useState } from 'react';
import {
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Typography,
} from '@mui/material';
import { DataGrid, GridRowSelectionModel, GridToolbarContainer, GridToolbarQuickFilter } from '@mui/x-data-grid';
import { listUsers, User } from '../client/http';
import useHandleError from '../common/error';

/**
 * It handles the addition of attendances. As it has to fetch all the users,
 * the user fetching logic is implemented here so that it wouldn't be fetched when the parent component
 * doesn't need this component.
 */
const AddAttendance = ({ applyAttendances }: { applyAttendances: (userIds: string[]) => Promise<void> }) => {
  const { handleError } = useHandleError();

  const [users, setUsers] = useState<User[]>([]);
  const [isLoadingUsers, setIsLoadingUsers] = useState(true);
  const [isApplyingAttendances, setIsApplyingAttendances] = useState(false);
  const [selectedUserIds, setSelectedUserIds] = useState<string[]>([]);
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);

  const [showOnlySelected, setShowOnlySelected] = useState(false);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        setIsLoadingUsers(true);
        // TODO(#177): Implement getAllUsers.
        const res = await listUsers(1, 9999999);
        setUsers(res.users);
      } catch (error) {
        handleError({
          error,
          messageAuth: 'Failed to fetch users. Contact the dev.',
          messageInternal: 'Failed to fetch users. Contact the dev.',
        });
      } finally {
        setIsLoadingUsers(false);
      }
    };
    fetchUsers();
  }, [handleError]);

  const handleSelectionChange = (newSelectionModel: GridRowSelectionModel) => {
    setSelectedUserIds(newSelectionModel.map((id) => id.toString()));
  };

  const handleApplyButtonClick = () => {
    setConfirmDialogOpen(true);
  };

  const handleConfirmDialogClose = () => {
    setConfirmDialogOpen(false);
  };

  const handleConfirmDialogConfirm = async () => {
    try {
      setIsApplyingAttendances(true);
      await applyAttendances(selectedUserIds);
      setConfirmDialogOpen(false);
    } catch (error) {
      handleError({
        error,
        messageAuth: 'Manual attendance application is restricted to admins.',
        messageInternal: 'Failed to apply attendances. Contact the dev.',
      });
    } finally {
      setIsApplyingAttendances(false);
    }
  };

  const columns = [
    { field: 'name', headerName: '이름', width: 150 },
    { field: 'generation', headerName: '기수', width: 150 },
  ];

  const usersToShow = showOnlySelected ? users.filter((user) => selectedUserIds.includes(user.id)) : users;

  return (
    <Box display="flex" flexDirection="column" gap={2} p={2}>
      <Box display="flex" justifyContent="flex-end" alignItems="center">
        <Button variant="outlined" onClick={handleApplyButtonClick}>
          <Typography variant="body2">Add new attendances</Typography>
        </Button>
      </Box>
      <DataGrid
        rows={usersToShow}
        columns={columns}
        loading={isLoadingUsers}
        disableColumnFilter
        initialState={{
          pagination: {
            paginationModel: { pageSize: 10, page: 0 },
          },
        }}
        checkboxSelection
        rowSelectionModel={selectedUserIds}
        onRowSelectionModelChange={handleSelectionChange}
        slots={{
          toolbar: () => CustomToolbar(showOnlySelected, () => setShowOnlySelected((prev) => !prev)),
        }}
        sx={{ width: '100%', height: '70vh' }}
      />
      <ConfirmDialog
        isConfirming={isApplyingAttendances}
        open={confirmDialogOpen}
        onClose={handleConfirmDialogClose}
        onConfirm={handleConfirmDialogConfirm}
      />
    </Box>
  );
};

const CustomToolbar = (showOnlySelected: boolean, handleButtonClick: () => void) => (
  <GridToolbarContainer sx={{ pt: 2, pr: 2, pl: 2 }}>
    <GridToolbarQuickFilter />
    <Button variant="outlined" onClick={handleButtonClick}>
      <Typography variant="body2">{showOnlySelected ? 'Show all' : 'Show selected'}</Typography>
    </Button>
  </GridToolbarContainer>
);

const ConfirmDialog = ({
  isConfirming,
  open,
  onClose,
  onConfirm,
}: {
  isConfirming: boolean;
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
}) => (
  <Dialog open={open} onClose={onClose}>
    <DialogTitle>출석 추가</DialogTitle>
    <DialogContent>
      <DialogContentText>정말로 출석을 추가하시겠습니까?</DialogContentText>
      <br />
      <DialogContentText sx={{ fontSize: '0.875rem' }}>
        * 이미 해당 세션에 출석한 유저가 포함된 경우, 해당 유저들은 제외된 후 출석이 추가됩니다.
      </DialogContentText>
    </DialogContent>
    <DialogActions>
      <Button onClick={onClose} disabled={isConfirming}>
        Cancel
      </Button>
      <Button onClick={onConfirm} color="primary" disabled={isConfirming}>
        {isConfirming ? <CircularProgress size={24} /> : 'Confirm'}
      </Button>
    </DialogActions>
  </Dialog>
);

export default AddAttendance;
