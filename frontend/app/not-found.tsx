'use client'

import Link from 'next/link'
import Navbar from '../components/Navbar'
import { Home, Search, ShoppingBag, HelpCircle, ArrowLeft } from 'lucide-react'

export default function NotFound() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50 to-indigo-100">
      <Navbar />
      
      <main className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center">
          {/* 404 Animation */}
          <div className="mb-12">
            <div className="relative">
              <div className="text-9xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent mb-4 animate-pulse">
                404
              </div>
              <div className="absolute -top-4 -right-4 w-8 h-8 bg-yellow-400 rounded-full animate-bounce"></div>
              <div className="absolute -bottom-4 -left-4 w-6 h-6 bg-pink-400 rounded-full animate-bounce delay-100"></div>
            </div>
            <h1 className="text-5xl font-bold text-gray-900 mb-6">Oops! Page Not Found</h1>
            <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
              The page you're looking for seems to have wandered off into the digital void. 
              Don't worry, even the best explorers sometimes take a wrong turn!
            </p>
          </div>

          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
            <Link 
              href="/" 
              className="group bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-4 rounded-2xl font-semibold hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl flex items-center justify-center"
            >
              <Home className="w-5 h-5 mr-2 group-hover:scale-110 transition-transform" />
              Go Home
            </Link>
            <Link 
              href="/products" 
              className="group bg-white/80 backdrop-blur-sm text-gray-800 px-8 py-4 rounded-2xl font-semibold hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl border border-white/20 flex items-center justify-center"
            >
              <ShoppingBag className="w-5 h-5 mr-2 group-hover:scale-110 transition-transform" />
              Browse Products
            </Link>
            <Link 
              href="/help" 
              className="group bg-white/80 backdrop-blur-sm text-gray-800 px-8 py-4 rounded-2xl font-semibold hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl border border-white/20 flex items-center justify-center"
            >
              <HelpCircle className="w-5 h-5 mr-2 group-hover:scale-110 transition-transform" />
              Get Help
            </Link>
          </div>

          {/* Popular Pages Grid */}
          <div className="bg-white/80 backdrop-blur-sm rounded-3xl shadow-xl p-8 border border-white/20">
            <h2 className="text-2xl font-bold mb-8 text-gray-900">Popular Destinations</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
              <Link 
                href="/products" 
                className="group p-4 rounded-2xl hover:bg-blue-50 transition-colors border border-transparent hover:border-blue-200"
              >
                <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform">
                  <ShoppingBag className="w-6 h-6 text-blue-600" />
                </div>
                <h3 className="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">Products</h3>
                <p className="text-sm text-gray-600 mt-1">Browse our catalog</p>
              </Link>

              <Link 
                href="/categories" 
                className="group p-4 rounded-2xl hover:bg-green-50 transition-colors border border-transparent hover:border-green-200"
              >
                <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform">
                  <Search className="w-6 h-6 text-green-600" />
                </div>
                <h3 className="font-semibold text-gray-900 group-hover:text-green-600 transition-colors">Categories</h3>
                <p className="text-sm text-gray-600 mt-1">Shop by category</p>
              </Link>

              <Link 
                href="/about" 
                className="group p-4 rounded-2xl hover:bg-purple-50 transition-colors border border-transparent hover:border-purple-200"
              >
                <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform">
                  <Home className="w-6 h-6 text-purple-600" />
                </div>
                <h3 className="font-semibold text-gray-900 group-hover:text-purple-600 transition-colors">About Us</h3>
                <p className="text-sm text-gray-600 mt-1">Learn our story</p>
              </Link>

              <Link 
                href="/contact" 
                className="group p-4 rounded-2xl hover:bg-orange-50 transition-colors border border-transparent hover:border-orange-200"
              >
                <div className="w-12 h-12 bg-orange-100 rounded-xl flex items-center justify-center mx-auto mb-3 group-hover:scale-110 transition-transform">
                  <HelpCircle className="w-6 h-6 text-orange-600" />
                </div>
                <h3 className="font-semibold text-gray-900 group-hover:text-orange-600 transition-colors">Contact</h3>
                <p className="text-sm text-gray-600 mt-1">Get in touch</p>
              </Link>
            </div>
          </div>

          {/* Back Button */}
          <div className="mt-8">
            <button 
              onClick={() => window.history.back()}
              className="inline-flex items-center text-gray-600 hover:text-gray-900 transition-colors"
            >
              <ArrowLeft className="w-4 h-4 mr-2" />
              Go Back
            </button>
          </div>
        </div>
      </main>
    </div>
  )
}
