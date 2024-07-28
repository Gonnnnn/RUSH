import { Outlet } from 'react-router-dom';
import { Box, Container } from '@mui/material';
import Sidebar, { SIDEBAR_WIDTH } from './Sidebar';

export const HEADER_HEIGHT = 64;

const Layout = () => (
  <>
    {/* <AppBar
        sx={{
          width: `calc(100% - ${SIDEBAR_WIDTH}px)`,
          height: HEADER_HEIGHT,
          boxShadow: 'none',
          bgcolor: 'transparent',
        }}
      /> */}

    <Box sx={{ display: 'flex', minHeight: '100%' }}>
      <Box sx={{ width: SIDEBAR_WIDTH }}>
        <Sidebar />
      </Box>

      <Container
        sx={{
          py: `${HEADER_HEIGHT + 16}px`,
          boxSizing: 'border-box',
          width: `calc(100% - ${SIDEBAR_WIDTH}px)`,
        }}
      >
        <Outlet />
      </Container>
    </Box>
  </>
);

export default Layout;
