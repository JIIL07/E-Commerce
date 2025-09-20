'use client'

import { useState, useEffect } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'

export default function OrderSuccessPage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const [paymentIntent, setPaymentIntent] = useState<string | null>(null)

  useEffect(() => {
    const paymentIntentId = searchParams.get('paymentIntent')
    if (paymentIntentId) {
      setPaymentIntent(paymentIntentId)
    }
  }, [searchParams])

  return (
    <div className="max-w-2xl mx-auto p-6">
      <div className="text-center">
        <div className="mb-6">
          <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
            <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Payment Successful!</h1>
          <p className="text-lg text-gray-600">
            Thank you for your purchase. Your order has been confirmed.
          </p>
        </div>

        {paymentIntent && (
          <div className="bg-gray-100 rounded-lg p-4 mb-6">
            <h3 className="font-semibold text-gray-900 mb-2">Payment Details</h3>
            <p className="text-sm text-gray-600">
              Payment Intent ID: <span className="font-mono">{paymentIntent}</span>
            </p>
          </div>
        )}

        <div className="space-y-4">
          <button
            onClick={() => router.push('/')}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700"
          >
            Continue Shopping
          </button>
          <button
            onClick={() => router.push('/orders')}
            className="w-full bg-gray-200 text-gray-800 py-2 px-4 rounded-md hover:bg-gray-300"
          >
            View My Orders
          </button>
        </div>
      </div>
    </div>
  )
}
