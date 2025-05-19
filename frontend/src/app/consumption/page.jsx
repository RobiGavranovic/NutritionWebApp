"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import IngredientConsumption from "@/app/consumption/IngredientConsumption";
import TodaysConsumptionHistory from "@/app/consumption/TodaysConsumptionHistory";
import ConsumptionAnalitics from "@/app/consumption/ConsumptionAnalitics";
import TopBar from "@/app/components/TopBar";

export default function ConsumptionPage() {
  const [username, setUsername] = useState("");
  const [ingredients, setIngredients] = useState([]);
  const router = useRouter();
  const [todaysConsumption, setTodaysConsumption] = useState([]);

  const handleNewConsumption = (newItem) => {
    setTodaysConsumption((prev) => [...prev, newItem]);
  };

  const handleDeleteConsumption = (id) => {
    setTodaysConsumption((prev) => prev.filter((item) => item.ID !== id));
  };

  useEffect(() => {
    // Get User Data
    fetch("http://localhost:3200/profile", {
      method: "GET",
      credentials: "include",
    })
      .then(async (res) => {
        if (!res.ok) {
          router.push("/login");
          return;
        }
        const data = await res.json();
        setUsername(data.Username);
      })
      .catch((err) => router.push("/login"));

    // Get All Ingredients
    const tokenResponse = JSON.parse(localStorage.getItem("googleToken"));
    fetch("http://localhost:3200/getAllIngredients", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ tokenResponse }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (!Array.isArray(data)) {
          return;
        }
        const names = data.map((item) => item.name); // Map to string array
        setIngredients(names);
      })
      .catch((err) => console.error("Failed to load ingredients", err));

    // Fetch Today's Consumption
    fetch("http://localhost:3200/consumption/today", {
      method: "POST",
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        if (Array.isArray(data.consumptions)) {
          setTodaysConsumption(data.consumptions);
        }
      });
  }, []);

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      {/* Header */}
      <TopBar
        navbarOptions={["RECIPES", "PROFILE", "LOGOUT"]}
        subtitle="Search for any recipes by name or origin"
      />

      {/* Page Content */}
      <div className="px-6 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="flex justify-center lg:justify-end">
            <IngredientConsumption
              ingredientOptions={ingredients}
              onConsume={handleNewConsumption}
            />
          </div>
          <div className="flex justify-center lg:justify-start">
            <TodaysConsumptionHistory
              items={todaysConsumption}
              onDelete={handleDeleteConsumption}
            />
          </div>
        </div>
        <div className="mt-6">
          <ConsumptionAnalitics />
        </div>
      </div>
    </main>
  );
}
