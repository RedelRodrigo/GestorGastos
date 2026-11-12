import type { Gasto, GastoInput, ResumenItem } from "./types";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch("/api" + path, {
    headers: { "Content-Type": "application/json" },
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
