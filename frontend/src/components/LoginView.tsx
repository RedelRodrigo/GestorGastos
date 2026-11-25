import React, { useEffect, useRef, useState } from "react";
import { loginApi, registroApi, guardarSesionUsuario } from "../api";
import { Wallet, Mail, User, Lock, ArrowRight } from "lucide-react";
import FormAlert from "./FormAlert";

declare global {
  interface Window {
    turnstile?: {
      render: (
        container: HTMLElement,
        options: {
          sitekey: string;
          callback?: (token: string) => void;
          "expired-callback"?: () => void;
          "error-callback"?: () => void;
        },
      ) => string;
      reset: (widgetId?: string) => void;
      remove: (widgetId?: string) => void;
    };
  }
}

interface LoginViewProps {
  onLoginSuccess: () => void;
}

const SITE_KEY = import.meta.env.VITE_TURNSTILE_SITE_KEY ?? "";

export default function LoginView({ onLoginSuccess }: LoginViewProps) {
  const [modo, setModo] = useState<"login" | "registro">("login");
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [captchaToken, setCaptchaToken] = useState("");
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);

  const captchaContainerRef = useRef<HTMLDivElement | null>(null);
  const widgetIdRef = useRef<string | undefined>(undefined);

  useEffect(() => {
    let cancelado = false;

    function montarWidget() {
      if (cancelado || !captchaContainerRef.current || !window.turnstile)
        return;
      if (!SITE_KEY) {
        setError("Falta configurar VITE_TURNSTILE_SITE_KEY en el entorno.");
        return;
      }

      if (widgetIdRef.current) {
        window.turnstile.remove(widgetIdRef.current);
        widgetIdRef.current = undefined;
      }

      widgetIdRef.current = window.turnstile.render(
        captchaContainerRef.current,
        {
          sitekey: SITE_KEY,
          callback: (token: string) => {
            setCaptchaToken(token);
            setError("");
          },
          "expired-callback": () => {
            setCaptchaToken("");
          },
          "error-callback": () => {
            setCaptchaToken("");
            setError("Error al cargar la verificación de seguridad.");
          },
        },
      );
    }

    if (window.turnstile) {
      montarWidget();
    } else {
      const timer = window.setInterval(() => {
        if (window.turnstile) {
          window.clearInterval(timer);
          montarWidget();
        }
      }, 150);

      return () => {
        cancelado = true;
        window.clearInterval(timer);
      };
    }

    return () => {
      cancelado = true;
      if (widgetIdRef.current && window.turnstile) {
        window.turnstile.remove(widgetIdRef.current);
        widgetIdRef.current = undefined;
      }
    };
  }, [modo]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (!captchaToken) {
      setError("Completá la verificación de seguridad antes de continuar.");
      return;
    }

    setCargando(true);

    try {
      const usuario =
        modo === "login"
          ? await loginApi(email.trim(), password, captchaToken)
          : await registroApi(
              email.trim(),
              username.trim(),
              password,
              captchaToken,
            );

      guardarSesionUsuario(usuario);
      onLoginSuccess();
    } catch (err) {
      setError((err as Error).message);
      setCaptchaToken("");
      if (widgetIdRef.current && window.turnstile) {
        window.turnstile.reset(widgetIdRef.current);
      }
    } finally {
      setCargando(false);
    }
  }

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 flex items-center justify-center p-4">
      <div className="max-w-md w-full bg-zinc-900 border border-zinc-800 rounded-3xl p-8 shadow-2xl space-y-6">
        <div className="text-center space-y-2">
          <div className="inline-flex bg-emerald-500/10 p-3 rounded-2xl border border-emerald-500/20 text-emerald-400 mb-2">
            <Wallet className="w-8 h-8" />
          </div>
          <h1 className="text-2xl font-bold tracking-tight">
            Gestor de Gastos
          </h1>
          <p className="text-sm text-zinc-400">
            {modo === "login"
              ? "Iniciá sesión para administrar tus finanzas"
              : "Creá tu cuenta en segundos"}
          </p>
        </div>

        <div className="flex bg-zinc-950 p-1 rounded-xl border border-zinc-800">
          <button
            type="button"
            onClick={() => {
              setModo("login");
              setCaptchaToken("");
              setError("");
            }}
            className={`flex-1 py-2 text-xs font-semibold rounded-lg transition-all ${
              modo === "login"
                ? "bg-zinc-800 text-zinc-100 shadow"
                : "text-zinc-400 hover:text-zinc-200"
            }`}
          >
            Iniciar Sesión
          </button>
          <button
            type="button"
            onClick={() => {
              setModo("registro");
              setCaptchaToken("");
              setError("");
            }}
            className={`flex-1 py-2 text-xs font-semibold rounded-lg transition-all ${
              modo === "registro"
                ? "bg-zinc-800 text-zinc-100 shadow"
                : "text-zinc-400 hover:text-zinc-200"
            }`}
          >
            Registrarse
          </button>
        </div>

        {error && <FormAlert message={error} />}

        <form onSubmit={handleSubmit} className="space-y-4">
          {modo === "registro" && (
            <div>
              <label className="block text-xs font-medium text-zinc-400 mb-1">
                Nombre de Usuario
              </label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-zinc-500">
                  <User className="w-4 h-4" />
                </span>
                <input
                  type="text"
                  required
                  minLength={2}
                  maxLength={40}
                  placeholder="Ej. Totti Redel"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  className="w-full bg-zinc-950 border border-zinc-800 rounded-xl pl-10 pr-4 py-2.5 text-zinc-100 focus:outline-none focus:ring-2 focus:ring-emerald-500 text-sm"
                />
              </div>
            </div>
          )}

          <div>
            <label className="block text-xs font-medium text-zinc-400 mb-1">
              Correo Electrónico
            </label>
            <div className="relative">
              <span className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-zinc-500">
                <Mail className="w-4 h-4" />
              </span>
              <input
                type="email"
                required
                placeholder="tu@correo.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-zinc-950 border border-zinc-800 rounded-xl pl-10 pr-4 py-2.5 text-zinc-100 focus:outline-none focus:ring-2 focus:ring-emerald-500 text-sm"
              />
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-zinc-400 mb-1">
              Contraseña{" "}
              {modo === "registro" && (
                <span className="text-zinc-500">(mínimo 6 caracteres)</span>
              )}
            </label>
            <div className="relative">
              <span className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-zinc-500">
                <Lock className="w-4 h-4" />
              </span>
              <input
                type="password"
                required
                minLength={modo === "registro" ? 6 : 1}
                maxLength={72}
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-zinc-950 border border-zinc-800 rounded-xl pl-10 pr-4 py-2.5 text-zinc-100 focus:outline-none focus:ring-2 focus:ring-emerald-500 text-sm"
              />
            </div>
          </div>

          <div className="flex justify-center pt-1">
            <div ref={captchaContainerRef} />
          </div>

          <button
            type="submit"
            disabled={cargando}
            className="w-full bg-emerald-600 hover:bg-emerald-500 text-zinc-950 font-semibold py-3 px-4 rounded-xl transition-all flex items-center justify-center gap-2 shadow-lg shadow-emerald-900/20 disabled:opacity-50 text-sm mt-2"
          >
            <span>
              {cargando
                ? "Procesando..."
                : modo === "login"
                  ? "Ingresar al Sistema"
                  : "Crear Cuenta"}
            </span>
            <ArrowRight className="w-4 h-4" />
          </button>
        </form>
      </div>
    </div>
  );
}
