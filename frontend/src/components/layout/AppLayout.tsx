import type { ReactNode } from "react";
import { NavLink } from "react-router-dom";

interface AppLayoutProps {
  children: ReactNode;
}

const linkClassName = ({ isActive }: { isActive: boolean }): string =>
  isActive ? "main-nav-link active" : "main-nav-link";

export function AppLayout({ children }: AppLayoutProps) {
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
        </nav>
      </header>

      <main className="page-main">{children}</main>
    </div>
  );
}
