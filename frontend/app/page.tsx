'use client'

import { useState, useEffect, useCallback, useMemo } from 'react'
import Link from 'next/link'
import Navbar from '../components/Navbar'
import HeroSection from '../components/HeroSection'
import CategoryGrid from '../components/CategoryGrid'
import TestimonialSection from '../components/TestimonialSection'
import FeaturesSection from '../components/FeaturesSection'
import StatsSection from '../components/StatsSection'
import BrandsSection from '../components/BrandsSection'
import BlogSection from '../components/BlogSection'
import FAQSection from '../components/FAQSection'
import NewsletterSection from '../components/NewsletterSection'
import Footer from '../components/Footer'
import WishlistButton from '../components/WishlistButton'
import OptimizedImage from '../components/OptimizedImage'
import { 
  ShoppingCart, 
  Star, 
  Eye, 
  ArrowRight, 
  Facebook, 
  Instagram, 
  Twitter, 
  Mail,
  TrendingUp,
  Users,
  Award,
  Heart,
  Plus
} from 'lucide-react'

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

interface Testimonial {
  id: string
  name: string
  role: string
  company: string
  content: string
  rating: number
  avatar: string
  verified: boolean
  product?: string
}

export default function Home() {
  const [featuredProducts, setFeaturedProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

  const fetchFeaturedProducts = useCallback(async () => {
    try {
      setLoading(true)
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'
      const response = await fetch(`${apiUrl}/api/products/featured`)
      
      if (response.ok) {
        const data = await response.json()
        setFeaturedProducts(data)
      } else {
        // Mock data for development
        const mockProducts: Product[] = [
          {
            id: '1',
            name: 'Premium Wireless Headphones',
            price: 299.99,
            image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
            category: { name: 'Electronics' },
            stock: 15,
            averageRating: 4.8,
            reviewCount: 124
          },
          {
            id: '2',
            name: 'Smart Fitness Watch',
            price: 199.99,
            image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
            category: { name: 'Electronics' },
            stock: 8,
            averageRating: 4.6,
            reviewCount: 89
          },
          {
            id: '3',
            name: 'Organic Cotton T-Shirt',
            price: 29.99,
            image: 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
            category: { name: 'Clothing' },
            stock: 25,
            averageRating: 4.7,
            reviewCount: 67
          },
          {
            id: '4',
            name: 'Minimalist Desk Lamp',
            price: 89.99,
            image: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
            category: { name: 'Home' },
            stock: 12,
            averageRating: 4.9,
            reviewCount: 156
          }
        ]
        setFeaturedProducts(mockProducts)
      }
    } catch (error) {
      console.error('Error fetching featured products:', error)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchFeaturedProducts()
  }, [fetchFeaturedProducts])

  const memoizedProducts = useMemo(() => featuredProducts, [featuredProducts])

  const categories: Category[] = [
    {
      id: 'electronics',
      name: 'Electronics',
      description: 'Latest gadgets and tech accessories',
      image: 'https://images.unsplash.com/photo-1498049794561-7780e7231661?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      productCount: 245,
      isTrending: true,
      rating: 4.8,
      reviewCount: 1250
    },
    {
      id: 'clothing',
      name: 'Fashion',
      description: 'Trendy clothing and accessories',
      image: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      productCount: 189,
      isTrending: true,
      rating: 4.6,
      reviewCount: 890
    },
    {
      id: 'home',
      name: 'Home & Garden',
      description: 'Beautiful home decor and furniture',
      image: 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      productCount: 156,
      rating: 4.7,
      reviewCount: 567
    },
    {
      id: 'sports',
      name: 'Sports & Fitness',
      description: 'Equipment for active lifestyle',
      image: 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      productCount: 98,
      rating: 4.5,
      reviewCount: 234
    },
    {
      id: 'books',
      name: 'Books & Media',
      description: 'Educational and entertainment content',
      image: 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      productCount: 67,
      rating: 4.9,
      reviewCount: 345
    },
    {
      id: 'beauty',
      name: 'Beauty & Health',
      description: 'Personal care and wellness products',
      image: 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      productCount: 134,
      isTrending: true,
      rating: 4.6,
      reviewCount: 678
    }
  ]

  const testimonials: Testimonial[] = [
    {
      id: '1',
      name: 'Sarah Johnson',
      role: 'Marketing Manager',
      company: 'TechCorp',
      content: 'The shopping experience here is absolutely fantastic! Fast shipping, great customer service, and the product quality exceeded my expectations. I\'ve been a loyal customer for over 2 years now.',
      rating: 5,
      avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b786?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80',
      verified: true,
      product: 'Wireless Headphones'
    },
    {
      id: '2',
      name: 'Michael Chen',
      role: 'Software Engineer',
      company: 'StartupXYZ',
      content: 'I love how easy it is to find exactly what I need. The search functionality is amazing and the recommendations are spot-on. Plus, the return policy gives me confidence in my purchases.',
      rating: 5,
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80',
      verified: true,
      product: 'Smart Watch'
    },
    {
      id: '3',
      name: 'Emily Rodriguez',
      role: 'Designer',
      company: 'Creative Studio',
      content: 'The customer support team is incredible! They helped me resolve an issue within minutes and were so friendly and professional. This is how customer service should be done.',
      rating: 5,
      avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80',
      verified: true,
      product: 'Design Software'
    }
  ]

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100">
      <Navbar />
      
      <main>
        {/* Hero Section */}
        <HeroSection />

        {/* Featured Products */}
        <section className="py-20 bg-white/50 backdrop-blur-sm">
          <div className="container mx-auto px-4">
            <div className="text-center mb-16">
              <div className="inline-flex items-center px-4 py-2 bg-blue-100 rounded-full text-sm font-medium text-blue-600 mb-4">
                <Star className="w-4 h-4 mr-2 fill-current" />
                Featured Products
              </div>
              <h2 className="text-responsive-3xl font-bold text-gray-900 mb-4">
                Handpicked for You
              </h2>
              <p className="text-responsive-lg text-gray-600 max-w-2xl mx-auto">
                Discover our carefully curated selection of premium products that our customers love
              </p>
            </div>

            {loading ? (
              <div className="flex justify-center items-center py-16">
                <div className="spinner w-12 h-12"></div>
              </div>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
                {memoizedProducts.map((product) => (
                  <div key={product.id} className="group card-hover">
                    <Link href={`/products/${product.id}`} className="block">
                      <div className="relative aspect-square overflow-hidden">
                        <OptimizedImage
                          src={product.image}
                          alt={product.name}
                          className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
                        />
                        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                        <div className="absolute top-4 left-4 px-3 py-1 bg-white/90 backdrop-blur-sm rounded-full text-sm font-medium text-gray-700 shadow-lg">
                          {product.category.name}
                        </div>
                        <div className="absolute top-4 right-4">
                          <WishlistButton productId={product.id} />
                        </div>
                        {product.stock < 10 && (
                          <div className="absolute bottom-4 left-4 px-3 py-1 bg-red-500 text-white rounded-full text-xs font-medium">
                            Only {product.stock} left
                          </div>
                        )}
                      </div>
                    </Link>
                    <div className="p-6">
                      <h3 className="font-bold text-xl mb-3 line-clamp-2 text-gray-900 group-hover:text-blue-600 transition-colors">
                        <Link href={`/products/${product.id}`}>{product.name}</Link>
                      </h3>
                      <div className="flex items-center mb-4">
                        <div className="flex text-yellow-400">
                          {[...Array(5)].map((_, i) => (
                            <Star 
                              key={i} 
                              size={20}
                              className={i < Math.floor(product.averageRating) ? 'fill-current text-yellow-400' : 'text-gray-300'}
                            />
                          ))}
                        </div>
                        <span className="text-sm text-gray-500 ml-2">({product.reviewCount})</span>
                      </div>
                      <div className="flex items-center justify-between mb-6">
                        <span className="text-3xl font-bold text-gray-900">${product.price.toFixed(2)}</span>
                        <Link 
                          href={`/products/${product.id}`}
                          className="text-blue-600 hover:text-blue-700 transition-colors"
                        >
                          <Eye className="w-6 h-6" />
                        </Link>
                      </div>
                      <div className="flex gap-3">
                        <Link 
                          href={`/products/${product.id}`}
                          className="flex-1 bg-gray-100 text-gray-700 py-3 px-4 rounded-2xl font-medium hover:bg-gray-200 transition-colors text-center"
                        >
                          View Details
                        </Link>
                        <button className="bg-blue-600 text-white p-3 rounded-2xl hover:bg-blue-700 transition-colors">
                          <Plus className="w-6 h-6" />
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}

            <div className="text-center mt-12">
              <Link
                href="/products"
                className="group inline-flex items-center px-8 py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-2xl hover:from-blue-700 hover:to-purple-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-1"
              >
                View All Products
                <ArrowRight className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
              </Link>
            </div>
          </div>
        </section>

        {/* Categories */}
        <CategoryGrid categories={categories} />

        {/* Features */}
        <FeaturesSection />

        {/* Stats */}
        <StatsSection />

        {/* Brands */}
        <BrandsSection />

        {/* Testimonials */}
        <TestimonialSection testimonials={testimonials} />

        {/* Blog */}
        <BlogSection />

        {/* FAQ */}
        <FAQSection />

        {/* Newsletter */}
        <NewsletterSection />
      </main>

      {/* Footer */}
      <Footer />
    </div>
  )
}