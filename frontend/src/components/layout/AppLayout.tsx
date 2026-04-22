import type { ReactNode } from "react";
import { NavLink } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";

interface AppLayoutProps {
  children: ReactNode;
}

const linkClassName = ({ isActive }: { isActive: boolean }): string =>
  isActive ? "main-nav-link active" : "main-nav-link";

const modeButtonClassName = (isActive: boolean): string =>
  isActive ? "mode-switch-button active" : "mode-switch-button";

export function AppLayout({ children }: AppLayoutProps) {
  const { user, mode, setMode, logout } = useAuth();

  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="topbar-brand">
          <NavLink to="/" className="brand-link">
            GetAPet
          </NavLink>
          <p className="brand-subtitle">С любовью к животным</p>
        </div>

        <nav className="main-nav" aria-label="Основная навигация">
          <NavLink to="/" className={linkClassName} end>
            Каталог
          </NavLink>

          {user && mode === "seller" ? (
            <NavLink to="/my-pets" className={linkClassName}>
              Мои объявления
            </NavLink>
          ) : null}

          {user ? (
            <NavLink to="/requests" className={linkClassName}>
              Мои заявки
            </NavLink>
          ) : null}

          {user ? (
            <NavLink to="/profile" className={linkClassName}>
              Профиль
            </NavLink>
          ) : (
            <>
              <NavLink to="/login" className={linkClassName}>
                Войти
              </NavLink>
              <NavLink to="/register" className={linkClassName}>
                Регистрация
              </NavLink>
            </>
          )}
        </nav>

        <div className="topbar-user">
          {user ? (
            <>
              <p>{user.fio}</p>

              <div className="topbar-actions">
                <div className="mode-switch" role="group" aria-label="Режим работы">
                  <button
                    type="button"
                    className={modeButtonClassName(mode === "buyer")}
                    onClick={() => setMode("buyer")}
                  >
                    Покупатель
                  </button>
                  <button
                    type="button"
                    className={modeButtonClassName(mode === "seller")}
                    onClick={() => setMode("seller")}
                  >
                    Продавец
                  </button>
                </div>

                <button type="button" className="danger-outline-button" onClick={logout}>
                  Выйти
                </button>
              </div>
            </>
          ) : (
            <p>Гостевой режим</p>
          )}
        </div>
      </header>

      <main className="page-main">{children}</main>
    </div>
  );
}
