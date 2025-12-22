import { cookies } from "next/headers";
import { redirect } from "next/navigation";

export default async function HomePage() {
  const cookieStore = await cookies(); // await karena async
  const token = cookieStore.get("access_token")?.value;

  if (token) {
    redirect("/dashboard");
  }

  return <h1>Welcome! Please <a href="/login">login</a>.</h1>;
}