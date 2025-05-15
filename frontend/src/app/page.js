import RecipesClient from "@/app/components/RecipesClient";
import TopBar from "@/app/components/TopBar";

export default function Home() {
  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      {/* Header & Navbar (unchanged) */}
      <TopBar
        navbarOptions={["REGISTER", "LOGIN"]}
        subtitle="The all-in-one app for cooking, nutrition, and health."
      />

      {/* Recipes UI (Client-Sided) */}
      <RecipesClient />
    </main>
  );
}
