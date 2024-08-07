import { useLocation, useNavigate } from 'react-router-dom';
import { Stack, Typography } from '@mui/material';
import { useHeader } from './Layout';
import GoogleSignInButton from './common/GoogleSignInButton';

const SignIn = () => {
  useHeader({ newTitle: 'Sign-In' });
  const location = useLocation();
  const navigate = useNavigate();
  const { from } = location.state || { from: { pathname: '/' } };

  return (
    <Stack justifyContent="center" alignItems="center" spacing={2} sx={{ height: '100vh' }}>
      <Typography variant="h4">Sign-In required for this page</Typography>
      <GoogleSignInButton
        callBack={() => {
          navigate(from, { replace: true });
        }}
        text="Sign in with Google"
      />
    </Stack>
  );
};

export default SignIn;
