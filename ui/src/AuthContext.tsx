import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import Cookies, { CookieAttributes } from 'js-cookie';
import { debounce } from 'lodash';
import { checkAuth, signIn } from './client/http';

interface AuthContextType {
  authenticated: boolean;
  isLoading: boolean;
  login: (token: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
  authenticated: false,
  isLoading: true,
  login: async () => {},
  logout: () => {},
});

const cookieName = 'rush-auth';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const verifyAuth = useCallback(async () => {
    try {
      // TODO(#105): Refresh the token if it's expired.
      const isAuthenticated = await checkAuth();
      setAuthenticated(isAuthenticated);
    } catch (error) {
      // TODO(#65): Handle this error globally.
      // eslint-disable-next-line no-console
      console.error('Error verifying authentication:', error);
      logout();
    } finally {
      setIsLoading(false);
    }
  }, []);

  const debouncedVerifyAuth = useMemo(() => debounce(verifyAuth, 100, { leading: true, trailing: true }), [verifyAuth]);

  useEffect(() => {
    debouncedVerifyAuth();

    const handleVisibilityChange = () => {
      // When the tab is visible, verify auth.
      if (!document.hidden) {
        debouncedVerifyAuth();
      }
    };

    const handlePageShow = (event: PageTransitionEvent) => {
      // When it's not the first load, verify auth.
      if (event.persisted) {
        debouncedVerifyAuth();
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    window.addEventListener('pageshow', handlePageShow);

    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      window.removeEventListener('pageshow', handlePageShow);
      debouncedVerifyAuth.cancel();
    };
  }, [debouncedVerifyAuth]);

  const login = useCallback(
    async (token: string) => {
      if (token) {
        const rushToken = await signIn(token);
        Cookies.set(cookieName, rushToken, getCookieOptions());
        setAuthenticated(true);
      }
    },
    [setAuthenticated],
  );

  const logout = () => {
    Cookies.remove(cookieName, getCookieOptions());
    setAuthenticated(false);
  };

  const value = useMemo(() => ({ authenticated, isLoading, login, logout }), [authenticated, isLoading, login]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

const hostName = new URL(import.meta.env.VITE_SERVER_ENDPOINT).hostname;
const env = import.meta.env.VITE_ENV;

const getCookieOptions = () => {
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

export const useAuth = () => useContext(AuthContext);
