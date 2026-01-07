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
import { useAuth } from '@/context/AuthContext'


export default function Sidebar() {
  const pathname = usePathname()
  const router = useRouter()
  const auth = useAuth()

  if (!auth) {
    return null;
  }

  if (auth.loading) {
    return (
      <aside className={styles.sidebar}>
        <div className={styles.brand}>Barberin.</div>
        <nav className={styles.menu}>
          <p>Loading...</p>
        </nav>
      </aside>
    );
  } 


  const handleLogout = async () => {
    try {
      await auth.logout()
      window.location.href = "/login"
    } catch (err) {
      console.error(err)
      alert("Terjadi kesalahan saat logout")
    }
  }

  // Menu items berdasarkan role
  const getMenuItems = () => {
    const userRole = auth.user?.role
    console.log('Current user:', auth.user); // Debug logging
    console.log('User role:', userRole); // Debug logging

    if (userRole === 'SUPERADMIN') {
      return [
        { name: 'Dashboard', href: '/dashboard', icon: <LayoutDashboard size={18} /> },
        { name: 'Mitra', href: '/mitra', icon: <Handshake size={18} /> },
        { name: 'Service & Product', href: '/services', icon: <ShoppingCart size={18} /> },
        { name: 'Stock Management', href: '/stock', icon: <Box size={18} /> },
      ]
    } else if (userRole === 'MITRA') {
      return [
        { name: 'Dashboard', href: '/dashboard', icon: <LayoutDashboard size={18} /> },
        { name: 'Outlets', href: '/outlets', icon: <Handshake size={18} /> },
        { name: 'Barbers', href: '/barbers', icon: <ShoppingCart size={18} /> },
        { name: 'Service & Product', href: '/services', icon: <ShoppingCart size={18} /> },
        { name: 'Stock Management', href: '/stock-outlet', icon: <Box size={18} /> },
      ]
    } else if (userRole === 'BARBER') {
      return [
        { name: 'Dashboard', href: '/dashboard', icon: <LayoutDashboard size={18} /> },
        { name: 'Services', href: '/services-barber', icon: <ShoppingCart size={18} /> },
        { name: 'Product Sales', href: '/sales-barber', icon: <ShoppingCart size={18} /> },
      ]
    }
    return []
  }

  const menuItems = getMenuItems()

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
