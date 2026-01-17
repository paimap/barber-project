// lib/fetchClient.ts

export async function fetchClient(url: string, options: RequestInit = {}) {
  // Ambil token dari browser cookie secara manual
  const token = document.cookie
    .split("; ")
    .find((row) => row.startsWith("access_token="))
    ?.split("=")[1];

  const headers = {
    ...options.headers,
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };

  return fetch(url, { ...options, headers });
}