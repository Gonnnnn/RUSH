import { NavLink } from 'react-router-dom';
import { EmergencyRecordingOutlined, PersonOutlined, RunCircleOutlined } from '@mui/icons-material';
import { Box, ListItemButton, Stack, Typography } from '@mui/material';

export const BOTTOM_NAV_HEIGHT = 56;

const BottomNavigation = () => (
  <Box
    sx={{
      position: 'fixed',
      bottom: 0,
      left: 0,
      right: 0,
      height: BOTTOM_NAV_HEIGHT,
      bgcolor: 'background.paper',
      borderTop: (theme) => `solid 1px ${theme.palette.divider}`,
      display: 'flex',
      alignItems: 'center',
      px: 2,
    }}
  >
    <NavigationButton title="Me" path="/admin/me" icon={<PersonOutlined />} />
    <NavigationButton title="Sessions" path="/admin/sessions" icon={<RunCircleOutlined />} />
    <NavigationButton title="Exception" path="/admin/exceptions" icon={<EmergencyRecordingOutlined />} />
  </Box>
);

const NavigationButton = ({ title, path, icon }: { title: string; path: string; icon: JSX.Element }) => (
  <ListItemButton
    component={NavLink}
    to={path}
    sx={{
      flexDirection: 'column',
      alignItems: 'center',
      '&.active': {
        color: 'primary.main',
      },
      flex: 1,
    }}
  >
    <Stack sx={{ width: 24, height: 24, mb: 0.5 }}>{icon}</Stack>
    <Typography sx={{ fontSize: '0.75rem' }}>{title}</Typography>
  </ListItemButton>
);

export default BottomNavigation;
