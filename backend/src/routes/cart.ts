import { Router } from 'express'
import { prisma } from '../lib/prisma'
import { authenticateToken, AuthRequest } from '../middleware/auth'
import { CartItemWithProduct } from '../types'

const router = Router()

router.get('/', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const cartItems = await prisma.cartItem.findMany({
      where: { userId: req.user!.id },
      include: {
        product: {
          include: {
            category: true,
            reviews: {
              select: { rating: true }
            }
          }
        }
      },
      orderBy: { id: 'desc' }
    })

    const cartWithTotals = cartItems.map((item: any) => ({
      ...item,
      product: {
        ...item.product,
        averageRating: item.product.reviews.length > 0 
          ? item.product.reviews.reduce((sum: number, review: any) => sum + review.rating, 0) / item.product.reviews.length
          : 0,
        reviewCount: item.product.reviews.length
      },
      totalPrice: Number(item.product.price) * item.quantity
    }))

    const subtotal = cartWithTotals.reduce((sum, item) => sum + item.totalPrice, 0)
    const itemCount = cartWithTotals.reduce((sum, item) => sum + item.quantity, 0)

    res.json({
      items: cartWithTotals,
      summary: {
        itemCount,
        subtotal,
        tax: subtotal * 0.1,
        shipping: subtotal > 100 ? 0 : 10,
        total: subtotal + (subtotal * 0.1) + (subtotal > 100 ? 0 : 10)
      }
    })
  } catch (error) {
    console.error('Error fetching cart:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.post('/', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { productId, quantity = 1 } = req.body

    if (!productId) {
      return res.status(400).json({ message: 'Product ID is required' })
    }

    const product = await prisma.product.findUnique({
      where: { id: productId }
    })

    if (!product) {
      return res.status(404).json({ message: 'Product not found' })
    }

    if (!product.inStock || product.stock < quantity) {
      return res.status(400).json({ message: 'Product is out of stock' })
    }

    const existingItem = await prisma.cartItem.findUnique({
      where: {
        userId_productId: {
          userId: req.user!.id,
          productId
        }
      }
    })

    if (existingItem) {
      const newQuantity = existingItem.quantity + quantity
      
      if (newQuantity > product.stock) {
        return res.status(400).json({ message: 'Not enough stock available' })
      }

      const updatedItem = await prisma.cartItem.update({
        where: { id: existingItem.id },
        data: { quantity: newQuantity },
        include: {
          product: {
            include: {
              category: true,
              reviews: {
                select: { rating: true }
              }
            }
          }
        }
      })

      res.json({
        message: 'Cart item updated',
        item: {
          ...updatedItem,
          product: {
            ...updatedItem.product,
            averageRating: updatedItem.product.reviews.length > 0 
              ? updatedItem.product.reviews.reduce((sum: number, review: any) => sum + review.rating, 0) / updatedItem.product.reviews.length
              : 0,
            reviewCount: updatedItem.product.reviews.length
          },
          totalPrice: Number(updatedItem.product.price) * updatedItem.quantity
        }
      })
    } else {
      const newItem = await prisma.cartItem.create({
        data: {
          userId: req.user!.id,
          productId,
          quantity
        },
        include: {
          product: {
            include: {
              category: true,
              reviews: {
                select: { rating: true }
              }
            }
          }
        }
      })

      res.status(201).json({
        message: 'Item added to cart',
        item: {
          ...newItem,
          product: {
            ...newItem.product,
            averageRating: newItem.product.reviews.length > 0 
              ? newItem.product.reviews.reduce((sum: number, review: any) => sum + review.rating, 0) / newItem.product.reviews.length
              : 0,
            reviewCount: newItem.product.reviews.length
          },
          totalPrice: Number(newItem.product.price) * newItem.quantity
        }
      })
    }
  } catch (error) {
    console.error('Error adding to cart:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.put('/:itemId', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { itemId } = req.params
    const { quantity } = req.body

    if (!quantity || quantity < 1) {
      return res.status(400).json({ message: 'Valid quantity is required' })
    }

    const cartItem = await prisma.cartItem.findFirst({
      where: {
        id: itemId,
        userId: req.user!.id
      },
      include: { product: true }
    })

    if (!cartItem) {
      return res.status(404).json({ message: 'Cart item not found' })
    }

    if (quantity > cartItem.product.stock) {
      return res.status(400).json({ message: 'Not enough stock available' })
    }

    const updatedItem = await prisma.cartItem.update({
      where: { id: itemId },
      data: { quantity },
      include: {
        product: {
          include: {
            category: true,
            reviews: {
              select: { rating: true }
            }
          }
        }
      }
    })

    res.json({
      message: 'Cart item updated',
      item: {
        ...updatedItem,
        product: {
          ...updatedItem.product,
          averageRating: updatedItem.product.reviews.length > 0 
            ? updatedItem.product.reviews.reduce((sum: number, review: any) => sum + review.rating, 0) / updatedItem.product.reviews.length
            : 0,
          reviewCount: updatedItem.product.reviews.length
        },
        totalPrice: Number(updatedItem.product.price) * updatedItem.quantity
      }
    })
  } catch (error) {
    console.error('Error updating cart item:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.delete('/:itemId', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { itemId } = req.params

    const cartItem = await prisma.cartItem.findFirst({
      where: {
        id: itemId,
        userId: req.user!.id
      }
    })

    if (!cartItem) {
      return res.status(404).json({ message: 'Cart item not found' })
    }

    await prisma.cartItem.delete({
      where: { id: itemId }
    })

    res.json({ message: 'Item removed from cart' })
  } catch (error) {
    console.error('Error removing cart item:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.delete('/', authenticateToken, async (req: AuthRequest, res) => {
  try {
    await prisma.cartItem.deleteMany({
      where: { userId: req.user!.id }
    })

    res.json({ message: 'Cart cleared' })
  } catch (error) {
    console.error('Error clearing cart:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

export default router
