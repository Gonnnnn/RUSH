import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';
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

  useEffect(() => {
    const verifyAuth = async () => {
      try {
        const isAuthenticated = await checkAuth();
        setAuthenticated(isAuthenticated);
      } catch (error) {
        // TODO(#65): Handle this error globally.
        // eslint-disable-next-line no-console
        console.error('Error verifying authentication:', error);
        logout();
      }
    };

    verifyAuth();
  }, []);

  const login = useCallback(
    async (token: string) => {
      if (token) {
        try {
          const cookie = await signIn(token);
          document.cookie = `${cookieName}=${cookie}`;
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
    document.cookie = `${cookieName}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`;

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
