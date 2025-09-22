'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { 
  X, 
  Menu, 
  User, 
  ShoppingCart, 
  Heart, 
  Settings, 
  LogOut,
  Search,
  ChevronDown
} from 'lucide-react'

interface MobileMenuProps {
  isLoggedIn: boolean
  user: any
  cartCount: number
  onLogout: () => void
}

export default function MobileMenu({ isLoggedIn, user, cartCount, onLogout }: MobileMenuProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [expandedSection, setExpandedSection] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const router = useRouter()

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      router.push(`/search?q=${encodeURIComponent(searchQuery.trim())}`)
      setIsOpen(false)
    }
  }

  const toggleSection = (section: string) => {
    setExpandedSection(expandedSection === section ? null : section)
  }

  const closeMenu = () => {
    setIsOpen(false)
    setExpandedSection(null)
  }

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = 'unset'
    }

    return () => {
      document.body.style.overflow = 'unset'
    }
  }, [isOpen])

  return (
    <>
      {/* Menu Button */}
      <button
        onClick={() => setIsOpen(true)}
        className="lg:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors"
        aria-label="Open menu"
      >
        <Menu className="w-6 h-6 text-gray-700" />
      </button>

      {/* Mobile Menu Overlay */}
      {isOpen && (
        <div className="fixed inset-0 z-50 lg:hidden">
          {/* Backdrop */}
          <div 
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            onClick={closeMenu}
          />
          
          {/* Menu Panel */}
          <div className="absolute right-0 top-0 h-full w-80 max-w-[90vw] bg-white shadow-2xl">
            <div className="flex flex-col h-full">
              {/* Header */}
              <div className="flex items-center justify-between p-6 border-b border-gray-200">
                <h2 className="text-xl font-bold text-gray-900">Menu</h2>
                <button
                  onClick={closeMenu}
                  className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
                  aria-label="Close menu"
                >
                  <X className="w-6 h-6 text-gray-700" />
                </button>
              </div>

              {/* Search */}
              <div className="p-6 border-b border-gray-200">
                <form onSubmit={handleSearch} className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                  <input
                    type="text"
                    placeholder="Search products..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </form>
              </div>

              {/* User Section */}
              {isLoggedIn ? (
                <div className="p-6 border-b border-gray-200">
                  <div className="flex items-center space-x-3 mb-4">
                    <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center">
                      <User className="w-6 h-6 text-white" />
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">{user?.name || user?.email}</p>
                      <p className="text-sm text-gray-600">Welcome back!</p>
                    </div>
                  </div>
                  
                  <div className="flex space-x-2">
                    <Link
                      href="/cart"
                      onClick={closeMenu}
                      className="flex-1 flex items-center justify-center space-x-2 bg-blue-600 text-white py-3 px-4 rounded-xl hover:bg-blue-700 transition-colors"
                    >
                      <ShoppingCart className="w-5 h-5" />
                      <span>Cart</span>
                      {cartCount > 0 && (
                        <span className="bg-white text-blue-600 text-xs px-2 py-1 rounded-full font-bold">
                          {cartCount}
                        </span>
                      )}
                    </Link>
                    <Link
                      href="/wishlist"
                      onClick={closeMenu}
                      className="flex-1 flex items-center justify-center space-x-2 bg-gray-100 text-gray-700 py-3 px-4 rounded-xl hover:bg-gray-200 transition-colors"
                    >
                      <Heart className="w-5 h-5" />
                      <span>Wishlist</span>
                    </Link>
                  </div>
                </div>
              ) : (
                <div className="p-6 border-b border-gray-200">
                  <div className="space-y-3">
                    <Link
                      href="/login"
                      onClick={closeMenu}
                      className="block w-full bg-blue-600 text-white py-3 px-4 rounded-xl text-center font-semibold hover:bg-blue-700 transition-colors"
                    >
                      Sign In
                    </Link>
                    <Link
                      href="/register"
                      onClick={closeMenu}
                      className="block w-full bg-gray-100 text-gray-700 py-3 px-4 rounded-xl text-center font-semibold hover:bg-gray-200 transition-colors"
                    >
                      Sign Up
                    </Link>
                  </div>
                </div>
              )}

              {/* Navigation */}
              <div className="flex-1 overflow-y-auto">
                <div className="p-6 space-y-2">
                  {/* Main Navigation */}
                  <Link
                    href="/products"
                    onClick={closeMenu}
                    className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                  >
                    <span className="font-medium text-gray-900">All Products</span>
                    <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                  </Link>

                  <Link
                    href="/categories"
                    onClick={closeMenu}
                    className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                  >
                    <span className="font-medium text-gray-900">Categories</span>
                    <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                  </Link>

                  {isLoggedIn && (
                    <>
                      <Link
                        href="/orders"
                        onClick={closeMenu}
                        className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                      >
                        <span className="font-medium text-gray-900">My Orders</span>
                        <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                      </Link>

                      <Link
                        href="/profile"
                        onClick={closeMenu}
                        className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                      >
                        <span className="font-medium text-gray-900">Profile</span>
                        <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                      </Link>

                      <Link
                        href="/settings"
                        onClick={closeMenu}
                        className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                      >
                        <span className="font-medium text-gray-900">Settings</span>
                        <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                      </Link>
                    </>
                  )}

                  {/* Support Section */}
                  <div className="pt-4 border-t border-gray-200">
                    <Link
                      href="/help"
                      onClick={closeMenu}
                      className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                    >
                      <span className="font-medium text-gray-900">Help Center</span>
                      <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                    </Link>

                    <Link
                      href="/contact"
                      onClick={closeMenu}
                      className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                    >
                      <span className="font-medium text-gray-900">Contact Us</span>
                      <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                    </Link>

                    <Link
                      href="/about"
                      onClick={closeMenu}
                      className="flex items-center justify-between p-4 rounded-xl hover:bg-gray-100 transition-colors"
                    >
                      <span className="font-medium text-gray-900">About Us</span>
                      <ChevronDown className="w-5 h-5 text-gray-400 rotate-[-90deg]" />
                    </Link>
                  </div>
                </div>
              </div>

              {/* Footer */}
              {isLoggedIn && (
                <div className="p-6 border-t border-gray-200">
                  <button
                    onClick={() => {
                      onLogout()
                      closeMenu()
                    }}
                    className="flex items-center space-x-3 w-full p-4 rounded-xl hover:bg-red-50 transition-colors text-red-600"
                  >
                    <LogOut className="w-5 h-5" />
                    <span className="font-medium">Sign Out</span>
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  )
}
