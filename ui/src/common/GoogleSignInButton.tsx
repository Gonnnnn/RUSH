import { Button, Typography } from '@mui/material';
import { AxiosError } from 'axios';
import { signInWithPopup } from 'firebase/auth';
import GoogleLogo from '../assets/google_logo.svg';
import { useAuth } from '../auth';
import { useSnackbar } from '../contexts/snackbar';
import { auth, provider } from '../firebase';

const GoogleSignInButton = ({ text = '', callBack }: { text?: string; callBack?: () => void }) => {
  const { login } = useAuth();
  const { showInfo, showError } = useSnackbar();

  const handleGoogleSignIn = async () => {
    try {
      const credential = await signInWithPopup(auth, provider);
      const idToken = await credential.user.getIdToken();
      await login(idToken);
      callBack?.();
      showInfo('Successfully signed in with Google.');
    } catch (error) {
      if (error instanceof AxiosError) {
        if (error.response?.status === 404) {
          showError('You are not registered in the system. Please contact the administrator.');
          return;
        }
      }

      // TODO(#23): Handle error after centralizing the error handler.
      // eslint-disable-next-line no-console
      console.error('Error signing in with Google:', error);
      showError('Failed to sign in with Google. Please contact dev if thie error persists.');
    }
  };

  return (
    <Button
      size="large"
      color="inherit"
      variant="outlined"
      sx={{
        borderColor: 'rgba(0, 0, 0, 0.23)',
        bgcolor: 'white',
        boxShadow: '0px 2px 4px rgba(0, 0, 0, 0.1)',
        display: 'flex',
        columnGap: 2,
        padding: '8px 8px',
      }}
      onClick={() => handleGoogleSignIn()}
    >
      <img src={GoogleLogo} width={20} height={20} alt="Google logo" />
      <Typography variant="body1">{text}</Typography>
    </Button>
  );
};

export default GoogleSignInButton;
