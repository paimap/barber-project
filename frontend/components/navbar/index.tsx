'use client'

import styles from './Navbar.module.css'
import {
  LayoutDashboard,
  Handshake,
  ShoppingCart,
  AlignEndHorizontal,
  LogOut,
} from 'lucide-react'

export default function Sidebar() {
  return (
    <aside className={styles.sidebar}>
      <div className={styles.brand}>Barberin.</div>

      <nav className={styles.menu}>
        <div className={styles.item}>
          <LayoutDashboard size={18} />
          <span>Dashboard</span>
        </div>
        <div className={styles.item}>
          <Handshake size={18} />
          <span>Mitra</span>
        </div>
        <div className={styles.item}>
          <ShoppingCart size={18} />
          <span>Service & Product</span>
        </div>
        <div className={styles.item}>
          <AlignEndHorizontal size={18} />
          <span>Leaderboard</span>
        </div>
      </nav>

      <div className={styles.logout}>
        <div className={styles.item}>
          <LogOut size={18} />
          <span>Logout</span>
        </div>
      </div>
    </aside>
  )
}
