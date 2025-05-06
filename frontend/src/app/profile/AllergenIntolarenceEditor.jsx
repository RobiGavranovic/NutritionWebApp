import { motion } from "framer-motion";
import { Check, Loader2, Plus, X } from "lucide-react";

function AllergenIntoleranceEditor({
  label,
  options,
  selectedItems,
  setSelectedItems,
  selectValue,
  setSelectValue,
  buttonState,
  setButtonState,
  endpoint,
  payloadKey,
}) {
  const handleAdd = () => {
    if (selectValue && !selectedItems.includes(selectValue)) {
      setSelectedItems([...selectedItems, selectValue]);
      setSelectValue("");
    }
  };

  const handleRemove = (item) => {
    setSelectedItems(selectedItems.filter((i) => i !== item));
  };

  const handleApply = async () => {
    const tokenResponse = JSON.parse(localStorage.getItem("googleToken"));
    setButtonState("loading");
    try {
      await fetch(endpoint, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ tokenResponse, [payloadKey]: selectedItems }),
      });
      setButtonState("success");
      setTimeout(() => setButtonState("idle"), 1500);
    } catch {
      setButtonState("error");
      setTimeout(() => setButtonState("idle"), 1500);
    }
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
      <label className="block text-lg font-medium mb-2">{label}</label>

      <div className="flex gap-2 mb-4">
        <select
          className="border border-gray-300 p-2 rounded flex-grow"
          value={selectValue}
          onChange={(e) => setSelectValue(e.target.value)}
        >
          <option value="">-- Choose --</option>
          {options.map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
        <button
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 rounded"
          onClick={handleAdd}
        >
          <Plus />
        </button>
      </div>

      <div className="flex flex-wrap gap-2 mb-4">
        {selectedItems.map((item) => (
          <div
            key={item}
            className="flex items-center bg-gray-200 rounded-full px-3 py-1 text-sm"
          >
            {item}
            <button
              className="ml-2 text-red-500"
              onClick={() => handleRemove(item)}
            >
              <X />
            </button>
          </div>
        ))}
      </div>

      <motion.button
        onClick={handleApply}
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

export default AllergenIntoleranceEditor;
