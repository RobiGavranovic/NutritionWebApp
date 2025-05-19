"use client";

import { useState } from "react";

export default function RecipesSearch({ onResults }) {
  const [nameQuery, setNameQuery] = useState("");
  const [originQuery, setOriginQuery] = useState("");

  const handleSearchByName = async () => {
    if (!nameQuery.trim()) return;
    try {
      const res = await fetch(`http://localhost:3200/recipes/search/name/${nameQuery}`);
      const data = await res.json();
      onResults(data.meals || []);
    } catch (err) {
      onResults([]);
    }
  };

  const handleSearchByOrigin = async () => {
    if (!originQuery.trim()) return;
    try {
      const res = await fetch(`http://localhost:3200/recipes/search/origin/${originQuery}`);
      const data = await res.json();
      onResults(data.meals || []);
    } catch (err) {
      onResults([]);
    }
  };

  return (
    <div className="flex flex-col md:flex-row justify-center items-center gap-6 my-8 px-4">
      {/* Search by Name */}
      <div className="flex flex-col w-full max-w-xs">
        <label className="text-sm font-medium text-gray-700 mb-1 text-center">
          Search by name
        </label>
        <input
          type="text"
          placeholder="e.g. Pasta"
          value={nameQuery}
          onChange={(e) => setNameQuery(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSearchByName()}
          className="border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
        />
        <button
          onClick={handleSearchByName}
          className="mt-2 bg-blue-500 text-white py-1 px-3 rounded hover:bg-blue-600 text-sm"
        >
          Search
        </button>
      </div>

      {/* Search by Origin */}
      <div className="flex flex-col w-full max-w-xs">
        <label className="text-sm font-medium text-gray-700 mb-1 text-center">
          Search by origin
        </label>
        <input
          type="text"
          placeholder="e.g. Italian"
          value={originQuery}
          onChange={(e) => setOriginQuery(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSearchByOrigin()}
          className="border border-gray-300 rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
        />
        <button
          onClick={handleSearchByOrigin}
          className="mt-2 bg-blue-500 text-white py-1 px-3 rounded hover:bg-blue-600 text-sm"
        >
          Search
        </button>
      </div>
    </div>
  );
}
