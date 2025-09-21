import Link from 'next/link'
import Navbar from '../components/Navbar'

export default function NotFound() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <main className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto text-center">
          <div className="mb-8">
            <div className="text-9xl font-bold text-blue-600 mb-4">404</div>
            <h1 className="text-4xl font-bold text-gray-900 mb-4">Page Not Found</h1>
            <p className="text-xl text-gray-600 mb-8">
              Sorry, we couldn't find the page you're looking for. 
              It might have been moved, deleted, or you entered the wrong URL.
            </p>
          </div>

          <div className="space-y-4">
            <Link 
              href="/" 
              className="inline-block bg-blue-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-blue-700 transition"
            >
              Go Home
            </Link>
            <div className="text-center">
              <Link 
                href="/products" 
                className="text-blue-600 hover:text-blue-700 font-medium"
              >
                Browse Products
              </Link>
              <span className="mx-2 text-gray-400">•</span>
              <Link 
                href="/contact" 
                className="text-blue-600 hover:text-blue-700 font-medium"
              >
                Contact Support
              </Link>
            </div>
          </div>

          <div className="mt-12">
            <h2 className="text-lg font-semibold mb-4">Popular Pages</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <Link href="/products" className="text-blue-600 hover:text-blue-700 text-sm">
                Products
              </Link>
              <Link href="/categories" className="text-blue-600 hover:text-blue-700 text-sm">
                Categories
              </Link>
              <Link href="/about" className="text-blue-600 hover:text-blue-700 text-sm">
                About Us
              </Link>
              <Link href="/contact" className="text-blue-600 hover:text-blue-700 text-sm">
                Contact
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
