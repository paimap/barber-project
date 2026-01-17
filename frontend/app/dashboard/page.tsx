import { redirect } from "next/navigation";
import DashboardClient from "./DashboardClient";
import { fetchWithAuth } from "@/lib/fetchWithAuth";

async function getDashboardData() {
  const baseUrl = "http://localhost:8080/api/admin/dashboard";

  try {
    const [summaryRes, chartRes, distRes, leaderboardRes] = await Promise.all([
      fetchWithAuth(`${baseUrl}/summary`, { cache: "no-store" }),
      fetchWithAuth(`${baseUrl}/chart`, { cache: "no-store" }),
      fetchWithAuth(`${baseUrl}/distribution`, { cache: "no-store" }),
      fetchWithAuth(`${baseUrl}/leaderboard`, { cache: "no-store" }),
    ]);

    if (summaryRes.status === 401) redirect("/login");

    const summaryJson = await summaryRes.json().catch(() => ({}));
    const chartJson = await chartRes.json().catch(() => ({}));
    const distJson = await distRes.json().catch(() => ({}));
    const leaderboardJson = await leaderboardRes.json().catch(() => ({}));

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
      },
    };
  } catch (err) {
    console.error("Dashboard Server Error:", err);
    return null;
  }
}

export default async function DashboardPage() {
  const data = await getDashboardData();

  if (!data) {
    return (
      <div className="flex h-screen items-center justify-center">
        <p className="text-red-500">
          Gagal memuat dashboard. Pastikan backend berjalan.
        </p>
      </div>
    );
  }

  return <DashboardClient initialData={data} />;
}
