"use client";

import { GoogleOAuthProvider, useGoogleLogin } from "@react-oauth/google";
import { motion, AnimatePresence } from "framer-motion";
import { useState } from "react";
import { useRouter } from "next/navigation";

export default function Login() {
  return (
    <GoogleOAuthProvider clientId={process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID}>
      <MainContent />
    </GoogleOAuthProvider>
  );
}

function MainContent() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const router = useRouter();

  const login = useGoogleLogin({
    onSuccess: async (tokenResponse) => {
      try {
        const res = await fetch("http://localhost:3200/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(tokenResponse),
          credentials: "include",
        });

        if (res.ok) {
          setIsAuthenticated(true);
          router.push("/profile");
          localStorage.setItem("googleToken", JSON.stringify(tokenResponse));
        } else {
          const errorData = await res.json();
          setErrorMessage(errorData.error || "Login failed. Try again.");
        }
      } catch (err) {
        setErrorMessage("An error occurred during login.");
      }
    },
    onError: () => setErrorMessage("Google login failed."),
  });

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      {/* Top Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 h-full border-primary border-t border-b">
        <div className="flex flex-col justify-center p-5 sm:p-5 md:p-5 lg:p-16 border-r border-primary">
          <div className="flex items-center justify-center w-full">
            <h1 className="text-5xl sm:text-6xl md:text-7xl lg:text-7xl xl:text-10xl font-bold leading-none text-center">
              YourCompany
            </h1>
          </div>
        </div>

        <div className="sm:mt-0 md:mt-1 md:ml-5 md:mr-5 flex flex-col justify-start lg:pt-8">
          <nav className="text-2xl flex justify-center border-b border-primary lg:pb-4 font-medium">
            <a href="/register" className="hover:underline">
              Register
            </a>
          </nav>

          <div className="mt-6 mb-6 border-primary pb-0">
            <h2 className="text-2xl font-light px-8 text-center">
              Log in with your Google account to continue your journey.
            </h2>
          </div>
        </div>
      </div>

      {/* AnimatePresence Zone */}
      <div className="flex flex-col items-center justify-start space-y-8">
        <AnimatePresence mode="wait">
          {!isAuthenticated && (
            <motion.button
              key="login-button"
              onClick={() => login()}
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.8 }}
              transition={{ duration: 0.4 }}
              className="flex items-center gap-3 border border-gray-300 py-2 px-6 rounded-full shadow-sm hover:shadow-md transition bg-white mt-6"
            >
              <img
                src="/google_logo.svg"
                alt="Google Icon"
                className="w-6 h-6"
              />
              <span className="text-gray-700 text-lg font-medium">
                Log in with Google
              </span>
            </motion.button>
          )}

          {errorMessage && (
            <motion.div
              key="login-error"
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.8 }}
              transition={{ duration: 0.4 }}
              className="bg-red-100 border border-red-400 text-red-700 px-6 py-4 rounded relative mt-10 text-center"
            >
              {errorMessage}
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </main>
  );
}
