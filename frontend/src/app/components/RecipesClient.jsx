"use client";

import { useState } from "react";
import RecipesSearch from "./RecipesSearch";
import MealCardsSection from "./MealCardsSection";

export default function RecipesClient() {
  const [searchResults, setSearchResults] = useState(null);

  return (
    <>
      <RecipesSearch onResults={setSearchResults} />
      {searchResults && searchResults.length === 0 ? (
        <p className="text-center text-gray-500 mt-4">No meals found.</p>
      ) : (
        <MealCardsSection searchResults={searchResults} />
      )}
    </>
  );
}
