import { Button } from '@mui/material';
import { signInWithPopup } from 'firebase/auth';
import GoogleLogo from '../assets/google_logo.svg';
import { useAuth } from '../auth/context';
import { auth, provider } from '../firebase';

const GoogleSignInButton = ({ text = '' }: { text?: string }) => {
  const { login } = useAuth();
  const handleGoogleSignIn = async () => {
    try {
      const credential = await signInWithPopup(auth, provider);
      const idToken = await credential.user.getIdToken();
      await login(idToken);
    } catch (error) {
      // TODO(#23): Handle error after centralizing the error handler.
      // eslint-disable-next-line no-console
      console.error('Error signing in with Google:', error);
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
        padding: '8px 0',
      }}
      onClick={() => handleGoogleSignIn()}
    >
      <img src={GoogleLogo} width={20} height={20} alt="Google logo" />
      {text}
    </Button>
  );
};

export default GoogleSignInButton;
