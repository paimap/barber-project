import "./globals.css";
import Sidebar from '../components/navbar/index';
import { cookies } from "next/headers";
import { AuthProvider } from "@/context/AuthContext";

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const cookieStore = await cookies();
  const token = cookieStore.get("access_token")?.value;
  const isLoggedIn = Boolean(token);

  return (
    <html lang="en">
      <body>
        <AuthProvider>
          <div className="layout-container">
            {isLoggedIn && <Sidebar />}
            <main className={`content ${isLoggedIn ? 'content-with-sidebar' : 'content-full'}`}>
              {children}
            </main>
          </div>
        </AuthProvider>
      </body>
    </html>
  );
}