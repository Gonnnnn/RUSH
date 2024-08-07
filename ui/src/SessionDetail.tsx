import { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { CalendarTodayOutlined, StarBorderRounded } from '@mui/icons-material';
import { Container, Typography, Paper, Box, Button, CircularProgress, Stack, Grid } from '@mui/material';
import { TimeIcon } from '@mui/x-date-pickers';
import { AxiosError } from 'axios';
import { QRCodeCanvas } from 'qrcode.react';
import { useHeader } from './Layout';
import { useSnackbar } from './SnackbarContext';
import { Session, closeSession, createSessionForm, getSession } from './client/http';
import { toYYslashMMslashDDspaceHHcolonMM, toYYYY년MM월DD일HH시MM분 } from './common/date';

const SessionDetail = () => {
  useHeader({ newTitle: 'Session Detail' });
  const navigate = useNavigate();
  const { showWarning, showError } = useSnackbar();
  const { id } = useParams();

  const [session, setSession] = useState<Session | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isCreatingForm, setIsCreatingForm] = useState(false);
  const [isClosingSession, setIsClosingSession] = useState(false);
  const qrRef = useRef<HTMLDivElement>(null);
  const qrSizePx = 128;

  useEffect(() => {
    if (!id) {
      navigate('/sessions');
      return;
    }

    const fetchSession = async () => {
      try {
        setIsLoading(true);
        const fetchedSession = await getSession(id);
        setSession(fetchedSession);
      } catch (error) {
        console.error(error);
        navigate('/sessions');
      } finally {
        setIsLoading(false);
      }
    };

    fetchSession();
  }, [navigate, id]);

  if (!id) {
    navigate('/sessions');
    return null;
  }

  const handleQrCodeCreateClick = async () => {
    try {
      setIsCreatingForm(true);
      await createSessionForm(id);
      const updatedSession = await getSession(id);
      setSession(updatedSession);
    } catch (error) {
      handleError(
        error,
        'Form creation is restricted to authenticated users',
        'Failed to create a form. Contact the administrator.',
      );
    } finally {
      setIsCreatingForm(false);
    }
  };

  const handleCloseSessionBtnClick = async () => {
    try {
      setIsClosingSession(true);
      await closeSession(id);
      const updatedSession = await getSession(id);
      setSession(updatedSession);
    } catch (error) {
      handleError(
        error,
        'Form closing is restricted to authenticated users',
        'Failed to close the form. Contact the administrator.',
      );
    } finally {
      setIsClosingSession(false);
    }
  };

  const handleError = (error: unknown, warningMessage: string, errorMessage: string) => {
    if (error instanceof AxiosError && error.response?.status === 401) {
      showWarning(warningMessage);
    } else {
      showError(errorMessage);
    }
  };

  if (isLoading || !session) {
    return (
      <Container>
        <Box display="flex" justifyContent="center" alignItems="center">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  return (
    <Container>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 3 }}>
        <Button variant="outlined" onClick={() => navigate('/sessions')} sx={{ alignSelf: 'flex-start' }}>
          Back
        </Button>
      </Box>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
        <SessionInfo session={session} />
        <AttendanceQrPanel
          session={session}
          qrRef={qrRef}
          qrSizePx={qrSizePx}
          isCreatingForm={isCreatingForm}
          isClosingSession={isClosingSession}
          onCreateQRCode={handleQrCodeCreateClick}
          onCloseSession={handleCloseSessionBtnClick}
        />
      </Box>
    </Container>
  );
};

const SessionInfo = ({ session }: { session: Session }) => (
  <Paper sx={{ p: 2 }} elevation={4}>
    <Stack spacing={2}>
      <Typography variant="h6">{session.name}</Typography>
      <Paper sx={{ p: 1 }} variant="outlined">
        <Typography variant="body2" color={session.description ? 'initial' : 'text.secondary'}>
          {session.description ? session.description : 'No description'}
        </Typography>
      </Paper>
      <Grid container spacing={2}>
        <Grid item xs={12} sm={6}>
          <Box display="flex" alignItems="center">
            <CalendarTodayOutlined sx={{ mr: 1 }} color="primary" />
            <Typography variant="body2">시작 시각: {toYYYY년MM월DD일HH시MM분(session.startsAt)}</Typography>
          </Box>
        </Grid>
        <Grid item xs={12} sm={6}>
          <Box display="flex" alignItems="center">
            <StarBorderRounded sx={{ mr: 1 }} color="primary" />
            <Typography variant="body2">출석 점수: {session.score}점</Typography>
          </Box>
        </Grid>
      </Grid>
      <Box display="flex" gap={1} alignItems="center" justifyContent="flex-end">
        <TimeIcon color="action" style={{ width: 16, height: 16 }} />
        <Typography variant="body2" color="text.secondary">
          Created at: {toYYslashMMslashDDspaceHHcolonMM(session.createdAt)}
        </Typography>
      </Box>
    </Stack>
  </Paper>
);

const AttendanceQrPanel = ({
  session,
  qrRef,
  qrSizePx,
  isCreatingForm,
  isClosingSession,
  onCreateQRCode,
  onCloseSession,
}: {
  session: Session;
  qrRef: React.RefObject<HTMLDivElement>;
  qrSizePx: number;
  isCreatingForm: boolean;
  isClosingSession: boolean;
  onCreateQRCode: () => void;
  onCloseSession: () => void;
}) => (
  <Paper sx={{ p: 2, mb: 3 }} elevation={4}>
    {session.googleFormUri ? (
      <>
        <Typography variant="h6" gutterBottom>
          출석 QR
        </Typography>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2, my: 2 }}>
          <div ref={qrRef}>
            <QRCodeCanvas value={session.googleFormUri} size={qrSizePx} />
          </div>
          <Grid container spacing={2} justifyContent="center">
            <Grid item xs={12} sm={4}>
              <Button
                variant="outlined"
                fullWidth
                onClick={() => onQrDownload(qrRef, qrSizePx, toYYYY년MM월DD일HH시MM분(session.startsAt))}
              >
                QR 다운로드
              </Button>
            </Grid>
            <Grid item xs={12} sm={4}>
              <Button variant="outlined" fullWidth onClick={() => window.open(session.googleFormUri, '_blank')}>
                Google form 열기 (제출용)
              </Button>
            </Grid>
            {!session.isClosed && (
              <Grid item xs={12} sm={4}>
                <Button variant="outlined" fullWidth onClick={onCloseSession} disabled={isClosingSession}>
                  {isClosingSession ? <CircularProgress size={24} /> : '출석 반영'}
                </Button>
              </Grid>
            )}
          </Grid>
        </Box>
      </>
    ) : (
      <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2, my: 2 }}>
        <Typography variant="h6">No google form attached yet!</Typography>
        <Button variant="contained" onClick={onCreateQRCode} disabled={isCreatingForm}>
          {isCreatingForm ? <CircularProgress size={24} /> : 'Create QR code'}
        </Button>
      </Box>
    )}
  </Paper>
);

const onQrDownload = (qrRef: React.RefObject<HTMLDivElement>, qrSize: number, text: string) => {
  if (!qrRef.current) return;

  const canvas = qrRef.current.querySelector('canvas');
  if (!canvas) return;

  const newCanvas = document.createElement('canvas');
  const ctx = newCanvas.getContext('2d');
  if (!ctx) return;

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
  a.download = `${text.replace(/ /g, '_')}.png`;
  a.click();
};

export default SessionDetail;
