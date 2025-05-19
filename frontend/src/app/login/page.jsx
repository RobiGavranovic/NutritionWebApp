"use client";

import { GoogleOAuthProvider, useGoogleLogin } from "@react-oauth/google";
import { motion, AnimatePresence } from "framer-motion";
import { useState } from "react";
import { useRouter } from "next/navigation";
import TopBar from "@/app/components/TopBar";

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
      <TopBar
        navbarOptions={["REGISTER"]}
        subtitle="Login using google account"
      />

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
