import express from 'express'
import cors from 'cors'
import helmet from 'helmet'
import morgan from 'morgan'
import dotenv from 'dotenv'
import productsRouter from './routes/products'
import categoriesRouter from './routes/categories'
import authRouter from './routes/auth'
import cartRouter from './routes/cart'
import ordersRouter from './routes/orders'
import reviewsRouter from './routes/reviews'
import paymentsRouter from './routes/payments'
import wishlistRouter from './routes/wishlist'

dotenv.config()

const app = express()
const PORT = process.env.PORT || 5000

app.use(helmet())
app.use(cors({
  origin: process.env.FRONTEND_URL || 'http://localhost:3000',
  credentials: true
}))
app.use(morgan('combined'))
app.use(express.json({ limit: '10mb' }))
app.use(express.urlencoded({ extended: true, limit: '10mb' }))

app.get('/', (req, res) => {
  res.json({
    message: 'E-Commerce API Server',
    version: '1.0.0',
    status: 'running',
    endpoints: {
      products: '/api/products',
      categories: '/api/categories',
      auth: '/api/auth',
      cart: '/api/cart',
      orders: '/api/orders',
      reviews: '/api/reviews',
      payments: '/api/payments',
      health: '/api/health'
    }
  })
})

app.get('/api/health', (req, res) => {
  res.json({
    status: 'healthy',
    timestamp: new Date().toISOString(),
    uptime: process.uptime()
  })
})

app.use('/api/auth', authRouter)
app.use('/api/products', productsRouter)
app.use('/api/categories', categoriesRouter)
app.use('/api/cart', cartRouter)
app.use('/api/orders', ordersRouter)
app.use('/api/reviews', reviewsRouter)
app.use('/api/payments', paymentsRouter)
app.use('/api/wishlist', wishlistRouter)

app.use((err: any, req: express.Request, res: express.Response, next: express.NextFunction) => {
  console.error(err.stack)
  res.status(500).json({
    message: 'Something went wrong!',
    error: process.env.NODE_ENV === 'development' ? err.message : {}
  })
})

app.use('*', (req, res) => {
  res.status(404).json({
    message: 'Route not found',
    path: req.originalUrl
  })
})

app.listen(PORT, () => {
  console.log(`🚀 Server running on port ${PORT}`)
  console.log(`📊 Health check: http://localhost:${PORT}/api/health`)
  console.log(`🔐 Auth API: http://localhost:${PORT}/api/auth`)
  console.log(`🛍️  Products API: http://localhost:${PORT}/api/products`)
  console.log(`📂 Categories API: http://localhost:${PORT}/api/categories`)
  console.log(`🛒 Cart API: http://localhost:${PORT}/api/cart`)
  console.log(`📦 Orders API: http://localhost:${PORT}/api/orders`)
  console.log(`⭐ Reviews API: http://localhost:${PORT}/api/reviews`)
  console.log(`💳 Payments API: http://localhost:${PORT}/api/payments`)
})
