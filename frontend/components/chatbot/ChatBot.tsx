"use client";

import React, { useState, useRef, useEffect } from "react";
import { usePathname } from "next/navigation";
import { MessageCircle, X, Send, Bot, Sparkles } from "lucide-react";
import ReactMarkdown from "react-markdown";
import styles from "./ChatBot.module.css";
import { useAuth } from "@/context/AuthContext";

interface ChatLog {
  role: "user" | "model";
  content: string;
}

export default function ChatBot() {
  const { user, loading: authLoading } = useAuth();
  const [isOpen, setIsOpen] = useState(false);
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);
  const pathname = usePathname();
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const defaultMsg: ChatLog = { 
    role: "model", 
    content: "Halo! Saya asisten AI Barberin. Ada yang bisa saya bantu dengan data bisnis Anda hari ini?" 
  };
  
  const [chatLog, setChatLog] = useState<ChatLog[]>([defaultMsg]);

  useEffect(() => {
    if (authLoading || user?.role !== "SUPERADMIN") return;

    const fetchHistory = async () => {
      try {
        const response = await fetch("/api/chatbot");
        const data = await response.json();
        if (data.history && data.history.length > 0) {
          const formattedHistory = data.history.map((h: any) => ({
            role: h.role,
            content: h.Content || h.content 
          }));
          setChatLog(formattedHistory);
        }
      } catch (err) {
        console.error("Gagal memuat history chat:", err);
      }
    };
    fetchHistory();
  }, [authLoading, user]);

  useEffect(() => {
    if (isOpen) {
      messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }
  }, [chatLog, isOpen]);

  const resolveContext = (path: string) => {
    const outletMatch = path.match(/\/mitra\/\d+\/outlet\/(\d+)/);
    if (outletMatch) return `outlet_${outletMatch[1]}`;
    const mitraMatch = path.match(/\/mitra\/(\d+)$/);
    if (mitraMatch) return `mitra_${mitraMatch[1]}`;
    if (path === "/" || path === "/dashboard") return "dashboard";
    return path.split("/").filter(Boolean).pop() || "dashboard";
  };

  const handleSend = async () => {
    if (!message.trim() || loading) return;
    const userMsg = message;
    const currentContext = resolveContext(pathname);
    setMessage("");
    setChatLog((prev) => [...prev, { role: "user", content: userMsg }]);
    setLoading(true);

    try {
      const response = await fetch("/api/chatbot", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ message: userMsg, page_context: currentContext }),
      });
      const data = await response.json();
      if (response.status === 429) {
        setChatLog((prev) => [...prev, { role: "model", content: data.error || "Kuota chat harian habis." }]);
        return;
      }
      if (!response.ok) throw new Error("Server error");
      if (data.reply) {
        setChatLog((prev) => [...prev, { role: "model", content: data.reply }]);
      }
    } catch (error) {
      setChatLog((prev) => [...prev, { role: "model", content: "Maaf, terjadi gangguan koneksi." }]);
    } finally {
      setLoading(false);
    }
  };

  if (authLoading || user?.role !== "SUPERADMIN") return null;

  return (
    <div className={styles.chatbotWrapper}>
      <button className={styles.chatbotToggle} onClick={() => setIsOpen(!isOpen)}>
        {isOpen ? <X size={28} /> : <MessageCircle size={28} />}
      </button>

      {isOpen && (
        <div className={styles.chatbotWindow}>
          <div className={styles.chatbotHeader}>
            <div className={styles.headerTitle}>
              <div className={styles.iconCircle}><Bot size={20} color="#0ea5e9" /></div>
              <div>
                <h3>BARBERIN AI</h3>
                <p>Context: {resolveContext(pathname).replace('_', ' #')}</p>
              </div>
            </div>
            <Sparkles size={16} className={styles.sparkleIcon} />
          </div>

          <div className={styles.messagesContainer}>
            {chatLog.map((chat, idx) => (
              <div 
                key={idx} 
                className={`${styles.message} ${styles[chat.role]} ${chat.content.includes("Kuota") ? styles.errorMessage : ""}`}
              >
                {/* Solusi Error: Bungkus ReactMarkdown dengan div className */}
                <div className={styles.markdownWrapper}>
                  <ReactMarkdown>{chat.content}</ReactMarkdown>
                </div>
              </div>
            ))}
            {loading && (
              <div className={`${styles.message} ${styles.model}`}>
                <span className={styles.typingIndicator}>Sedang berpikir...</span>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>

          <div className={styles.inputArea}>
            <input
              type="text"
              placeholder="Tanya laporan bisnis..."
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSend()}
              disabled={loading}
            />
            <button className={styles.sendBtn} onClick={handleSend} disabled={loading || !message.trim()}>
              <Send size={20} />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}