"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Loader2, Check, X } from "lucide-react";

export default function IngredientConsumption({ ingredientOptions }) {
  const [searchText, setSearchText] = useState("");
  const [filteredOptions, setFilteredOptions] = useState([]);
  const [selectedIngredient, setSelectedIngredient] = useState("");
  const [imageVisible, setImageVisible] = useState(true);
  const [weight, setWeight] = useState("");
  const [status, setStatus] = useState("idle");
  const [consumedCalories, setConsumedCalories] = useState(null);

  const handleInputChange = (e) => {
    const input = e.target.value;
    setSearchText(input);
    if (input.trim() === "") {
      setFilteredOptions([]);
      return;
    }
    const matches = ingredientOptions.filter((item) =>
      item.toLowerCase().startsWith(input.toLowerCase())
    );
    setFilteredOptions(matches.slice(0, 10));
  };

  const handleSelect = (name) => {
    setSelectedIngredient(name);
    setSearchText(name);
    setFilteredOptions([]);
    setImageVisible(true);
  };

  const handleConsume = async () => {
    setImageVisible(false);
    setSelectedIngredient("");
    setStatus("loading");

    const tokenResponse = JSON.parse(localStorage.getItem("googleToken"));
    const payload = {
      tokenResponse,
      ingredient: selectedIngredient,
      weight: parseFloat(weight),
    };

    try {
      const response = await fetch("http://localhost:3200/consume", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(payload),
      });

      const data = await response.json();

      if (!response.ok) {
        setStatus("error");
        return;
      }

      setConsumedCalories(data.caloriesConsumed);
      setStatus("success");
    } catch (err) {
      setStatus("error");
    } finally {
      // Reset Inputs
      setTimeout(() => setStatus("idle"), 1500);
      setSearchText("");
      setSelectedIngredient("");
      setFilteredOptions([]);
      setWeight("");
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-md p-6 border border-gray-200 max-w-2xl mx-auto mt-6">
      <div className="mb-4 relative">
        <input
          type="text"
          placeholder="Search for an ingredient..."
          value={searchText}
          onChange={handleInputChange}
          className="w-full p-3 border border-gray-300 rounded text-base"
        />
        {filteredOptions.length > 0 && (
          <ul className="absolute z-10 bg-white border border-gray-300 w-full rounded mt-1 max-h-48 overflow-auto shadow">
            {filteredOptions.map((option) => (
              <li
                key={option}
                onClick={() => handleSelect(option)}
                className="px-4 py-2 hover:bg-blue-100 cursor-pointer"
              >
                {option}
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="grid grid-cols-3 gap-4 items-start">
        {/* Image */}
        <div className="col-span-2">
          <div className="w-full h-48 bg-gray-100 rounded flex items-center justify-center overflow-hidden">
            <AnimatePresence>
              {imageVisible &&
                selectedIngredient &&
                (() => {
                  const imageUrl =
                    selectedIngredient === "Hazelnuts"
                      ? "https://upload.wikimedia.org/wikipedia/commons/e/e4/Hazelnuts_%28Corylus_avellana%29_-_whole_with_kernels.jpg"
                      : `https://www.themealdb.com/images/ingredients/${selectedIngredient.replaceAll(
                          " ",
                          "_"
                        )}.png`;

                  return (
                    <motion.img
                      key="ingredient-image"
                      src={imageUrl}
                      alt={selectedIngredient}
                      initial={{ opacity: 1, x: 0 }}
                      exit={{ opacity: 0, x: -100 }}
                      transition={{ duration: 0.5 }}
                      className="max-h-full max-w-full object-contain"
                    />
                  );
                })()}
            </AnimatePresence>
          </div>
        </div>

        <div className="flex flex-col justify-between h-full gap-4">
          {/* Label + Input */}
          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium text-gray-700">
              Weight (g)
            </label>
            <input
              type="number"
              placeholder="0"
              value={weight}
              onChange={(e) => setWeight(e.target.value)}
              className="p-2 border border-gray-300 rounded text-sm"
            />
          </div>

          {/* Calories Display */}
          <AnimatePresence>
            {consumedCalories !== null && (
              <motion.div
                key="calorie-info"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: 10 }}
                className="text-sm text-center bg-blue-100 text-blue-800 px-3 py-2 rounded"
              >
                Consumed: {Math.round(consumedCalories)} kcal
              </motion.div>
            )}
          </AnimatePresence>

          {/* Consume Button */}
          <motion.button
            onClick={handleConsume}
            disabled={
              !selectedIngredient ||
              isNaN(parseFloat(weight)) ||
              parseFloat(weight) <= 0 ||
              status === "loading"
            }
            whileTap={{ scale: 0.97 }}
            className={`
              text-white rounded py-2 px-4 text-sm shadow flex items-center justify-center gap-2
              ${
                status === "idle" || status === "loading"
                  ? "bg-green-600 hover:bg-green-700"
                  : ""
              }
              ${status === "success" ? "bg-green-500" : ""}
              ${status === "error" ? "bg-red-500" : ""}
              disabled:opacity-50`}
          >
            {status === "idle" && "CONSUME"}
            {status === "loading" && (
              <>
                <Loader2 className="animate-spin" size={16} /> Processing
              </>
            )}
            {status === "success" && (
              <>
                <Check size={16} /> Consumed
              </>
            )}
            {status === "error" && (
              <>
                <X size={16} /> Error
              </>
            )}
          </motion.button>
        </div>
      </div>
    </div>
  );
}
