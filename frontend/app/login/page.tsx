"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Mail, Lock } from "lucide-react";
import { useAuth } from "@/context/AuthContext"; // Pastikan path import benar
import styles from "./LoginPage.module.css";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  
  const router = useRouter();
  const { user, refetchUser, loading: authLoading } = useAuth();

  // Redirect otomatis jika user sudah terdeteksi login oleh Context
  useEffect(() => {
    if (user) {
      if (user.role === "BARBER") {
        window.location.href = "/services-barber";
      } else {
        window.location.href = "/dashboard";
      }
    }
  }, [user, router]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);

    try {
      const res = await fetch("/api/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (res.ok) {
        // 1. Ambil data user terbaru ke dalam Context
        await refetchUser();
        
        // 2. Refresh router agar server component mendapatkan state cookie terbaru
        router.refresh();
        
        // Catatan: useEffect di atas akan menangani redirect setelah user terisi
      } else {
        const errorData = await res.json();
        alert(errorData.message || "Login gagal - Cek kembali email dan kata sandi Anda");
      }
    } catch (error) {
      console.error("Login error:", error);
      alert("Terjadi kesalahan koneksi");
    } finally {
      setLoading(false);
    }
  }

  // Jika sedang mengecek status auth, tampilkan loading agar tidak flicker
  if (authLoading) return <div className={styles.container}>Memuat...</div>;

  return (
    <div className={styles.container}>
      <div className={styles.loginCard}>
        <div className={styles.header}>
          <h1>Selamat Datang Kembali</h1>
          <p>Silakan masukkan detail akun Anda untuk masuk</p>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.inputGroup}>
            <label>Alamat Email</label>
            <div className={styles.inputWrapper}>
              <Mail size={18} />
              <input
                className={styles.input}
                type="email"
                placeholder="nama@perusahaan.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
          </div>

          <div className={styles.inputGroup}>
            <label>Kata Sandi</label>
            <div className={styles.inputWrapper}>
              <Lock size={18} />
              <input
                className={styles.input}
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
          </div>

          <button 
            type="submit" 
            className={styles.loginButton}
            disabled={loading}
          >
            {loading ? "Masuk..." : "Masuk"}
          </button>
        </form>

        <p className={styles.footerText}>
          Belum punya akun? Hubungi Admin
        </p>
      </div>
    </div>
  );
}