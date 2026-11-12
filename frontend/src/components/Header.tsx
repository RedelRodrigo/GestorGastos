import { Wallet } from "lucide-react";
import TotalCard from "./TotalCard";
import UserBadge from "./UserBadge";
import { obtenerUsuarioActual } from "../api";

interface HeaderProps {
  total: number;
  onLogout: () => void;
}

export default function Header({ total, onLogout }: HeaderProps) {
  const usuario = obtenerUsuarioActual();

  return (
    <header className="flex flex-col lg:flex-row lg:items-center justify-between gap-4 border-b border-zinc-800/80 pb-6">
      <div className="flex items-center gap-3">
        <div className="bg-emerald-500/10 p-3 rounded-2xl border border-emerald-500/20 text-emerald-400">
          <Wallet className="w-8 h-8" />
        </div>
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-zinc-100">
            Gestor de Gastos
          </h1>
          <p className="text-sm text-zinc-400">
            Control financiero personal de alto rendimiento
          </p>
        </div>
      </div>

      <div className="flex flex-wrap items-center gap-3">
        <UserBadge usuario={usuario} onLogout={onLogout} />
        <TotalCard total={total} />
      </div>
    </header>
  );
}
