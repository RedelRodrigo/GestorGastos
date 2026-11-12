interface TotalCardProps {
  total: number;
}

export default function TotalCard({ total }: TotalCardProps) {
  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-2xl px-6 py-3 flex items-center gap-4 shadow-inner">
      <div>
        <span className="text-xs text-zinc-400 block font-medium uppercase tracking-wider">
          Gasto Total Acumulado
        </span>
        <span className="text-2xl font-bold text-emerald-400 tabular-nums">
          $
          {Number(total ?? 0).toLocaleString("es-AR", {
            minimumFractionDigits: 2,
          })}
        </span>
      </div>
    </div>
  );
}
