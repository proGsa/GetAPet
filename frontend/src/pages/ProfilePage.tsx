import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { usersApi } from "../api/users";
import { AlertMessage } from "../components/ui/AlertMessage";
import { readSessionPassword, saveSessionPassword } from "../context/authSession";
import { useAuth } from "../hooks/useAuth";
import type { User, UserUpsertPayload } from "../types/user";
import { getErrorMessage } from "../utils/error";

interface ProfileFormState {
  fio: string;
  telephone_number: string;
  city: string;
  user_login: string;
  user_description: string;
}

const mapUserToForm = (user: User): ProfileFormState => ({
  fio: user.fio,
  telephone_number: user.telephone_number,
  city: user.city,
  user_login: user.user_login,
  user_description: user.user_description,
});

export function ProfilePage() {
  const navigate = useNavigate();
  const { user, token, setUser, logout } = useAuth();

  const [form, setForm] = useState<ProfileFormState | null>(user ? mapUserToForm(user) : null);
  const [newPassword, setNewPassword] = useState("");
  const [confirmNewPassword, setConfirmNewPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    setForm(user ? mapUserToForm(user) : null);
    setNewPassword("");
    setConfirmNewPassword("");
  }, [user]);

  if (!user || !token || !form) {
    return (
      <section className="page-content narrow-page">
        <EmptyProfile />
      </section>
    );
  }

  const setField = (key: keyof ProfileFormState, value: string) => {
    setForm((current) => (current ? { ...current, [key]: value } : current));
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setMessage(null);

    const isPasswordChange = Boolean(newPassword || confirmNewPassword);
    if (isPasswordChange && newPassword !== confirmNewPassword) {
      setMessage("Пароли не совпадают.");
      return;
    }

    if (newPassword.length > 0 && newPassword.length < 6) {
      setMessage("Пароль должен содержать минимум 6 символов.");
      return;
    }

    const fallbackPassword = readSessionPassword();
    const passwordForRequest = isPasswordChange ? newPassword : (fallbackPassword ?? "");

    if (!passwordForRequest) {
      setMessage("Не удалось определить пароль для сохранения. Войдите заново или задайте новый пароль.");
      return;
    }

    const payload: UserUpsertPayload = {
      ...form,
      status: user.status,
      user_password: passwordForRequest,
    };

    setIsSubmitting(true);

    try {
      const updatedUser = await usersApi.update(user.id, payload, token);
      setUser(updatedUser);
      saveSessionPassword(passwordForRequest);
      setNewPassword("");
      setConfirmNewPassword("");
      setMessage("Профиль обновлен.");
    } catch (updateError) {
      setMessage(getErrorMessage(updateError, "Не удалось обновить профиль"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async () => {
    const shouldDelete = window.confirm("Удалить профиль безвозвратно?");
    if (!shouldDelete) {
      return;
    }

    setMessage(null);
    setIsSubmitting(true);

    try {
      await usersApi.remove(user.id, token);
      logout();
      navigate("/", { replace: true });
    } catch (deleteError) {
      setMessage(getErrorMessage(deleteError, "Не удалось удалить профиль"));
      setIsSubmitting(false);
    }
  };

  const handleLogout = () => {
    logout();
    navigate("/", { replace: true });
  };

  return (
    <section className="page-content narrow-page">
      <div className="page-title-row">
        <h1>Профиль</h1>
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

          <div className="profile-row-break" aria-hidden="true" />

          <label>
            Логин
            <input
              value={form.user_login}
              onChange={(event) => setField("user_login", event.target.value)}
              required
            />
          </label>

          <label>
            Новый пароль
            <input
              type="password"
              value={newPassword}
              onChange={(event) => setNewPassword(event.target.value)}
              autoComplete="new-password"
            />
          </label>

          <label>
            Подтвердите новый пароль
            <input
              type="password"
              value={confirmNewPassword}
              onChange={(event) => setConfirmNewPassword(event.target.value)}
              autoComplete="new-password"
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

          <div className="button-row profile-button-row wide-label">
            <button
              type="submit"
              className="primary-button profile-action-button profile-save-button"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Сохранение..." : "Сохранить изменения"}
            </button>

            <div className="profile-right-actions">
              <button
                type="button"
                className="danger-outline-button profile-action-button"
                disabled={isSubmitting}
                onClick={handleLogout}
              >
                Выйти
              </button>
              <button
                type="button"
                className="danger-button profile-action-button"
                disabled={isSubmitting}
                onClick={() => {
                  void handleDelete();
                }}
              >
                Удалить профиль
              </button>
            </div>
          </div>
        </form>

        {message ? <AlertMessage>{message}</AlertMessage> : null}
      </article>
    </section>
  );
}

function EmptyProfile() {
  return (
    <article className="panel">
      <h1>Профиль</h1>
      <p>Вы просматриваете сайт в гостевом режиме.</p>
      <p>
        <Link to="/login">Войдите</Link> или <Link to="/register">зарегистрируйтесь</Link>, чтобы
        управлять профилем.
      </p>
    </article>
  );
}
