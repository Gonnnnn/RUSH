import { AxiosResponse } from 'axios';
import { replaceCookieHeader, setAuthToken } from '../auth';

const setToken = (response: AxiosResponse) => {
  const newToken = response.headers[replaceCookieHeader];
  if (newToken) {
    setAuthToken(newToken);
  }
  return response;
};

export default setToken;
