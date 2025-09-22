'use client'

import { useState } from 'react'
import { Minus, Plus } from 'lucide-react'

interface QuantitySelectorProps {
  value: number
  min?: number
  max?: number
  onChange: (value: number) => void
  disabled?: boolean
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export default function QuantitySelector({
  value,
  min = 1,
  max = 99,
  onChange,
  disabled = false,
  size = 'md',
  className = ''
}: QuantitySelectorProps) {
  const [localValue, setLocalValue] = useState(value)

  const sizeClasses = {
    sm: {
      button: 'w-8 h-8',
      input: 'w-12 h-8 text-sm',
      icon: 'w-3 h-3'
    },
    md: {
      button: 'w-10 h-10',
      input: 'w-16 h-10 text-base',
      icon: 'w-4 h-4'
    },
    lg: {
      button: 'w-12 h-12',
      input: 'w-20 h-12 text-lg',
      icon: 'w-5 h-5'
    }
  }

  const handleDecrease = () => {
    if (localValue > min && !disabled) {
      const newValue = localValue - 1
      setLocalValue(newValue)
      onChange(newValue)
    }
  }

  const handleIncrease = () => {
    if (localValue < max && !disabled) {
      const newValue = localValue + 1
      setLocalValue(newValue)
      onChange(newValue)
    }
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = parseInt(e.target.value) || min
    
    if (inputValue >= min && inputValue <= max) {
      setLocalValue(inputValue)
      onChange(inputValue)
    }
  }

  const handleBlur = () => {
    if (localValue < min) {
      setLocalValue(min)
      onChange(min)
    } else if (localValue > max) {
      setLocalValue(max)
      onChange(max)
    }
  }

  return (
    <div className={`flex items-center border border-gray-300 rounded-lg overflow-hidden ${className}`}>
      <button
        type="button"
        onClick={handleDecrease}
        disabled={disabled || localValue <= min}
        className={`
          ${sizeClasses[size].button}
          flex items-center justify-center
          bg-gray-100 hover:bg-gray-200
          disabled:bg-gray-50 disabled:cursor-not-allowed
          transition-colors duration-200
          border-r border-gray-300
        `}
        aria-label="Decrease quantity"
      >
        <Minus className={`${sizeClasses[size].icon} text-gray-600`} />
      </button>

      <input
        type="number"
        value={localValue}
        onChange={handleInputChange}
        onBlur={handleBlur}
        min={min}
        max={max}
        disabled={disabled}
        className={`
          ${sizeClasses[size].input}
          text-center font-medium
          border-0 focus:outline-none focus:ring-0
          disabled:bg-gray-50 disabled:cursor-not-allowed
          bg-white
        `}
        aria-label="Quantity"
      />

      <button
        type="button"
        onClick={handleIncrease}
        disabled={disabled || localValue >= max}
        className={`
          ${sizeClasses[size].button}
          flex items-center justify-center
          bg-gray-100 hover:bg-gray-200
          disabled:bg-gray-50 disabled:cursor-not-allowed
          transition-colors duration-200
          border-l border-gray-300
        `}
        aria-label="Increase quantity"
      >
        <Plus className={`${sizeClasses[size].icon} text-gray-600`} />
      </button>
    </div>
  )
}
