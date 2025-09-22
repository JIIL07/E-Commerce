'use client'

import { useState, useEffect } from 'react'
import { Award, Star, Shield, Zap } from 'lucide-react'

interface Brand {
  id: string
  name: string
  logo: string
  description: string
  rating: number
  products: number
  isVerified: boolean
  isPremium: boolean
}

interface BrandsSectionProps {
  title?: string
  subtitle?: string
  brands?: Brand[]
  showStats?: boolean
}

export default function BrandsSection({ 
  title = "Trusted Brands",
  subtitle = "Shop from the world's most trusted brands",
  brands,
  showStats = true
}: BrandsSectionProps) {
  const [currentIndex, setCurrentIndex] = useState(0)

  const defaultBrands: Brand[] = [
    {
      id: 'apple',
      name: 'Apple',
      logo: 'https://images.unsplash.com/photo-1611186871348-b1ce696e52c9?ixlib=rb-4.0.3&auto=format&fit=crop&w=200&q=80',
      description: 'Innovation and quality in every product',
      rating: 4.9,
      products: 45,
      isVerified: true,
      isPremium: true
    },
    {
      id: 'samsung',
      name: 'Samsung',
      logo: 'https://images.unsplash.com/photo-1610945265064-0e34e5519bbf?ixlib=rb-4.0.3&auto=format&fit=crop&w=200&q=80',
      description: 'Leading technology and innovation',
      rating: 4.8,
      products: 67,
      isVerified: true,
      isPremium: true
    },
    {
      id: 'nike',
      name: 'Nike',
      logo: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?ixlib=rb-4.0.3&auto=format&fit=crop&w=200&q=80',
      description: 'Just Do It - Premium athletic wear',
      rating: 4.7,
      products: 89,
      isVerified: true,
      isPremium: true
    },
    {
      id: 'adidas',
      name: 'Adidas',
      logo: 'https://images.unsplash.com/photo-1595950653106-6c9ebd614d3a?ixlib=rb-4.0.3&auto=format&fit=crop&w=200&q=80',
      description: 'Impossible is nothing',
      rating: 4.6,
      products: 72,
      isVerified: true,
      isPremium: true
    },
    {
      id: 'sony',
      name: 'Sony',
      logo: 'https://images.unsplash.com/photo-1606107557195-0e29a4b5b4aa?ixlib=rb-4.0.3&auto=format&fit=crop&w=200&q=80',
      description: 'Be moved by our technology',
      rating: 4.8,
      products: 56,
      isVerified: true,
      isPremium: true
    },
    {
      id: 'lg',
      name: 'LG',
      logo: 'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?ixlib=rb-4.0.3&auto=format&fit=crop&w=200&q=80',
      description: 'Life\'s Good with LG',
      rating: 4.5,
      products: 43,
      isVerified: true,
      isPremium: false
    }
  ]

  const displayBrands = brands || defaultBrands

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % displayBrands.length)
    }, 4000)

    return () => clearInterval(timer)
  }, [displayBrands.length])

  return (
    <section className="py-20 bg-gradient-to-br from-gray-50 to-blue-50">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-16">
          <div className="inline-flex items-center px-4 py-2 bg-blue-100 rounded-full text-sm font-medium text-blue-600 mb-4">
            <Award className="w-4 h-4 mr-2" />
            Premium Partners
          </div>
          
          <h2 className="text-responsive-3xl font-bold text-gray-900 mb-4">
            {title}
          </h2>
          
          <p className="text-responsive-lg text-gray-600 max-w-2xl mx-auto">
            {subtitle}
          </p>
        </div>

        {/* Featured Brand */}
        <div className="max-w-4xl mx-auto mb-16">
          <div className="bg-white rounded-3xl shadow-2xl p-8 md:p-12 relative overflow-hidden">
            {/* Background Pattern */}
            <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-blue-100 to-purple-100 rounded-full -translate-y-16 translate-x-16"></div>
            <div className="absolute bottom-0 left-0 w-24 h-24 bg-gradient-to-tr from-green-100 to-blue-100 rounded-full translate-y-12 -translate-x-12"></div>
            
            <div className="relative z-10 grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
              {/* Brand Info */}
              <div>
                <div className="flex items-center mb-6">
                  <img
                    src={displayBrands[currentIndex].logo}
                    alt={displayBrands[currentIndex].name}
                    className="w-20 h-20 rounded-2xl object-cover mr-4"
                  />
                  <div>
                    <h3 className="text-2xl font-bold text-gray-900">
                      {displayBrands[currentIndex].name}
                    </h3>
                    <div className="flex items-center">
                      <div className="flex text-yellow-400 mr-2">
                        {[...Array(5)].map((_, i) => (
                          <Star
                            key={i}
                            size={16}
                            className={`${
                              i < Math.floor(displayBrands[currentIndex].rating)
                                ? 'text-yellow-400 fill-current'
                                : 'text-gray-300'
                            }`}
                          />
                        ))}
                      </div>
                      <span className="text-sm text-gray-600">
                        {displayBrands[currentIndex].rating}
                      </span>
                    </div>
                  </div>
                </div>
                
                <p className="text-gray-600 mb-6">
                  {displayBrands[currentIndex].description}
                </p>
                
                <div className="flex items-center space-x-6 mb-6">
                  <div className="text-center">
                    <div className="text-2xl font-bold text-gray-900">
                      {displayBrands[currentIndex].products}
                    </div>
                    <div className="text-sm text-gray-600">Products</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-gray-900">
                      {displayBrands[currentIndex].rating}
                    </div>
                    <div className="text-sm text-gray-600">Rating</div>
                  </div>
                </div>
                
                <div className="flex items-center space-x-4">
                  {displayBrands[currentIndex].isVerified && (
                    <div className="flex items-center px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium">
                      <Shield className="w-4 h-4 mr-1" />
                      Verified
                    </div>
                  )}
                  {displayBrands[currentIndex].isPremium && (
                    <div className="flex items-center px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
                      <Zap className="w-4 h-4 mr-1" />
                      Premium
                    </div>
                  )}
                </div>
              </div>
              
              {/* Brand Image */}
              <div className="relative">
                <div className="aspect-square rounded-2xl overflow-hidden">
                  <img
                    src={displayBrands[currentIndex].logo}
                    alt={displayBrands[currentIndex].name}
                    className="w-full h-full object-cover"
                  />
                </div>
                <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent rounded-2xl"></div>
              </div>
            </div>
          </div>
        </div>

        {/* Brands Grid */}
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-6">
          {displayBrands.map((brand, index) => (
            <div
              key={brand.id}
              className={`group relative p-6 bg-white rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 cursor-pointer ${
                index === currentIndex ? 'ring-2 ring-blue-500 scale-105' : 'hover:scale-105'
              }`}
              onClick={() => setCurrentIndex(index)}
            >
              <div className="text-center">
                <div className="relative mb-4">
                  <img
                    src={brand.logo}
                    alt={brand.name}
                    className="w-16 h-16 rounded-xl object-cover mx-auto"
                  />
                  {brand.isVerified && (
                    <div className="absolute -top-1 -right-1 w-5 h-5 bg-green-500 rounded-full flex items-center justify-center">
                      <Shield className="w-3 h-3 text-white" />
                    </div>
                  )}
                </div>
                
                <h4 className="font-semibold text-gray-900 mb-2 text-sm">
                  {brand.name}
                </h4>
                
                <div className="flex items-center justify-center mb-2">
                  <Star className="w-4 h-4 text-yellow-400 fill-current mr-1" />
                  <span className="text-sm text-gray-600">{brand.rating}</span>
                </div>
                
                <p className="text-xs text-gray-500">
                  {brand.products} products
                </p>
              </div>
              
              {/* Hover Effect */}
              <div className="absolute inset-0 bg-gradient-to-r from-blue-500/5 to-purple-500/5 rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>
            </div>
          ))}
        </div>

        {/* Stats */}
        {showStats && (
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mt-16 pt-16 border-t border-gray-200">
            <div className="text-center">
              <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <Award className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-2">500+</h3>
              <p className="text-gray-600">Brand Partners</p>
            </div>
            
            <div className="text-center">
              <div className="w-16 h-16 bg-gradient-to-r from-green-500 to-green-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <Star className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-2">4.8</h3>
              <p className="text-gray-600">Average Rating</p>
            </div>
            
            <div className="text-center">
              <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-purple-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <Shield className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-2">100%</h3>
              <p className="text-gray-600">Verified Brands</p>
            </div>
            
            <div className="text-center">
              <div className="w-16 h-16 bg-gradient-to-r from-orange-500 to-orange-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
                <Zap className="w-8 h-8 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-gray-900 mb-2">24/7</h3>
              <p className="text-gray-600">Support</p>
            </div>
          </div>
        )}
      </div>
    </section>
  )
}
