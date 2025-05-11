"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import IngredientConsumption from "@/app/consumption/IngredientConsumption";

export default function ConsumptionPage() {
  const [username, setUsername] = useState("");
  const [ingredients, setIngredients] = useState([]);
  const router = useRouter();

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
      .catch(() => router.push("/login"));

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
        const names = data.map((item) => item.name); // Map to string array
        setIngredients(names);
      })
      .catch((err) => console.error("Failed to load ingredients", err));
  }, []);

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      {/* Header */}
      <div className="grid grid-cols-1 lg:grid-cols-2 h-full border-primary border-t border-b">
        <div className="flex flex-col justify-center p-5 lg:p-16 border-r border-primary">
          <h1 className="text-5xl font-bold text-center">YourCompany</h1>
        </div>

        <div className="flex flex-col justify-start lg:pt-8">
          <nav className="text-2xl flex justify-around items-center border-b border-primary pb-4 font-medium px-4">
            <button
              onClick={() => router.push("/recipes")}
              className="hover:underline"
            >
              RECIPES
            </button>
            <button
              onClick={() => router.refresh()}
              className="hover:underline"
            >
              {username}
            </button>
            <button
              onClick={async () => {
                await fetch("http://localhost:3200/logout", {
                  method: "POST",
                  credentials: "include",
                });
                router.push("/login");
              }}
              className="hover:underline"
            >
              LOGOUT
            </button>
          </nav>

          <div className="mt-6 mb-6 pb-0">
            <h2 className="text-2xl font-light px-8 text-center">
              Welcome, {username}!
            </h2>
          </div>
        </div>
      </div>

      {/* Page Content */}
      <div className="px-6 py-8">
        <IngredientConsumption ingredientOptions={ingredients} />
      </div>
    </main>
  );
}
