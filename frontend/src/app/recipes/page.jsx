import RecipesClient from "@/app/recipes/RecipesClient";
import TopBar from "@/app/components/TopBar";

export default function Home() {

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      {/* Header & Navbar */}
      <TopBar
        navbarOptions={["CONSUMPTION", "PROFILE", "LOGOUT"]}
        subtitle="Search for any recipes by name or origin"
      />

      {/* Recipes UI (Client-Sided) */}
      <RecipesClient />
    </main>
  );
}
