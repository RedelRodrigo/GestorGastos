import type { ResumenItem } from "../types";

interface CategoryProgressBarProps {
  item: ResumenItem;
  porcentaje: number;
}

export default function CategoryProgressBar({
  item,
  porcentaje,
}: CategoryProgressBarProps) {
  return (
    <div className="space-y-1">
      <div className="flex justify-between text-sm">
        <span className="font-medium text-zinc-200">{item.categoria}</span>
        <span className="text-emerald-400 font-semibold tabular-nums">
          $
          {(item?.total ?? 0).toLocaleString("es-AR", {
            minimumFractionDigits: 2,
          })}
        </span>
      </div>
      <div className="w-full bg-zinc-950 rounded-full h-2 overflow-hidden border border-zinc-800/60">
        <div
          className="bg-emerald-500 h-full rounded-full transition-all duration-500"
          style={{ width: `${porcentaje}%` }}
        />
      </div>
    </div>
  );
}
