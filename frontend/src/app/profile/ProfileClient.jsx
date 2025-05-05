'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function ProfileClient() {
  const [loading, setLoading] = useState(true);
  const [profileData, setProfileData] = useState(null);
  const router = useRouter();

  useEffect(() => {
    fetch("http://localhost:3200/profile", {
      method: "GET",
      credentials: "include",
    })
      .then(async (res) => {
        if (!res.ok) {
            console.log("ehhh" + res)
          router.push("/login");
          return;
        }
        console.log("meh" + res)
        const data = await res.json();
        setProfileData(data);
      })
      .catch(() => {
        console.log("lehhh" + res)
        router.push("/login");
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p className="text-center mt-10">Loading profile...</p>;

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      {/* Top Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 h-full border-primary border-t border-b">
        {/* Left Side */}
        <div className="flex flex-col justify-center p-5 lg:p-16 border-r border-primary">
          <div className="flex items-center justify-center w-full">
            <h1 className="text-5xl font-bold text-center">YourCompany</h1>
          </div>
        </div>

        {/* Right Side */}
        <div className="flex flex-col justify-start lg:pt-8">
          <nav className="text-2xl flex justify-center border-b border-primary pb-4 font-medium">
            <a href="/logout" className="hover:underline">Logout</a>
          </nav>

          <div className="mt-6 mb-6 pb-0">
            <h2 className="text-2xl font-light px-8 text-center">
              Welcome, {profileData?.Username}!
            </h2>
          </div>
        </div>
      </div>

      <div className="flex flex-col items-center mt-10">
        <p className="text-lg text-gray-700">
          Allergens: {profileData?.Allergens?.join(', ') || 'None'}
        </p>
        <p className="text-lg text-gray-700">
          Intolerances: {profileData?.Intolerances?.join(', ') || 'None'}
        </p>
      </div>
    </main>
  );
}
