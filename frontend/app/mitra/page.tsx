import { redirect } from "next/navigation";
import MitraClient from "./MitraClient";
import { fetchWithAuth } from "@/lib/fetchWithAuth";
import { MitraType, MitraSummaryType } from "./types";

export default async function MitraPage() {
  try {
    const baseUrl = "http://localhost:8080/api/admin/mitra";

    const [resMitra, resSummary] = await Promise.all([
      fetchWithAuth(baseUrl),
      fetchWithAuth(`${baseUrl}/summary`)
    ]);

    if (!resMitra.ok || !resSummary.ok) redirect("/login");

    const mitraJson: { mitra: MitraType[] } = await resMitra.json();
    const summaryJson: { data: MitraSummaryType } = await resSummary.json();

    // ✅ FORMAT DI SERVER
    const revenueFormatted = new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      maximumFractionDigits: 0,
    }).format(summaryJson.data.today_revenue);

    return (
      <MitraClient 
        mitraData={mitraJson.mitra || []} 
        summaryData={{
          ...summaryJson.data,
          today_revenue_formatted: revenueFormatted, // ✅ kirim string
        }}
      />
    );
  } catch (err) {
    console.error(err);
    redirect("/login");
  }
}
