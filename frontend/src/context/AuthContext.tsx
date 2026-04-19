import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type PropsWithChildren,
} from "react";
import { usersApi } from "../api/users";
import type { LoginPayload, User, UserUpsertPayload } from "../types/user";
import { clearSessionPassword, saveSessionPassword } from "./authSession";

export type UserMode = "buyer" | "seller";

interface AuthState {
  user: User | null;
  token: string | null;
  mode: UserMode;
}

interface AuthContextValue {
  user: User | null;
  token: string | null;
  mode: UserMode;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (payload: LoginPayload) => Promise<void>;
  register: (payload: UserUpsertPayload) => Promise<User>;
  logout: () => void;
  refreshUser: () => Promise<User | null>;
  setUser: (nextUser: User | null) => void;
  setMode: (nextMode: UserMode) => void;
}

const STORAGE_KEY = "getapet.auth.v2";

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

const isUserMode = (value: unknown): value is UserMode => value === "buyer" || value === "seller";

const defaultModeForUser = (user: User | null): UserMode => (user?.status === "seller" ? "seller" : "buyer");

const isUserShape = (value: unknown): value is User => {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const candidate = value as Partial<User>;
  return typeof candidate.id === "string" && typeof candidate.user_login === "string";
};

const isAuthStateShape = (value: unknown): value is Partial<AuthState> => {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const candidate = value as Partial<AuthState>;
  if (candidate.token !== null && candidate.token !== undefined && typeof candidate.token !== "string") {
    return false;
  }

  if (candidate.user !== null && candidate.user !== undefined && !isUserShape(candidate.user)) {
    return false;
  }

  if (candidate.mode !== undefined && candidate.mode !== null && !isUserMode(candidate.mode)) {
    return false;
  }

  return true;
};

const readStoredAuth = (): AuthState => {
  const rawValue = localStorage.getItem(STORAGE_KEY);
  if (!rawValue) {
    return { user: null, token: null, mode: "buyer" };
  }

  try {
    const parsed = JSON.parse(rawValue) as unknown;
    if (!isAuthStateShape(parsed)) {
      return { user: null, token: null, mode: "buyer" };
    }

    const parsedUser = parsed.user ?? null;
    return {
      user: parsedUser,
      token: parsed.token ?? null,
      mode: isUserMode(parsed.mode) ? parsed.mode : defaultModeForUser(parsedUser),
    };
  } catch {
    return { user: null, token: null, mode: "buyer" };
  }
};

export function AuthProvider({ children }: PropsWithChildren) {
  const [authState, setAuthState] = useState<AuthState>(() => readStoredAuth());
  const [isLoading, setIsLoading] = useState(false);

  const persistAuth = useCallback((nextState: AuthState) => {
    setAuthState(nextState);

    if (!nextState.user || !nextState.token) {
      localStorage.removeItem(STORAGE_KEY);
      return;
    }

    localStorage.setItem(STORAGE_KEY, JSON.stringify(nextState));
  }, []);

  const performLogin = useCallback(
    async ({ user_login, user_password }: LoginPayload): Promise<User> => {
      const loginResult = await usersApi.login({ user_login, user_password });
      const currentUser = await usersApi.getById(loginResult.user_id, loginResult.token);
      persistAuth({
        user: currentUser,
        token: loginResult.token,
        mode: defaultModeForUser(currentUser),
      });
      saveSessionPassword(user_password);
      return currentUser;
    },
    [persistAuth],
  );

  const login = useCallback(
    async (payload: LoginPayload) => {
      setIsLoading(true);
      try {
        await performLogin(payload);
      } finally {
        setIsLoading(false);
      }
    },
    [performLogin],
  );

  const register = useCallback(
    async (payload: UserUpsertPayload) => {
      setIsLoading(true);
      try {
        await usersApi.register(payload);
        return await performLogin({
          user_login: payload.user_login,
          user_password: payload.user_password,
        });
      } finally {
        setIsLoading(false);
      }
    },
    [performLogin],
  );

  const logout = useCallback(() => {
    persistAuth({ user: null, token: null, mode: "buyer" });
    clearSessionPassword();
  }, [persistAuth]);

  const refreshUser = useCallback(async () => {
    if (!authState.user || !authState.token) {
      return null;
    }

    const updatedUser = await usersApi.getById(authState.user.id, authState.token);
    persistAuth({
      user: updatedUser,
      token: authState.token,
      mode: authState.mode,
    });
    return updatedUser;
  }, [authState.mode, authState.token, authState.user, persistAuth]);

  const setUser = useCallback(
    (nextUser: User | null) => {
      if (!nextUser) {
        persistAuth({ user: null, token: null, mode: "buyer" });
        return;
      }

      persistAuth({
        user: nextUser,
        token: authState.token,
        mode: authState.mode,
      });
    },
    [authState.mode, authState.token, persistAuth],
  );

  const setMode = useCallback(
    (nextMode: UserMode) => {
      if (!authState.user || !authState.token) {
        return;
      }

      persistAuth({
        user: authState.user,
        token: authState.token,
        mode: nextMode,
      });
    },
    [authState.token, authState.user, persistAuth],
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      user: authState.user,
      token: authState.token,
      mode: authState.mode,
      isLoading,
      isAuthenticated: Boolean(authState.user && authState.token),
      login,
      register,
      logout,
      refreshUser,
      setUser,
      setMode,
    }),
    [
      authState.mode,
      authState.token,
      authState.user,
      isLoading,
      login,
      logout,
      refreshUser,
      register,
      setMode,
      setUser,
    ],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export const useAuthContext = (): AuthContextValue => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuthContext must be used inside AuthProvider");
  }

  return context;
};
