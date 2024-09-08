import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Container, TextField, Typography } from '@mui/material';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import { AxiosError } from 'axios';
import dayjs from 'dayjs';
import { useSnackbar } from './SnackbarContext';
import { createSession } from './client/http';

const SessionCreate = () => {
  const navigate = useNavigate();
  const { showWarning, showError } = useSnackbar();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [startsAt, setStartsAt] = useState(dayjs().add(1, 'hour').startOf('minute'));
  const [score, setScore] = useState(0);
  const [isFormValid, setIsFormValid] = useState(false);

  useEffect(() => {
    const isNameValid = name.trim().length > 0;
    const isStartsAtValid = startsAt.isAfter(dayjs());
    setIsFormValid(isNameValid && isStartsAtValid);
  }, [name, startsAt]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!startsAt || !isFormValid) return;
    try {
      const id = await createSession(name, description, new Date(startsAt.toISOString()), score);
      navigate(`/sessions/${id}`);
    } catch (error: unknown) {
      if (!(error instanceof AxiosError)) {
        showError('An unexpected error occurred. Please contact the administrator.');
        return;
      }

      const status = error.response?.status;
      switch (status) {
        case 401:
          showWarning('Session creation is restricted to admin users');
          break;
        case 403:
          showWarning('Session creation is restricted to admin users');
          break;
        default:
          showError('An unexpected error occurred. Please contact the administrator.');
          break;
      }
    }
  };

  const handleDateTimeChange = (newValue: dayjs.Dayjs | null) => {
    if (newValue === null) return;
    // Truncate seconds and milliseconds because sessions are created in minute precision.
    setStartsAt(newValue.startOf('minute'));
  };

  return (
    <Container>
      <Typography variant="h4" sx={{ mb: 3 }}>
        Create Session
      </Typography>
      <form onSubmit={handleSubmit}>
        <TextField
          label="Name"
          name="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          fullWidth
          sx={{ mb: 2 }}
          error={name.trim().length === 0}
          helperText={name.trim().length === 0 ? 'Name is required' : ''}
        />
        <TextField
          label="Description"
          name="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          fullWidth
          sx={{ mb: 2 }}
        />
        <DateTimePicker
          label="Starts At"
          value={startsAt}
          onChange={handleDateTimeChange}
          minutesStep={10}
          shouldDisableDate={(date) => date.isBefore(dayjs(), 'day')}
          views={['year', 'month', 'day', 'hours', 'minutes']}
          format="YYYY/MM/DD HH:mm"
          skipDisabled
          sx={{ mb: 2, width: '100%' }}
          slotProps={{
            textField: {
              helperText: startsAt.isBefore(dayjs()) ? 'Start time must be in the future' : '',
              error: startsAt.isBefore(dayjs()),
            },
          }}
        />
        <TextField
          type="number"
          label="Score"
          name="score"
          value={score}
          onChange={(e) => setScore(parseInt(e.target.value, 10))}
          fullWidth
          sx={{ mb: 2 }}
        />
        <Box sx={{ textAlign: 'center' }}>
          <Button type="submit" variant="contained" sx={{ mt: 2 }} onClick={handleSubmit} disabled={!isFormValid}>
            Create
          </Button>
        </Box>
      </form>
    </Container>
  );
};

export default SessionCreate;
