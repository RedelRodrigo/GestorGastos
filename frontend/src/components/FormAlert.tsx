import { AlertCircle } from "lucide-react";

interface FormAlertProps {
  message: string;
}

export default function FormAlert({ message }: FormAlertProps) {
  return (
    <div
      role="alert"
      className="flex items-center gap-2 bg-rose-500/10 border border-rose-500/20 text-rose-400 p-3 rounded-xl text-sm"
    >
      <AlertCircle className="w-4 h-4 shrink-0" />
      <span>{message}</span>
    </div>
  );
}
