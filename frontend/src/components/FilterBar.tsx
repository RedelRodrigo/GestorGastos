import { Filter } from "lucide-react";
import type { ResumenItem } from "../types";

interface FilterBarProps {
  filtroCategoria: string;
  setFiltroCategoria: (cat: string) => void;
  categoriasDisponibles: ResumenItem[];
}

export default function FilterBar({
  filtroCategoria,
  setFiltroCategoria,
  categoriasDisponibles,
}: FilterBarProps) {
  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-4 flex flex-wrap items-center justify-between gap-4">
      <div className="flex items-center gap-2 text-zinc-300 text-sm font-medium">
        <Filter className="w-4 h-4 text-emerald-400" />
        <span>Filtrar movimientos:</span>
      </div>
      <div className="flex items-center gap-2">
        <select
          value={filtroCategoria}
          onChange={(e) => setFiltroCategoria(e.target.value)}
          className="bg-zinc-950 border border-zinc-800 text-zinc-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
          <option value="">Todas las categorías</option>
          {categoriasDisponibles?.map((r) => (
            <option key={r.categoria} value={r.categoria}>
              {r.categoria}
            </option>
          ))}
        </select>
        {filtroCategoria && (
          <button
            onClick={() => setFiltroCategoria("")}
            className="text-xs text-zinc-400 hover:text-zinc-200 underline px-2 py-1 transition-colors"
          >
            Limpiar filtro
          </button>
        )}
      </div>
    </div>
  );
}
