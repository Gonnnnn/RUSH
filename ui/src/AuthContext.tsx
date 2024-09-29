import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { debounce } from 'lodash';
import { removeAuthToken, setAuthToken } from './client/auth';
import { checkAuth, signIn } from './client/http';

interface AuthContextType {
  authenticated: boolean;
  isLoading: boolean;
  login: (token: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
  // The indicator of whether the user is authenticated. It is refreshed when the page is loaded, users open the tab again.
  // It is also updated when the user logs in or out.
  authenticated: false,
  // The indicator of whether the authentication is being checked. It is provided outside so that the app can show a loading indicator
  // to initialize the authentication state before the page is loaded.
  isLoading: true,
  // The function to log in the user. Components should use this to log in the user.
  login: async () => {},
  // The function to log out the user. Components should use this to log out the user.
  logout: () => {},
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const verifyAuth = useCallback(async () => {
    try {
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
        setAuthToken(await signIn(token));
        setAuthenticated(true);
      }
    },
    [setAuthenticated],
  );

  const logout = () => {
    removeAuthToken();
    setAuthenticated(false);
  };

  const value = useMemo(() => ({ authenticated, isLoading, login, logout }), [authenticated, isLoading, login]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);
