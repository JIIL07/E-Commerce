'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Navbar from '../../components/Navbar'
import { 
  Heart,
  ShoppingCart,
  Trash2,
  ArrowLeft,
  Star,
  Package,
  Plus,
  Eye,
  Sparkles,
  Loader2
} from 'lucide-react'

interface WishlistItem {
  id: string
  product: {
    id: string
    name: string
    price: number
    images: string[]
    category: {
      name: string
    }
    averageRating: number
    reviewCount: number
    stock: number
  }
  addedAt: string
}

export default function WishlistPage() {
  const [wishlistItems, setWishlistItems] = useState<WishlistItem[]>([])
  const [loading, setLoading] = useState(true)
  const [removing, setRemoving] = useState<string | null>(null)
  const [addingToCart, setAddingToCart] = useState<string | null>(null)
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    fetchWishlist()
  }, [router])

  const fetchWishlist = async () => {
    try {
      const token = localStorage.getItem('token')
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'
      
      const response = await fetch(`${apiUrl}/api/wishlist`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        setWishlistItems(data.items || [])
      } else {
        console.error('Failed to fetch wishlist:', response.status)
        setWishlistItems([])
      }
    } catch (error) {
      console.error('Error fetching wishlist:', error)
      setWishlistItems([])
    } finally {
      setLoading(false)
    }
  }

  const removeFromWishlist = async (itemId: string) => {
    setRemoving(itemId)
    try {
      const token = localStorage.getItem('token')
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'
      
      // Find the product ID from the item
      const item = wishlistItems.find(item => item.id === itemId)
      if (!item) return

      const response = await fetch(`${apiUrl}/api/wishlist/${item.product.id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (response.ok) {
        setWishlistItems(prev => prev.filter(item => item.id !== itemId))
      }
    } catch (error) {
      console.error('Error removing from wishlist:', error)
    } finally {
      setRemoving(null)
    }
  }

  const addToCart = async (productId: string) => {
    setAddingToCart(productId)
    try {
      const token = localStorage.getItem('token')
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'
      
      const response = await fetch(`${apiUrl}/api/cart`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ productId, quantity: 1 })
      })

      if (response.ok) {
        // Remove from wishlist after adding to cart
        setWishlistItems(prev => prev.filter(item => item.product.id !== productId))
      }
    } catch (error) {
      console.error('Error adding to cart:', error)
    } finally {
      setAddingToCart(null)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-pink-50 via-purple-50 to-indigo-50">
        <Navbar />
        <div className="flex items-center justify-center min-h-[60vh]">
          <div className="text-center">
            <div className="relative">
              <div className="w-16 h-16 border-4 border-pink-200 border-t-pink-600 rounded-full animate-spin"></div>
              <div className="absolute inset-0 w-16 h-16 border-4 border-transparent border-t-purple-600 rounded-full animate-spin" style={{animationDelay: '0.5s'}}></div>
            </div>
            <p className="mt-4 text-gray-600">Loading your wishlist...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-pink-50 via-purple-50 to-indigo-50">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Link 
                href="/products" 
                className="flex items-center space-x-2 text-gray-600 hover:text-pink-600 transition-colors"
              >
                <ArrowLeft className="w-5 h-5" />
                <span>Back to Products</span>
              </Link>
            </div>
            <div className="text-right">
              <h1 className="text-3xl font-bold text-gray-900 flex items-center">
                <Heart className="w-8 h-8 text-pink-500 mr-3" />
                My Wishlist
              </h1>
              <p className="text-gray-600">{wishlistItems.length} {wishlistItems.length === 1 ? 'item' : 'items'}</p>
            </div>
          </div>
        </div>

        {wishlistItems.length === 0 ? (
          /* Empty Wishlist */
          <div className="text-center py-20">
            <div className="w-32 h-32 bg-gradient-to-r from-pink-100 to-purple-100 rounded-full flex items-center justify-center mx-auto mb-8">
              <Heart className="w-16 h-16 text-pink-400" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Your wishlist is empty</h2>
            <p className="text-gray-600 mb-8 max-w-md mx-auto">
              Start adding products you love to your wishlist and they'll appear here!
            </p>
            <Link 
              href="/products"
              className="inline-flex items-center px-8 py-4 bg-gradient-to-r from-pink-500 to-purple-600 text-white rounded-2xl font-bold text-lg hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl"
            >
              <Sparkles className="w-6 h-6 mr-3" />
              Start Shopping
              <ShoppingCart className="w-6 h-6 ml-3" />
            </Link>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {wishlistItems.map((item, index) => (
              <div 
                key={item.id}
                className="group bg-white/80 backdrop-blur-sm rounded-3xl overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-500 hover:scale-105"
                style={{animationDelay: `${index * 100}ms`}}
              >
                {/* Product Image */}
                <div className="relative h-64 overflow-hidden">
                  {item.product.images.length > 0 ? (
                    <img 
                      src={item.product.images[0]} 
                      alt={item.product.name}
                      className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                    />
                  ) : (
                    <div className="w-full h-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
                      <Package className="w-16 h-16 text-gray-400" />
                    </div>
                  )}
                  
                  {/* Overlay */}
                  <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                  
                  {/* Quick Actions */}
                  <div className="absolute top-4 right-4 flex flex-col gap-2 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-x-4 group-hover:translate-x-0">
                    <button className="w-10 h-10 bg-white/90 backdrop-blur-sm rounded-full flex items-center justify-center hover:bg-white transition-colors">
                      <Eye className="w-5 h-5 text-gray-600" />
                    </button>
                    <button
                      onClick={() => removeFromWishlist(item.id)}
                      disabled={removing === item.id}
                      className="w-10 h-10 bg-red-50 rounded-full flex items-center justify-center hover:bg-red-100 transition-colors group disabled:opacity-50"
                    >
                      {removing === item.id ? (
                        <Loader2 className="w-4 h-4 animate-spin text-red-500" />
                      ) : (
                        <Trash2 className="w-4 h-4 text-red-500 group-hover:scale-110 transition-transform" />
                      )}
                    </button>
                  </div>
                  
                  {/* Category Badge */}
                  <div className="absolute top-4 left-4">
                    <span className="px-3 py-1 bg-white/90 backdrop-blur-sm rounded-full text-xs font-medium text-gray-700">
                      {item.product.category.name}
                    </span>
                  </div>

                  {/* Heart Icon */}
                  <div className="absolute bottom-4 right-4">
                    <div className="w-10 h-10 bg-pink-500 rounded-full flex items-center justify-center">
                      <Heart className="w-5 h-5 text-white fill-current" />
                    </div>
                  </div>
                </div>
                
                {/* Product Info */}
                <div className="p-6">
                  <h3 className="font-bold text-lg mb-2 line-clamp-2 group-hover:text-pink-600 transition-colors">
                    {item.product.name}
                  </h3>
                  
                  {/* Rating */}
                  <div className="flex items-center mb-3">
                    <div className="flex text-yellow-400">
                      {[...Array(5)].map((_, i) => (
                        <Star 
                          key={i} 
                          size={16}
                          className={i < Math.floor(item.product.averageRating) ? 'fill-current text-yellow-400' : 'text-gray-300'}
                        />
                      ))}
                    </div>
                    <span className="text-sm text-gray-500 ml-2">({item.product.reviewCount})</span>
                  </div>
                  
                  {/* Price and Stock */}
                  <div className="flex items-center justify-between mb-4">
                    <span className="text-2xl font-bold text-gray-900">${item.product.price}</span>
                    <span className={`text-sm px-2 py-1 rounded-full ${
                      item.product.stock > 10 ? 'bg-green-100 text-green-800' : 
                      item.product.stock > 0 ? 'bg-yellow-100 text-yellow-800' : 
                      'bg-red-100 text-red-800'
                    }`}>
                      {item.product.stock > 0 ? `${item.product.stock} in stock` : 'Out of stock'}
                    </span>
                  </div>

                  {/* Actions */}
                  <div className="flex space-x-3">
                    <Link
                      href={`/products/${item.product.id}`}
                      className="flex-1 bg-gray-100 text-gray-700 py-3 px-4 rounded-2xl font-medium hover:bg-gray-200 transition-colors text-center"
                    >
                      View Details
                    </Link>
                    <button
                      onClick={() => addToCart(item.product.id)}
                      disabled={addingToCart === item.product.id || item.product.stock === 0}
                      className="flex-1 bg-gradient-to-r from-pink-500 to-purple-600 text-white py-3 px-4 rounded-2xl font-medium hover:scale-105 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
                    >
                      {addingToCart === item.product.id ? (
                        <Loader2 className="w-4 h-4 animate-spin" />
                      ) : (
                        <>
                          <Plus className="w-4 h-4 mr-2" />
                          Add to Cart
                        </>
                      )}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Quick Actions */}
        {wishlistItems.length > 0 && (
          <div className="mt-12 text-center">
            <div className="bg-white/80 backdrop-blur-sm rounded-3xl p-8 shadow-lg border border-white/20 max-w-2xl mx-auto">
              <h3 className="text-xl font-bold text-gray-900 mb-4">Quick Actions</h3>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <button className="px-6 py-3 bg-gradient-to-r from-pink-500 to-purple-600 text-white rounded-2xl font-medium hover:scale-105 transition-all duration-300">
                  Add All to Cart
                </button>
                <button className="px-6 py-3 bg-gray-100 text-gray-700 rounded-2xl font-medium hover:bg-gray-200 transition-colors">
                  Share Wishlist
                </button>
                <button className="px-6 py-3 bg-red-100 text-red-600 rounded-2xl font-medium hover:bg-red-200 transition-colors">
                  Clear All
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
