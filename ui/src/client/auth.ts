import Cookies, { CookieAttributes } from 'js-cookie';

export const replaceCookieHeader = 'X-Replace-Cookie';

export const setAuthToken = (token: string) => {
  Cookies.set(authCookieName, token, getAuthCookieOptions());
};

export const removeAuthToken = () => {
  Cookies.remove(authCookieName, getAuthCookieOptions());
};

const BASE_URL = import.meta.env.VITE_SERVER_ENDPOINT;
const hostName = new URL(BASE_URL).hostname;
const env = import.meta.env.VITE_ENV;
const authCookieName = 'rush-auth';

const getAuthCookieOptions = () => {
  const cookieOptions: CookieAttributes = {
    expires: 30,
    domain: hostName,
    path: '/',
  };
  if (env === 'local') {
    return cookieOptions;
  }
  cookieOptions.secure = true;
  cookieOptions.sameSite = 'Strict';
  return cookieOptions;
};
