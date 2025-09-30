import { useEffect, useRef, useState, useCallback } from 'react'

interface WebSocketMessage {
  type: string
  data: any
  timestamp: number
  user_id?: string
}

interface UseWebSocketOptions {
  url: string
  onMessage?: (message: WebSocketMessage) => void
  onOpen?: () => void
  onClose?: () => void
  onError?: (error: Event) => void
  reconnectInterval?: number
  maxReconnectAttempts?: number
}

export function useWebSocket({
  url,
  onMessage,
  onOpen,
  onClose,
  onError,
  reconnectInterval = 5000,
  maxReconnectAttempts = 5
}: UseWebSocketOptions) {
  const [isConnected, setIsConnected] = useState(false)
  const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'disconnected' | 'error'>('disconnected')
  const ws = useRef<WebSocket | null>(null)
  const reconnectAttempts = useRef(0)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>()

  const connect = useCallback(() => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      return
    }

    setConnectionStatus('connecting')
    
    try {
      ws.current = new WebSocket(url)
      
      ws.current.onopen = () => {
        setIsConnected(true)
        setConnectionStatus('connected')
        reconnectAttempts.current = 0
        onOpen?.()
      }
      
      ws.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          onMessage?.(message)
        } catch (error) {
          console.error('Error parsing WebSocket message:', error)
        }
      }
      
      ws.current.onclose = () => {
        setIsConnected(false)
        setConnectionStatus('disconnected')
        onClose?.()
        
        // Attempt to reconnect
        if (reconnectAttempts.current < maxReconnectAttempts) {
          reconnectAttempts.current++
          reconnectTimeoutRef.current = setTimeout(() => {
            connect()
          }, reconnectInterval)
        }
      }
      
      ws.current.onerror = (error) => {
        setConnectionStatus('error')
        onError?.(error)
      }
    } catch (error) {
      setConnectionStatus('error')
      onError?.(error as Event)
    }
  }, [url, onMessage, onOpen, onClose, onError, reconnectInterval, maxReconnectAttempts])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }
    
    if (ws.current) {
      ws.current.close()
      ws.current = null
    }
    
    setIsConnected(false)
    setConnectionStatus('disconnected')
  }, [])

  const sendMessage = useCallback((message: any) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message))
    } else {
      console.warn('WebSocket is not connected')
    }
  }, [])

  const sendPing = useCallback(() => {
    sendMessage({
      type: 'ping',
      data: { timestamp: Date.now() },
      timestamp: Date.now()
    })
  }, [sendMessage])

  useEffect(() => {
    connect()
    
    return () => {
      disconnect()
    }
  }, [connect, disconnect])

  return {
    isConnected,
    connectionStatus,
    sendMessage,
    sendPing,
    connect,
    disconnect
  }
}

// Hook for real-time notifications
export function useNotifications() {
  const [notifications, setNotifications] = useState<Array<{
    id: string
    type: string
    title: string
    message: string
    icon?: string
    priority?: string
    category?: string
    timestamp: number
    read: boolean
  }>>([])

  const { sendMessage } = useWebSocket({
    url: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:5000/ws',
    onMessage: (message) => {
      if (message.type === 'notification') {
        const notification = {
          id: `${message.timestamp}-${Math.random()}`,
          type: message.data.type || 'info',
          title: message.data.title,
          message: message.data.message,
          icon: message.data.icon,
          priority: message.data.priority || 'medium',
          category: message.data.category,
          timestamp: message.timestamp,
          read: false
        }
        
        setNotifications(prev => [notification, ...prev])
        
        // Auto-remove low priority notifications after 5 seconds
        if (notification.priority === 'low') {
          setTimeout(() => {
            setNotifications(prev => prev.filter(n => n.id !== notification.id))
          }, 5000)
        }
      }
    }
  })

  const markAsRead = (id: string) => {
    setNotifications(prev => 
      prev.map(notification => 
        notification.id === id 
          ? { ...notification, read: true }
          : notification
      )
    )
  }

  const removeNotification = (id: string) => {
    setNotifications(prev => prev.filter(notification => notification.id !== id))
  }

  const clearAll = () => {
    setNotifications([])
  }

  const unreadCount = notifications.filter(n => !n.read).length

  return {
    notifications,
    unreadCount,
    markAsRead,
    removeNotification,
    clearAll
  }
}

