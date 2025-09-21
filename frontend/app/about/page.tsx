import Link from 'next/link'
import Navbar from '../../components/Navbar'

export default function AboutPage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      
      <main className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold mb-4">About Us</h1>
            <p className="text-xl text-gray-600">
              Your trusted partner in online shopping
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-12 mb-16">
            <div>
              <h2 className="text-2xl font-semibold mb-4">Our Story</h2>
              <p className="text-gray-700 leading-relaxed mb-4">
                Founded in 2024, our e-commerce platform was born from a simple idea: 
                to make online shopping more accessible, secure, and enjoyable for everyone. 
                We believe that technology should enhance the shopping experience, not complicate it.
              </p>
              <p className="text-gray-700 leading-relaxed">
                Today, we're proud to serve thousands of customers worldwide, offering 
                a curated selection of high-quality products at competitive prices, 
                backed by exceptional customer service.
              </p>
            </div>
            
            <div>
              <h2 className="text-2xl font-semibold mb-4">Our Mission</h2>
              <p className="text-gray-700 leading-relaxed mb-4">
                To provide a seamless, secure, and satisfying online shopping experience 
                that connects customers with the products they love, while supporting 
                businesses and communities around the world.
              </p>
              <p className="text-gray-700 leading-relaxed">
                We're committed to innovation, sustainability, and building lasting 
                relationships with our customers and partners.
              </p>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow-md p-8 mb-16">
            <h2 className="text-2xl font-semibold mb-6 text-center">Our Values</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              <div className="text-center">
                <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-blue-600 text-2xl">🔒</span>
                </div>
                <h3 className="text-lg font-semibold mb-2">Security</h3>
                <p className="text-gray-600">
                  Your data and payments are protected with industry-leading security measures.
                </p>
              </div>
              
              <div className="text-center">
                <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-green-600 text-2xl">⭐</span>
                </div>
                <h3 className="text-lg font-semibold mb-2">Quality</h3>
                <p className="text-gray-600">
                  We carefully curate our product selection to ensure the highest standards.
                </p>
              </div>
              
              <div className="text-center">
                <div className="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <span className="text-purple-600 text-2xl">🤝</span>
                </div>
                <h3 className="text-lg font-semibold mb-2">Service</h3>
                <p className="text-gray-600">
                  Our dedicated support team is here to help you every step of the way.
                </p>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-12 mb-16">
            <div>
              <h2 className="text-2xl font-semibold mb-4">Why Choose Us?</h2>
              <ul className="space-y-3">
                <li className="flex items-start">
                  <span className="text-green-500 mr-3 mt-1">✓</span>
                  <span className="text-gray-700">Free shipping on orders over $100</span>
                </li>
                <li className="flex items-start">
                  <span className="text-green-500 mr-3 mt-1">✓</span>
                  <span className="text-gray-700">30-day return policy</span>
                </li>
                <li className="flex items-start">
                  <span className="text-green-500 mr-3 mt-1">✓</span>
                  <span className="text-gray-700">Secure payment processing</span>
                </li>
                <li className="flex items-start">
                  <span className="text-green-500 mr-3 mt-1">✓</span>
                  <span className="text-gray-700">24/7 customer support</span>
                </li>
                <li className="flex items-start">
                  <span className="text-green-500 mr-3 mt-1">✓</span>
                  <span className="text-gray-700">Mobile-friendly platform</span>
                </li>
              </ul>
            </div>
            
            <div>
              <h2 className="text-2xl font-semibold mb-4">Get in Touch</h2>
              <div className="space-y-4">
                <div className="flex items-center">
                  <span className="text-gray-500 mr-3">📧</span>
                  <span className="text-gray-700">support@ecommerce.com</span>
                </div>
                <div className="flex items-center">
                  <span className="text-gray-500 mr-3">📞</span>
                  <span className="text-gray-700">+1 (555) 123-4567</span>
                </div>
                <div className="flex items-center">
                  <span className="text-gray-500 mr-3">📍</span>
                  <span className="text-gray-700">123 Commerce St, City, State 12345</span>
                </div>
                <div className="flex items-center">
                  <span className="text-gray-500 mr-3">🕒</span>
                  <span className="text-gray-700">Mon-Fri: 9AM-6PM EST</span>
                </div>
              </div>
            </div>
          </div>

          <div className="text-center">
            <h2 className="text-2xl font-semibold mb-4">Ready to Start Shopping?</h2>
            <p className="text-gray-600 mb-6">
              Join thousands of satisfied customers and discover amazing products today.
            </p>
            <div className="space-x-4">
              <Link href="/products" className="bg-blue-600 text-white px-8 py-3 rounded-lg hover:bg-blue-700 transition">
                Browse Products
              </Link>
              <Link href="/contact" className="bg-gray-200 text-gray-800 px-8 py-3 rounded-lg hover:bg-gray-300 transition">
                Contact Us
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
