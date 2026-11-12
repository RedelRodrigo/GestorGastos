import React, { useState } from "react";
import { crearGasto, editarGasto } from "../api";
import type { Gasto } from "../types";
import { PlusCircle, Save, X } from "lucide-react";
import FormAlert from "./FormAlert";

interface GastoFormProps {
  inicial: Gasto | null;
  onGuardado: () => void;
  onCancelar: () => void;
}

const CATEGORIAS_SUGERIDAS = [
  "Comida",
  "Transporte",
  "Servicios",
  "Ocio",
  "Salud",
  "Otros",
];

export default function GastoForm({
  inicial,
  onGuardado,
  onCancelar,
}: GastoFormProps) {
  const [monto, setMonto] = useState(inicial ? String(inicial.monto) : "");
  const [categoria, setCategoria] = useState(inicial?.categoria ?? "");
  const [fecha, setFecha] = useState(
    inicial?.fecha ?? new Date().toISOString().slice(0, 10),
  );
  const [descripcion, setDescripcion] = useState(inicial?.descripcion ?? "");
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);

  async function enviar(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setCargando(true);

    try {
      const datos = { monto: Number(monto), categoria, fecha, descripcion };

      if (inicial) {
        await editarGasto(inicial.id, datos);
      } else {
        await crearGasto(datos);
        // Limpiamos el formulario si es un nuevo registro
        setMonto("");
        setCategoria("");
        setDescripcion("");
        setFecha(new Date().toISOString().slice(0, 10));
      }

      // Notificamos al componente padre para que recargue los datos y actualice la tabla/resumen
      onGuardado();
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setCargando(false);
    }
  }

  return (
    <form
      onSubmit={enviar}
      className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 shadow-xl space-y-4"
    >
      <div className="flex items-center justify-between border-b border-zinc-800 pb-3">
        <h2 className="text-lg font-semibold text-zinc-100">
          {inicial ? "Editar Gasto" : "Nuevo Registro"}
        </h2>
        {inicial && (
          <button
            type="button"
            onClick={onCancelar}
            className="text-zinc-400 hover:text-zinc-200 transition-colors p-1 rounded-lg focus:outline-none focus:ring-2 focus:ring-emerald-500"
            aria-label="Cancelar edición"
          >
            <X className="w-5 h-5" />
          </button>
        )}
      </div>

      {error && <FormAlert message={error} />}

      <div>
        <label
          htmlFor="monto"
          className="block text-xs font-medium text-zinc-400 mb-1"
        >
          Monto ($)
        </label>
        <input
          id="monto"
          name="monto"
          type="number"
          step="0.01"
          placeholder="0.00"
          value={monto}
          onChange={(e) => setMonto(e.target.value)}
          required
          className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-2.5 text-zinc-100 focus:outline-none focus:ring-2 focus:ring-emerald-500 tabular-nums"
        />
      </div>

      <div>
        <label
          htmlFor="categoria"
          className="block text-xs font-medium text-zinc-400 mb-1"
        >
          Categoría
        </label>
        <input
          id="categoria"
          name="categoria"
          list="sugerencias-cat"
          placeholder="Ej. Comida, Transporte..."
          value={categoria}
          onChange={(e) => setCategoria(e.target.value)}
          required
          className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-2.5 text-zinc-100 focus:outline-none focus:ring-2 focus:ring-emerald-500"
        />
        <datalist id="sugerencias-cat">
          {CATEGORIAS_SUGERIDAS?.map((cat) => (
            <option key={cat} value={cat} />
          ))}
        </datalist>
      </div>

      <div>
        <label
          htmlFor="fecha"
          className="block text-xs font-medium text-zinc-400 mb-1"
        >
          Fecha
        </label>
        <input
          id="fecha"
          name="fecha"
          type="date"
          value={fecha}
          onChange={(e) => setFecha(e.target.value)}
          required
          className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-2.5 text-zinc-100 focus:outline-none focus:ring-2 focus:ring-emerald-500"
        />
      </div>

      <div>
        <label
          htmlFor="descripcion"
          className="block text-xs font-medium text-zinc-400 mb-1"
        >
          Descripción (Opcional)
        </label>
        <input
          id="descripcion"
          name="descripcion"
          type="text"
          placeholder="Detalle adicional..."
          value={descripcion}
          onChange={(e) => setDescripcion(e.target.value)}
          className="w-full bg-zinc-950 border border-zinc-800 rounded-xl px-4 py-2.5 text-zinc-100 focus:outline-none focus:ring-2 focus:ring-emerald-500"
        />
      </div>

      <div className="flex gap-3 pt-2">
        <button
          type="submit"
          disabled={cargando}
          className="flex-1 bg-emerald-600 hover:bg-emerald-500 text-zinc-950 font-semibold py-2.5 px-4 rounded-xl transition-all flex items-center justify-center gap-2 shadow-lg shadow-emerald-900/20 disabled:opacity-50"
        >
          {inicial ? (
            <Save className="w-4 h-4" />
          ) : (
            <PlusCircle className="w-4 h-4" />
          )}
          <span>{inicial ? "Guardar Cambios" : "Agregar Gasto"}</span>
        </button>

        {inicial && (
          <button
            type="button"
            onClick={onCancelar}
            className="bg-zinc-800 hover:bg-zinc-700 text-zinc-300 font-medium py-2.5 px-4 rounded-xl transition-colors"
          >
            Cancelar
          </button>
        )}
      </div>
    </form>
  );
}
