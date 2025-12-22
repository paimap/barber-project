import { cookies } from "next/headers";

export async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const token = (await cookies()).get("access_token");
  if (!token) throw new Error("No token found");

  const res = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      Authorization: `Bearer ${token.value}`,
    },
    cache: "no-store", 
  });

  return res;
}
