import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Container, Input, TextField, Typography } from '@mui/material';
import { createSession } from './client/http';

const SessionCreate = () => {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const { dateInputValue, onDateInputValueChange, date: startsAt } = useDateTimeLocalInput();
  const [score, setScore] = useState(0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const id = await createSession(name, description, startsAt, score);
      navigate(`/sessions/${id}`);
    } catch (error) {
      console.error(error);
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
        <Input
          type="datetime-local"
          name="startsAt"
          value={dateInputValue}
          onChange={(e) => onDateInputValueChange(e.target.value)}
          fullWidth
          sx={{ mb: 2 }}
        />
        <Input
          type="number"
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

const useDateTimeLocalInput = (initialDate = new Date()) => {
  const [date, setDate] = useState(initialDate);

  const formatDateToInput = (newDate: Date) => {
    const pad = (number: number) => number.toString().padStart(2, '0');

    const year = newDate.getFullYear();
    const month = pad(newDate.getMonth() + 1);
    const day = pad(newDate.getDate());
    const hours = pad(newDate.getHours());
    const minutes = pad(newDate.getMinutes());

    return `${year}-${month}-${day}T${hours}:${minutes}`;
  };

  const handleChange = (newInputValue: string) => {
    setDate(new Date(newInputValue));
  };

  return {
    dateInputValue: formatDateToInput(date),
    onDateInputValueChange: handleChange,
    date,
  };
};

export default SessionCreate;
