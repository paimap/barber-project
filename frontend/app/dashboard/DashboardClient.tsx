"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import styles from "./Dashboard.module.css";

type Props = {
  user: {
    user_id: number;
    email: string;
  };
};

export default function DashboardClient({ user }: Props) {
  const [count, setCount] = useState(0);
  const router = useRouter();

  const handleLogout = async () => {
    const res = await fetch("/api/auth/logout", { method: "POST" });
    if (res.ok) {
      router.push("/login");
    } else {
      alert("Logout failed");
    }
  };

  return (
    <div>
        <div className={styles.ujang}>

        </div>
        <h1>Welcome, {user.email}</h1>
        <p>Counter: {count}</p>
        <button onClick={() => setCount(count + 1)}>Increment</button>
        <br />
        <button onClick={handleLogout} style={{ marginTop: "10px" }}>
          Logout
        </button>
    </div>
  );
}
