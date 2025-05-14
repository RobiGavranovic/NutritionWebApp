"use client";

import { useEffect, useRef, useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import MealCard from "./MealCard";

export default function MealCardsSection({ searchResults }) {
  const [meals, setMeals] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedMeal, setSelectedMeal] = useState(null);
  const observer = useRef(null);

  const fetchMeals = async (amount = 10) => {
    if (loading || searchResults) return; // If Search Result: Dont Load

    setLoading(true);
    try {
      const res = await fetch(`http://localhost:3200/getRandomMeals/${amount}`);
      const data = await res.json();

      if (!data || !Array.isArray(data.Meals)) {
        return;
      }

      const seenIDs = new Set(meals.map((m) => m.IDMeal));
      const newMeals = data.Meals.filter((meal) => !seenIDs.has(meal.IDMeal));

      // Retry once if no unique meals found
      if (newMeals.length === 0) {
        const retryRes = await fetch(
          `http://localhost:3200/getRandomMeals/${amount}`
        );
        const retryData = await retryRes.json();
        if (retryData && Array.isArray(retryData.Meals)) {
          const retryMeals = retryData.Meals.filter(
            (meal) => !seenIDs.has(meal.IDMeal)
          );
          setMeals((prev) => [...prev, ...retryMeals]);
        }
      } else {
        setMeals((prev) => [...prev, ...newMeals]);
      }
    } catch (err) {
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!searchResults) {
      fetchMeals(20); // If No Search Has Been Made
    } else {
      setMeals(searchResults); // Show Only Search Results
    }
  }, [searchResults]);

  // Observe The Last Card To Detect When To Load More
  const lastMealRef = useCallback(
    (node) => {
      if (loading) return;
      if (observer.current) observer.current.disconnect();
      observer.current = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting) {
          fetchMeals(10);
        }
      });
      if (node) observer.current.observe(node);
    },
    [loading, meals]
  );

  return (
    <div className="relative">
      <div className="flex flex-wrap justify-center gap-8 mt-8">
        {meals.map((meal, index) => {
          const key = meal.IDMeal || `${meal.StrMeal}-${index}`; // fallback key
          return (
            <motion.div
              key={key}
              layoutId={key}
              onClick={() => setSelectedMeal(meal)}
              ref={index === meals.length - 1 ? lastMealRef : null}
            >
              <MealCard meal={meal} />
            </motion.div>
          );
        })}
      </div>

      {loading && (
        <div className="h-12 flex justify-center items-center">
          <p>Loading more meals...</p>
        </div>
      )}

      <AnimatePresence>
        {selectedMeal && (
          <motion.div
            layoutId={selectedMeal.IDMeal}
            className="fixed top-0 left-0 w-screen h-screen flex items-center justify-center bg-white/70 z-50"
            onClick={() => setSelectedMeal(null)}
          >
            <div
              className="p-4 bg-white shadow-lg max-w-3xl w-full h-[90vh] overflow-y-auto rounded-lg"
              style={{ WebkitOverflowScrolling: "touch" }}
            >
              <img
                src={selectedMeal.StrMealThumb}
                alt={selectedMeal.StrMeal}
                className="w-full h-auto rounded max-w-md mx-auto"
              />
              <h2 className="text-2xl font-bold mt-4">
                {selectedMeal.StrMeal}
              </h2>
              <p className="mt-2 text-gray-600">{selectedMeal.StrArea}</p>
              <p className="mt-2 font-semibold text-black">Ingredients:</p>
              <div className="grid grid-cols-2 gap-x-8 gap-y-2 mb-4">
                {Array.from({ length: 20 }, (_, index) => {
                  const measure = selectedMeal[`StrMeasure${index + 1}`];
                  const ingredient = selectedMeal[`StrIngredient${index + 1}`];
                  if (ingredient && ingredient.trim() !== "") {
                    return (
                      <p key={index} className="text-gray-600">
                        {measure} {ingredient}
                      </p>
                    );
                  }
                  return null;
                })}
              </div>
              <p className="mt-2 text-gray-600">
                {selectedMeal.StrInstructions}
              </p>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
