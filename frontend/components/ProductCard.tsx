'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Eye, ShoppingCart, Heart } from 'lucide-react'
import OptimizedImage from './OptimizedImage'
import WishlistButton from './WishlistButton'
import StarRating from './StarRating'

interface Product {
  id: string
  name: string
  price: number
  image: string
  category: {
    name: string
  }
  stock: number
  averageRating: number
  reviewCount: number
}

interface ProductCardProps {
  product: Product
  onAddToCart?: (productId: string) => void
  className?: string
}

export default function ProductCard({ product, onAddToCart, className = '' }: ProductCardProps) {
  const [isHovered, setIsHovered] = useState(false)

  const handleAddToCart = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (onAddToCart) {
      onAddToCart(product.id)
    }
  }

  return (
    <div 
      className={`group bg-white rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden ${className}`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <Link href={`/products/${product.id}`} className="block">
        {/* Product Image */}
        <div className="relative aspect-square overflow-hidden">
          <OptimizedImage
            src={product.image}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
          />
          
          {/* Category Badge */}
          <div className="absolute top-4 left-4 px-3 py-1 bg-white/90 backdrop-blur-sm rounded-full text-sm font-medium text-gray-700 shadow-lg">
            {product.category.name}
          </div>
          
          {/* Wishlist Button */}
          <div className="absolute top-4 right-4">
            <WishlistButton productId={product.id} />
          </div>
          
          {/* Stock Warning */}
          {product.stock < 10 && (
            <div className="absolute bottom-4 left-4 px-3 py-1 bg-red-500 text-white rounded-full text-xs font-medium">
              Only {product.stock} left
            </div>
          )}
          
          {/* Quick Actions Overlay */}
          <div className={`
            absolute inset-0 bg-black/50 flex items-center justify-center
            transition-opacity duration-300
            ${isHovered ? 'opacity-100' : 'opacity-0'}
          `}>
            <div className="flex space-x-3">
              <button
                onClick={handleAddToCart}
                className="bg-white text-gray-800 p-3 rounded-full hover:bg-gray-100 transition-colors shadow-lg"
                title="Add to cart"
              >
                <ShoppingCart className="w-5 h-5" />
              </button>
              <Link
                href={`/products/${product.id}`}
                className="bg-white text-gray-800 p-3 rounded-full hover:bg-gray-100 transition-colors shadow-lg"
                title="View details"
              >
                <Eye className="w-5 h-5" />
              </Link>
            </div>
          </div>
        </div>

        {/* Product Info */}
        <div className="p-6">
          <h3 className="font-bold text-xl mb-3 line-clamp-2 text-gray-900 group-hover:text-blue-600 transition-colors">
            {product.name}
          </h3>
          
          {/* Rating */}
          <div className="flex items-center mb-4">
            <StarRating 
              rating={product.averageRating} 
              size="sm" 
              showValue={false}
            />
            <span className="text-sm text-gray-500 ml-2">
              ({product.reviewCount})
            </span>
          </div>
          
          {/* Price and Actions */}
          <div className="flex items-center justify-between">
            <span className="text-3xl font-bold text-gray-900">
              ${product.price.toFixed(2)}
            </span>
            
            <button
              onClick={handleAddToCart}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors font-medium"
            >
              Add to Cart
            </button>
          </div>
        </div>
      </Link>
    </div>
  )
}
