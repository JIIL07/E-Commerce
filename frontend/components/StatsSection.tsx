'use client'

import { useState, useEffect } from 'react'
import { Users, ShoppingBag, Star, Award, TrendingUp, Globe } from 'lucide-react'

interface Stat {
  id: string
  value: string
  label: string
  icon: React.ComponentType<any>
  color: string
  description?: string
}

interface StatsSectionProps {
  title?: string
  subtitle?: string
  stats?: Stat[]
  animated?: boolean
}

export default function StatsSection({ 
  title = "Our Impact",
  subtitle = "Numbers that speak for themselves",
  stats,
  animated = true
}: StatsSectionProps) {
  const [counters, setCounters] = useState<Record<string, number>>({})

  const defaultStats: Stat[] = [
    {
      id: 'customers',
      value: '50000',
      label: 'Happy Customers',
      icon: Users,
      color: 'from-blue-500 to-blue-600',
      description: 'Satisfied customers worldwide'
    },
    {
      id: 'products',
      value: '10000',
      label: 'Products Sold',
      icon: ShoppingBag,
      color: 'from-green-500 to-green-600',
      description: 'Quality products delivered'
    },
    {
      id: 'rating',
      value: '4.9',
      label: 'Average Rating',
      icon: Star,
      color: 'from-yellow-500 to-yellow-600',
      description: 'Based on customer reviews'
    },
    {
      id: 'countries',
      value: '50',
      label: 'Countries Served',
      icon: Globe,
      color: 'from-purple-500 to-purple-600',
      description: 'Global shipping network'
    },
    {
      id: 'growth',
      value: '150',
      label: 'Growth Rate',
      icon: TrendingUp,
      color: 'from-pink-500 to-pink-600',
      description: 'Year over year growth'
    },
    {
      id: 'awards',
      value: '25',
      label: 'Industry Awards',
      icon: Award,
      color: 'from-indigo-500 to-indigo-600',
      description: 'Recognition and achievements'
    }
  ]

  const displayStats = stats || defaultStats

  useEffect(() => {
    if (!animated) return

    const animateCounters = () => {
      displayStats.forEach((stat) => {
        const numericValue = parseFloat(stat.value.replace(/[^\d.]/g, ''))
        if (isNaN(numericValue)) return

        const duration = 2000 // 2 seconds
        const steps = 60
        const increment = numericValue / steps
        let current = 0

        const timer = setInterval(() => {
          current += increment
          if (current >= numericValue) {
            current = numericValue
            clearInterval(timer)
          }

          setCounters(prev => ({
            ...prev,
            [stat.id]: current
          }))
        }, duration / steps)
      })
    }

    // Start animation after a short delay
    const timer = setTimeout(animateCounters, 500)
    return () => clearTimeout(timer)
  }, [displayStats, animated])

  const formatValue = (stat: Stat, value: number) => {
    if (stat.id === 'rating') {
      return value.toFixed(1)
    }
    if (stat.id === 'growth') {
      return `${Math.round(value)}%`
    }
    if (value >= 1000) {
      return `${(value / 1000).toFixed(0)}K+`
    }
    return Math.round(value).toString()
  }

  return (
    <section className="py-20 bg-gradient-to-br from-gray-900 via-blue-900 to-purple-900 relative overflow-hidden">
      {/* Background Effects */}
      <div className="absolute inset-0">
        <div className="absolute top-0 left-0 w-full h-full">
          <div className="absolute top-20 left-20 w-72 h-72 bg-blue-500/10 rounded-full blur-3xl animate-float"></div>
          <div className="absolute top-40 right-20 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl animate-float" style={{ animationDelay: '1s' }}></div>
          <div className="absolute bottom-20 left-1/3 w-80 h-80 bg-pink-500/10 rounded-full blur-3xl animate-float" style={{ animationDelay: '2s' }}></div>
        </div>
      </div>

      <div className="container mx-auto px-4 relative z-10">
        {/* Header */}
        <div className="text-center mb-16">
          <div className="inline-flex items-center px-4 py-2 bg-white/10 backdrop-blur-sm rounded-full text-sm font-medium text-white mb-4">
            <TrendingUp className="w-4 h-4 mr-2" />
            Our Success Story
          </div>
          
          <h2 className="text-responsive-3xl font-bold text-white mb-4">
            {title}
          </h2>
          
          <p className="text-responsive-lg text-blue-100 max-w-2xl mx-auto">
            {subtitle}
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          {displayStats.map((stat, index) => {
            const currentValue = animated ? (counters[stat.id] || 0) : parseFloat(stat.value.replace(/[^\d.]/g, ''))
            
            return (
              <div
                key={stat.id}
                className="group relative"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <div className="bg-white/10 backdrop-blur-md rounded-3xl p-8 border border-white/20 hover:bg-white/20 transition-all duration-300 hover:scale-105">
                  {/* Icon */}
                  <div className={`w-16 h-16 bg-gradient-to-r ${stat.color} rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300`}>
                    <stat.icon className="w-8 h-8 text-white" />
                  </div>

                  {/* Value */}
                  <div className="mb-4">
                    <span className="text-4xl font-bold text-white">
                      {animated ? formatValue(stat, currentValue) : stat.value}
                    </span>
                    {stat.id === 'rating' && <span className="text-2xl text-yellow-400 ml-1">★</span>}
                    {stat.id === 'growth' && <span className="text-2xl text-green-400 ml-1">%</span>}
                  </div>

                  {/* Label */}
                  <h3 className="text-xl font-semibold text-white mb-2">
                    {stat.label}
                  </h3>

                  {/* Description */}
                  {stat.description && (
                    <p className="text-blue-100 text-sm">
                      {stat.description}
                    </p>
                  )}

                  {/* Progress Bar */}
                  <div className="mt-6">
                    <div className="w-full bg-white/20 rounded-full h-2">
                      <div 
                        className={`h-2 bg-gradient-to-r ${stat.color} rounded-full transition-all duration-2000 ease-out`}
                        style={{ 
                          width: animated ? `${Math.min((currentValue / parseFloat(stat.value.replace(/[^\d.]/g, ''))) * 100, 100)}%` : '100%' 
                        }}
                      ></div>
                    </div>
                  </div>
                </div>

                {/* Hover Effect */}
                <div className="absolute inset-0 bg-gradient-to-r from-blue-500/5 to-purple-500/5 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>
              </div>
            )
          })}
        </div>

        {/* Bottom CTA */}
        <div className="text-center mt-16">
          <div className="bg-white/10 backdrop-blur-md rounded-3xl p-8 border border-white/20 max-w-2xl mx-auto">
            <h3 className="text-2xl font-bold text-white mb-4">
              Join Our Growing Community
            </h3>
            <p className="text-blue-100 mb-6">
              Be part of our success story and experience the best in online shopping
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button className="btn-primary">
                Start Shopping
              </button>
              <button className="btn-secondary bg-white/20 text-white border-white/30 hover:bg-white/30">
                Learn More
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
