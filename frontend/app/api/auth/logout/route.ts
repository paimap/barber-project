import { NextResponse } from "next/server";

export async function POST() {
  const response = NextResponse.json({ success: true });

  response.cookies.set({
    name: "access_token",
    value: "",
    maxAge: 0,
    httpOnly: true,
    path: "/",
    sameSite: "lax",
  });

  return response;
}
