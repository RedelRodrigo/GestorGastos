import { User, LogOut } from "lucide-react";
import type { Usuario } from "../types";

interface UserBadgeProps {
  usuario: Usuario | null;
  onLogout: () => void;
}

export default function UserBadge({ usuario, onLogout }: UserBadgeProps) {
  if (!usuario) return null;

  return (
    <div className="flex items-center gap-3 bg-zinc-900 border border-zinc-800 rounded-2xl px-4 py-2">
      <div className="bg-emerald-500/10 p-2 rounded-xl text-emerald-400">
        <User className="w-4 h-4" />
      </div>
      <div className="text-left">
        <span className="text-[10px] text-zinc-400 block uppercase tracking-wider">
          Conectado como
        </span>
        <span className="text-sm font-semibold text-zinc-200">
          {usuario.username}{" "}
          <span className="text-xs text-zinc-500 font-normal">
            ({usuario.email})
          </span>
        </span>
      </div>
      <button
        onClick={onLogout}
        title="Cerrar sesión"
        className="ml-2 p-2 bg-zinc-800 hover:bg-rose-500/20 hover:text-rose-400 text-zinc-400 rounded-xl transition-colors"
      >
        <LogOut className="w-4 h-4" />
      </button>
    </div>
  );
}
