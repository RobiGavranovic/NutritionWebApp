"use client";

import { GoogleOAuthProvider, useGoogleLogin } from "@react-oauth/google";
import UserProfileSetup from "../components/UserProfileSetUp";
import { motion, AnimatePresence } from "framer-motion"; // 👈 import motion + AnimatePresence
import { useState } from "react";
import TopBar from "@/app/components/TopBar";

export default function Register() {
  return (
    <GoogleOAuthProvider clientId={process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID}>
      <MainContent />
    </GoogleOAuthProvider>
  );
}

async function onFinishedUserSetUp(
  profileData,
  setProfileFinished,
  setErrorMessage
) {
  setTimeout(() => {
    setProfileFinished(true);
  }, 400);

  const res = await fetch("http://localhost:3200/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      tokenResponse: profileData.tokenResponse,
      username: profileData.username,
      allergens: profileData.allergens,
      intolerances: profileData.intolerances,
      gender: profileData.gender,
    }),
  });

  if (res.ok) {
  } else {
    const errorData = await res.json();
    // Instruction : Go Login
    setTimeout(() => {
      if (errorData.error === "User already exists, please login instead") {
        setErrorMessage("Account already exists. Please log in!");
      }
    }, 300);
  }
}

function MainContent() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [googleTokenResponse, setGoogleTokenResponse] = useState(null);
  const [profileFinished, setProfileFinished] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  const login = useGoogleLogin({
    onSuccess: async (TokenResponse) => {
      setGoogleTokenResponse(TokenResponse);
      setIsAuthenticated(true);
    },
    onError: () => {},
  });

  return (
    <main className="min-h-screen bg-gray-100 text-gray-900">
      {/* Top Layout */}
      <TopBar
        navbarOptions={["LOGIN"]}
        subtitle="Register using a google account"
      />

      {/* AnimatePresence zone */}
      <div className="flex flex-col items-center justify-start space-y-8">
        <AnimatePresence mode="wait">
          {!isAuthenticated && (
            <motion.button
              key="sign-up-button"
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
                Sign up with Google
              </span>
            </motion.button>
          )}

          {isAuthenticated && !profileFinished && (
            <motion.div
              key="user-profile-setup"
              initial={{ opacity: 0, y: 50 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: 50 }}
              transition={{ duration: 0.5 }}
              className="w-full flex justify-center"
            >
              <UserProfileSetup
                googleTokenResponse={googleTokenResponse}
                onFinish={(finalData) =>
                  onFinishedUserSetUp(
                    finalData,
                    setProfileFinished,
                    setErrorMessage
                  )
                }
              />
            </motion.div>
          )}

          {errorMessage && (
            <motion.div
              key="error-message"
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
