'use client'

import { useState, useEffect } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import Navbar from '../../components/Navbar'
import { 
  CheckCircle, 
  Package, 
  Truck, 
  Mail, 
  ArrowRight, 
  ShoppingBag,
  Calendar,
  CreditCard,
  Download,
  Home
} from 'lucide-react'

export default function OrderSuccessPage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const [paymentIntent, setPaymentIntent] = useState<string | null>(null)
  const [orderId, setOrderId] = useState<string | null>(null)

  useEffect(() => {
    const paymentIntentId = searchParams.get('paymentIntent')
    const orderIdParam = searchParams.get('orderId')
    
    if (paymentIntentId) {
      setPaymentIntent(paymentIntentId)
    }
    if (orderIdParam) {
      setOrderId(orderIdParam)
    }
  }, [searchParams])

  return (
    <div className="min-h-screen bg-gradient-to-br from-green-50 via-blue-50 to-purple-50">
      <Navbar />
      
      <main className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto">
          {/* Success Header */}
          <div className="text-center mb-12">
            <div className="w-32 h-32 bg-gradient-to-r from-green-500 to-emerald-600 rounded-full flex items-center justify-center mx-auto mb-8 shadow-2xl animate-pulse">
              <CheckCircle className="w-20 h-20 text-white" />
            </div>
            <h1 className="text-5xl font-bold mb-4 text-gray-900">
              Order <span className="text-green-600">Successful!</span>
            </h1>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Thank you for your purchase! Your order has been confirmed and will be processed shortly.
            </p>
          </div>

          {/* Order Details Card */}
          <div className="bg-white/80 backdrop-blur-sm rounded-3xl shadow-xl p-8 mb-8 border border-white/20">
            <h2 className="text-2xl font-bold mb-6 text-gray-900">Order Details</h2>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div className="flex items-center">
                  <Package className="w-6 h-6 text-blue-600 mr-3" />
                  <div>
                    <p className="text-sm text-gray-500">Order ID</p>
                    <p className="font-semibold text-gray-900">
                      {orderId || 'ORD-2024-001'}
                    </p>
                  </div>
                </div>
                
                <div className="flex items-center">
                  <Calendar className="w-6 h-6 text-green-600 mr-3" />
                  <div>
                    <p className="text-sm text-gray-500">Order Date</p>
                    <p className="font-semibold text-gray-900">
                      {new Date().toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      })}
                    </p>
                  </div>
                </div>
              </div>
              
              <div className="space-y-4">
                <div className="flex items-center">
                  <CreditCard className="w-6 h-6 text-purple-600 mr-3" />
                  <div>
                    <p className="text-sm text-gray-500">Payment ID</p>
                    <p className="font-semibold text-gray-900 font-mono text-sm">
                      {paymentIntent || 'pi_1234567890'}
                    </p>
                  </div>
                </div>
                
                <div className="flex items-center">
                  <Truck className="w-6 h-6 text-orange-600 mr-3" />
                  <div>
                    <p className="text-sm text-gray-500">Estimated Delivery</p>
                    <p className="font-semibold text-gray-900">
                      {new Date(Date.now() + 3 * 24 * 60 * 60 * 1000).toLocaleDateString('en-US', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric'
                      })}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Next Steps */}
          <div className="bg-white/80 backdrop-blur-sm rounded-3xl shadow-xl p-8 mb-8 border border-white/20">
            <h2 className="text-2xl font-bold mb-6 text-gray-900">What's Next?</h2>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="text-center">
                <div className="w-16 h-16 bg-blue-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <Mail className="w-8 h-8 text-blue-600" />
                </div>
                <h3 className="font-bold text-lg mb-2">Confirmation Email</h3>
                <p className="text-gray-600 text-sm">
                  You'll receive an order confirmation email with all the details within the next few minutes.
                </p>
              </div>
              
              <div className="text-center">
                <div className="w-16 h-16 bg-green-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <Package className="w-8 h-8 text-green-600" />
                </div>
                <h3 className="font-bold text-lg mb-2">Order Processing</h3>
                <p className="text-gray-600 text-sm">
                  We'll prepare your order and send you tracking information once it ships.
                </p>
              </div>
              
              <div className="text-center">
                <div className="w-16 h-16 bg-purple-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <Truck className="w-8 h-8 text-purple-600" />
                </div>
                <h3 className="font-bold text-lg mb-2">Fast Delivery</h3>
                <p className="text-gray-600 text-sm">
                  Your order will be delivered within 3-5 business days with tracking updates.
                </p>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <button
              onClick={() => router.push('/orders')}
              className="group bg-gradient-to-r from-blue-600 to-purple-600 text-white p-6 rounded-2xl hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-bold text-lg mb-2">View My Orders</h3>
                  <p className="text-blue-100 text-sm">Track your order status</p>
                </div>
                <ArrowRight className="w-6 h-6 group-hover:translate-x-1 transition-transform" />
              </div>
            </button>
            
            <button
              onClick={() => router.push('/products')}
              className="group bg-white/80 backdrop-blur-sm text-gray-800 p-6 rounded-2xl hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl border border-white/20"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-bold text-lg mb-2">Continue Shopping</h3>
                  <p className="text-gray-600 text-sm">Discover more products</p>
                </div>
                <ShoppingBag className="w-6 h-6 group-hover:scale-110 transition-transform" />
              </div>
            </button>
            
            <button
              onClick={() => router.push('/')}
              className="group bg-white/80 backdrop-blur-sm text-gray-800 p-6 rounded-2xl hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl border border-white/20"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-bold text-lg mb-2">Back to Home</h3>
                  <p className="text-gray-600 text-sm">Return to homepage</p>
                </div>
                <Home className="w-6 h-6 group-hover:scale-110 transition-transform" />
              </div>
            </button>
          </div>

          {/* Support Section */}
          <div className="mt-12 text-center">
            <div className="bg-gradient-to-r from-blue-600 to-purple-600 rounded-3xl p-8 text-white">
              <h2 className="text-2xl font-bold mb-4">Need Help?</h2>
              <p className="text-blue-100 mb-6">
                Our customer support team is here to help with any questions about your order.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <button
                  onClick={() => router.push('/contact')}
                  className="bg-white/20 backdrop-blur-sm text-white py-3 px-6 rounded-xl font-medium hover:bg-white/30 transition-colors"
                >
                  Contact Support
                </button>
                <button
                  onClick={() => router.push('/help')}
                  className="bg-white/20 backdrop-blur-sm text-white py-3 px-6 rounded-xl font-medium hover:bg-white/30 transition-colors"
                >
                  Help Center
                </button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
