import { Router } from 'express'
import { prisma } from '../lib/prisma'
import { authenticateToken, requireAdmin, AuthRequest } from '../middleware/auth'
import { ProductWithRating, ProductWithCategory } from '../types'

const router = Router()

router.get('/', async (req, res) => {
  try {
    const page = parseInt(req.query.page as string) || 1
    const limit = parseInt(req.query.limit as string) || 10
    const category = req.query.category as string
    const search = req.query.search as string
    const featured = req.query.featured === 'true'
    const sortBy = req.query.sortBy as string || 'createdAt'
    const sortOrder = req.query.sortOrder as string || 'desc'

    const skip = (page - 1) * limit

    const where: any = {}
    
    if (category) {
      where.category = { slug: category }
    }
    
    if (search) {
      where.OR = [
        { name: { contains: search, mode: 'insensitive' } },
        { description: { contains: search, mode: 'insensitive' } }
      ]
    }
    
    if (featured) {
      where.featured = true
    }

    const [products, total] = await Promise.all([
      prisma.product.findMany({
        where,
        include: {
          category: true,
          reviews: {
            select: {
              rating: true
            }
          }
        },
        orderBy: { [sortBy]: sortOrder },
        skip,
        take: limit
      }),
      prisma.product.count({ where })
    ])

    const productsWithRating: ProductWithRating[] = products.map((product: any) => ({
      ...product,
      averageRating: product.reviews.length > 0 
        ? product.reviews.reduce((sum: number, review: any) => sum + review.rating, 0) / product.reviews.length
        : 0,
      reviewCount: product.reviews.length
    }))

    res.json({
      products: productsWithRating,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit)
      }
    })
  } catch (error) {
    console.error('Error fetching products:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.get('/:id', async (req, res) => {
  try {
    const product = await prisma.product.findUnique({
      where: { id: req.params.id },
      include: {
        category: true,
        reviews: {
          include: {
            user: {
              select: { name: true, image: true }
            }
          },
          orderBy: { createdAt: 'desc' }
        }
      }
    })

    if (!product) {
      return res.status(404).json({ message: 'Product not found' })
    }

    const averageRating = product.reviews.length > 0 
      ? product.reviews.reduce((sum: number, review: any) => sum + review.rating, 0) / product.reviews.length
      : 0

    res.json({
      ...product,
      averageRating,
      reviewCount: product.reviews.length
    })
  } catch (error) {
    console.error('Error fetching product:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.post('/', authenticateToken, requireAdmin, async (req: AuthRequest, res) => {
  try {
    const { name, description, price, comparePrice, images, categoryId, stock, featured } = req.body

    const slug = name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')

    const product = await prisma.product.create({
      data: {
        name,
        slug,
        description,
        price: parseFloat(price),
        comparePrice: comparePrice ? parseFloat(comparePrice) : null,
        images: Array.isArray(images) ? images : [],
        categoryId,
        stock: parseInt(stock) || 0,
        featured: Boolean(featured)
      },
      include: {
        category: true
      }
    })

    res.status(201).json(product)
  } catch (error) {
    console.error('Error creating product:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.put('/:id', authenticateToken, requireAdmin, async (req: AuthRequest, res) => {
  try {
    const { name, description, price, comparePrice, images, categoryId, stock, featured, inStock } = req.body

    const updateData: any = {
      description,
      price: parseFloat(price),
      comparePrice: comparePrice ? parseFloat(comparePrice) : null,
      images: Array.isArray(images) ? images : [],
      categoryId,
      stock: parseInt(stock) || 0,
      featured: Boolean(featured),
      inStock: Boolean(inStock)
    }

    if (name) {
      updateData.name = name
      updateData.slug = name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
    }

    const product = await prisma.product.update({
      where: { id: req.params.id },
      data: updateData,
      include: {
        category: true
      }
    })

    res.json(product)
  } catch (error) {
    console.error('Error updating product:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

router.delete('/:id', authenticateToken, requireAdmin, async (req: AuthRequest, res) => {
  try {
    await prisma.product.delete({
      where: { id: req.params.id }
    })

    res.json({ message: 'Product deleted successfully' })
  } catch (error) {
    console.error('Error deleting product:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

export default router
