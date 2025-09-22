'use client'

import Link from 'next/link'
import { ArrowRight, TrendingUp, Star, Users } from 'lucide-react'

interface Category {
  id: string
  name: string
  description: string
  image: string
  productCount: number
  isTrending?: boolean
  rating?: number
  reviewCount?: number
}

interface CategoryGridProps {
  categories: Category[]
  title?: string
  subtitle?: string
  showViewAll?: boolean
  maxItems?: number
}

export default function CategoryGrid({ 
  categories, 
  title = "Shop by Category",
  subtitle = "Discover products organized by category",
  showViewAll = true,
  maxItems = 6
}: CategoryGridProps) {
  const displayCategories = categories.slice(0, maxItems)

  return (
    <section className="py-20 bg-gradient-to-br from-white via-blue-50/30 to-purple-50/30">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-16">
          <div className="inline-flex items-center px-4 py-2 bg-blue-100 rounded-full text-sm font-medium text-blue-600 mb-4">
            <TrendingUp className="w-4 h-4 mr-2" />
            Popular Categories
          </div>
          
          <h2 className="text-responsive-3xl font-bold text-gray-900 mb-4">
            {title}
          </h2>
          
          <p className="text-responsive-lg text-gray-600 max-w-2xl mx-auto">
            {subtitle}
          </p>
        </div>

        {/* Categories Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8 mb-12">
          {displayCategories.map((category, index) => (
            <Link
              key={category.id}
              href={`/categories/${category.id}`}
              className="group relative overflow-hidden rounded-3xl bg-white shadow-lg hover:shadow-2xl transition-all duration-500 transform hover:-translate-y-2"
              style={{ animationDelay: `${index * 100}ms` }}
            >
              {/* Image Container */}
              <div className="relative aspect-[4/3] overflow-hidden">
                <img
                  src={category.image}
                  alt={category.name}
                  className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                />
                
                {/* Overlay */}
                <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                
                {/* Trending Badge */}
                {category.isTrending && (
                  <div className="absolute top-4 left-4 px-3 py-1 bg-gradient-to-r from-orange-500 to-red-500 text-white text-xs font-bold rounded-full flex items-center">
                    <TrendingUp className="w-3 h-3 mr-1" />
                    Trending
                  </div>
                )}
                
                {/* Product Count */}
                <div className="absolute top-4 right-4 px-3 py-1 bg-white/90 backdrop-blur-sm text-gray-700 text-sm font-medium rounded-full">
                  {category.productCount} items
                </div>
              </div>

              {/* Content */}
              <div className="p-6">
                <h3 className="text-xl font-bold text-gray-900 mb-2 group-hover:text-blue-600 transition-colors">
                  {category.name}
                </h3>
                
                <p className="text-gray-600 text-sm mb-4 line-clamp-2">
                  {category.description}
                </p>

                {/* Rating */}
                {category.rating && (
                  <div className="flex items-center mb-4">
                    <div className="flex items-center">
                      {[...Array(5)].map((_, i) => (
                        <Star
                          key={i}
                          size={16}
                          className={`${
                            i < Math.floor(category.rating!)
                              ? 'text-yellow-400 fill-current'
                              : 'text-gray-300'
                          }`}
                        />
                      ))}
                    </div>
                    <span className="text-sm text-gray-500 ml-2">
                      {category.rating.toFixed(1)} ({category.reviewCount} reviews)
                    </span>
                  </div>
                )}

                {/* CTA */}
                <div className="flex items-center justify-between">
                  <span className="text-blue-600 font-medium group-hover:text-blue-700 transition-colors">
                    Explore Category
                  </span>
                  <ArrowRight className="w-5 h-5 text-blue-600 group-hover:translate-x-1 transition-transform" />
                </div>
              </div>

              {/* Hover Effect */}
              <div className="absolute inset-0 bg-gradient-to-r from-blue-600/5 to-purple-600/5 opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>
            </Link>
          ))}
        </div>

        {/* View All Button */}
        {showViewAll && categories.length > maxItems && (
          <div className="text-center">
            <Link
              href="/categories"
              className="group inline-flex items-center px-8 py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-2xl hover:from-blue-700 hover:to-purple-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-1"
            >
              View All Categories
              <ArrowRight className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
            </Link>
          </div>
        )}
      </div>
    </section>
  )
}
