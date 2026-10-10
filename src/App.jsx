import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./auth";
import { Layout, NotFound } from "./components/Shell";
import AuthPage from "./pages/AuthPage";
import Notes from "./pages/Notes";

function Private({ children }) {
  const { user } = useAuth();
  return user ? children : <Navigate to="/masuk" replace />;
}

function Guest({ children }) {
  const { user } = useAuth();
  return user ? <Navigate to="/" replace /> : children;
}

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/masuk" element={<Guest><AuthPage mode="login" /></Guest>} />
        <Route path="/daftar" element={<Guest><AuthPage mode="register" /></Guest>} />
        <Route path="/" element={<Private><Notes /></Private>} />
        <Route path="*" element={<NotFound />} />
      </Route>
    </Routes>
  );
}
