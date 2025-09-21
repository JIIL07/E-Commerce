const { PrismaClient } = require('@prisma/client')

const prisma = new PrismaClient()

async function main() {
  console.log('🌱 Starting database seeding...')

  // Create categories
  const electronics = await prisma.category.upsert({
    where: { slug: 'electronics' },
    update: {},
    create: {
      name: 'Electronics',
      slug: 'electronics',
      description: 'Latest electronics and gadgets',
      image: 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=500'
    }
  })

  const clothing = await prisma.category.upsert({
    where: { slug: 'clothing' },
    update: {},
    create: {
      name: 'Clothing',
      slug: 'clothing',
      description: 'Fashion and apparel for everyone',
      image: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=500'
    }
  })

  const books = await prisma.category.upsert({
    where: { slug: 'books' },
    update: {},
    create: {
      name: 'Books',
      slug: 'books',
      description: 'Books and literature',
      image: 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=500'
    }
  })

  const home = await prisma.category.upsert({
    where: { slug: 'home' },
    update: {},
    create: {
      name: 'Home & Garden',
      slug: 'home',
      description: 'Home decor and garden supplies',
      image: 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500'
    }
  })

  console.log('✅ Categories created')

  // Create sample products
  const products = [
    {
      name: 'iPhone 15 Pro',
      description: 'Latest iPhone with advanced camera system and A17 Pro chip',
      price: 999.99,
      stock: 50,
      featured: true,
      categoryId: electronics.id,
      images: [
        'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=500',
        'https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=500'
      ]
    },
    {
      name: 'MacBook Pro 16"',
      description: 'Powerful laptop for professionals with M3 Max chip',
      price: 2499.99,
      stock: 25,
      featured: true,
      categoryId: electronics.id,
      images: [
        'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=500',
        'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=500'
      ]
    },
    {
      name: 'Samsung Galaxy S24',
      description: 'Premium Android smartphone with AI features',
      price: 899.99,
      stock: 40,
      featured: true,
      categoryId: electronics.id,
      images: [
        'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500',
        'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500'
      ]
    },
    {
      name: 'Designer T-Shirt',
      description: 'Premium cotton t-shirt with modern design',
      price: 49.99,
      stock: 100,
      featured: true,
      categoryId: clothing.id,
      images: [
        'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500',
        'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500'
      ]
    },
    {
      name: 'Classic Jeans',
      description: 'Comfortable and stylish denim jeans',
      price: 79.99,
      stock: 75,
      featured: true,
      categoryId: clothing.id,
      images: [
        'https://images.unsplash.com/photo-1542272604-787c3835535d?w=500',
        'https://images.unsplash.com/photo-1542272604-787c3835535d?w=500'
      ]
    },
    {
      name: 'JavaScript: The Good Parts',
      description: 'Essential guide to JavaScript programming',
      price: 29.99,
      stock: 200,
      featured: true,
      categoryId: books.id,
      images: [
        'https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?w=500',
        'https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?w=500'
      ]
    },
    {
      name: 'Modern Furniture Set',
      description: 'Contemporary living room furniture set',
      price: 1299.99,
      stock: 15,
      featured: true,
      categoryId: home.id,
      images: [
        'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500',
        'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=500'
      ]
    },
    {
      name: 'Garden Tools Kit',
      description: 'Complete set of professional garden tools',
      price: 89.99,
      stock: 30,
      featured: true,
      categoryId: home.id,
      images: [
        'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=500',
        'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=500'
      ]
    }
  ]

  for (const product of products) {
    const slug = product.name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
    await prisma.product.upsert({
      where: { slug: slug },
      update: {},
      create: {
        ...product,
        slug: slug
      }
    })
  }

  console.log('✅ Products created')

  // Create sample user
  const user = await prisma.user.upsert({
    where: { email: 'demo@example.com' },
    update: {},
    create: {
      name: 'Demo User',
      email: 'demo@example.com',
      password: '$2a$10$K7L1OJ45/4Y2nIvhRVpCe.FSmhDdWoXehVzJByJp8Yz59dK3k2K2C', // password: password
      role: 'USER'
    }
  })

  console.log('✅ Demo user created')

  console.log('🎉 Database seeding completed!')
}

main()
  .catch((e) => {
    console.error('❌ Error during seeding:', e)
  })
  .then(async () => {
    await prisma.$disconnect()
  })
