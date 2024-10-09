import { useState } from 'react';
import { CalendarTodayOutlined, StarBorderRounded } from '@mui/icons-material';
import { TextField, Typography, Paper, Stack, Box, Grid } from '@mui/material';
import { TimeIcon } from '@mui/x-date-pickers';
import { Session } from '../client/http';
import { toYYslashMMslashDDspaceHHcolonMM, toYYYY년MM월DD일HH시MM분 } from '../common/date';

type EditableSessionData = {
  name: string;
  description: string;
  score: number;
  startsAt: Date;
};

type InputNames = keyof EditableSessionData;

const SessionInfo = ({ session }: { session: Session }) => {
  const [isEditing] = useState(false);
  const [sessionData, setSessionData] = useState<EditableSessionData>({
    name: session.name,
    description: session.description,
    score: session.score,
    startsAt: session.startsAt,
  });
  const [errors, setErrors] = useState({
    name: false,
    score: false,
  });

  const handleInputChange = (inputName: InputNames, newValue: string) => {
    const prevData = sessionData;
    const newData = {
      ...prevData,
      [inputName]: inputName === 'score' ? Number(newValue) : newValue,
    };
    setSessionData(newData);
    validateFields(newData);
  };

  const validateFields = (newData: EditableSessionData) => {
    let isValid = true;
    const newErrors = { name: false, score: false };

    if (!newData.name) {
      newErrors.name = true;
      isValid = false;
    }

    if (newData.score <= 0) {
      newErrors.score = true;
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  return (
    <Paper sx={{ p: 2 }} elevation={4}>
      <Stack spacing={2}>
        {isEditing ? (
          <>
            <TextField
              label="Session Name"
              value={sessionData.name}
              onChange={(e) => handleInputChange('name', e.target.value)}
              variant="outlined"
              required
              fullWidth
              error={errors.name}
              helperText={errors.name ? 'Session name is required' : ''}
            />
            <TextField
              label="Description"
              value={sessionData.description || ''}
              onChange={(e) => handleInputChange('description', e.target.value)}
              multiline
              rows={3}
              variant="outlined"
              fullWidth
            />
            <TextField
              label="Starts At"
              type="datetime-local"
              value={sessionData.startsAt}
              onChange={(e) => handleInputChange('startsAt', e.target.value)}
              variant="outlined"
              required
              fullWidth
            />
            <TextField
              label="Score"
              type="number"
              value={sessionData.score}
              onChange={(e) => handleInputChange('score', e.target.value)}
              variant="outlined"
              inputProps={{ min: 0 }}
              required
              fullWidth
              error={errors.score}
              helperText={errors.score ? 'Score must be greater than 0' : ''}
            />
          </>
        ) : (
          <>
            <Typography variant="h6">{session.name}</Typography>
            <Paper sx={{ p: 1 }} variant="outlined">
              <Typography variant="body2" color={session.description ? 'initial' : 'text.secondary'}>
                {session.description ? session.description : 'No description'}
              </Typography>
            </Paper>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Box display="flex" alignItems="center">
                  <CalendarTodayOutlined sx={{ mr: 1 }} color="primary" />
                  <Typography variant="body2">시작 시간: {toYYYY년MM월DD일HH시MM분(session.startsAt)}</Typography>
                </Box>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Box display="flex" alignItems="center">
                  <StarBorderRounded sx={{ mr: 1 }} color="primary" />
                  <Typography variant="body2">출석 점수: {session.score}점</Typography>
                </Box>
              </Grid>
            </Grid>
          </>
        )}
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Box display="flex" gap={1} alignItems="center" justifyContent="center">
            <TimeIcon color="action" style={{ width: 16, height: 16 }} />
            <Typography variant="body2" color="text.secondary">
              Created at {toYYslashMMslashDDspaceHHcolonMM(session.createdAt)}
            </Typography>
          </Box>
          {/* TODO(#188): Add the button to toggle. */}
        </Box>
      </Stack>
    </Paper>
  );
};

export default SessionInfo;
