import { motion } from "framer-motion";
import { Check, Loader2, X } from "lucide-react";

function CalorieCalculator({
  age,
  setAge,
  height,
  setHeight,
  weight,
  setWeight,
  calories,
  buttonState,
  onCalculate,
}) {
  const handleCalculate = async () => {
    // Update Backend Info + Get Daily Calories Calculated
    onCalculate();
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
      <div className="grid grid-cols-3 gap-4">
        {/* Left: Inputs (1/3) */}
        <div className="space-y-4 col-span-1">
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Age
            </label>
            <input
              type="number"
              value={age}
              onChange={(e) => setAge(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Height (cm)
            </label>
            <input
              type="number"
              value={height}
              onChange={(e) => setHeight(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">
              Weight (kg)
            </label>
            <input
              type="number"
              value={weight}
              onChange={(e) => setWeight(e.target.value)}
              className="w-full p-2 border border-gray-300 rounded"
            />
          </div>
        </div>

        {/* Right: Output + Button (2/3) */}
        <div className="col-span-2 flex flex-col justify-between">
          <div className="text-xl font-medium text-center py-4 bg-gray-100 rounded mb-4 flex-grow flex items-center justify-center">
            {calories ? `${calories} kcal/day` : "0"}
          </div>

          <motion.button
            onClick={handleCalculate}
            disabled={buttonState === "loading"}
            className="w-full px-4 py-2 rounded bg-blue-600 text-white flex justify-center items-center gap-2 font-medium shadow-md hover:bg-blue-700 transition"
            whileTap={{ scale: 0.98 }}
          >
            {buttonState === "idle" && "Calculate"}
            {buttonState === "loading" && (
              <>
                <Loader2 className="animate-spin" size={20} /> Calculating
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
      </div>
    </div>
  );
}

export default CalorieCalculator;
