// proxy.ts
import { NextResponse, NextRequest } from "next/server";

export function proxy(request: NextRequest) {
  const token = request.cookies.get('access_token')?.value;
  const { pathname } = request.nextUrl;

  // 1. Jika TIDAK ada token dan mencoba akses selain halaman login
  if (!token && pathname !== '/login') {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // 2. Jika ADA token dan mencoba akses halaman login
  if (token && pathname === '/login') {
    return NextResponse.redirect(new URL('/', request.url));
  }

  // 3. Redirect halaman utama "/" ke dashboard jika sudah login, atau login jika belum
  if (pathname === '/') {
    if (token) return NextResponse.redirect(new URL('/dashboard', request.url));
    return NextResponse.redirect(new URL('/login', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/dashboard/:path*", "/"],
};
