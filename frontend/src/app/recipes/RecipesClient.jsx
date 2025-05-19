"use client";

import { useState, useEffect } from "react";
import RecipesSearch from "@/app/components/RecipesSearch";
import MealCardsSection from "@/app/components/MealCardsSection";

export default function RecipesClient() {
  const [searchResults, setSearchResults] = useState(null);
  const [userAllergens, setUserAllergens] = useState([]);
  const [userIntolerances, setUserIntolerances] = useState([]);

  useEffect(() => {
    fetch("http://localhost:3200/profile", {
      method: "GET",
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        setUserAllergens(data.Allergens || []);
        setUserIntolerances(data.Intolerances || []);
      });
  }, []);

  return (
    <>
      <RecipesSearch onResults={setSearchResults} />
      {searchResults && searchResults.length === 0 ? (
        <p className="text-center text-gray-500 mt-4">No meals found.</p>
      ) : (
        <MealCardsSection
          searchResults={searchResults}
          userAllergens={userAllergens}
          userIntolerances={userIntolerances}
        />
      )}
    </>
  );
}
