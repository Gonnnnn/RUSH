import { createContext, ReactNode, useContext, useEffect, useMemo, useState } from 'react';

const DEFAULT_TITLE = 'RU:SH';

const HeaderContext = createContext<{
  pageTitle: string;
  setPageTitle: (pageTitle: string) => void;
}>({ pageTitle: DEFAULT_TITLE, setPageTitle: () => {} });

const HeaderProvider = ({ children }: { children: ReactNode }) => {
  const [pageTitle, setPageTitle] = useState(DEFAULT_TITLE);

  const value = useMemo(() => ({ pageTitle, setPageTitle }), [pageTitle, setPageTitle]);
  return <HeaderContext.Provider value={value}>{children}</HeaderContext.Provider>;
};

const useHeader = ({ newTitle }: { newTitle?: string } = {}) => {
  const { pageTitle, setPageTitle } = useContext(HeaderContext);

  useEffect(() => {
    if (newTitle) {
      setPageTitle(newTitle);
    }
    return () => {
      setPageTitle(DEFAULT_TITLE);
    };
  }, [newTitle, setPageTitle]);

  return { pageTitle, setPageTitle };
};

export { HeaderProvider, useHeader };
