import { AxiosError } from 'axios';
import { useSnackbar } from '../SnackbarContext';

const useHandleError = () => {
  const { showWarning, showError } = useSnackbar();

  const handleError = ({
    error,
    messageAuth,
    messageInternal,
  }: {
    error: unknown;
    messageAuth: string;
    messageInternal: string;
  }) => {
    if (!(error instanceof AxiosError)) {
      showError(messageInternal);
      return;
    }

    const status = error.response?.status;
    switch (status) {
      case 401:
        showWarning(messageAuth);
        break;
      case 403:
        showWarning(messageAuth);
        break;
      default:
        showError(messageInternal);
        break;
    }
  };

  return { handleError };
};

export default useHandleError;
