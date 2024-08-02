import { createContext, ReactNode, useContext, useMemo, useState } from 'react';

interface AuthContextType {
  authenticated: boolean;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
  authenticated: false,
  login: () => {},
  logout: () => {},
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [authenticated, setAuthenticated] = useState(false);

  const login = (token: string) => {
    if (token) {
      setAuthenticated(true);
      return;
    }
    setAuthenticated(false);
  };

  const logout = () => {
    setAuthenticated(false);
  };

  const value = useMemo(() => ({ authenticated, login, logout }), [authenticated]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => useContext(AuthContext);
