import { TrendingDown } from "lucide-react";
import type { ResumenItem } from "../types";
import CategoryProgressBar from "./CategoryProgressBar";

interface CategorySummaryProps {
  resumen: ResumenItem[];
  total: number;
}

export default function CategorySummary({
  resumen,
  total,
}: CategorySummaryProps) {
  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-4">
      <h3 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider flex items-center gap-2">
        <TrendingDown className="w-4 h-4 text-emerald-400" />
        Resumen por Categoría
      </h3>

      {resumen?.length === 0 ? (
        <p className="text-sm text-zinc-500 italic">
          No hay datos registrados aún.
        </p>
      ) : (
        <div className="space-y-3">
          {resumen?.map((item) => {
            const porcentaje = total > 0 ? (item.total / total) * 100 : 0;
            return (
              <CategoryProgressBar
                key={item.categoria}
                item={item}
                porcentaje={porcentaje}
              />
            );
          })}
        </div>
      )}
    </div>
  );
}
