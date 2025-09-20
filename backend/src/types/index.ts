// Base types for API responses
export interface ProductWithRating {
  id: string
  name: string
  slug: string
  description: string | null
  price: number
  comparePrice: number | null
  images: string[]
  inStock: boolean
  stock: number
  featured: boolean
  categoryId: string
  createdAt: Date
  updatedAt: Date
  averageRating: number
  reviewCount: number
  category?: {
    id: string
    name: string
    slug: string
    description: string | null
    image: string | null
  }
  reviews?: Array<{
    id: string
    rating: number
    comment: string | null
    createdAt: Date
  }>
}

export interface CategoryWithProducts {
  id: string
  name: string
  slug: string
  description: string | null
  image: string | null
  createdAt: Date
  updatedAt: Date
  products: ProductWithRating[]
  _count: {
    products: number
  }
}

export interface ProductWithCategory {
  id: string
  name: string
  slug: string
  description: string | null
  price: number
  comparePrice: number | null
  images: string[]
  inStock: boolean
  stock: number
  featured: boolean
  categoryId: string
  createdAt: Date
  updatedAt: Date
  category: {
    id: string
    name: string
    slug: string
    description: string | null
    image: string | null
  }
  reviews: Array<{
    id: string
    rating: number
    comment: string | null
    createdAt: Date
  }>
}

export interface ReviewWithUser {
  id: string
  userId: string
  productId: string
  rating: number
  comment: string | null
  createdAt: Date
  updatedAt: Date
  user: {
    name: string | null
    image: string | null
  }
}

export interface OrderWithItems {
  id: string
  userId: string
  status: string
  total: number
  subtotal: number
  tax: number
  shipping: number
  shippingAddress: any
  billingAddress: any
  paymentIntent: string | null
  createdAt: Date
  updatedAt: Date
  orderItems: Array<{
    id: string
    orderId: string
    productId: string
    quantity: number
    price: number
    product: ProductWithRating
  }>
  user: {
    id: string
    email: string
    name: string | null
    image: string | null
  }
}

export interface CartItemWithProduct {
  id: string
  userId: string
  productId: string
  quantity: number
  product: ProductWithRating
}

export interface ApiResponse<T> {
  data?: T
  message: string
  error?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    pages: number
  }
}
