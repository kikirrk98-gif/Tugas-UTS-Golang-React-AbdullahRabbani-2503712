import axios from "axios";

// Di production Go menyajikan React dari origin yang sama, jadi baseURL cukup "/api".
const api = axios.create({ baseURL: import.meta.env.VITE_API_URL || "/api" });

api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Token kedaluwarsa atau tidak valid: bersihkan sesi dan kembali ke login.
api.interceptors.response.use(
  (res) => res,
  (err) => {
    const url = err.config?.url ?? "";
    if (err.response?.status === 401 && !url.startsWith("/auth/")) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.dispatchEvent(new Event("auth:logout"));
    }
    return Promise.reject(err);
  }
);

export const errMsg = (err) =>
  err.response?.data?.message ?? "Tidak bisa terhubung ke server. Pastikan backend Go berjalan.";

export default api;
