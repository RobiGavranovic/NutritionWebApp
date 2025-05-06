import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { jwtVerify } from "jose";
import ProfileClient from "./ProfileClient";

export default async function ProfilePage() {
  const cookieStore = await cookies();
  const token = cookieStore?.get("session_token")?.value;

  if (!token) {
    redirect("/login");
  }

  try {
    await jwtVerify(
      token,
      new TextEncoder().encode(process.env.JWT_PRIVATE_KEY)
    );
  } catch (e) {
    redirect("/login");
  }

  return <ProfileClient />;
}
