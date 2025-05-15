"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function TopBar({ navbarOptions = [], subtitle = "" }) {
  const router = useRouter();
  const [username, setUsername] = useState(null);

  useEffect(() => {
    if (navbarOptions.includes("PROFILE")) {
      fetch("http://localhost:3200/profile", {
        method: "GET",
        credentials: "include",
      })
        .then((res) => (res.ok ? res.json() : null))
        .then((data) => {
          if (data && data.Username) setUsername(data.Username);
        })
        .catch(() => {});
    }
  }, [navbarOptions]);

  const handleClick = async (option) => {
    switch (option) {
      case "REGISTER":
        router.push("/register");
        break;
      case "LOGIN":
        router.push("/login");
        break;
      case "RECIPES":
        router.push("/recipes");
        break;
      case "CONSUMPTION":
        router.push("/consumption");
        break;
      case "PROFILE":
        router.push("/profile");
        break;
      case "LOGOUT":
        await fetch("http://localhost:3200/logout", {
          method: "POST",
          credentials: "include",
        });
        router.push("/login");
        break;
      default:
        break;
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 h-full border-primary border-t border-b">
      {/* Company Name */}
      <div className="flex flex-col justify-center p-5 lg:p-12 border-r border-primary">
        <h1 className="text-5xl font-bold text-center">YourCompany</h1>
      </div>

      {/* Navbar and Subtitle */}
      <div className="flex flex-col justify-start lg:pt-8">
        <nav className="sm:text-xl lg:text-2xl flex justify-around items-center border-b border-primary pb-4 font-medium px-4">
          {navbarOptions.map((option, i) => (
            <button
              key={i}
              onClick={() => handleClick(option)}
              className="hover:underline"
            >
              {option === "PROFILE" && username ? username : option}
            </button>
          ))}
        </nav>

        {subtitle && (
          <div className="mt-6 mb-6 pb-0">
            <h2 className="text-2xl font-light px-8 text-center">{subtitle}</h2>
          </div>
        )}
      </div>
    </div>
  );
}
