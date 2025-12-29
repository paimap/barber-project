import { redirect } from "next/navigation";
import MitraClient from "./MitraClient";
import { fetchWithAuth } from "@/lib/fetchWithAuth";
import { MitraType } from './types';

export default async function MitraPage() {
  try {
    const res = await fetchWithAuth("http://localhost:8080/api/admin/mitra");
    if (!res.ok) redirect("/login");

    const mitraData: { mitra: MitraType[] } = await res.json();

    return <MitraClient mitraData={mitraData.mitra} />;
  } catch (err) {
    redirect("/login");
  }

  return <MitraClient mitraData={[]} />;
}
