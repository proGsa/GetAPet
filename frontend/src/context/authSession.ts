const SESSION_PASSWORD_KEY = "getapet.auth.session-password.v1";

export const saveSessionPassword = (password: string) => {
  if (!password) {
    return;
  }

  sessionStorage.setItem(SESSION_PASSWORD_KEY, password);
};

export const readSessionPassword = (): string | null => {
  return sessionStorage.getItem(SESSION_PASSWORD_KEY);
};

export const clearSessionPassword = () => {
  sessionStorage.removeItem(SESSION_PASSWORD_KEY);
};
