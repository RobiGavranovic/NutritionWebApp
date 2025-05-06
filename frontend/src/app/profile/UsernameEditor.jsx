import { motion } from "framer-motion";
import { Check, Loader2, X } from "lucide-react";

function UsernameEditor({ label, value, setValue, onApply, buttonState }) {
  return (
    <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
      <label className="block text-lg font-medium mb-2">{label}</label>
      <input
        type="text"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        className="w-full p-2 border border-gray-300 rounded mb-4 focus:outline-none focus:ring-2 focus:ring-blue-400"
      />
      <motion.button
        onClick={onApply}
        disabled={buttonState === "loading"}
        className="w-full px-4 py-2 rounded bg-blue-600 text-white flex justify-center items-center gap-2 font-medium shadow-md hover:bg-blue-700 transition"
        whileTap={{ scale: 0.98 }}
      >
        {buttonState === "idle" && "Apply"}
        {buttonState === "loading" && (
          <>
            <Loader2 className="animate-spin" size={20} /> Processing
          </>
        )}
        {buttonState === "success" && (
          <>
            <Check size={20} /> Done
          </>
        )}
        {buttonState === "error" && (
          <>
            <X size={20} /> Error
          </>
        )}
      </motion.button>
    </div>
  );
}

export default UsernameEditor;
