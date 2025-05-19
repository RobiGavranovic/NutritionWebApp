import { motion } from "framer-motion";
import { Check, Loader2, X } from "lucide-react";
import { useState } from "react";

function DailyCalorieGoal({
  goal,
  setGoal,
  goalType,
  setGoalType,
  buttonState,
  onSetGoal,
}) {
  const handleSetGoal = async () => {
    onSetGoal();
  };

  return (
    <div className="bg-white p-6 rounded-xl shadow-lg border border-gray-200">
      <label className="block text-lg font-medium mb-2">
        Set Daily Calorie Goal
      </label>

      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700">
            Calories (kcal/day)
          </label>
          <input
            type="number"
            value={goal}
            onChange={(e) => setGoal(e.target.value)}
            className="w-full p-2 border border-gray-300 rounded"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700">
            Goal Type
          </label>
          <select
            value={goalType}
            onChange={(e) => setGoalType(e.target.value)}
            className="w-full p-2 border border-gray-300 rounded"
          >
            <option value="maintain">Maintain Weight</option>
            <option value="gain">Gain Weight</option>
            <option value="lose">Lose Weight</option>
          </select>
        </div>

        <motion.button
          onClick={handleSetGoal}
          disabled={buttonState === "loading"}
          className="w-full px-4 py-2 rounded bg-blue-600 text-white flex justify-center items-center gap-2 font-medium shadow-md hover:bg-blue-700 transition"
          whileTap={{ scale: 0.98 }}
        >
          {buttonState === "idle" && "Set Goal"}
          {buttonState === "loading" && (
            <>
              <Loader2 className="animate-spin" size={20} /> Saving
            </>
          )}
          {buttonState === "success" && (
            <>
              <Check size={20} /> Saved
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
  );
}

export default DailyCalorieGoal;
