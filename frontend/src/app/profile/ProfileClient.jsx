"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import UsernameEditor from "@/app/profile/UsernameEditor";
import AllergenIntoleranceEditor from "@/app/profile/AllergenIntolarenceEditor";

// #TODO all options
const allergenOptions = [
  "Dairy",
  "Peanuts",
  "Tree nuts",
  "Soy",
  "Eggs",
  "Wheat",
  "Fish",
  "Shellfish",
];

// #TODO all options
const intoleranceOptions = ["Dairy"];

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

  const router = useRouter();

  useEffect(() => {
    // Get Profile Data
    fetch("http://localhost:3200/profile", {
      method: "GET",
      credentials: "include",
    })
      .then(async (res) => {
        if (!res.ok) {
          router.push("/login");
          return;
        }
        // Update Data On FE W/O Refresh
        const data = await res.json();
        setUsername(data.Username);
        setUsernameInput(data.Username);
        setAllergens(data.Allergens || []);
        setIntolerances(data.Intolerances || []);
      })
      .catch(() => {
        router.push("/login");
      })
      .finally(() => setLoading(false));
  }, []);

  // Send PUT Request To URL To Change DB Data
  const updateField = async (url, payloadKey, data, setState, onSuccess) => {
    setState("loading");
    const tokenResponse = JSON.parse(localStorage.getItem("googleToken"));
    try {
      await fetch(url, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ tokenResponse, [payloadKey]: data }),
      });
      setState("success");
      if (onSuccess) onSuccess();
      setTimeout(() => setState("idle"), 1500);
    } catch {
      setState("error");
      setTimeout(() => setState("idle"), 1500);
    }
  };

  if (loading) return <p className="text-center mt-10">Loading profile...</p>;

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      <div className="grid grid-cols-1 lg:grid-cols-2 h-full border-primary border-t border-b">
        {/* Left Side */}
        <div className="flex flex-col justify-center p-5 lg:p-16 border-r border-primary">
          <div className="flex items-center justify-center w-full">
            <h1 className="text-5xl font-bold text-center">YourCompany</h1>
          </div>
        </div>

        {/* Right Side */}
        <div className="flex flex-col justify-start lg:pt-8">
          {/* Navbar */}
          <nav className="text-2xl flex justify-between items-center border-b border-primary pb-4 font-medium px-4">
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
              Logout
            </button>
            <button
              onClick={() => router.refresh()}
              className="hover:underline"
            >
              {username}
            </button>
          </nav>

          {/* Text Under navbar */}
          <div className="mt-6 mb-6 pb-0">
            <h2 className="text-2xl font-light px-8 text-center">
              Welcome, {username}!
            </h2>
          </div>
        </div>
      </div>

      <div className="px-8 mt-2 space-y-6">
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
    </main>
  );
}
