'use client'

import { useState, useEffect } from 'react'
import Toast from './Toast'

interface ToastMessage {
  id: string
  title: string
  message?: string
  type: 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

export default function ToastContainer() {
  const [toasts, setToasts] = useState<ToastMessage[]>([])

  useEffect(() => {
    // Listen for custom toast events
    const handleToast = (event: CustomEvent) => {
      const { title, message, type, duration } = event.detail
      addToast(title, message, type, duration)
    }

    window.addEventListener('toast', handleToast as EventListener)
    return () => window.removeEventListener('toast', handleToast as EventListener)
  }, [])

  const addToast = (title: string, message?: string, type: ToastMessage['type'] = 'info', duration = 5000) => {
    const id = Math.random().toString(36).substr(2, 9)
    const newToast: ToastMessage = { id, title, message, type, duration }
    
    setToasts(prev => [...prev, newToast])
  }

  const removeToast = (id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id))
  }

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2">
      {toasts.map(toast => (
        <Toast
          key={toast.id}
          id={toast.id}
          title={toast.title}
          message={toast.message}
          type={toast.type}
          duration={toast.duration}
          onClose={removeToast}
        />
      ))}
    </div>
  )
}

// Utility functions to show toasts
export const showToast = (title: string, message?: string, type: ToastMessage['type'] = 'info', duration = 5000) => {
  const event = new CustomEvent('toast', {
    detail: { title, message, type, duration }
  })
  window.dispatchEvent(event)
}

export const showSuccess = (title: string, message?: string, duration = 5000) => showToast(title, message, 'success', duration)
export const showError = (title: string, message?: string, duration = 5000) => showToast(title, message, 'error', duration)
export const showWarning = (title: string, message?: string, duration = 5000) => showToast(title, message, 'warning', duration)
export const showInfo = (title: string, message?: string, duration = 5000) => showToast(title, message, 'info', duration)
