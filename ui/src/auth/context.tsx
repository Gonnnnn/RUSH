import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { checkAuth, signIn } from '../client/http';

interface AuthContextType {
  authenticated: boolean;
  login: (token: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
  authenticated: false,
  login: async () => {},
  logout: () => {},
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

  const logout = () => {
    document.cookie = '';
    setAuthenticated(false);
  };

  const value = useMemo(() => ({ authenticated, login, logout }), [authenticated, login]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);
