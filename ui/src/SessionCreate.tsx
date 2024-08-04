import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Container, TextField, Typography } from '@mui/material';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import { AxiosError } from 'axios';
import dayjs from 'dayjs';
import { useSnackbar } from './SnackbarContex';
import { createSession } from './client/http';

const SessionCreate = () => {
  const navigate = useNavigate();
  const { showWarning, showError } = useSnackbar();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [startsAt, setStartsAt] = useState(dayjs());
  const [score, setScore] = useState(0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!startsAt) return;
    try {
      const id = await createSession(name, description, new Date(startsAt.toISOString()), score);
      navigate(`/sessions/${id}`);
    } catch (error: unknown) {
      if (error instanceof AxiosError && error.response?.status === 401) {
        showWarning('Session creation is restricted to authenticated users');
      } else {
        showError('Failed to create a form. Contact the administrator.');
      }
    }
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
          onChange={(newValue) => {
            console.log(newValue);
            setStartsAt(newValue ?? dayjs());
          }}
          minutesStep={10}
          shouldDisableDate={(date) => date.isBefore(dayjs(), 'day')}
          views={['year', 'month', 'day', 'hours', 'minutes']}
          format="YYYY/MM/DD HH:mm"
          skipDisabled
          sx={{ mb: 2, width: '100%' }}
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
          <Button type="submit" variant="contained" sx={{ mt: 2 }} onClick={handleSubmit}>
            Create
          </Button>
        </Box>
      </form>
    </Container>
  );
};

export default SessionCreate;
