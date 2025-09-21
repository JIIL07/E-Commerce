import { Router } from 'express'
import { PrismaClient } from '../../../database/node_modules/@prisma/client'
import { authenticateToken } from '../middleware/auth'

const router = Router()
const prisma = new PrismaClient()

// Get user's wishlist
router.get('/', authenticateToken, async (req, res) => {
  try {
    const userId = (req as any).user.id

    const wishlistItems = await prisma.wishlistItem.findMany({
      where: { userId },
      include: {
        product: {
          include: {
            category: true,
            reviews: {
              select: {
                rating: true
              }
            }
          }
        }
      },
      orderBy: { addedAt: 'desc' }
    })

    const itemsWithRating = wishlistItems.map(item => ({
      ...item,
      product: {
        ...item.product,
        averageRating: item.product.reviews.length > 0 
          ? item.product.reviews.reduce((sum, review) => sum + review.rating, 0) / item.product.reviews.length
          : 0,
        reviewCount: item.product.reviews.length
      }
    }))

    res.json({ items: itemsWithRating })
  } catch (error) {
    console.error('Error fetching wishlist:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

// Add item to wishlist
router.post('/', authenticateToken, async (req, res) => {
  try {
    const userId = (req as any).user.id
    const { productId } = req.body

    if (!productId) {
      return res.status(400).json({ message: 'Product ID is required' })
    }

    // Check if product exists
    const product = await prisma.product.findUnique({
      where: { id: productId }
    })

    if (!product) {
      return res.status(404).json({ message: 'Product not found' })
    }

    // Check if already in wishlist
    const existingItem = await prisma.wishlistItem.findUnique({
      where: {
        userId_productId: {
          userId,
          productId
        }
      }
    })

    if (existingItem) {
      return res.status(400).json({ message: 'Product already in wishlist' })
    }

    const wishlistItem = await prisma.wishlistItem.create({
      data: {
        userId,
        productId
      },
      include: {
        product: {
          include: {
            category: true
          }
        }
      }
    })

    res.status(201).json({ item: wishlistItem })
  } catch (error) {
    console.error('Error adding to wishlist:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

// Remove item from wishlist
router.delete('/:productId', authenticateToken, async (req, res) => {
  try {
    const userId = (req as any).user.id
    const { productId } = req.params

    const wishlistItem = await prisma.wishlistItem.findUnique({
      where: {
        userId_productId: {
          userId,
          productId
        }
      }
    })

    if (!wishlistItem) {
      return res.status(404).json({ message: 'Item not found in wishlist' })
    }

    await prisma.wishlistItem.delete({
      where: {
        userId_productId: {
          userId,
          productId
        }
      }
    })

    res.json({ message: 'Item removed from wishlist' })
  } catch (error) {
    console.error('Error removing from wishlist:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

// Check if product is in wishlist
router.get('/check/:productId', authenticateToken, async (req, res) => {
  try {
    const userId = (req as any).user.id
    const { productId } = req.params

    const wishlistItem = await prisma.wishlistItem.findUnique({
      where: {
        userId_productId: {
          userId,
          productId
        }
      }
    })

    res.json({ inWishlist: !!wishlistItem })
  } catch (error) {
    console.error('Error checking wishlist:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

export default router
