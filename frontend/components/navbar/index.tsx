'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import styles from './Navbar.module.css'
import {
  LayoutDashboard,
  Handshake,
  ShoppingCart,
  AlignEndHorizontal,
  LogOut,
  Box,
} from 'lucide-react'




export default function Sidebar() {
  const pathname = usePathname()
  const router = useRouter()


  const handleLogout = async () => {
    try {
      const res = await fetch("/api/auth/logout", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
      });

      if (res.ok) {
        window.location.href = "/login";
      } else {
        alert("Logout gagal - Silakan coba lagi");
      }
    } catch (err) {
      console.error(err);
      alert("Terjadi kesalahan saat logout");
    }
  }


  const menuItems = [
    { name: 'Dashboard', href: '/dashboard', icon: <LayoutDashboard size={18} /> },
    { name: 'Mitra', href: '/mitra', icon: <Handshake size={18} /> },
    { name: 'Service & Product', href: '/services', icon: <ShoppingCart size={18} /> },
    { name: 'Stock Management', href: '/stock', icon: <Box size={18} /> },
  ]

  return (
    <aside className={styles.sidebar}>
      <div className={styles.brand}>Barberin.</div>

      <nav className={styles.menu}>
        {menuItems.map((item) => {
          const isActive = pathname === item.href
          return (
            <Link 
              key={item.href} 
              href={item.href} 
              className={`${styles.item} ${isActive ? styles.active : ''}`}
            >
              {item.icon}
              <span>{item.name}</span>
            </Link>
          )
        })}
      </nav>

      <div className={styles.logout}>
        <button className={styles.logoutButton} onClick={handleLogout}>
          <LogOut size={18} />
          <span>Logout</span>
        </button>
      </div>
    </aside>
  )
}
