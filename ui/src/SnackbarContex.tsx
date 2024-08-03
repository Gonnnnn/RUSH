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
  showInfo: (message: string) => void;
  showError: (message: string) => void;
  showSuccess: (message: string) => void;
  showWarning: (message: string) => void;
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

  const showInfo = useCallback((msg: string) => showMessage(msg, SnackbarMessageType.info), [showMessage]);
  const showError = useCallback((msg: string) => showMessage(msg, SnackbarMessageType.error), [showMessage]);
  const showSuccess = useCallback((msg: string) => showMessage(msg, SnackbarMessageType.success), [showMessage]);
  const showWarning = useCallback((msg: string) => showMessage(msg, SnackbarMessageType.warning), [showMessage]);

  const handleClose = () => {
    setOpen(false);
  };

  const value = useMemo(
    () => ({ showMessage, showInfo, showError, showSuccess, showWarning }),
    [showMessage, showInfo, showError, showSuccess, showWarning],
  );

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
