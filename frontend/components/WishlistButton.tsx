'use client'

import { useState, useEffect } from 'react'
import { Heart } from 'lucide-react'

interface WishlistButtonProps {
  productId: string
  className?: string
  size?: 'sm' | 'md' | 'lg'
}

export default function WishlistButton({ productId, className = '', size = 'md' }: WishlistButtonProps) {
  const [isInWishlist, setIsInWishlist] = useState(false)
  const [loading, setLoading] = useState(false)
  const [isLoggedIn, setIsLoggedIn] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem('token')
    setIsLoggedIn(!!token)
    
    if (token) {
      checkWishlistStatus()
    }
  }, [productId])

  const checkWishlistStatus = async () => {
    try {
      const token = localStorage.getItem('token')
      if (!token) return

      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'
      const response = await fetch(`${apiUrl}/api/wishlist/check/${productId}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        setIsInWishlist(data.inWishlist)
      }
    } catch (error) {
      console.error('Error checking wishlist status:', error)
    }
  }

  const toggleWishlist = async () => {
    if (!isLoggedIn) {
      // Redirect to login
      window.location.href = '/login'
      return
    }

    setLoading(true)
    try {
      const token = localStorage.getItem('token')
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'

      if (isInWishlist) {
        // Remove from wishlist
        const response = await fetch(`${apiUrl}/api/wishlist/${productId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })

        if (response.ok) {
          setIsInWishlist(false)
        }
      } else {
        // Add to wishlist
        const response = await fetch(`${apiUrl}/api/wishlist`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ productId })
        })

        if (response.ok) {
          setIsInWishlist(true)
        }
      }
    } catch (error) {
      console.error('Error toggling wishlist:', error)
    } finally {
      setLoading(false)
    }
  }

  const sizeClasses = {
    sm: 'w-8 h-8',
    md: 'w-10 h-10',
    lg: 'w-12 h-12'
  }

  const iconSizes = {
    sm: 16,
    md: 20,
    lg: 24
  }

  return (
    <button
      onClick={toggleWishlist}
      disabled={loading}
      className={`
        ${sizeClasses[size]}
        ${isInWishlist 
          ? 'bg-pink-500 text-white' 
          : 'bg-white/90 backdrop-blur-sm text-gray-600 hover:bg-pink-50 hover:text-pink-500'
        }
        rounded-full flex items-center justify-center transition-all duration-300 
        hover:scale-110 active:scale-95 shadow-lg hover:shadow-xl
        disabled:opacity-50 disabled:cursor-not-allowed
        ${className}
      `}
      title={isInWishlist ? 'Remove from wishlist' : 'Add to wishlist'}
    >
      <Heart 
        size={iconSizes[size]}
        className={isInWishlist ? 'fill-current' : ''}
      />
    </button>
  )
}
