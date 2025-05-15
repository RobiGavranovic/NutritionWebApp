"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import UsernameEditor from "@/app/profile/UsernameEditor";
import AllergenIntoleranceEditor from "@/app/profile/AllergenIntolarenceEditor";
import CalorieCalculator from "@/app/profile/CalorieCalculator";
import DailyCalorieGoal from "@/app/profile/DailyCalorieGoal";
import TopBar from "@/app/components/TopBar";

const allergenOptions = [
  "Celery",
  "Gluten",
  "Crustaceans",
  "Eggs",
  "Fish",
  "Lupin",
  "Milk",
  "Molluscs",
  "Mustard",
  "Nuts",
  "Peanuts",
  "Sesame seeds",
  "Soya",
  "Sulphur dioxide",
];

const intoleranceOptions = [
  "Lactose",
  "Gluten",
  "Histamine",
  "Caffeine",
  "Alcohol",
  "Sulphites",
  "Salicylates",
  "Monosodium glutamate",
];

export default function ProfileClient() {
  const [loading, setLoading] = useState(true);
  const [username, setUsername] = useState("");
  const [usernameInput, setUsernameInput] = useState("");
  const [allergens, setAllergens] = useState([]);
  const [intolerances, setIntolerances] = useState([]);
  const [usernameButtonState, setUsernameButtonState] = useState("idle");
  const [allergenSelect, setAllergenSelect] = useState("");
  const [intoleranceSelect, setIntoleranceSelect] = useState("");
  const [allergenButtonState, setAllergenButtonState] = useState("idle");
  const [intoleranceButtonState, setIntoleranceButtonState] = useState("idle");

  // Calorie Calculator
  const [age, setAge] = useState("");
  const [height, setHeight] = useState("");
  const [weight, setWeight] = useState("");
  const [calories, setCalories] = useState(null);
  const [calorieButtonState, setCalorieButtonState] = useState("idle");

  // Daily Calorie Goal
  const [dailyGoal, setDailyGoal] = useState("");
  const [goalType, setGoalType] = useState("");
  const [goalButtonState, setGoalButtonState] = useState("idle");

  const router = useRouter();

  useEffect(() => {
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
        setUsernameInput(data.Username);
        setAllergens(data.Allergens || []);
        setIntolerances(data.Intolerances || []);
        setAge(data.Age || "");
        setHeight(data.Height || "");
        setWeight(data.Weight || "");
        setDailyGoal(data.DailyCalorieGoal || "");
      })
      .catch(() => {
        router.push("/login");
      })
      .finally(() => setLoading(false));
  }, []);

  const updateField = async (url, payloadKey, data, setState, onSuccess) => {
    setState("loading");
    try {
      await fetch(url, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ [payloadKey]: data }),
      });
      setState("success");
      if (onSuccess) onSuccess();
      setTimeout(() => setState("idle"), 1500);
    } catch {
      setState("error");
      setTimeout(() => setState("idle"), 1500);
    }
  };

  const updatePersonalInfo = async () => {
    if (!age || parseInt(age) <= 0 || !height || parseInt(height) <= 0) {
      alert("Please enter valid age and height values before calculating.");
      return;
    }

    setCalorieButtonState("loading");
    try {
      var res = await fetch(
        "http://localhost:3200/profile/updatePersonalInfo",
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({
            age: parseInt(age),
            height: parseInt(height),
            weight: parseInt(weight),
          }),
        }
      );
      const data = await res.json();
      setCalories(data.dailyCalorieGoal);
      setCalorieButtonState("success");
      setTimeout(() => setCalorieButtonState("idle"), 1500);
    } catch {
      setCalorieButtonState("error");
      setTimeout(() => setCalorieButtonState("idle"), 1500);
    }
  };

  const updateDailyCalorieGoal = async () => {
    setGoalButtonState("loading");
    try {
      await fetch("http://localhost:3200/profile/updateDailyCalorieGoal", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({
          dailyCalorieGoal: parseInt(dailyGoal),
          dailyGoalType: goalType,
        }),
      });
      setGoalButtonState("success");
      setTimeout(() => setGoalButtonState("idle"), 1500);
    } catch {
      setGoalButtonState("error");
      setTimeout(() => setGoalButtonState("idle"), 1500);
    }
  };

  if (loading) return <p className="text-center mt-10">Loading profile...</p>;

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      <TopBar
        navbarOptions={["RECIPES", "CONSUMPTION", "PROFILE", "LOGOUT"]}
        subtitle={"Welcome, " + username + "!"}
      />

      <div className="grid grid-cols-1 lg:grid-cols-2 h-full border-primary border-t border-b">
        <div className="space-y-6 mt-6 ml-3 mr-3">
          <UsernameEditor
            label="Edit Username"
            value={usernameInput}
            setValue={setUsernameInput}
            onApply={() =>
              updateField(
                "http://localhost:3200/profile/updateUsername",
                "username",
                usernameInput,
                setUsernameButtonState,
                () => setUsername(usernameInput)
              )
            }
            buttonState={usernameButtonState}
          />

          <AllergenIntoleranceEditor
            label="Edit Allergens"
            options={allergenOptions}
            selectedItems={allergens}
            setSelectedItems={setAllergens}
            selectValue={allergenSelect}
            setSelectValue={setAllergenSelect}
            buttonState={allergenButtonState}
            setButtonState={setAllergenButtonState}
            endpoint="http://localhost:3200/profile/updateAllergens"
            payloadKey="allergens"
          />

          <AllergenIntoleranceEditor
            label="Edit Intolerances"
            options={intoleranceOptions}
            selectedItems={intolerances}
            setSelectedItems={setIntolerances}
            selectValue={intoleranceSelect}
            setSelectValue={setIntoleranceSelect}
            buttonState={intoleranceButtonState}
            setButtonState={setIntoleranceButtonState}
            endpoint="http://localhost:3200/profile/updateIntolerances"
            payloadKey="intolerances"
          />
        </div>

        <div className="space-y-6 mt-6 ml-3 mr-3">
          <CalorieCalculator
            age={age}
            setAge={setAge}
            height={height}
            setHeight={setHeight}
            weight={weight}
            setWeight={setWeight}
            calories={calories}
            buttonState={calorieButtonState}
            onCalculate={updatePersonalInfo}
          />

          <DailyCalorieGoal
            goal={dailyGoal}
            setGoal={setDailyGoal}
            goalType={goalType}
            setGoalType={setGoalType}
            buttonState={goalButtonState}
            onSetGoal={updateDailyCalorieGoal}
          />
        </div>
      </div>
    </main>
  );
}
