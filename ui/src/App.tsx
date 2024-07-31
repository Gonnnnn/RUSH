import { BrowserRouter } from 'react-router-dom';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import dayjs from 'dayjs';
import 'dayjs/locale/ko';
import AppRoutes from './Routes';

dayjs.locale('ko');

const App = () => (
  <LocalizationProvider dateAdapter={AdapterDayjs}>
    <BrowserRouter>
      <AppRoutes />
    </BrowserRouter>
  </LocalizationProvider>
);
export default App;
