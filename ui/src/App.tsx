import { useEffect, useState } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Box } from '@mui/material';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import dayjs from 'dayjs';
import 'dayjs/locale/ko';
import { AuthProvider } from './AuthContext';
import AppRoutes from './Routes';
import { SnackbarProvider } from './SnackbarContex';
import Logo from './assets/logo.svg';

dayjs.locale('ko');

const App = () => {
  const [isLoading, setIsLoading] = useState(true);
  const [isFadingOut, setIsFadingOut] = useState(false);
  const fadeOutTimeMillis = 1000;

  useEffect(() => {
    const loadingTimer = setTimeout(() => {
      setIsLoading(false);
    }, fadeOutTimeMillis);

    setIsFadingOut(true);

    return () => clearTimeout(loadingTimer);
  }, []);

  return (
    <SnackbarProvider>
      <AuthProvider>
        <LocalizationProvider dateAdapter={AdapterDayjs}>
          <BrowserRouter>
            <AppRoutes />
          </BrowserRouter>
        </LocalizationProvider>
      </AuthProvider>
      {isLoading && (
        <Box
          onClick={() => setIsLoading(false)}
          position="fixed"
          top={0}
          left={0}
          width="100%"
          height="100%"
          display="flex"
          justifyContent="center"
          alignItems="center"
          bgcolor="#ffffff"
          zIndex={9999}
          sx={{
            opacity: isFadingOut ? 0 : 1,
            transition: `opacity ${fadeOutTimeMillis}ms cubic-bezier(0.6, 0, 1, 1)`,
          }}
        >
          <img src={Logo} alt="logo" width={256} />
        </Box>
      )}
    </SnackbarProvider>
  );
};
export default App;
