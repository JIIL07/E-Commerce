import { Router } from 'express'
import express from 'express'
import Stripe from 'stripe'
import { stripe, createPaymentIntent, retrievePaymentIntent, confirmPaymentIntent, cancelPaymentIntent } from '../lib/stripe'
import { prisma } from '../lib/prisma'
import { authenticateToken, AuthRequest } from '../middleware/auth'

const router = Router()

router.post('/create-payment-intent', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { amount, currency = 'usd', orderId } = req.body

    if (!amount || amount <= 0) {
      return res.status(400).json({ message: 'Valid amount is required' })
    }

    const paymentIntent = await createPaymentIntent(amount, currency)

    if (orderId) {
      await prisma.order.update({
        where: { id: orderId },
        data: { paymentIntent: paymentIntent.id }
      })
    }

    res.json({
      clientSecret: paymentIntent.client_secret,
      paymentIntentId: paymentIntent.id
    })
  } catch (error) {
    console.error('Error creating payment intent:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.get('/payment-intent/:id', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { id } = req.params
    const paymentIntent = await retrievePaymentIntent(id)

    res.json(paymentIntent)
  } catch (error) {
    console.error('Error retrieving payment intent:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.post('/confirm-payment', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { paymentIntentId, orderId } = req.body

    if (!paymentIntentId) {
      return res.status(400).json({ message: 'Payment intent ID is required' })
    }

    const paymentIntent = await confirmPaymentIntent(paymentIntentId)

    if (orderId && paymentIntent.status === 'succeeded') {
      await prisma.order.update({
        where: { id: orderId },
        data: { 
          status: 'PROCESSING',
          paymentIntent: paymentIntent.id
        }
      })
    }

    res.json({
      status: paymentIntent.status,
      paymentIntent
    })
  } catch (error) {
    console.error('Error confirming payment:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.post('/cancel-payment', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { paymentIntentId, orderId } = req.body

    if (!paymentIntentId) {
      return res.status(400).json({ message: 'Payment intent ID is required' })
    }

    const paymentIntent = await cancelPaymentIntent(paymentIntentId)

    if (orderId) {
      await prisma.order.update({
        where: { id: orderId },
        data: { 
          status: 'CANCELLED',
          paymentIntent: paymentIntent.id
        }
      })
    }

    res.json({
      status: paymentIntent.status,
      paymentIntent
    })
  } catch (error) {
    console.error('Error cancelling payment:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.post('/webhook', express.raw({ type: 'application/json' }), async (req, res) => {
  const sig = req.headers['stripe-signature']
  const endpointSecret = process.env.STRIPE_WEBHOOK_SECRET

  if (!endpointSecret) {
    return res.status(400).send('Webhook secret not configured')
  }

  let event

  try {
    event = stripe.webhooks.constructEvent(req.body, sig as string, endpointSecret)
  } catch (err) {
    console.error('Webhook signature verification failed:', err)
    return res.status(400).send('Webhook signature verification failed')
  }

  try {
    switch (event.type) {
      case 'payment_intent.succeeded':
        const paymentIntent = event.data.object as Stripe.PaymentIntent
        await handlePaymentSuccess(paymentIntent)
        break
      
      case 'payment_intent.payment_failed':
        const failedPayment = event.data.object as Stripe.PaymentIntent
        await handlePaymentFailure(failedPayment)
        break
      
      default:
        console.log(`Unhandled event type: ${event.type}`)
    }

    res.json({ received: true })
  } catch (error) {
    console.error('Error processing webhook:', error)
    res.status(500).json({ error: 'Webhook processing failed' })
  }
})

async function handlePaymentSuccess(paymentIntent: Stripe.PaymentIntent) {
  try {
    const order = await prisma.order.findFirst({
      where: { paymentIntent: paymentIntent.id }
    })

    if (order) {
      await prisma.order.update({
        where: { id: order.id },
        data: { status: 'PROCESSING' }
      })
    }
  } catch (error) {
    console.error('Error handling payment success:', error)
  }
}

async function handlePaymentFailure(paymentIntent: Stripe.PaymentIntent) {
  try {
    const order = await prisma.order.findFirst({
      where: { paymentIntent: paymentIntent.id }
    })

    if (order) {
      await prisma.order.update({
        where: { id: order.id },
        data: { status: 'CANCELLED' }
      })
    }
  } catch (error) {
    console.error('Error handling payment failure:', error)
  }
}

export default router
