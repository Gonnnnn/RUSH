import { useEffect, useState } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Box } from '@mui/material';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import dayjs from 'dayjs';
import 'dayjs/locale/ko';
import AppRoutes from './Routes';
import Logo from './assets/logo.svg';
import { AuthProvider } from './auth/context';

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
    <>
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
    </>
  );
};
export default App;
