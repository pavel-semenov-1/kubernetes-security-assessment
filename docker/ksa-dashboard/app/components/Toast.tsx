import { useEffect } from "react";

interface ToastProps {
  message: string;
  onClose: () => void;
  duration?: number;
}

// Toast for notifications
export default function Toast({ message, onClose, duration = 3000 }: ToastProps) {
  useEffect(() => {
    const timer = setTimeout(onClose, duration);
    return () => clearTimeout(timer);
  }, [onClose, duration]);

  return (
    <div className="fixed bottom-5 right-5 bg-gray-900 text-white px-5 py-3 rounded-lg shadow-lg transition-opacity opacity-100 animate-fade-in-out">
      {message}
    </div>
  );
}