// Hook for real-time cart updates
export function useCartWebSocket() {
  const { sendMessage } = useWebSocket({
    url: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:5000/ws',
    onMessage: (message) => {
      if (message.type === 'cart_updated') {
        // Trigger cart refresh
        window.dispatchEvent(new CustomEvent('cartUpdated', { 
          detail: message.data 
        }))
      }
    }
  })

  const trackAddToCart = (productId: string, quantity: number) => {
    sendMessage({
      type: 'add_to_cart',
      data: {
        product_id: productId,
        quantity: quantity
      },
      timestamp: Date.now()
    })
  }

  const trackRemoveFromCart = (productId: string) => {
    sendMessage({
      type: 'remove_from_cart',
      data: {
        product_id: productId
      },
      timestamp: Date.now()
    })
  }

  const trackCheckoutStart = (total: number, itemCount: number) => {
    sendMessage({
      type: 'checkout_start',
      data: {
        total: total,
        item_count: itemCount
      },
      timestamp: Date.now()
    })
  }

  const trackCheckoutComplete = (orderId: string, total: number) => {
    sendMessage({
      type: 'checkout_complete',
      data: {
        order_id: orderId,
        total: total
      },
      timestamp: Date.now()
    })
  }

  return {
    trackAddToCart,
    trackRemoveFromCart,
    trackCheckoutStart,
    trackCheckoutComplete
  }
}

// Hook for real-time product updates
export function useProductWebSocket() {
  const { sendMessage } = useWebSocket({
    url: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:5000/ws',
    onMessage: (message) => {
      if (message.type === 'product_update') {
        // Trigger product refresh
        window.dispatchEvent(new CustomEvent('productUpdated', { 
          detail: message.data 
        }))
      } else if (message.type === 'price_alert') {
        // Show price change notification
        window.dispatchEvent(new CustomEvent('priceAlert', { 
          detail: message.data 
        }))
      } else if (message.type === 'stock_alert') {
        // Show stock alert
        window.dispatchEvent(new CustomEvent('stockAlert', { 
          detail: message.data 
        }))
      }
    }
  })

  const trackProductView = (productId: string) => {
    sendMessage({
      type: 'product_view',
      data: {
        product_id: productId
      },
      timestamp: Date.now()
    })
  }

  const trackProductReview = (productId: string, rating: number) => {
    sendMessage({
      type: 'product_review',
      data: {
        product_id: productId,
        rating: rating
      },
      timestamp: Date.now()
    })
  }

  const trackWishlistAdd = (productId: string) => {
    sendMessage({
      type: 'wishlist_add',
      data: {
        product_id: productId
      },
      timestamp: Date.now()
    })
  }

  const trackWishlistRemove = (productId: string) => {
    sendMessage({
      type: 'wishlist_remove',
      data: {
        product_id: productId
      },
      timestamp: Date.now()
    })
  }

  return {
    trackProductView,
    trackProductReview,
    trackWishlistAdd,
    trackWishlistRemove
  }
}

// Hook for real-time order updates
export function useOrderWebSocket() {
  const { sendMessage } = useWebSocket({
    url: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:5000/ws',
    onMessage: (message) => {
      if (message.type === 'order_update') {
        // Trigger order refresh
        window.dispatchEvent(new CustomEvent('orderUpdated', { 
          detail: message.data 
        }))
      }
    }
  })

  return {
    // Order updates are typically sent from the server
    // This hook is mainly for listening to updates
  }
}

// Hook for real-time analytics (admin only)
export function useAnalyticsWebSocket() {
  const [analytics, setAnalytics] = useState<any>(null)
  const [realtimeStats, setRealtimeStats] = useState<any>(null)

  useWebSocket({
    url: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:5000/ws',
    onMessage: (message) => {
      if (message.type === 'analytics_update') {
        setAnalytics(message.data.metrics)
      } else if (message.type === 'realtime_stats') {
        setRealtimeStats(message.data.stats)
      }
    }
  })

  return {
    analytics,
    realtimeStats
  }
}
