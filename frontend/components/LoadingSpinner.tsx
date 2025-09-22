'use client'

import { useState, useEffect } from 'react'

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl'
  color?: 'blue' | 'white' | 'gray' | 'pink' | 'green' | 'purple'
  text?: string
  fullScreen?: boolean
  variant?: 'spinner' | 'dots' | 'pulse' | 'bounce'
  className?: string
}

export default function LoadingSpinner({ 
  size = 'md', 
  color = 'blue', 
  text,
  fullScreen = false,
  variant = 'spinner',
  className = ''
}: LoadingSpinnerProps) {
  const [dots, setDots] = useState('')

  useEffect(() => {
    const interval = setInterval(() => {
      setDots(prev => prev.length >= 3 ? '' : prev + '.')
    }, 500)

    return () => clearInterval(interval)
  }, [])

  const getSizeClasses = () => {
    switch (size) {
      case 'sm':
        return 'w-4 h-4'
      case 'md':
        return 'w-8 h-8'
      case 'lg':
        return 'w-12 h-12'
      case 'xl':
        return 'w-16 h-16'
      default:
        return 'w-8 h-8'
    }
  }

  const getColorClasses = () => {
    switch (color) {
      case 'blue':
        return 'border-blue-600 bg-blue-600'
      case 'white':
        return 'border-white bg-white'
      case 'gray':
        return 'border-gray-600 bg-gray-600'
      case 'pink':
        return 'border-pink-600 bg-pink-600'
      case 'green':
        return 'border-green-600 bg-green-600'
      case 'purple':
        return 'border-purple-600 bg-purple-600'
      default:
        return 'border-blue-600 bg-blue-600'
    }
  }

  const renderSpinner = () => {
    switch (variant) {
      case 'dots':
        return (
          <div className="flex space-x-1">
            {[0, 1, 2].map((i) => (
              <div
                key={i}
                className={`${getSizeClasses()} ${getColorClasses().split(' ')[1]} rounded-full animate-bounce`}
                style={{ animationDelay: `${i * 0.1}s` }}
              />
            ))}
          </div>
        )
      
      case 'pulse':
        return (
          <div className={`${getSizeClasses()} ${getColorClasses().split(' ')[1]} rounded-full animate-pulse`} />
        )
      
      case 'bounce':
        return (
          <div className="flex space-x-1">
            {[0, 1, 2].map((i) => (
              <div
                key={i}
                className={`w-2 h-2 ${getColorClasses().split(' ')[1]} rounded-full animate-bounce`}
                style={{ animationDelay: `${i * 0.1}s` }}
              />
            ))}
          </div>
        )
      
      default:
        return (
          <div className={`
            ${getSizeClasses()} 
            ${getColorClasses().split(' ')[0]} 
            border-2 border-t-transparent rounded-full animate-spin
            ${className}
          `} />
        )
    }
  }

  const spinner = (
    <div className="flex flex-col items-center justify-center">
      {renderSpinner()}
      {text && (
        <p className="mt-2 text-sm text-gray-600 animate-pulse">
          {text}{dots}
        </p>
      )}
    </div>
  )

  if (fullScreen) {
    return (
      <div className="fixed inset-0 bg-white bg-opacity-75 flex items-center justify-center z-50 backdrop-blur-sm">
        <div className="bg-white rounded-lg p-8 shadow-xl">
          {spinner}
        </div>
      </div>
    )
  }

  return spinner
}

// Page Loading Component
export function PageLoader({ text = 'Loading...' }: { text?: string }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <LoadingSpinner size="xl" variant="spinner" text={text} />
      </div>
    </div>
  )
}

// Button Loading Component
export function ButtonLoader({ size = 'sm' }: { size?: 'sm' | 'md' }) {
  return <LoadingSpinner size={size} color="white" variant="spinner" />
}
