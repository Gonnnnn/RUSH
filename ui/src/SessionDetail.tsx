import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
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
  Box,
  Button,
} from '@mui/material';
import QrCode from './assets/qr.png';
import { Session, getSession } from './client/http';

const SessionDetail = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const attendanceUploaded = id === '1';
  const [session, setSession] = useState<Session>();
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!id) {
      alert('Invalid session ID');
      navigate('/sessions');
      return;
    }

    const init = async () => {
      try {
        setIsLoading(true);
        setSession(await getSession(id));
      } catch (error) {
        console.error(error);
        navigate('/sessions');
      } finally {
        setIsLoading(false);
      }
    };

    init();
  }, [navigate, id]);

  if (isLoading || !session) {
    return (
      <Container>
        <Typography variant="h4" sx={{ mb: 3 }}>
          Session Detail
        </Typography>
        <Typography>Loading...</Typography>
      </Container>
    );
  }

  return (
    <Container>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
          gap: 1,
          mb: 3,
        }}
      >
        <Typography variant="h4">Session Detail</Typography>
        <Button variant="outlined" onClick={() => navigate('/sessions')} sx={{ alignSelf: 'flex-start' }}>
          Back
        </Button>
      </Box>
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="h6">Details</Typography>
        <Typography>Name: {session.name}</Typography>
        <Typography>Description: {session.description}</Typography>
        <Typography>Start Time: {session.startsAt.toISOString()}</Typography>
        <Typography>Score: {session.score}</Typography>
        <Typography>Created At: {session.createdAt.toISOString()}</Typography>
      </Paper>
      <Typography variant="h6" sx={{ mb: 5 }}>
        User attendance
      </Typography>
      {attendanceUploaded ? (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>User Name</TableCell>
                <TableCell>User Gen</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {session.joinningUsers.map((user: string) => (
                <TableRow key={user}>
                  <TableCell>{user}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      ) : (
        <>
          <Paper sx={{ p: 2, mb: 3 }}>
            <Typography variant="h6">QR to the form</Typography>
            <Typography>Scan the QR code to fill the attendance form</Typography>
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
              <img src={QrCode} alt="QR Code" height="200" />
            </Box>
          </Paper>
          <Paper sx={{ p: 2, mb: 3 }}>
            <Typography variant="h6">Upload Attendance</Typography>
            <Typography>Drag and drop your attendance file here</Typography>
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
              <Button variant="contained" color="primary">
                Upload
              </Button>
            </Box>
          </Paper>
        </>
      )}
    </Container>
  );
};

export default SessionDetail;
