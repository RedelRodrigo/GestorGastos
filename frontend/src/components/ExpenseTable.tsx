import ExpenseTableRow from "./ExpenseTableRow";
import type { Gasto } from "../types";
import { Receipt } from "lucide-react";

interface ExpenseTableProps {
  gastos: Gasto[];
  onEditar: (g: Gasto) => void;
  onEliminar: (id: number) => void;
}

export default function ExpenseTable({
  gastos,
  onEditar,
  onEliminar,
}: ExpenseTableProps) {
  if (gastos?.length === 0) {
    return (
      <div className="bg-zinc-900 border border-zinc-800 rounded-3xl p-12 text-center space-y-3">
        <div className="inline-flex bg-zinc-800 p-3 rounded-2xl text-zinc-400">
          <Receipt className="w-6 h-6" />
        </div>
        <p className="text-sm font-medium text-zinc-300">
          No hay gastos registrados
        </p>
        <p className="text-xs text-zinc-500">
          Agregá tu primer movimiento utilizando el formulario.
        </p>
      </div>
    );
  }

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-3xl overflow-hidden shadow-lg">
      <div className="overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-zinc-800 text-[11px] uppercase tracking-wider text-zinc-400 bg-zinc-950/50">
              <th className="py-3 px-4 font-semibold">Fecha</th>
              <th className="py-3 px-4 font-semibold">Categoría</th>
              <th className="py-3 px-4 font-semibold">Descripción</th>
              <th className="py-3 px-4 font-semibold text-right">Monto</th>
              <th className="py-3 px-4 font-semibold text-right">Acciones</th>
            </tr>
          </thead>
          <tbody>
            {gastos?.map((gasto) => (
              <ExpenseTableRow
                key={gasto.id}
                gasto={gasto}
                onEditar={onEditar}
                onEliminar={onEliminar}
              />
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
