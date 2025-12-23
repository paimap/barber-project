import "./globals.css";
import Sidebar from '../components/navbar/index'

export const metadata = {
  title: "My App",
  description: "Next.js App Router Example",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <div className="layout-container">
          <Sidebar />
          <main className="content">
            {children}
          </main>
        </div>
      </body>
    </html>
  );
}
