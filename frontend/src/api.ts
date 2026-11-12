import type { Gasto, GastoInput, ResumenItem, Usuario } from "./types";

export const obtenerUsuarioActual = (): Usuario | null => {
  const data = localStorage.getItem("gestor_sesion_usuario");
  return data ? JSON.parse(data) : null;
};

export const guardarSesionUsuario = (usuario: Usuario) => {
  localStorage.setItem("gestor_sesion_usuario", JSON.stringify(usuario));
};

export const cerrarSesionUsuario = () => {
  localStorage.removeItem("gestor_sesion_usuario");
};

export const registroApi = async (
  email: string,
  username: string,
  password: string,
): Promise<Usuario> => {
  const res = await fetch("/api/auth/registro", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, username, password }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error ?? "Error al registrar cuenta");
  return body.data;
};

export const loginApi = async (
  email: string,
  password: string,
): Promise<Usuario> => {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error ?? "Credenciales incorrectas");
  return body.data;
};

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const usuario = obtenerUsuarioActual();
  const identificador = usuario ? usuario.email : "invitado";

  const res = await fetch("/api" + path, {
    headers: {
      "Content-Type": "application/json",
      "X-User-ID": identificador,
    },
    ...init,
  });

  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(
      (body as { error?: string }).error ?? `Error HTTP: ${res.status}`,
    );
  }
  return body as T;
}

export const listarGastos = (query = "") =>
  request<{ data: Gasto[] }>(`/gastos${query}`).then((r) => r.data);

export const crearGasto = (g: GastoInput) =>
  request<{ data: Gasto }>("/gastos", {
    method: "POST",
    body: JSON.stringify(g),
  });

export const editarGasto = (id: number, g: GastoInput) =>
  request<{ data: Gasto }>(`/gastos/${id}`, {
    method: "PUT",
    body: JSON.stringify(g),
  });

export const borrarGasto = (id: number) =>
  request<{ mensaje: string }>(`/gastos/${id}`, {
    method: "DELETE",
  });

export const obtenerResumen = () =>
  request<{ data: ResumenItem[]; total: number }>("/gastos/resumen");
