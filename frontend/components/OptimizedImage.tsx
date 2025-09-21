'use client'

import { useState } from 'react'
import { Package } from 'lucide-react'

interface OptimizedImageProps {
  src: string
  alt: string
  className?: string
  fallbackIcon?: React.ReactNode
}

export default function OptimizedImage({ 
  src, 
  alt, 
  className = '', 
  fallbackIcon 
}: OptimizedImageProps) {
  const [imageError, setImageError] = useState(false)
  const [imageLoading, setImageLoading] = useState(true)

  const handleImageLoad = () => {
    setImageLoading(false)
  }

  const handleImageError = () => {
    setImageError(true)
    setImageLoading(false)
  }

  if (imageError) {
    return (
      <div className={`w-full h-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center ${className}`}>
        {fallbackIcon || <Package className="w-16 h-16 text-gray-400" />}
      </div>
    )
  }

  return (
    <div className={`relative ${className}`}>
      {imageLoading && (
        <div className="absolute inset-0 bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
          <div className="w-8 h-8 border-2 border-gray-300 border-t-blue-500 rounded-full animate-spin"></div>
        </div>
      )}
      <img 
        src={src} 
        alt={alt}
        className={`w-full h-full object-cover transition-opacity duration-300 ${
          imageLoading ? 'opacity-0' : 'opacity-100'
        }`}
        onLoad={handleImageLoad}
        onError={handleImageError}
        loading="lazy"
      />
    </div>
  )
}
