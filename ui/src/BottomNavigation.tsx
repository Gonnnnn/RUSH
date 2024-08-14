import { NavLink } from 'react-router-dom';
import { CheckCircleOutline, GroupOutlined, PersonOutlined, RunCircleOutlined } from '@mui/icons-material';
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
      justifyContent: 'space-around',
      alignItems: 'center',
      px: 2,
    }}
  >
    <NavigationButton title="Me" path="/me" icon={<PersonOutlined />} />
    <NavigationButton title="Sessions" path="/sessions" icon={<RunCircleOutlined />} />
    <NavigationButton title="Users" path="/users" icon={<GroupOutlined />} />
    <NavigationButton title="Attendance" path="/attendances" icon={<CheckCircleOutline />} />
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
    }}
  >
    <Stack sx={{ width: 24, height: 24, mb: 0.5 }}>{icon}</Stack>
    <Typography sx={{ fontSize: '0.75rem' }}>{title}</Typography>
  </ListItemButton>
);

export default BottomNavigation;
