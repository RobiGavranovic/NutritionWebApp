'use client';

import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import MealCard from './MealCard';

{/* We Pass Meals From page.js */}
export default function MealCardsSection({meals}) {
    const [selectedMeal, setSelectedMeal] = useState(null);

    return (
        <div className="relative">
        {/* Meals Cards */}
        <div className="flex flex-wrap justify-center gap-8 mt-8">
          {meals.map((meal) => (
            <motion.div 
              key={meal.IDMeal} 
              layoutId={meal.IDMeal} 
              onClick={() => setSelectedMeal(meal)}
            >
              <MealCard meal={meal}/>
            </motion.div>
          ))}
        </div>
  
        {/* Expanded Meal View */}
        <AnimatePresence>
          {selectedMeal && (
            <motion.div
              layoutId={selectedMeal.IDMeal}
              className="fixed top-0 left-0 w-screen h-screen flex items-center justify-center bg-white/70 z-50"
              onClick={() => setSelectedMeal(null)}
            >
              <div className="p-8 bg-white shadow-lg max-w-3xl w-full md:h-[600px] overflow-y-auto rounded-lg">
                <img src={selectedMeal.StrMealThumb} alt={selectedMeal.StrMeal} className="w-full h-auto rounded max-w-md mx-auto" />
                
                <h2 className="text-2xl font-bold mt-4">{selectedMeal.StrMeal}</h2>
                  <p className="mt-2 text-gray-600">{selectedMeal.StrArea}</p>
                  <p className="mt-2 font-semibold text-black">Ingredients:</p>

                  <div className="grid grid-cols-2 gap-x-8 gap-y-2 mb-4 border-primary">

                    {Array.from({ length: 20 }, (_, index) => {
                      const measure = selectedMeal[`StrMeasure${index + 1}`];
                      const ingredient = selectedMeal[`StrIngredient${index + 1}`];
                      
                      {/* The Empty Ones Dont Get Shown */}
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
                  <p className="mt-2 text-gray-600">{selectedMeal.StrInstructions}</p>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    );
}
