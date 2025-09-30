'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

export default function Header() {
  const pathname = usePathname()

  return (
    <header className="bg-white shadow-sm border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link href="/" className="text-2xl font-bold text-blue-600">
              E-Commerce
            </Link>
          </div>
          
          <nav className="hidden md:flex space-x-8">
            <Link 
              href="/" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                pathname === '/' 
                  ? 'bg-blue-100 text-blue-700' 
                  : 'text-gray-700 hover:text-blue-600'
              }`}
            >
              Home
            </Link>
            <Link 
              href="/products" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                pathname === '/products' 
                  ? 'bg-blue-100 text-blue-700' 
                  : 'text-gray-700 hover:text-blue-600'
              }`}
            >
              Products
            </Link>
            <Link 
              href="/categories" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                pathname === '/categories' 
                  ? 'bg-blue-100 text-blue-700' 
                  : 'text-gray-700 hover:text-blue-600'
              }`}
            >
              Categories
            </Link>
            <Link 
              href="/login" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                pathname === '/login' 
                  ? 'bg-blue-100 text-blue-700' 
                  : 'text-gray-700 hover:text-blue-600'
              }`}
            >
              Login
            </Link>
            <Link 
              href="/register" 
              className={`px-3 py-2 rounded-md text-sm font-medium ${
                pathname === '/register' 
                  ? 'bg-blue-100 text-blue-700' 
                  : 'text-gray-700 hover:text-blue-600'
              }`}
            >
              Register
            </Link>
          </nav>
        </div>
      </div>
    </header>
  )
}