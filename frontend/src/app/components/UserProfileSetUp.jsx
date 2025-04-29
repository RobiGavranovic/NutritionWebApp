'use client';

import { Tabs, TabsList, TabsTrigger, TabsContent } from '@radix-ui/react-tabs';
import { motion, AnimatePresence } from 'framer-motion';
import { useState } from 'react';
import { ArrowLeft, ArrowRight, X, Plus } from 'lucide-react';

const steps = ['username', 'allergens', 'intolerances'];

export default function UserProfileSetup({ googleTokenResponse, onFinish }) {
    const [currentIndex, setCurrentIndex] = useState(0);
    const [username, setUsername] = useState('');
    const [allergens, setAllergens] = useState('');
    const [intolerances, setIntolerances] = useState('');

    const currentTab = steps[currentIndex];
    const [errorUsername, setErrorUsername] = useState(false);

    const handleNext = () => {
        if (currentTab === 'username' && username.trim() === '') {
            setErrorUsername(true);
            return; // block moving forward
        }

        setErrorUsername(false);

        if (currentIndex < steps.length - 1) {
            setCurrentIndex(currentIndex + 1);
        }
    };

    const handlePrev = () => {
        if (currentIndex > 0) {
            setCurrentIndex(currentIndex - 1);
        }
    };

    const allergenOptions = [
        "Peanuts",
        "Tree Nuts",
        "Dairy",
        "Eggs",
        "Fish",
        "Shellfish",
        "Soy",
        "Wheat",
        "Gluten",
        "Sesame",
    ];
    const intoleranceOptions = [
        "Lactose",
        "Gluten",
        "Fructose",
        "Histamine",
        "Sorbitol",
        "Salicylates",
    ];  

    const [addedAllergens, setAddedAllergens] = useState([]);
    const [addedIntolerances, setAddedIntolerances] = useState([]);

    const [isExiting, setIsExiting] = useState(false);

    return (
        <motion.div
          initial={{ opacity: 1, y: 0 }}
          animate={isExiting ? { opacity: 0, y: 50 } : { opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
          className="flex flex-col items-center justify-start w-full bg-gray-100 mt-6"
        >
            <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">

                {/* Tabs Indicator */}
                <Tabs defaultValue="username" value={currentTab}>
                    <TabsList className="flex justify-around mb-8 relative border-b border-gray-300">
                        {steps.map((tab) => (
                        <TabsTrigger
                          key={tab}
                          value={tab}
                          className={`pb-2 relative pointer-events-none ${
                            currentTab === tab ? 'text-blue-600' : 'text-gray-400'
                          }`}
                        >
                            {tab.charAt(0).toUpperCase() + tab.slice(1)}
                            {currentTab === tab && (
                            <motion.div
                              layoutId="underline"
                              className="absolute left-0 right-0 bottom-0 h-0.5 bg-blue-600 rounded-full"
                            />
                            )}
                        </TabsTrigger>
                        ))}
                    </TabsList>

                    {/* Content Section */}
                    <TabsContent value="username" className="flex flex-col gap-4">
                        <h2 className="text-2xl font-bold text-center">Insert your Username</h2>
                        <input
                          type="text"
                          placeholder="Username"
                          className="border border-gray-300 p-2 rounded w-full"
                          value={username}
                          onChange={(e) => setUsername(e.target.value)}
                        />

                        {/* Empty Username Warning */}
                        {errorUsername && (
                            <p className="text-red-500 text-center text-sm">
                                Please insert a username before proceeding.
                            </p>
                        )}
                    </TabsContent>

                    <TabsContent value="allergens" className="flex flex-col gap-4">
                    <h2 className="text-2xl font-bold text-center">Select your Allergens</h2>

                    <div className="flex gap-2">

                        <select
                          className="border border-gray-300 p-2 rounded flex-grow"
                          value={allergens}
                          onChange={(e) => setAllergens(e.target.value)}
                        >
                            <option value="">-- Choose an allergen --</option>
                            {allergenOptions.map((option) => (
                                <option key={option} value={option}>
                                    {option}
                                </option>
                            ))}
                        </select>

                        <button
                          className="bg-blue-500 hover:bg-blue-600 text-white px-4 rounded"
                          onClick={() => {
                            if (allergens && !addedAllergens.includes(allergens)) {
                                setAddedAllergens([...addedAllergens, allergens]);
                                setAllergens(''); // Reset select after adding
                            }
                          }}
                        >
                            <Plus/>
                        </button>
                    </div>

                    {/* List of Added Allergens */}
                    <div className="flex flex-wrap gap-2 mt-4">
                        {addedAllergens.map((allergen) => (
                            <div
                              key={allergen}
                              className="flex items-center bg-gray-200 rounded-full px-3 py-1 text-sm"
                            >
                                {allergen}
                                <button
                                  className="ml-2 text-red-500"
                                  onClick={() =>
                                    setAddedAllergens(addedAllergens.filter((item) => item !== allergen))
                                }
                                >
                                <X/>
                                </button>
                            </div>
                        ))}
                    </div>
                    </TabsContent>


                    <TabsContent value="intolerances" className="flex flex-col gap-4">
                        <h2 className="text-2xl font-bold text-center">Select your Intolerances</h2>

                        <div className="flex gap-2">
                            <select
                              className="border border-gray-300 p-2 rounded flex-grow"
                              value={intolerances}
                              onChange={(e) => setIntolerances(e.target.value)}
                            >
                            <option value="">-- Choose an intolerance --</option>

                            {intoleranceOptions.map((option) => (
                                <option key={option} value={option}>
                                    {option}
                                </option>
                            ))}

                            </select>

                            <button
                              className="bg-blue-500 hover:bg-blue-600 text-white px-4 rounded"
                              onClick={() => {
                                if (intolerances && !addedIntolerances.includes(intolerances)) {
                                    setAddedIntolerances([...addedIntolerances, intolerances]);
                                    setIntolerances('');
                                }
                            }}
                            >
                            <Plus/>
                            </button>
                        </div>

                        {/* List of Added Intolerances */}
                        <div className="flex flex-wrap gap-2 mt-4">
                            {addedIntolerances.map((intolerance) => (
                                <div
                                key={intolerance}
                                className="flex items-center bg-gray-200 rounded-full px-3 py-1 text-sm"
                                >
                                    {intolerance}
                                    <button
                                      onClick={() =>
                                        setAddedIntolerances(addedIntolerances.filter((item) => item !== intolerance))
                                      }
                                      className="ml-2 text-red-500"
                                    >
                                        <X/>
                                    </button>
                                </div>
                            ))}
                        </div>

                        <button
                          onClick={() => {
                            setIsExiting(true);
                            const finalData = {
                                tokenResponse: googleTokenResponse,
                                username: username,
                                allergens: addedAllergens,
                                intolerances: addedIntolerances,
                            };
                            onFinish(finalData);
                          }}
                          className="bg-green-500 hover:bg-green-600 text-white py-2 px-6 rounded-full mt-4"
                        >
                        Finish
                        </button>
                    </TabsContent>

                </Tabs>

                {/* Arrows */}
                <div className="flex justify-between items-center mt-8 relative w-full">

                    {/* Left Arrow container */}
                    <div className="w-12 flex justify-start">
                        <AnimatePresence>
                            {currentIndex > 0 && (
                                <motion.button
                                  key="left-arrow"
                                  initial={{ opacity: 0, scale: 0.8 }}
                                  animate={{ opacity: 1, scale: 1 }}
                                  exit={{ opacity: 0, scale: 0.8 }}
                                  transition={{ duration: 0.3 }}
                                  onClick={handlePrev}
                                  className="p-2 bg-gray-200 rounded-full"
                                >
                                    <ArrowLeft size={24}/>
                                </motion.button>
                            )}
                        </AnimatePresence>
                    </div>

                    {/* Right Arrow container */}
                    <div className="w-12 flex justify-end">
                        <AnimatePresence>
                            {currentIndex < steps.length - 1 && (
                                <motion.button
                                  key="right-arrow"
                                  initial={{ opacity: 0, scale: 0.8 }}
                                  animate={{ opacity: 1, scale: 1 }}
                                  exit={{ opacity: 0, scale: 0.8 }}
                                  transition={{ duration: 0.3 }}
                                  onClick={handleNext}
                                  className="p-2 bg-gray-200 rounded-full"
                                >
                                    <ArrowRight size={24}/>
                                </motion.button>
                            )}
                        </AnimatePresence>
                    </div>

                </div>

            </div>
        </motion.div>
    );
}
