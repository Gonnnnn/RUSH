import { useState } from 'react';
import { Button, CircularProgress } from '@mui/material';
import { useAuth } from '../AuthContext';
import { useSnackbar } from '../SnackbarContex';
import GoogleLogo from '../assets/google_logo.svg';

const GoogleSignOutButton = ({ text = 'Sign Out' }: { text?: string }) => {
  const { logout } = useAuth();
  const { showInfo } = useSnackbar();
  const [isLoggingOut, setIsLoggingOut] = useState(false);

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
        padding: '8px',
        fontSize: '14px',
        lineHeight: '20px',
        fontWeight: 400,
      }}
      onClick={() => {
        setIsLoggingOut(true);
        logout();
        setIsLoggingOut(false);
        showInfo('Successfully signed out.');
      }}
    >
      {isLoggingOut ? (
        <CircularProgress
          size={20}
          sx={{
            color: 'black',
          }}
        />
      ) : (
        <img src={GoogleLogo} alt="Google Logo" height={20} width={20} />
      )}
      {text}
    </Button>
  );
};

export default GoogleSignOutButton;
