'use client'

import { useState, useEffect } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import Checkout from '../../components/Checkout'

export default function CheckoutPage() {
  const searchParams = useSearchParams()
  const router = useRouter()
  const [orderData, setOrderData] = useState<{
    orderId: string
    amount: number
  } | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const orderId = searchParams.get('orderId')
    const amount = searchParams.get('amount')

    if (!orderId || !amount) {
      setError('Missing order information')
      setLoading(false)
      return
    }

    setOrderData({
      orderId,
      amount: parseInt(amount)
    })
    setLoading(false)
  }, [searchParams])

  const handlePaymentSuccess = (paymentIntent: any) => {
    router.push(`/order-success?paymentIntent=${paymentIntent.id}`)
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error || !orderData) {
    return (
      <div className="max-w-2xl mx-auto p-6">
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error || 'Invalid order data'}
        </div>
        <button
          onClick={() => router.push('/')}
          className="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Return to Home
        </button>
      </div>
    )
  }

  return (
    <Checkout
      orderId={orderData.orderId}
      amount={orderData.amount}
      onPaymentSuccess={handlePaymentSuccess}
    />
  )
}
