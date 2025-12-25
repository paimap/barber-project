import DashboardClient from "./DashboardClient";
import { fetchWithAuth } from "@/lib/fetchWithAuth";
import { redirect } from "next/navigation";

export default async function DashboardPage() {
  return <DashboardClient />;
  // try {
  //   // const res = await fetchWithAuth("http://localhost:8080/api/profile");
  //   // if (!res.ok) redirect("/login");
  //   // const user = await res.json();

  //   return <DashboardClient />;
  // } catch (err) {
  //   redirect("/login");
  // }
}
