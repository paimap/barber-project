import { redirect } from "next/navigation";
import { cookies } from "next/headers"; // WAJIB DIIMPOR
import DashboardClient from "./DashboardClient";
import { fetchWithAuth } from "@/lib/fetchWithAuth";

async function getProfile() {
  const res = await fetchWithAuth("http://localhost:8080/api/me")

  if (!res.ok) return null;
  return res.json();
}

async function getDashboardData(baseUrl: string) {
  try {
    const [summaryRes, chartRes, distRes, leaderboardRes] = await Promise.all([
      fetchWithAuth(`${baseUrl}/summary`, { cache: "no-store" }),
      fetchWithAuth(`${baseUrl}/chart`, { cache: "no-store" }),
      fetchWithAuth(`${baseUrl}/distribution`, { cache: "no-store" }),
      fetchWithAuth(`${baseUrl}/leaderboard`, { cache: "no-store" }),
    ]);

    // Error handling yang lebih rapi
    if (summaryRes.status === 401) return null;

    const [summaryJson, chartJson, distJson, leaderboardJson] = await Promise.all([
      summaryRes.json().catch(() => ({})),
      chartRes.json().catch(() => ({})),
      distRes.json().catch(() => ({})),
      leaderboardRes.json().catch(() => ({})),
    ]);

    return {
      summary: {
        total_revenue: Number(summaryJson.data?.total_revenue ?? 0),
        profit_20: Number(summaryJson.data?.profit_20 ?? 0),
        service_revenue: Number(summaryJson.data?.service_revenue ?? 0),
        product_revenue: Number(summaryJson.data?.product_revenue ?? 0),
      },
      chart: chartJson.data ?? [],
      distribution: {
        products: distJson.data?.products ?? [],
        services: distJson.data?.services ?? [],
      },
      leaderboard: {
        top_outlets: leaderboardJson.data?.top_outlets ?? [],
        top_mitras: leaderboardJson.data?.top_mitras ?? [],
        top_barbers: leaderboardJson.data?.top_barbers ?? [],
      },
    };
  } catch (err) {
    console.error("Dashboard Server Error:", err);
    return null;
  }
}

export default async function DashboardPage() {
  const profileData = await getProfile();

  if (!profileData || !profileData.data) {
    redirect("/login");
  }

  const id = profileData.data.id;
  const roleRaw = profileData.data.role;
  const rolePath = roleRaw === "SUPERADMIN" ? "admin" : roleRaw.toLowerCase();

  const baseUrl = rolePath === "admin" ? `http://localhost:8080/api/${rolePath}/dashboard` : `http://localhost:8080/api/${rolePath}/dashboard/${id}`;
  const data = await getDashboardData(baseUrl);

  if (!data) {
    redirect("/login");
  }

  return <DashboardClient initialData={data} />;
}