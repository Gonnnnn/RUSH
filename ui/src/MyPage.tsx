import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Typography, Paper, Box, Button, CircularProgress } from '@mui/material';
import { User, getUser } from './client/http';

const MyPage = () => {
  const navigate = useNavigate();
  const [user, setUser] = useState<User>();
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const init = async () => {
      try {
        setIsLoading(true);
        // TODO(#42): Fetch the user ID from the auth.
        setUser(await getUser(''));
      } catch (error) {
        console.error(error);
        navigate('/');
      } finally {
        setIsLoading(false);
      }
    };

    init();
  }, [navigate]);

  if (isLoading || !user) {
    return (
      <Container>
        <Typography variant="h4" sx={{ mb: 3 }}>
          My Page
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
        <Typography variant="h4">My Page</Typography>
        <Button variant="outlined" onClick={() => navigate('/')} sx={{ alignSelf: 'flex-start' }}>
          Back
        </Button>
      </Box>
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="h6">Details</Typography>
        <Typography>Name: {user.name}</Typography>
        <Typography>Generation: {user.generation}</Typography>
      </Paper>
    </Container>
  );
};

export default MyPage;
