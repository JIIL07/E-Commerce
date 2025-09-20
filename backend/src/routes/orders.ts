import { Router } from 'express'
import { prisma } from '../lib/prisma'
import { authenticateToken, requireAdmin, AuthRequest } from '../middleware/auth'
import { OrderWithItems } from '../types'

const router = Router()

router.get('/', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const page = parseInt(req.query.page as string) || 1
    const limit = parseInt(req.query.limit as string) || 10
    const status = req.query.status as string
    const skip = (page - 1) * limit

    const where: any = {}
    
    if (req.user!.role !== 'ADMIN') {
      where.userId = req.user!.id
    }
    
    if (status) {
      where.status = status
    }

    const [orders, total] = await Promise.all([
      prisma.order.findMany({
        where,
        include: {
          orderItems: {
            include: {
              product: {
                include: {
                  category: true
                }
              }
            }
          },
          user: {
            select: {
              id: true,
              name: true,
              email: true
            }
          }
        },
        orderBy: { createdAt: 'desc' },
        skip,
        take: limit
      }),
      prisma.order.count({ where })
    ])

    res.json({
      orders,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit)
      }
    })
  } catch (error) {
    console.error('Error fetching orders:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.get('/:id', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const order = await prisma.order.findUnique({
      where: { id: req.params.id },
      include: {
        orderItems: {
          include: {
            product: {
              include: {
                category: true
              }
            }
          }
        },
        user: {
          select: {
            id: true,
            name: true,
            email: true
          }
        }
      }
    })

    if (!order) {
      return res.status(404).json({ message: 'Order not found' })
    }

    if (req.user!.role !== 'ADMIN' && order.userId !== req.user!.id) {
      return res.status(403).json({ message: 'Access denied' })
    }

    res.json(order)
  } catch (error) {
    console.error('Error fetching order:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.post('/', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { shippingAddress, billingAddress, paymentIntent } = req.body

    if (!shippingAddress || !billingAddress) {
      return res.status(400).json({ message: 'Shipping and billing addresses are required' })
    }

    const cartItems = await prisma.cartItem.findMany({
      where: { userId: req.user!.id },
      include: { product: true }
    })

    if (cartItems.length === 0) {
      return res.status(400).json({ message: 'Cart is empty' })
    }

    const subtotal = cartItems.reduce((sum, item) => sum + (Number(item.product.price) * item.quantity), 0)
    const tax = subtotal * 0.1
    const shipping = subtotal > 100 ? 0 : 10
    const total = subtotal + tax + shipping

    for (const item of cartItems) {
      if (item.product.stock < item.quantity) {
        return res.status(400).json({ 
          message: `Not enough stock for ${item.product.name}` 
        })
      }
    }

    const order = await prisma.$transaction(async (tx) => {
      const newOrder = await tx.order.create({
        data: {
          userId: req.user!.id,
          status: 'PENDING',
          subtotal,
          tax,
          shipping,
          total,
          shippingAddress,
          billingAddress,
          paymentIntent
        }
      })

      const orderItems = []
      for (const item of cartItems) {
        const orderItem = await tx.orderItem.create({
          data: {
            orderId: newOrder.id,
            productId: item.productId,
            quantity: item.quantity,
            price: item.product.price
          }
        })

        await tx.product.update({
          where: { id: item.productId },
          data: {
            stock: {
              decrement: item.quantity
            }
          }
        })

        orderItems.push(orderItem)
      }

      await tx.cartItem.deleteMany({
        where: { userId: req.user!.id }
      })

      return newOrder
    })

    const completeOrder = await prisma.order.findUnique({
      where: { id: order.id },
      include: {
        orderItems: {
          include: {
            product: {
              include: {
                category: true
              }
            }
          }
        },
        user: {
          select: {
            id: true,
            name: true,
            email: true
          }
        }
      }
    })

    res.status(201).json({
      message: 'Order created successfully',
      order: completeOrder
    })
  } catch (error) {
    console.error('Error creating order:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.put('/:id/status', authenticateToken, requireAdmin, async (req: AuthRequest, res) => {
  try {
    const { id } = req.params
    const { status } = req.body

    const validStatuses = ['PENDING', 'PROCESSING', 'SHIPPED', 'DELIVERED', 'CANCELLED', 'REFUNDED']
    
    if (!validStatuses.includes(status)) {
      return res.status(400).json({ 
        message: 'Invalid status. Must be one of: ' + validStatuses.join(', ') 
      })
    }

    const order = await prisma.order.update({
      where: { id },
      data: { status },
      include: {
        orderItems: {
          include: {
            product: true
          }
        },
        user: {
          select: {
            id: true,
            name: true,
            email: true
          }
        }
      }
    })

    res.json({
      message: 'Order status updated',
      order
    })
  } catch (error) {
    console.error('Error updating order status:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.put('/:id/payment', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { id } = req.params
    const { paymentIntent } = req.body

    const order = await prisma.order.findFirst({
      where: {
        id,
        userId: req.user!.id
      }
    })

    if (!order) {
      return res.status(404).json({ message: 'Order not found' })
    }

    const updatedOrder = await prisma.order.update({
      where: { id },
      data: { paymentIntent },
      include: {
        orderItems: {
          include: {
            product: {
              include: {
                category: true
              }
            }
          }
        },
        user: {
          select: {
            id: true,
            name: true,
            email: true
          }
        }
      }
    })

    res.json({
      message: 'Payment intent updated',
      order: updatedOrder
    })
  } catch (error) {
    console.error('Error updating payment intent:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.delete('/:id', authenticateToken, async (req: AuthRequest, res) => {
  try {
    const { id } = req.params

    const order = await prisma.order.findFirst({
      where: {
        id,
        userId: req.user!.id
      },
      include: {
        orderItems: true
      }
    })

    if (!order) {
      return res.status(404).json({ message: 'Order not found' })
    }

    if (order.status !== 'PENDING') {
      return res.status(400).json({ 
        message: 'Only pending orders can be cancelled' 
      })
    }

    await prisma.$transaction(async (tx) => {
      await tx.order.update({
        where: { id },
        data: { status: 'CANCELLED' }
      })

      for (const item of order.orderItems) {
        await tx.product.update({
          where: { id: item.productId },
          data: {
            stock: {
              increment: item.quantity
            }
          }
        })
      }
    })

    res.json({ message: 'Order cancelled successfully' })
  } catch (error) {
    console.error('Error cancelling order:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

export default router
