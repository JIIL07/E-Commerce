import { Router } from 'express'
import { prisma } from '../lib/prisma'
import { authenticateToken, AuthRequest } from '../middleware/auth'
import { ReviewWithUser } from '../types'

const router = Router()

router.get('/product/:productId', async (req, res) => {
  try {
    const page = parseInt(req.query.page as string) || 1
    const limit = parseInt(req.query.limit as string) || 10
    const skip = (page - 1) * limit

    const [reviews, total] = await Promise.all([
      prisma.review.findMany({
        where: { productId: req.params.productId },
        include: {
          user: {
            select: {
              name: true,
              image: true
            }
          }
        },
        orderBy: { createdAt: 'desc' },
        skip,
        take: limit
      }),
      prisma.review.count({
        where: { productId: req.params.productId }
      })
    ])

    // Calculate average rating
    const averageRating = reviews.length > 0 
      ? reviews.reduce((sum, review) => sum + review.rating, 0) / reviews.length
      : 0

    // Count ratings by star
    const ratingCounts = {
      5: reviews.filter(r => r.rating === 5).length,
      4: reviews.filter(r => r.rating === 4).length,
      3: reviews.filter(r => r.rating === 3).length,
      2: reviews.filter(r => r.rating === 2).length,
      1: reviews.filter(r => r.rating === 1).length
    }

    res.json({
      reviews,
      summary: {
        averageRating: Math.round(averageRating * 10) / 10,
        totalReviews: total,
        ratingCounts
      },
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit)
      }
    })
  } catch (error) {
    console.error('Error fetching reviews:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.get('/user', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const page = parseInt(req.query.page as string) || 1
    const limit = parseInt(req.query.limit as string) || 10
    const skip = (page - 1) * limit

    const [reviews, total] = await Promise.all([
      prisma.review.findMany({
        where: { userId: req.user!.id },
        include: {
          product: {
            select: {
              id: true,
              name: true,
              slug: true,
              images: true,
              price: true
            }
          }
        },
        orderBy: { createdAt: 'desc' },
        skip,
        take: limit
      }),
      prisma.review.count({
        where: { userId: req.user!.id }
      })
    ])

    res.json({
      reviews,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit)
      }
    })
  } catch (error) {
    console.error('Error fetching user reviews:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.post('/', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { productId, rating, comment } = req.body

    if (!productId || !rating) {
      return res.status(400).json({ message: 'Product ID and rating are required' })
    }

    if (rating < 1 || rating > 5) {
      return res.status(400).json({ message: 'Rating must be between 1 and 5' })
    }

    // Check if product exists
    const product = await prisma.product.findUnique({
      where: { id: productId }
    })

    if (!product) {
      return res.status(404).json({ message: 'Product not found' })
    }

    // Check if user has already reviewed this product
    const existingReview = await prisma.review.findUnique({
      where: {
        userId_productId: {
          userId: req.user!.id,
          productId
        }
      }
    })

    if (existingReview) {
      return res.status(400).json({ message: 'You have already reviewed this product' })
    }

    // Check if user has purchased this product (optional validation)
    const hasPurchased = await prisma.orderItem.findFirst({
      where: {
        productId,
        order: {
          userId: req.user!.id,
          status: {
            in: ['DELIVERED', 'SHIPPED']
          }
        }
      }
    })

    if (!hasPurchased) {
      return res.status(400).json({ 
        message: 'You can only review products you have purchased' 
      })
    }

    const review = await prisma.review.create({
      data: {
        userId: req.user!.id,
        productId,
        rating,
        comment: comment || null
      },
      include: {
        user: {
          select: {
            name: true,
            image: true
          }
        },
        product: {
          select: {
            id: true,
            name: true,
            slug: true
          }
        }
      }
    })

    res.status(201).json({
      message: 'Review created successfully',
      review
    })
  } catch (error) {
    console.error('Error creating review:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.put('/:id', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { id } = req.params
    const { rating, comment } = req.body

    if (rating && (rating < 1 || rating > 5)) {
      return res.status(400).json({ message: 'Rating must be between 1 and 5' })
    }

    const existingReview = await prisma.review.findFirst({
      where: {
        id,
        userId: req.user!.id
      }
    })

    if (!existingReview) {
      return res.status(404).json({ message: 'Review not found' })
    }

    const updateData: any = {}
    if (rating !== undefined) updateData.rating = rating
    if (comment !== undefined) updateData.comment = comment

    const review = await prisma.review.update({
      where: { id },
      data: updateData,
      include: {
        user: {
          select: {
            name: true,
            image: true
          }
        },
        product: {
          select: {
            id: true,
            name: true,
            slug: true
          }
        }
      }
    })

    res.json({
      message: 'Review updated successfully',
      review
    })
  } catch (error) {
    console.error('Error updating review:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.delete('/:id', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { id } = req.params

    const review = await prisma.review.findFirst({
      where: {
        id,
        userId: req.user!.id
      }
    })

    if (!review) {
      return res.status(404).json({ message: 'Review not found' })
    }

    await prisma.review.delete({
      where: { id }
    })

    res.json({ message: 'Review deleted successfully' })
  } catch (error) {
    console.error('Error deleting review:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.get('/stats/:productId', async (req, res) => {
  try {
    const productId = req.params.productId

    const reviews = await prisma.review.findMany({
      where: { productId },
      select: { rating: true }
    })

    if (reviews.length === 0) {
      return res.json({
        averageRating: 0,
        totalReviews: 0,
        ratingCounts: { 5: 0, 4: 0, 3: 0, 2: 0, 1: 0 }
      })
    }

    const averageRating = reviews.reduce((sum, review) => sum + review.rating, 0) / reviews.length
    const ratingCounts = {
      5: reviews.filter(r => r.rating === 5).length,
      4: reviews.filter(r => r.rating === 4).length,
      3: reviews.filter(r => r.rating === 3).length,
      2: reviews.filter(r => r.rating === 2).length,
      1: reviews.filter(r => r.rating === 1).length
    }

    res.json({
      averageRating: Math.round(averageRating * 10) / 10,
      totalReviews: reviews.length,
      ratingCounts
    })
  } catch (error) {
    console.error('Error fetching review stats:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

export default router
