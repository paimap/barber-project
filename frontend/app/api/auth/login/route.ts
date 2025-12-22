import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export async function POST(req: NextRequest) {
  const body = await req.json();

  const res = await fetch("http://localhost:8080/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  if (!res.ok) {
    return NextResponse.json(
      { error: "login failed" },
      { status: 401 }
    );
  }

  const data = await res.json();
  const token = data.token;

  const response = NextResponse.json({ success: true });

  response.cookies.set({
    name: "access_token",
    value: token,
    httpOnly: true,
    path: "/",
    sameSite: "lax",
    secure: false, // true kalau https
  });

  return response;
}
