import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import Cookies from 'js-cookie';
import { debounce } from 'lodash';
import { checkAuth, signIn } from './client/http';

interface AuthContextType {
  authenticated: boolean;
  login: (token: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  authenticated: false,
  login: async () => {},
  logout: async () => {},
});

const cookieName = 'rush-auth';

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false);

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
        try {
          const cookie = await signIn(token);
          Cookies.set(cookieName, cookie, {
            // Longer than the token expiration time so that it'll be sent again and refreshed.
            expires: 30,
            secure: import.meta.env.VITE_ENV !== 'local',
            domain: import.meta.env.VITE_SERVER_ENDPOINT,
            sameSite: import.meta.env.VITE_ENV === 'local' ? 'None' : 'Strict',
            path: '/',
          });
          setAuthenticated(true);
        } catch (error) {
          // TODO(#65): Handle this error globally.
          // eslint-disable-next-line no-console
          console.error(error);
          logout();
        }
      }
    },
    [setAuthenticated],
  );

  const logout = async () => {
    Cookies.remove(cookieName, { path: '/' });

    // Set timer to reload the page just to make users feel like they are logged out.
    // That's why the cookie is removed before the timer is initialized.
    return new Promise<void>((resolve) => {
      setTimeout(() => {
        setAuthenticated(false);
        resolve();
      }, 1000);
    });
  };

  const value = useMemo(() => ({ authenticated, login, logout }), [authenticated, login]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);
