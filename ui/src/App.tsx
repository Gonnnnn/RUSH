import { ReactNode, useEffect, useState } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Box } from '@mui/material';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import dayjs from 'dayjs';
import 'dayjs/locale/ko';
import { AuthProvider, useAuth } from './AuthContext';
import AppRoutes from './Routes';
import { SnackbarProvider } from './SnackbarContex';
import Logo from './assets/logo.svg';

dayjs.locale('ko');

const App = () => (
  <SnackbarProvider>
    <AuthProvider>
      <LocalizationProvider dateAdapter={AdapterDayjs}>
        <DataLoader>
          <BrowserRouter>
            <AppRoutes />
          </BrowserRouter>
        </DataLoader>
      </LocalizationProvider>
    </AuthProvider>
  </SnackbarProvider>
);

const DataLoader = ({ children }: { children: ReactNode }) => {
  const { isLoading: isAuthLoading } = useAuth();
  // To keep logo until necessary data is loaded.
  const [shouldKeepLogo, setShouldKeepLogo] = useState(true);
  const [isFadingOut, setIsFadingOut] = useState(false);
  const fadeOutTimeMillis = 1000;

  useEffect(() => {
    if (isAuthLoading) {
      return;
    }
    setIsFadingOut(true);
  }, [isAuthLoading]);

  useEffect(() => {
    if (!isFadingOut) {
      return;
    }
    setTimeout(() => {
      setShouldKeepLogo(false);
    }, fadeOutTimeMillis);
  }, [isFadingOut]);

  const handleLogoClick = () => {
    if (isAuthLoading) {
      return;
    }
    setShouldKeepLogo(false);
    setIsFadingOut(true);
  };

  return (
    <>
      {isAuthLoading ? null : children}
      {shouldKeepLogo && (
        <Box
          onClick={() => handleLogoClick()}
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
