import { useCallback, useEffect, useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { Container, Box, Button, CircularProgress } from '@mui/material';
import { useHeader } from '../../Layout';
import { Session } from '../../client/http/data';
import { getSession } from '../../client/http/default';
import useHandleError from '../../common/error';
import SessionInfo from '../common/SessionInfo';
import SessionAttendanceTable from './AttendanceTable';

const UserSessionDetail = () => {
  useHeader({ newTitle: 'Session Detail' });
  const navigate = useNavigate();
  const { state } = useLocation();
  const { handleError } = useHandleError();
  const { id } = useParams();

  const [session, setSession] = useState<Session | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const fetchSession = useCallback(
    async (sessionId: string) => {
      try {
        setIsLoading(true);
        const fetchedSession = await getSession(sessionId);
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
      </Box>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
        <SessionInfo session={session} />
        <SessionAttendanceTable sessionId={id} />
      </Box>
    </Container>
  );
};

export default UserSessionDetail;
