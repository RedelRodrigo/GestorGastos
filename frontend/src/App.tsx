import { useCallback, useEffect, useState } from "react";
import { borrarGasto, listarGastos, obtenerResumen } from "./api";
import GastoForm from "./components/GastoForm";
import type { Gasto, ResumenItem } from "./types";
import { Wallet, Trash2, Edit3, TrendingDown, Filter } from "lucide-react";

export default function App() {
  const [gastos, setGastos] = useState<Gasto[]>([]);
  const [resumen, setResumen] = useState<ResumenItem[]>([]);
  const [total, setTotal] = useState(0);
  const [editando, setEditando] = useState<Gasto | null>(null);
  const [error, setError] = useState("");

  // Filtros
  const [filtroCategoria, setFiltroCategoria] = useState("");

  const cargar = useCallback(async () => {
    try {
      const query = filtroCategoria
        ? `?categoria=${encodeURIComponent(filtroCategoria)}`
        : "";
      const [lista, res] = await Promise.all([
        listarGastos(query),
        obtenerResumen(),
      ]);
      setGastos(lista);
      setResumen(res.data);
      setTotal(res.total);
      setError("");
    } catch (e) {
      setError((e as Error).message);
    }
  }, [filtroCategoria]);

  useEffect(() => {
    cargar();
  }, [cargar]);

  async function eliminar(id: number) {
    if (!confirm("¿Estás seguro de eliminar este gasto?")) return;
    try {
      await borrarGasto(id);
      await cargar();
    } catch (e) {
      setError((e as Error).message);
    }
  }

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 p-4 md:p-8">
      <main className="max-w-7xl mx-auto space-y-8">
        {/* Header */}
        <header className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-800/80 pb-6">
          <div className="flex items-center gap-3">
            <div className="bg-emerald-500/10 p-3 rounded-2xl border border-emerald-500/20 text-emerald-400">
              <Wallet className="w-8 h-8" />
            </div>
            <div>
              <h1 className="text-2xl font-bold tracking-tight">
                Gestor de Gastos
              </h1>
              <p className="text-sm text-zinc-400">
                Control financiero personal de alto rendimiento
              </p>
            </div>
          </div>

          <div className="bg-zinc-900 border border-zinc-800 rounded-2xl px-6 py-3 flex items-center gap-4 shadow-inner">
            <div>
              <span className="text-xs text-zinc-400 block font-medium uppercase tracking-wider">
                Gasto Total Acumulado
              </span>
              <span className="text-2xl font-bold text-emerald-400 tabular-nums">
                ${total.toLocaleString("es-AR", { minimumFractionDigits: 2 })}
              </span>
            </div>
          </div>
        </header>

        {error && (
          <div
            role="alert"
            className="bg-rose-500/10 border border-rose-500/20 text-rose-400 p-4 rounded-xl text-sm"
          >
            {error}
          </div>
        )}

        {/* Layout Principal en 2 Columnas */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          {/* Columna Izquierda: Formulario y Resumen por Categoría */}
          <div className="lg:col-span-4 space-y-6">
            <GastoForm
              key={editando?.id ?? "nuevo"}
              inicial={editando}
              onGuardado={() => {
                setEditando(null);
                cargar();
              }}
              onCancelar={() => setEditando(null)}
            />

            {/* Tarjeta de Desglose por Categoría */}
            <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-4">
              <h3 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider flex items-center gap-2">
                <TrendingDown className="w-4 h-4 text-emerald-400" />
                Resumen por Categoría
              </h3>

              {resumen.length === 0 ? (
                <p className="text-sm text-zinc-500 italic">
                  No hay datos registrados aún.
                </p>
              ) : (
                <div className="space-y-3">
                  {resumen.map((item) => {
                    const porcentaje =
                      total > 0 ? (item.total / total) * 100 : 0;
                    return (
                      <div key={item.categoria} className="space-y-1">
                        <div className="flex justify-between text-sm">
                          <span className="font-medium text-zinc-200">
                            {item.categoria}
                          </span>
                          <span className="text-emerald-400 font-semibold tabular-nums">
                            $
                            {item.total.toLocaleString("es-AR", {
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
                  })}
                </div>
              )}
            </div>
          </div>

          {/* Columna Derecha: Filtros y Tabla de Movimientos */}
          <div className="lg:col-span-8 space-y-6">
            {/* Barra de Filtros */}
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
                  {resumen.map((r) => (
                    <option key={r.categoria} value={r.categoria}>
                      {r.categoria}
                    </option>
                  ))}
                </select>
                {filtroCategoria && (
                  <button
                    onClick={() => setFiltroCategoria("")}
                    className="text-xs text-zinc-400 hover:text-zinc-200 underline px-2 py-1"
                  >
                    Limpiar filtro
                  </button>
                )}
              </div>
            </div>

            {/* Tabla de Gastos */}
            <div className="bg-zinc-900 border border-zinc-800 rounded-2xl shadow-xl overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full text-left border-collapse">
                  <thead>
                    <tr className="border-b border-zinc-800 bg-zinc-950/50 text-xs font-semibold text-zinc-400 uppercase tracking-wider">
                      <th className="py-4 px-6">Fecha</th>
                      <th className="py-4 px-6">Categoría</th>
                      <th className="py-4 px-6">Descripción</th>
                      <th className="py-4 px-6 text-right">Monto</th>
                      <th className="py-4 px-6 text-center">Acciones</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-zinc-800/60 text-sm">
                    {gastos.length === 0 ? (
                      <tr>
                        <td
                          colSpan={5}
                          className="py-12 text-center text-zinc-500"
                        >
                          Todavía no hay gastos registrados con los filtros
                          seleccionados.
                        </td>
                      </tr>
                    ) : (
                      gastos.map((g) => (
                        <tr
                          key={g.id}
                          className="hover:bg-zinc-800/30 transition-colors group"
                        >
                          <td className="py-4 px-6 text-zinc-300 tabular-nums whitespace-nowrap">
                            {g.fecha}
                          </td>
                          <td className="py-4 px-6">
                            <span className="bg-zinc-800 text-zinc-200 px-3 py-1 rounded-full text-xs font-medium border border-zinc-700/50">
                              {g.categoria}
                            </span>
                          </td>
                          <td className="py-4 px-6 text-zinc-400 max-w-xs truncate">
                            {g.descripcion || (
                              <span className="italic text-zinc-600">
                                Sin descripción
                              </span>
                            )}
                          </td>
                          <td className="py-4 px-6 text-right font-semibold text-emerald-400 tabular-nums whitespace-nowrap">
                            $
                            {g.monto.toLocaleString("es-AR", {
                              minimumFractionDigits: 2,
                            })}
                          </td>
                          <td className="py-4 px-6 text-center whitespace-nowrap">
                            <div className="flex items-center justify-center gap-2">
                              <button
                                onClick={() => setEditando(g)}
                                className="p-2 bg-zinc-800 hover:bg-zinc-700 text-zinc-300 rounded-xl transition-colors focus:outline-none focus:ring-2 focus:ring-emerald-500"
                                aria-label={`Editar gasto ${g.id}`}
                              >
                                <Edit3 className="w-4 h-4" />
                              </button>
                              <button
                                onClick={() => eliminar(g.id)}
                                className="p-2 bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 rounded-xl transition-colors focus:outline-none focus:ring-2 focus:ring-rose-500"
                                aria-label={`Eliminar gasto ${g.id}`}
                              >
                                <Trash2 className="w-4 h-4" />
                              </button>
                            </div>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
