import { Outlet } from 'react-router-dom';
import { Box, Container, Typography, useMediaQuery, useTheme } from '@mui/material';
import { useAuth } from '../auth';
import GoogleSignInButton from '../common/GoogleSignInButton';
import GoogleSignOutButton from '../common/GoogleSignOutButton';
import { useHeader } from '../contexts/header';
import BottomNavigation, { BOTTOM_NAV_HEIGHT } from './BottomNavigation';
import Sidebar, { SIDEBAR_WIDTH } from './Sidebar';

const Layout = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const { pageTitle } = useHeader();
  const { authenticated } = useAuth();

  return isMobile ? (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', p: 2 }}>
        <Typography variant="h6">{pageTitle}</Typography>
        {authenticated ? <GoogleSignOutButton /> : <GoogleSignInButton />}
      </Box>

      <Container
        sx={{
          py: '16px',
          pb: `${BOTTOM_NAV_HEIGHT + 16}px`,
          boxSizing: 'border-box',
          width: '100%',
          flexGrow: 1,
        }}
      >
        <Outlet />
      </Container>

      <BottomNavigation />
    </Box>
  ) : (
    <Box sx={{ display: 'flex', flexDirection: 'row', minHeight: '100vh' }}>
      <Box sx={{ width: SIDEBAR_WIDTH }}>
        <Sidebar />
      </Box>

      <Container
        sx={{
          py: '80px',
          pb: '16px',
          boxSizing: 'border-box',
          width: `calc(100% - ${SIDEBAR_WIDTH}px)`,
          flexGrow: 1,
        }}
      >
        <Outlet />
      </Container>
    </Box>
  );
};

export default Layout;
