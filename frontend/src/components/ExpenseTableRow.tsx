import { Pencil, Trash2 } from "lucide-react";
import type { Gasto } from "../types";

interface ExpenseTableRowProps {
  gasto: Gasto;
  onEditar: (g: Gasto) => void;
  onEliminar: (id: number) => void;
}

export default function ExpenseTableRow({
  gasto,
  onEditar,
  onEliminar,
}: ExpenseTableRowProps) {
  return (
    <tr className="border-b border-zinc-800/50 hover:bg-zinc-800/30 transition-colors">
      <td className="py-3 px-4 text-sm text-zinc-300">{gasto.fecha}</td>
      <td className="py-3 px-4 text-sm font-medium text-zinc-100">
        {gasto.categoria}
      </td>
      <td className="py-3 px-4 text-sm text-zinc-400 max-w-xs truncate">
        {gasto.descripcion || "—"}
      </td>
      <td className="py-3 px-4 text-sm font-semibold text-emerald-400 text-right">
        $
        {(gasto?.monto ?? 0).toLocaleString("es-AR", {
          minimumFractionDigits: 2,
        })}
      </td>
      <td className="py-3 px-4 text-right space-x-2">
        <button
          onClick={() => onEditar(gasto)}
          className="p-1.5 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-lg transition-colors"
          title="Editar"
        >
          <Pencil className="w-4 h-4" />
        </button>
        <button
          onClick={() => onEliminar(gasto.id)}
          className="p-1.5 bg-zinc-800 hover:bg-rose-500/20 hover:text-rose-400 text-zinc-400 rounded-lg transition-colors"
          title="Eliminar"
        >
          <Trash2 className="w-4 h-4" />
        </button>
      </td>
    </tr>
  );
}
