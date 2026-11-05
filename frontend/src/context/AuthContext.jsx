import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react';
import { api, setUnauthorizedHandler, tokenStore } from '../api/client';

const AuthContext = createContext(null);
export const useAuth = () => useContext(AuthContext);

export function AuthProvider({ children }) {
  const [profile, setProfile] = useState(null); // { user, permissions }
  const [loading, setLoading] = useState(!!tokenStore.get());

  const logout = useCallback(() => { tokenStore.clear(); setProfile(null); }, []);

  useEffect(() => {
    setUnauthorizedHandler(logout);
    if (!tokenStore.get()) return;
    api('/auth/me').then(setProfile).catch(logout).finally(() => setLoading(false));
  }, [logout]);

  const login = useCallback(async (email, password, remember) => {
    const res = await api('/auth/login', { method: 'POST', body: { email, password, remember } });
    tokenStore.set(res.token, remember);
    setProfile(res.profile);
  }, []);

  const value = useMemo(() => ({
    profile, loading, login, logout,
    user: profile?.user,
    can: (code) => !!profile?.permissions.includes(code),
  }), [profile, loading, login, logout]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
