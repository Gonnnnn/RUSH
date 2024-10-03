import { useEffect, useRef, useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { CalendarTodayOutlined, StarBorderRounded } from '@mui/icons-material';
import { Container, Typography, Paper, Box, Button, CircularProgress, Stack, Grid } from '@mui/material';
import { TimeIcon } from '@mui/x-date-pickers';
import { AxiosError } from 'axios';
import { QRCodeCanvas } from 'qrcode.react';
import { useHeader } from './Layout';
import { useSnackbar } from './SnackbarContext';
import { Session, createSessionForm, deleteSession, getSession } from './client/http';
import { formatDateToMonthDate, toYYslashMMslashDDspaceHHcolonMM, toYYYY년MM월DD일HH시MM분 } from './common/date';

const SessionDetail = () => {
  useHeader({ newTitle: 'Session Detail' });
  const navigate = useNavigate();
  const { state } = useLocation();
  const { showWarning, showError } = useSnackbar();
  const { id } = useParams();

  const [session, setSession] = useState<Session | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isCreatingForm, setIsCreatingForm] = useState(false);
  const qrRef = useRef<HTMLDivElement>(null);
  const qrSizePx = 128;
  const qrDownloadSizePx = 512;

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

  const handleDeleteBtnClick = async () => {
    try {
      await deleteSession(id);
      navigate('/sessions');
    } catch (error) {
      handleError({
        error,
        messageAuth: 'Session deletion is restricted to admin users',
        messageInternal: 'Failed to delete the session. Contact the dev.',
      });
    }
  };

  const handleQrCodeCreateClick = async () => {
    try {
      setIsCreatingForm(true);
      await createSessionForm(id);
      const updatedSession = await getSession(id);
      setSession(updatedSession);
    } catch (error) {
      handleError({
        error,
        messageAuth: 'Form creation is restricted to authenticated users',
        messageInternal: 'Failed to create a form. Contact the dev.',
      });
    } finally {
      setIsCreatingForm(false);
    }
  };

  const handleError = ({
    error,
    messageAuth,
    messageInternal,
  }: {
    error: unknown;
    messageAuth: string;
    messageInternal: string;
  }) => {
    if (!(error instanceof AxiosError)) {
      showError(messageInternal);
      return;
    }

    const status = error.response?.status;
    switch (status) {
      case 401:
        showWarning(messageAuth);
        break;
      case 403:
        showWarning(messageAuth);
        break;
      default:
        showError(messageInternal);
        break;
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
      <Box sx={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', mb: 3 }}>
        <Button
          variant="outlined"
          onClick={() => {
            if (state?.from) {
              navigate(state.from);
              return;
            }
            navigate('/sessions');
          }}
          sx={{ alignSelf: 'flex-start' }}
        >
          Back
        </Button>
        <Button
          variant="outlined"
          color="error"
          onClick={handleDeleteBtnClick}
          disabled={session.attendanceStatus === 'applied'}
        >
          {session.attendanceStatus === 'applied' ? 'Applied already' : 'Delete'}
        </Button>
      </Box>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
        <SessionInfo session={session} />
        <AttendanceQrPanel
          session={session}
          qrRef={qrRef}
          qrSizePx={qrSizePx}
          qrDownloadSizePx={qrDownloadSizePx}
          isCreatingForm={isCreatingForm}
          onCreateQRCode={handleQrCodeCreateClick}
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
            <Typography variant="body2">시작 시간: {toYYYY년MM월DD일HH시MM분(session.startsAt)}</Typography>
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
  qrDownloadSizePx,
  isCreatingForm,
  onCreateQRCode,
}: {
  session: Session;
  qrRef: React.RefObject<HTMLDivElement>;
  qrSizePx: number;
  qrDownloadSizePx: number;
  isCreatingForm: boolean;
  onCreateQRCode: () => void;
}) => (
  <Paper sx={{ p: 2, mb: 3 }} elevation={4}>
    {session.googleFormUri ? (
      <>
        <Typography variant="h6" gutterBottom>
          출석 QR
        </Typography>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2, my: 2 }}>
          <QRCodeCanvas value={session.googleFormUri} size={qrSizePx} />
          {/* Make a hidden QR for download. The QR for display is too small that it breaks when resizing for downloading. */}
          <div ref={qrRef} style={{ display: 'None' }}>
            <QRCodeCanvas value={session.googleFormUri} size={qrDownloadSizePx} />
          </div>
          <Grid container spacing={2} justifyContent="center">
            <Grid item xs={12} sm={6}>
              <Button
                variant="outlined"
                fullWidth
                // to Month Day format
                onClick={() => onQrDownload(qrRef, qrDownloadSizePx, formatDateToMonthDate(new Date(session.startsAt)))}
              >
                QR 다운로드
              </Button>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Button variant="outlined" fullWidth onClick={() => window.open(session.googleFormUri, '_blank')}>
                Google form 열기 (제출용)
              </Button>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Button
                variant="outlined"
                fullWidth
                onClick={() => window.open(`https://docs.google.com/forms/d/${session.googleFormId}/edit`, '_blank')}
              >
                Google form 열기 (편집용)
              </Button>
            </Grid>
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

  const paddingPx = 64;
  const textSpacePx = 128;
  const newCanvasWidth = qrSize + paddingPx * 2;
  const newCanvasHeight = qrSize + paddingPx * 2 + textSpacePx;

  newCanvas.width = newCanvasWidth;
  newCanvas.height = newCanvasHeight;
  ctx.fillStyle = 'white';
  ctx.fillRect(0, 0, newCanvas.width, newCanvas.height);

  const qrYoffset = (newCanvasHeight - qrSize) / 2;
  const qrXoffset = (newCanvasWidth - qrSize) / 2;
  ctx.drawImage(canvas, qrXoffset, qrYoffset, qrSize, qrSize);

  ctx.font = '32px Helvetica';
  ctx.fillStyle = 'black';
  ctx.textAlign = 'center';
  ctx.fillText(text, newCanvasWidth / 2, Math.min(qrYoffset + qrSize + textSpacePx / 2, newCanvasHeight - 16));

  const a = document.createElement('a');
  a.href = newCanvas.toDataURL('image/png');
  a.download = `${text.replace(/ /g, '_')}.png`;
  a.click();
};

export default SessionDetail;
