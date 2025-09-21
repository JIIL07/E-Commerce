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

  const sports = await prisma.category.upsert({
    where: { slug: 'sports' },
    update: {},
    create: {
      name: 'Sports & Fitness',
      slug: 'sports',
      description: 'Sports equipment and fitness gear',
      image: 'https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=500'
    }
  })

  const beauty = await prisma.category.upsert({
    where: { slug: 'beauty' },
    update: {},
    create: {
      name: 'Beauty & Health',
      slug: 'beauty',
      description: 'Beauty products and health supplements',
      image: 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=500'
    }
  })

  console.log('✅ Categories created')

  // Create products
  const products = [
    // Electronics
    {
      name: 'Premium Wireless Headphones',
      slug: 'premium-wireless-headphones',
      description: 'High-quality wireless headphones with noise cancellation and 30-hour battery life',
      price: 299.99,
      comparePrice: 399.99,
      images: ['https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500'],
      stock: 15,
      featured: true,
      categoryId: electronics.id
    },
    {
      name: 'Smart Fitness Watch',
      slug: 'smart-fitness-watch',
      description: 'Advanced fitness tracking with heart rate monitor and GPS',
      price: 199.99,
      comparePrice: 249.99,
      images: ['https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=500'],
      stock: 8,
      featured: true,
      categoryId: electronics.id
    },
    {
      name: '4K Ultra HD TV',
      slug: '4k-ultra-hd-tv',
      description: '55-inch 4K Smart TV with HDR and built-in streaming apps',
      price: 899.99,
      comparePrice: 1199.99,
      images: ['https://images.unsplash.com/photo-1593359677879-a4bb92f829d1?w=500'],
      stock: 5,
      featured: true,
      categoryId: electronics.id
    },
    {
      name: 'Gaming Laptop',
      slug: 'gaming-laptop',
      description: 'High-performance gaming laptop with RTX graphics and RGB keyboard',
      price: 1299.99,
      comparePrice: 1599.99,
      images: ['https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=500'],
      stock: 3,
      featured: true,
      categoryId: electronics.id
    },
    {
      name: 'Wireless Charging Pad',
      slug: 'wireless-charging-pad',
      description: 'Fast wireless charging pad compatible with all Qi-enabled devices',
      price: 49.99,
      comparePrice: 69.99,
      images: ['https://images.unsplash.com/photo-1609592807905-4b0b1b5b5b5b?w=500'],
      stock: 25,
      featured: false,
      categoryId: electronics.id
    },
    {
      name: 'Bluetooth Speaker',
      slug: 'bluetooth-speaker',
      description: 'Portable Bluetooth speaker with 360-degree sound and waterproof design',
      price: 79.99,
      comparePrice: 99.99,
      images: ['https://images.unsplash.com/photo-1608043152269-423dbba4e7e1?w=500'],
      stock: 20,
      featured: false,
      categoryId: electronics.id
    },

    // Clothing
    {
      name: 'Organic Cotton T-Shirt',
      slug: 'organic-cotton-t-shirt',
      description: 'Comfortable organic cotton t-shirt in various colors',
      price: 29.99,
      comparePrice: 39.99,
      images: ['https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=500'],
      stock: 25,
      featured: true,
      categoryId: clothing.id
    },
    {
      name: 'Designer Jeans',
      slug: 'designer-jeans',
      description: 'Premium denim jeans with perfect fit and modern style',
      price: 89.99,
      comparePrice: 119.99,
      images: ['https://images.unsplash.com/photo-1542272604-787c3835535d?w=500'],
      stock: 12,
      featured: true,
      categoryId: clothing.id
    },
    {
      name: 'Winter Jacket',
      slug: 'winter-jacket',
      description: 'Warm and stylish winter jacket with water-resistant material',
      price: 149.99,
      comparePrice: 199.99,
      images: ['https://images.unsplash.com/photo-1551028719-00167b16eac5?w=500'],
      stock: 8,
      featured: false,
      categoryId: clothing.id
    },
    {
      name: 'Running Shoes',
      slug: 'running-shoes',
      description: 'Lightweight running shoes with advanced cushioning technology',
      price: 129.99,
      comparePrice: 159.99,
      images: ['https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=500'],
      stock: 15,
      featured: false,
      categoryId: clothing.id
    },

    // Books
    {
      name: 'Programming Fundamentals',
      slug: 'programming-fundamentals',
      description: 'Complete guide to programming concepts and best practices',
      price: 49.99,
      comparePrice: 69.99,
      images: ['https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=500'],
      stock: 30,
      featured: true,
      categoryId: books.id
    },
    {
      name: 'Business Strategy Guide',
      slug: 'business-strategy-guide',
      description: 'Essential strategies for modern business success',
      price: 39.99,
      comparePrice: 49.99,
      images: ['https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=500'],
      stock: 20,
      featured: false,
      categoryId: books.id
    },
    {
      name: 'Cookbook Collection',
      slug: 'cookbook-collection',
      description: 'Delicious recipes from around the world',
      price: 34.99,
      comparePrice: 44.99,
      images: ['https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?w=500'],
      stock: 18,
      featured: false,
      categoryId: books.id
    },

    // Home & Garden
    {
      name: 'Garden Tools Kit',
      slug: 'garden-tools-kit',
      description: 'Complete set of professional garden tools',
      price: 89.99,
      comparePrice: 119.99,
      images: ['https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=500'],
      stock: 10,
      featured: true,
      categoryId: home.id
    },
    {
      name: 'Smart Home Hub',
      slug: 'smart-home-hub',
      description: 'Control all your smart devices from one central hub',
      price: 199.99,
      comparePrice: 249.99,
      images: ['https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=500'],
      stock: 6,
      featured: false,
      categoryId: home.id
    },
    {
      name: 'Indoor Plant Set',
      slug: 'indoor-plant-set',
      description: 'Beautiful collection of air-purifying indoor plants',
      price: 59.99,
      comparePrice: 79.99,
      images: ['https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=500'],
      stock: 12,
      featured: false,
      categoryId: home.id
    },

    // Sports & Fitness
    {
      name: 'Yoga Mat Premium',
      slug: 'yoga-mat-premium',
      description: 'Non-slip yoga mat with extra cushioning for comfort',
      price: 39.99,
      comparePrice: 59.99,
      images: ['https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?w=500'],
      stock: 20,
      featured: false,
      categoryId: sports.id
    },
    {
      name: 'Dumbbell Set',
      slug: 'dumbbell-set',
      description: 'Adjustable dumbbell set for home workouts',
      price: 149.99,
      comparePrice: 199.99,
      images: ['https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=500'],
      stock: 7,
      featured: false,
      categoryId: sports.id
    },
    {
      name: 'Basketball',
      slug: 'basketball',
      description: 'Official size basketball with premium grip',
      price: 24.99,
      comparePrice: 34.99,
      images: ['https://images.unsplash.com/photo-1546519638-68e109498ffc?w=500'],
      stock: 15,
      featured: false,
      categoryId: sports.id
    },

    // Beauty & Health
    {
      name: 'Skincare Set',
      slug: 'skincare-set',
      description: 'Complete skincare routine with natural ingredients',
      price: 79.99,
      comparePrice: 99.99,
      images: ['https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=500'],
      stock: 14,
      featured: false,
      categoryId: beauty.id
    },
    {
      name: 'Vitamin C Serum',
      slug: 'vitamin-c-serum',
      description: 'Anti-aging vitamin C serum for radiant skin',
      price: 34.99,
      comparePrice: 44.99,
      images: ['https://images.unsplash.com/photo-1556228720-195a672e8a03?w=500'],
      stock: 22,
      featured: false,
      categoryId: beauty.id
    },
    {
      name: 'Protein Powder',
      slug: 'protein-powder',
      description: 'High-quality whey protein powder for muscle recovery',
      price: 49.99,
      comparePrice: 69.99,
      images: ['https://images.unsplash.com/photo-1593095948071-474c5cc2989d?w=500'],
      stock: 16,
      featured: false,
      categoryId: beauty.id
    }
  ]

  for (const product of products) {
    await prisma.product.upsert({
      where: { slug: product.slug },
      update: {},
      create: product
    })
  }

  console.log('✅ Products created')

  // Create demo user
  const demoUser = await prisma.user.upsert({
    where: { email: 'demo@example.com' },
    update: {},
    create: {
      name: 'Demo User',
      email: 'demo@example.com',
      password: '$2a$10$K7L1OJ45/4Y2nIvhRVpCe.FSmhDdWoXehVzJByJp8Yz59dK3k2K2K', // password123
      role: 'USER'
    }
  })

  console.log('✅ Demo user created')

  // Create some reviews
  const reviewProducts = await prisma.product.findMany({ take: 5 })
  
  for (const product of reviewProducts) {
    await prisma.review.create({
      data: {
        rating: Math.floor(Math.random() * 2) + 4, // 4-5 stars
        comment: 'Great product! Highly recommended.',
        userId: demoUser.id,
        productId: product.id
      }
    })
  }

  console.log('✅ Reviews created')
  console.log('🎉 Database seeding completed!')
}

main()
  .catch((e) => {
    console.error('❌ Error seeding database:', e)
    process.exit(1)
  })
  .then(async () => {
    await prisma.$disconnect()
  })