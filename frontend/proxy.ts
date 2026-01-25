// proxy.ts (Update sederhana)
import { NextResponse, NextRequest } from "next/server";

export function proxy(request: NextRequest) {
  const token = request.cookies.get('access_token')?.value;
  const { pathname } = request.nextUrl;

  // Jika tidak ada token dan coba akses dashboard/services
  if (!token && pathname !== '/login') {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // Jika sudah login tapi coba buka halaman login lagi
  if (token && pathname === '/login') {
    // Karena kita tidak tahu role di middleware (API terbatas), 
    // kita lempar ke / dulu, nanti diproses oleh AuthContext di sana.
    return NextResponse.redirect(new URL('/', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/dashboard/:path*", "/services-barber/:path*", "/"],
};