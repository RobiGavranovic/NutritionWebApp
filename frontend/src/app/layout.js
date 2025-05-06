import "./globals.css";

export const metadata = {
  title: "Nutrition Web App",
  description: "Manage your meals, calories, and analytics",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
