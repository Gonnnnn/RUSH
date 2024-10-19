import { createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { useAuth } from '../auth';
import { Role } from '../auth/role';

const AdminModeContext = createContext<{
  adminMode: boolean;
  setAdminMode: (adminMode: boolean) => void;
}>({ adminMode: false, setAdminMode: () => {} });

const AdminModeProvider = ({ children }: { children: ReactNode }) => {
  const { role } = useAuth();
  const [adminMode, setAdminMode] = useState(false);

  useEffect(() => {
    setAdminMode(role === Role.ADMIN);
  }, [role]);

  const setAdminModeIfAdmin = useCallback(
    (newAdminMode: boolean) => {
      if (role === Role.ADMIN) {
        setAdminMode(newAdminMode);
      }
    },
    [role, setAdminMode],
  );

  const value = useMemo(() => ({ adminMode, setAdminMode: setAdminModeIfAdmin }), [adminMode, setAdminModeIfAdmin]);
  return <AdminModeContext.Provider value={value}>{children}</AdminModeContext.Provider>;
};

const useAdminMode = () => useContext(AdminModeContext);

export { AdminModeProvider, useAdminMode };
