'use client'

import { useState } from 'react'
import { Star } from 'lucide-react'

interface StarRatingProps {
  rating: number
  maxRating?: number
  size?: 'sm' | 'md' | 'lg'
  interactive?: boolean
  onRatingChange?: (rating: number) => void
  showValue?: boolean
  className?: string
}

export default function StarRating({
  rating,
  maxRating = 5,
  size = 'md',
  interactive = false,
  onRatingChange,
  showValue = false,
  className = ''
}: StarRatingProps) {
  const [hoveredRating, setHoveredRating] = useState(0)

  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6'
  }

  const handleClick = (newRating: number) => {
    if (interactive && onRatingChange) {
      onRatingChange(newRating)
    }
  }

  const handleMouseEnter = (newRating: number) => {
    if (interactive) {
      setHoveredRating(newRating)
    }
  }

  const handleMouseLeave = () => {
    if (interactive) {
      setHoveredRating(0)
    }
  }

  const displayRating = hoveredRating || rating

  return (
    <div className={`flex items-center space-x-1 ${className}`}>
      <div className="flex items-center">
        {[...Array(maxRating)].map((_, index) => {
          const starValue = index + 1
          const isFilled = starValue <= displayRating
          
          return (
            <button
              key={index}
              type="button"
              onClick={() => handleClick(starValue)}
              onMouseEnter={() => handleMouseEnter(starValue)}
              onMouseLeave={handleMouseLeave}
              disabled={!interactive}
              className={`
                ${interactive ? 'cursor-pointer hover:scale-110' : 'cursor-default'}
                transition-transform duration-150
                ${interactive ? 'focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-1 rounded' : ''}
              `}
            >
              <Star
                className={`
                  ${sizeClasses[size]}
                  ${isFilled 
                    ? 'text-yellow-400 fill-current' 
                    : 'text-gray-300'
                  }
                  transition-colors duration-150
                `}
              />
            </button>
          )
        })}
      </div>
      
      {showValue && (
        <span className="text-sm text-gray-600 ml-2">
          {rating.toFixed(1)}/{maxRating}
        </span>
      )}
    </div>
  )
}
