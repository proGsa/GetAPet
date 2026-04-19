import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { AlertMessage } from "../components/ui/AlertMessage";
import { useAuth } from "../hooks/useAuth";
import type { UserRole, UserUpsertPayload } from "../types/user";
import { getErrorMessage } from "../utils/error";

const DEFAULT_ROLE: UserRole = "buyer";

const createDefaultPayload = (): UserUpsertPayload => ({
  fio: "",
  telephone_number: "",
  city: "",
  user_login: "",
  user_password: "",
  status: DEFAULT_ROLE,
  user_description: "",
});

export function RegisterPage() {
  const { user, register, isLoading } = useAuth();
  const navigate = useNavigate();

  const [form, setForm] = useState<UserUpsertPayload>(createDefaultPayload);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (user) {
      navigate("/profile", { replace: true });
    }
  }, [navigate, user]);

  const setField = (key: keyof UserUpsertPayload, value: string) => {
    setForm((current) => ({ ...current, [key]: value }));
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);

    try {
      await register({
        ...form,
        status: DEFAULT_ROLE,
      });
      navigate("/profile", { replace: true });
    } catch (registerError) {
      setError(getErrorMessage(registerError, "Не удалось зарегистрироваться"));
    }
  };

  return (
    <section className="page-content narrow-page">
      <div className="page-title-row">
        <h1>Регистрация</h1>
      </div>

      <article className="panel">
        <form className="form-grid" onSubmit={handleSubmit}>
          <label>
            ФИО
            <input value={form.fio} onChange={(event) => setField("fio", event.target.value)} required />
          </label>

          <label>
            Номер телефона
            <input
              value={form.telephone_number}
              onChange={(event) => setField("telephone_number", event.target.value)}
              required
              pattern="^(8[0-9]{10}|\\+7[0-9]{10})$"
              title="Формат: 8XXXXXXXXXX или +7XXXXXXXXXX"
            />
          </label>

          <label>
            Город
            <input value={form.city} onChange={(event) => setField("city", event.target.value)} />
          </label>

          <div className="auth-row-break" aria-hidden="true" />

          <label>
            Логин
            <input
              value={form.user_login}
              onChange={(event) => setField("user_login", event.target.value)}
              required
            />
          </label>

          <label>
            Пароль
            <input
              type="password"
              value={form.user_password}
              onChange={(event) => setField("user_password", event.target.value)}
              required
            />
          </label>

          <label className="wide-label">
            Описание
            <textarea
              value={form.user_description}
              rows={3}
              onChange={(event) => setField("user_description", event.target.value)}
            />
          </label>

          <button type="submit" className="primary-button" disabled={isLoading}>
            {isLoading ? "Создание..." : "Создать аккаунт"}
          </button>
        </form>

        {error ? <AlertMessage variant="error">{error}</AlertMessage> : null}

        <p className="hint-text">
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </p>
      </article>
    </section>
  );
}
