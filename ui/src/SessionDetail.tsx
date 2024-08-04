import { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { Container, Typography, Paper, Box, Button, CircularProgress } from '@mui/material';
import { AxiosError } from 'axios';
import { QRCodeCanvas } from 'qrcode.react';
import { useSnackbar } from './SnackbarContex';
import { Session, createSessionForm, getSession } from './client/http';
import toYYYY년MM월DD일HH시MM분 from './common/date';

const SessionDetail = () => {
  const navigate = useNavigate();
  const { showWarning, showError } = useSnackbar();
  const { id } = useParams();
  const [session, setSession] = useState<Session>();
  const [isLoading, setIsLoading] = useState(false);
  const [isCreatingForm, setIsCreatingForm] = useState(false);
  const qrRef = useRef<HTMLDivElement>(null);
  const qrSizePx = 128;

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
    } catch (error: unknown) {
      if (error instanceof AxiosError && error.response?.status === 401) {
        showWarning('Form creation is restricted to authenticated users');
      } else {
        showError('Failed to create a form. Contact the administrator.');
      }
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
        <Typography>Starts At: {toYYYY년MM월DD일HH시MM분(session.startsAt)}</Typography>
        <Typography>Score: {session.score}</Typography>
        <Typography>Created At: {toYYYY년MM월DD일HH시MM분(session.createdAt)}</Typography>
      </Paper>

      <Paper sx={{ p: 2, mb: 3 }}>
        {session.googleFormUri ? (
          <>
            <Typography variant="h6">QR code to the form</Typography>
            <Box
              sx={{
                display: 'flex',
                flexDirection: 'column',
                gap: 2,
                justifyContent: 'center',
                alignItems: 'center',
                mt: 2,
                mb: 2,
              }}
            >
              {/* TODO(#8): Replace the value with the actual form URL. */}
              <div ref={qrRef}>
                <QRCodeCanvas value={session.googleFormUri} size={qrSizePx} />
              </div>
              <Button
                variant="outlined"
                onClick={() => onQrDownload(qrRef, qrSizePx, toYYYY년MM월DD일HH시MM분(session.startsAt))}
              >
                Download QR code
              </Button>
              <Button variant="outlined" onClick={() => window.open(session.googleFormUri, '_blank')}>
                Open the form
              </Button>
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

// https://github.com/zpao/qrcode.react/issues/233
const onQrDownload = (qrRef: React.RefObject<HTMLDivElement>, qrSize: number, text: string) => {
  if (!qrRef.current) {
    return;
  }

  const canvas = qrRef.current.querySelector('canvas');
  if (!canvas) {
    return;
  }

  const newCanvas = document.createElement('canvas');
  const ctx = newCanvas.getContext('2d');
  if (!ctx) {
    return;
  }

  const newWidth = qrSize + 512;
  const newHeight = qrSize + 630;

  newCanvas.width = newWidth;
  newCanvas.height = newHeight;
  ctx.fillStyle = 'white';
  ctx.fillRect(0, 0, newCanvas.width, newCanvas.height);

  const newQrSize = 256;
  const qrYoffset = (newHeight - newQrSize) / 2;
  const qrXoffset = (newWidth - newQrSize) / 2;
  ctx.drawImage(canvas, qrXoffset, qrYoffset, 256, 256);

  ctx.font = '32px Arial';
  ctx.fillStyle = 'black';
  ctx.textAlign = 'center';
  ctx.fillText(text, newWidth / 2, Math.min(qrYoffset + newQrSize + 128, newHeight - 32));

  const a = document.createElement('a');
  a.href = newCanvas.toDataURL('image/png');
  // replace all spaces with underscores and add .png extension
  a.download = `${text.replace(/ /g, '_')}.png`;
  a.click();
};

export default SessionDetail;
