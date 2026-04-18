import { Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "./components/layout/AppLayout";
import { CatalogPage } from "./pages/CatalogPage";

function App() {
  return (
    <AppLayout>
      <Routes>
        <Route path="/" element={<CatalogPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </AppLayout>
  );
}

export default App;
