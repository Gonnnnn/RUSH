import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Container, Typography, Paper, Box, Button, CircularProgress } from '@mui/material';
import QRCode from 'qrcode.react';
import { Session, createSessionForm, getSession } from './client/http';

const SessionDetail = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [session, setSession] = useState<Session>();
  const [isLoading, setIsLoading] = useState(false);
  const [isCreatingForm, setIsCreatingForm] = useState(false);

  useEffect(() => {
    if (!id) {
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

  if (!id) {
    navigate('/sessions');
    return null;
  }

  const onQrCodeCreateClick = async () => {
    try {
      setIsCreatingForm(true);
      await createSessionForm(id);
      setSession(await getSession(id));
    } catch (error) {
      console.error(error);
    } finally {
      setIsCreatingForm(false);
    }
  };

  if (isLoading || !session) {
    return (
      <Container>
        <Typography variant="h4" sx={{ mb: 3 }}>
          Session Detail
        </Typography>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress />
        </Box>
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

      <Paper sx={{ p: 2, mb: 3 }}>
        {session.googleFormUri ? (
          <>
            <Typography variant="h6">QR code to the form</Typography>
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2, mb: 2 }}>
              {/* TODO(#8): Replace the value with the actual form URL. */}
              <QRCode value={session.googleFormUri} />
            </Box>
          </>
        ) : (
          <>
            <Typography variant="h6">No form is associated</Typography>
            <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2, mb: 2 }}>
              <Button variant="contained" onClick={onQrCodeCreateClick} disabled={isCreatingForm}>
                {isCreatingForm ? <CircularProgress /> : 'Create QR code'}
              </Button>
            </Box>
          </>
        )}
      </Paper>
    </Container>
  );
};

export default SessionDetail;
