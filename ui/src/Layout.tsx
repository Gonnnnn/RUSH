// Layout.tsx
import { Outlet } from 'react-router-dom';
import { Box, Container, useMediaQuery, useTheme } from '@mui/material';
import BottomNavigation, { BOTTOM_NAV_HEIGHT } from './BottomNavigation';
import Sidebar, { SIDEBAR_WIDTH } from './Sidebar';

export const HEADER_HEIGHT = 64;

const Layout = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  return (
    <Box sx={{ display: 'flex', flexDirection: isMobile ? 'column' : 'row', minHeight: '100vh' }}>
      {!isMobile && (
        <Box sx={{ width: SIDEBAR_WIDTH }}>
          <Sidebar />
        </Box>
      )}

      <Container
        sx={{
          py: isMobile ? `16px` : `${HEADER_HEIGHT + 16}px`,
          pb: isMobile ? `${BOTTOM_NAV_HEIGHT + 16}px` : '16px',
          boxSizing: 'border-box',
          width: isMobile ? '100%' : `calc(100% - ${SIDEBAR_WIDTH}px)`,
          flexGrow: 1,
        }}
      >
        <Outlet />
      </Container>

      {isMobile && <BottomNavigation />}
    </Box>
  );
};

export default Layout;
