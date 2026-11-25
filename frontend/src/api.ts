import type { Gasto, GastoInput, Usuario } from "./types";

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
  captcha: string,
): Promise<Usuario> => {
  const res = await fetch("/api/auth/registro", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, username, password, captcha }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error ?? "Error al registrar cuenta");
  return body.data;
};

export const loginApi = async (
  email: string,
  password: string,
  captcha: string,
): Promise<Usuario> => {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password, captcha }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(body.error ?? "Credenciales incorrectas");
  return body.data;
};

export const logoutApi = async (): Promise<void> => {
  await fetch("/api/auth/logout", {
    method: "POST",
    credentials: "include",
  }).catch(() => {});
  cerrarSesionUsuario();
};

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch("/api" + path, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
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

export const listarGastos = async (query = "") => {
  const res = await request<any>(`/gastos${query}`);

  if (Array.isArray(res)) return res;
  if (res && Array.isArray(res.data)) return res.data;

  return [];
};

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

export const obtenerResumen = async () => {
  const res = await request<any>("/gastos/resumen");

  if (Array.isArray(res)) {
    const totalCalculado = res.reduce(
      (acc, item) => acc + (item.total || 0),
      0,
    );
    return { data: res, total: totalCalculado };
  }

  const lista = Array.isArray(res?.data) ? res.data : [];
  const total =
    typeof res?.total === "number"
      ? res.total
      : lista.reduce((acc: number, item: any) => acc + (item.total || 0), 0);

  return {
    data: lista,
    total: total,
  };
};
