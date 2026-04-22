import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { AlertMessage } from "../components/ui/AlertMessage";
import { useAuth } from "../hooks/useAuth";
import { getErrorMessage } from "../utils/error";

export function LoginPage() {
  const { user, login, isLoading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const [userLogin, setUserLogin] = useState("");
  const [userPassword, setUserPassword] = useState("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (user) {
      navigate("/", { replace: true });
    }
  }, [navigate, user]);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);

    try {
      await login({ user_login: userLogin, user_password: userPassword });
      const targetPath =
        typeof location.state === "object" &&
        location.state !== null &&
        "from" in location.state &&
        typeof location.state.from === "string"
          ? location.state.from
          : "/";

      navigate(targetPath, { replace: true });
    } catch (loginError) {
      setError(getErrorMessage(loginError, "Не удалось выполнить вход"));
    }
  };

  return (
    <section className="page-content narrow-page">
      <div className="page-title-row">
        <h1>Вход</h1>
      </div>

      <article className="panel">
        <form className="auth-login-form" onSubmit={handleSubmit}>
          <label>
            Логин
            <input
              name="user_login"
              value={userLogin}
              onChange={(event) => setUserLogin(event.target.value)}
              required
              autoComplete="username"
            />
          </label>

          <label>
            Пароль
            <input
              name="user_password"
              type="password"
              value={userPassword}
              onChange={(event) => setUserPassword(event.target.value)}
              required
              autoComplete="current-password"
            />
          </label>

          <button type="submit" className="primary-button" disabled={isLoading}>
            {isLoading ? "Проверка..." : "Войти"}
          </button>

          <p className="hint-text auth-login-hint">
            Еще нет аккаунта? <Link to="/register">Зарегистрируйтесь</Link>
          </p>
        </form>

        {error ? <AlertMessage variant="error">{error}</AlertMessage> : null}
      </article>
    </section>
  );
}
