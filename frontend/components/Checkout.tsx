'use client'

import { useState, useEffect } from 'react'
import { loadStripe } from '@stripe/stripe-js'
import { Elements } from '@stripe/react-stripe-js'
import PaymentForm from './PaymentForm'

interface CheckoutProps {
  orderId: string
  amount: number
  onPaymentSuccess?: (paymentIntent: any) => void
}

export default function Checkout({ orderId, amount, onPaymentSuccess }: CheckoutProps) {
  const [stripePromise, setStripePromise] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const initializeStripe = async () => {
      const stripe = await loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY!)
      setStripePromise(stripe)
      setLoading(false)
    }

    initializeStripe()
  }, [])

  const handlePaymentSuccess = (paymentIntent: any) => {
    console.log('Payment succeeded:', paymentIntent)
    onPaymentSuccess?.(paymentIntent)
  }

  const handlePaymentError = (error: any) => {
    console.error('Payment failed:', error)
  }

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-3xl font-bold text-center mb-8">Checkout</h1>
      
      <div className="bg-white rounded-lg shadow-md p-6">
        <div className="mb-6">
          <h2 className="text-xl font-semibold mb-4">Order Summary</h2>
          <div className="space-y-2">
            <div className="flex justify-between">
              <span>Order ID:</span>
              <span className="font-mono">{orderId}</span>
            </div>
            <div className="flex justify-between text-lg font-semibold">
              <span>Total:</span>
              <span>${(amount / 100).toFixed(2)}</span>
            </div>
          </div>
        </div>

        <Elements stripe={stripePromise}>
          <PaymentForm
            amount={amount}
            orderId={orderId}
            onSuccess={handlePaymentSuccess}
            onError={handlePaymentError}
          />
        </Elements>
      </div>
    </div>
  )
}
