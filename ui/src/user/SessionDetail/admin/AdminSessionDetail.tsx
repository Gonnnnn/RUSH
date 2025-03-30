import { useCallback, useEffect, useRef, useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { Container, Box, Button, CircularProgress, Paper, Typography } from '@mui/material';
import { adminCreateSessionForm, adminDeleteSession, adminGetSession } from '../../../client/http/admin';
import { AdminSession, SessionAttendanceAppliedBy } from '../../../client/http/data';
import useHandleError from '../../../common/error';
import { useHeader } from '../../../contexts/header';
import SessionInfo from '../common/SessionInfo';
import AttendanceQrPanel from './AttendanceQrPanel';
import SessionAttendanceTable from './AttendanceTable';

const AdminSessionDetail = () => {
  useHeader({ newTitle: 'Session Detail' });
  const navigate = useNavigate();
  const { state } = useLocation();
  const { handleError } = useHandleError();
  const { id } = useParams();

  const [session, setSession] = useState<AdminSession | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const [isCreatingForm, setIsCreatingForm] = useState(false);
  const qrRef = useRef<HTMLDivElement>(null);
  const qrSizePx = 128;
  const qrDownloadSizePx = 512;

  const fetchSession = useCallback(
    async (sessionId: string) => {
      try {
        setIsLoading(true);
        const fetchedSession = await adminGetSession(sessionId);
        setSession(fetchedSession);
      } catch (error) {
        handleError({
          error,
          messageAuth: 'Session retrieval is restricted to authenticated users',
          messageInternal: 'Failed to retrieve the session. Contact the dev.',
        });
        navigate('/sessions');
      } finally {
        setIsLoading(false);
      }
    },
    [navigate, setSession, setIsLoading, handleError],
  );

  useEffect(() => {
    if (!id) {
      navigate('/sessions');
      return;
    }

    fetchSession(id);
  }, [navigate, id, fetchSession]);

  if (!id) {
    navigate('/sessions');
    return null;
  }

  const handleDeleteBtnClick = async () => {
    try {
      await adminDeleteSession(id);
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
      await adminCreateSessionForm(id);
      const updatedSession = await adminGetSession(id);
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
        <SessionAttendanceTable
          sessionId={id}
          reloadSession={() => fetchSession(id)}
          qrActivated={session.googleFormId !== '' || session.googleFormUri !== ''}
        />
        {session.attendanceAppliedBy === SessionAttendanceAppliedBy.Enum.unspecified ||
        session.attendanceAppliedBy === SessionAttendanceAppliedBy.Enum.form ? (
          <AttendanceQrPanel
            session={session}
            qrRef={qrRef}
            qrSizePx={qrSizePx}
            qrDownloadSizePx={qrDownloadSizePx}
            isCreatingForm={isCreatingForm}
            onCreateQRCode={handleQrCodeCreateClick}
          />
        ) : (
          <Paper sx={{ p: 2 }} elevation={4}>
            <Typography variant="h6" gutterBottom>
              출석 QR
            </Typography>
            <Typography variant="body1" gutterBottom>
              해당 세션은 Google form을 통해 출석처리 되지 않았습니다.
            </Typography>
          </Paper>
        )}
      </Box>
    </Container>
  );
};

export default AdminSessionDetail;
