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

  // 🟢 Agrega este log para ver exactamente qué responde el servidor de Go
  console.log(
    `[API Response] ${init?.method || "GET"} ${path} (Status: ${res.status}):`,
    body,
  );

  if (!res.ok) {
    throw new Error(
      (body as { error?: string }).error ?? `Error HTTP: ${res.status}`,
    );
  }
  return body as T;
}

export const listarGastos = async (query = "") => {
  const res = await request<any>(`/gastos${query}`);

  // Si 'res' es un array directamente
  if (Array.isArray(res)) return res;

  // Si viene envuelto en .data
  if (res && Array.isArray(res.data)) return res.data;

  // Si el servidor devolvió {} u otra cosa extraña, devolvemos array vacío
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

  // Si res es un array directamente
  if (Array.isArray(res)) {
    const totalCalculado = res.reduce(
      (acc, item) => acc + (item.total || 0),
      0,
    );
    return { data: res, total: totalCalculado };
  }

  // Si res tiene .data
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
