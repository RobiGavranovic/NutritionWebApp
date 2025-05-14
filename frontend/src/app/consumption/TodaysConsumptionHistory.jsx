"use client";

import { X } from "lucide-react";
import { useState, useEffect } from "react";

export default function TodaysConsumptionHistory({ items }) {
  const [consumedItems, setConsumedItems] = useState(items);

  useEffect(() => {
    if (Array.isArray(items)) {
      setConsumedItems(items);
    }
  }, [items]);

  const handleDelete = async (id) => {
    try {
      const res = await fetch(
        `http://localhost:3200/consumption/consume/${id}`,
        {
          method: "DELETE",
          credentials: "include",
        }
      );

      if (!res.ok) {
        console.error("Failed to delete item with ID:", id);
        return;
      }

      // Remove the item from state on success
      setConsumedItems((prev) => prev.filter((item) => item.ID !== id));
    } catch (err) {
      console.error("Error deleting item:", err);
    }
  };

  if (!Array.isArray(consumedItems)) {
    return (
      <div className="w-full max-w-2xl bg-white rounded-xl shadow-md p-6 border border-gray-200">
        <h2 className="text-xl font-semibold mb-4 text-center">
          Today's Consumption BLAAAA
        </h2>

        <div className="grid grid-cols-4 font-semibold border-b pb-2 mb-3 text-gray-700">
          <div>Name</div>
          <div>Weight (g)</div>
          <div>Calories</div>
          <div></div>
        </div>
      </div>
    );
  }

  if (consumedItems.length === 0) {
    return (
      <div className="w-full max-w-2xl bg-white rounded-xl shadow-md p-6 border border-gray-200 text-center">
        <h2 className="text-xl font-semibold mb-4">Today's Consumption</h2>
        <p className="text-gray-500">No ingredients consumed yet today.</p>
      </div>
    );
  }

  return (
    <div className="w-full max-w-2xl bg-white rounded-xl shadow-md p-6 border border-gray-200">
      <h2 className="text-xl font-semibold mb-4 text-center">
        Today's Consumption
      </h2>

      <div className="grid grid-cols-4 font-semibold border-b pb-2 mb-3 text-gray-700">
        <div>Name</div>
        <div>Weight (g)</div>
        <div>Calories</div>
        <div></div>
      </div>

      {consumedItems.map((item, index) => (
        <div
          key={item.id || `${item.Ingredient}-${index}`}
          className="grid grid-cols-4 items-center py-2 border-b last:border-none text-sm text-gray-800"
        >
          <div>{item.Ingredient}</div>
          <div>{item.Weight}</div>
          <div>{Math.round(parseInt(item.Calories))}</div>
          <div>
            <button
              className="text-red-500 hover:text-red-700"
              onClick={() => handleDelete(item.ID)}
            >
              <X size={18} />
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}
