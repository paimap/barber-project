// app/dashboard/page.tsx
import DashboardClient from "./DashboardClient";
import { fetchWithAuth } from "@/lib/fetchWithAuth";
import { redirect } from "next/navigation";

export default async function DashboardPage() {
  try {
    const res = await fetchWithAuth("http://localhost:8080/api/profile");
    if (!res.ok) redirect("/login");
    const user = await res.json();

    return <DashboardClient user={user} />;
  } catch (err) {
    redirect("/login");
  }
}
