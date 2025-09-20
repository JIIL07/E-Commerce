import { Router } from 'express'
import { prisma } from '../lib/prisma'
import { authenticateToken, requireAdmin, AuthRequest } from '../middleware/auth'
import { ProductWithRating, CategoryWithProducts } from '../types'

const router = Router()

// GET /api/categories - Get all categories
router.get('/', async (req, res) => {
  try {
    const categories = await prisma.category.findMany({
      include: {
        _count: {
          select: { products: true }
        }
      },
      orderBy: { name: 'asc' }
    })

    res.json(categories)
  } catch (error) {
    console.error('Error fetching categories:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

// GET /api/categories/:slug - Get single category with products
router.get('/:slug', async (req, res) => {
  try {
    const category = await prisma.category.findUnique({
      where: { slug: req.params.slug },
      include: {
        products: {
          include: {
            reviews: {
              select: { rating: true }
            }
          }
        },
        _count: {
          select: { products: true }
        }
      }
    })

    if (!category) {
      return res.status(404).json({ message: 'Category not found' })
    }

    // Calculate average rating for each product
    const productsWithRating: ProductWithRating[] = category.products.map((product: any) => ({
      ...product,
      averageRating: product.reviews.length > 0 
        ? product.reviews.reduce((sum: number, review: any) => sum + review.rating, 0) / product.reviews.length
        : 0,
      reviewCount: product.reviews.length
    }))

    res.json({
      ...category,
      products: productsWithRating
    })
  } catch (error) {
    console.error('Error fetching category:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

// POST /api/categories - Create new category (Admin only)
router.post('/', authenticateToken, requireAdmin, async (req: AuthRequest, res) => {
  try {
    const { name, description, image } = req.body

    // Generate slug from name
    const slug = name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')

    const category = await prisma.category.create({
      data: {
        name,
        slug,
        description,
        image
      }
    })

    res.status(201).json(category)
  } catch (error) {
    console.error('Error creating category:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

// PUT /api/categories/:id - Update category (Admin only)
router.put('/:id', authenticateToken, requireAdmin, async (req: AuthRequest, res) => {
  try {
    const { name, description, image } = req.body

    const updateData: any = {
      description,
      image
    }

    if (name) {
      updateData.name = name
      updateData.slug = name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
    }

    const category = await prisma.category.update({
      where: { id: req.params.id },
      data: updateData
    })

    res.json(category)
  } catch (error) {
    console.error('Error updating category:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

// DELETE /api/categories/:id - Delete category (Admin only)
router.delete('/:id', authenticateToken, requireAdmin, async (req: AuthRequest, res) => {
  try {
    // Check if category has products
    const productCount = await prisma.product.count({
      where: { categoryId: req.params.id }
    })

    if (productCount > 0) {
      return res.status(400).json({ 
        message: 'Cannot delete category with existing products' 
      })
    }

    await prisma.category.delete({
      where: { id: req.params.id }
    })

    res.json({ message: 'Category deleted successfully' })
  } catch (error) {
    console.error('Error deleting category:', error)
    res.status(500).json({ message: 'Internal server error' })
  }
})

export default router
