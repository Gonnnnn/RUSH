import { NavLink, useNavigate } from 'react-router-dom';
import { CheckCircleOutline, GroupOutlined, PersonOutlined, RunCircleOutlined } from '@mui/icons-material';
import ExpandMoreRoundedIcon from '@mui/icons-material/ExpandMoreRounded';
import { Box, ListItemButton, Stack, Typography } from '@mui/material';
import { alpha } from '@mui/material/styles';
import { useAuth } from './AuthContext';
import Logo from './assets/logo.svg';
import GoogleSignInButton from './common/GoogleSignInButton';
import GoogleSignOutButton from './common/GoogleSignOutButton';

export const SIDEBAR_WIDTH = 280;

const Sidebar = () => {
  const navigate = useNavigate();
  const { authenticated } = useAuth();
  return (
    <Stack
      spacing={3}
      sx={{
        width: SIDEBAR_WIDTH,
        mr: '1px',
        borderRight: (theme) => `dashed 1px ${theme.palette.divider}`,
        px: 2,
        py: 4,
        boxSizing: 'border-box',
        position: 'fixed',
        height: '100%',
      }}
    >
      <Box px={0.5} sx={{ display: 'flex', justifyContent: 'center' }}>
        <img
          src={Logo}
          alt="logo"
          width={256}
          onClick={() => {
            navigate('/');
          }}
          style={{ cursor: 'pointer' }}
        />
      </Box>

      <Stack>
        <NavigationButton title="Me" path="/me" icon={<PersonOutlined />} />
        <NavigationButton title="Sessions" path="/sessions" icon={<RunCircleOutlined />} />
        <NavigationButton title="Users" path="/users" icon={<GroupOutlined />} />
        <NavigationButton title="Attendance" path="/attendance" icon={<CheckCircleOutline />} />
      </Stack>
      {authenticated ? <GoogleSignOutButton /> : <GoogleSignInButton />}
    </Stack>
  );
};

export default Sidebar;

const NavigationButton = ({
  title,
  path,
  icon,
  expandable,
}: {
  title: string;
  path: string;
  icon?: JSX.Element;
  expandable?: boolean;
}) => (
  <ListItemButton
    key={title}
    component={NavLink}
    end
    to={path}
    sx={{
      minHeight: 44,
      borderRadius: '6px',
      typography: 'body2',
      color: 'text.secondary',
      fontWeight: 'fontWeightMedium',
      '&.active': {
        color: 'primary.main',
        fontWeight: 'fontWeightSemiBold',
        bgcolor: (theme) => alpha(theme.palette.primary.main, 0.08),
        '&:hover': {
          bgcolor: (theme) => alpha(theme.palette.primary.main, 0.16),
        },
      },
      marginBottom: 0.5,
    }}
  >
    <Stack sx={{ width: 24, height: 24, mr: 1.5 }}>{icon}</Stack>

    <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{ flex: 1 }}>
      <Typography
        sx={{
          width: 132,
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          whiteSpace: 'nowrap',
        }}
      >
        {title}
      </Typography>

      {expandable && <ExpandMoreRoundedIcon />}
    </Stack>
  </ListItemButton>
);
