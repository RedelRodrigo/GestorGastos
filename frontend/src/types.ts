export interface Gasto {
  id: number;
  monto: number;
  categoria: string;
  fecha: string; // AAAA-MM-DD
  descripcion: string;
}

export type GastoInput = Omit<Gasto, "id">;

export interface ResumenItem {
  categoria: string;
  total: number;
}
