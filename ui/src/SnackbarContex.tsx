import { ReactNode, createContext, useState, useContext, useCallback, useMemo } from 'react';
import Alert from '@mui/material/Alert';
import Snackbar from '@mui/material/Snackbar';

export enum SnackbarMessageType {
  error = 'error',
  warning = 'warning',
  info = 'info',
  success = 'success',
}

interface SnackbarContextType {
  showMessage: (message: string, type: SnackbarMessageType) => void;
}

const SnackbarContext = createContext<SnackbarContextType | undefined>(undefined);

export const useSnackbar = (): SnackbarContextType => {
  const context = useContext(SnackbarContext);
  if (!context) {
    throw new Error('useSnackbar must be used within a SnackbarProvider');
  }
  return context;
};

export const SnackbarProvider = ({ children }: { children: ReactNode }) => {
  const [open, setOpen] = useState(false);
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState<SnackbarMessageType>(SnackbarMessageType.info);
  const [lastTime, setLastTime] = useState(0);
  const throttleTimeMilli = 1000;

  const showMessage = useCallback(
    (msg: string, type: SnackbarMessageType) => {
      const now = Date.now();
      if (now - lastTime >= throttleTimeMilli) {
        setMessage(msg);
        setMessageType(type);
        setOpen(true);
        setLastTime(now);
      }
    },
    [lastTime],
  );

  const handleClose = () => {
    setOpen(false);
  };

  const value = useMemo(() => ({ showMessage }), [showMessage]);

  return (
    <SnackbarContext.Provider value={value}>
      {children}
      <Snackbar open={open} autoHideDuration={6000} onClose={handleClose}>
        <Alert onClose={handleClose} severity={messageType} sx={{ width: '100%' }}>
          {message}
        </Alert>
      </Snackbar>
    </SnackbarContext.Provider>
  );
};
