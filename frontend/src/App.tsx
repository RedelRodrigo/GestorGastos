import { useCallback, useEffect, useState } from "react";
import {
  borrarGasto,
  listarGastos,
  obtenerResumen,
  obtenerUsuarioActual,
} from "./api";
import Header from "./components/Header";
import CategorySummary from "./components/CategorySummary";
import FilterBar from "./components/FilterBar";
import ExpenseTable from "./components/ExpenseTable";
import GastoForm from "./components/GastoForm";
import LoginView from "./components/LoginView";
import type { Gasto, ResumenItem } from "./types";

export default function App() {
  const [sesionActiva, setSesionActiva] = useState<boolean>(
    () => !!obtenerUsuarioActual(),
  );
  const [gastos, setGastos] = useState<Gasto[]>([]);
  const [resumen, setResumen] = useState<ResumenItem[]>([]);
  const [total, setTotal] = useState(0);
  const [editando, setEditando] = useState<Gasto | null>(null);
  const [error, setError] = useState("");
  const [filtroCategoria, setFiltroCategoria] = useState("");

  // Función dedicada exclusivamente a recargas manuales tras mutaciones (POST, PUT, DELETE)
  const refrescarDatos = useCallback(async () => {
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

  // Efecto de sincronización: vive autocontenido con cleanup para evitar race conditions
  useEffect(() => {
    if (!sesionActiva) return;

    let ignorado = false;

    async function sincronizar() {
      try {
        const query = filtroCategoria
          ? `?categoria=${encodeURIComponent(filtroCategoria)}`
          : "";
        const [lista, res] = await Promise.all([
          listarGastos(query),
          obtenerResumen(),
        ]);
        if (!ignorado) {
          setGastos(lista);
          setResumen(res.data);
          setTotal(res.total);
          setError("");
        }
      } catch (e) {
        if (!ignorado) {
          setError((e as Error).message);
        }
      }
    }

    sincronizar();

    return () => {
      ignorado = true;
    };
  }, [sesionActiva, filtroCategoria]);

  async function eliminar(id: number) {
    if (!confirm("¿Estás seguro de eliminar este gasto?")) return;
    try {
      await borrarGasto(id);
      await refrescarDatos();
    } catch (e) {
      setError((e as Error).message);
    }
  }

  function handleLogout() {
    setSesionActiva(false);
    setGastos([]);
    setResumen([]);
    setTotal(0);
    setEditando(null);
    setFiltroCategoria("");
    setError("");
  }

  if (!sesionActiva) {
    return <LoginView onLoginSuccess={() => setSesionActiva(true)} />;
  }

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 p-4 md:p-8">
      <main className="max-w-7xl mx-auto space-y-8">
        <Header total={total} onLogout={handleLogout} />

        {error && (
          <div
            role="alert"
            className="bg-rose-500/10 border border-rose-500/20 text-rose-400 p-4 rounded-xl text-sm"
          >
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          <div className="lg:col-span-4 space-y-6">
            <GastoForm
              key={editando?.id ?? "nuevo"}
              inicial={editando}
              onGuardado={() => {
                setEditando(null);
                refrescarDatos();
              }}
              onCancelar={() => setEditando(null)}
            />
            <CategorySummary resumen={resumen} total={total} />
          </div>

          <div className="lg:col-span-8 space-y-6">
            <FilterBar
              filtroCategoria={filtroCategoria}
              setFiltroCategoria={setFiltroCategoria}
              categoriasDisponibles={resumen}
            />
            <ExpenseTable
              gastos={gastos}
              onEditar={(g) => setEditando(g)}
              onEliminar={eliminar}
            />
          </div>
        </div>
      </main>
    </div>
  );
}
