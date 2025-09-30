'use client'

import React, { useState, useEffect } from 'react'
import Link from 'next/link'

interface Category {
  id: string
  name: string
  description: string
}

interface Product {
  id: string
  name: string
  description: string
  price: number
  images: string[]
  stock: number
  category_id: string
}

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedCategory, setSelectedCategory] = useState<string>('')

  useEffect(() => {
    fetchData()
  }, [])

  const fetchData = async () => {
    try {
      const [categoriesRes, productsRes] = await Promise.all([
        fetch('/api/categories'),
        fetch('/api/products')
      ])

      if (categoriesRes.ok) {
        const categoriesData = await categoriesRes.json()
        setCategories(categoriesData.categories || [])
      }

      if (productsRes.ok) {
        const productsData = await productsRes.json()
        setProducts(productsData.data || [])
      }
    } catch (error) {
      console.error('Error fetching data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getProductsByCategory = (categoryId: string) => {
    return products.filter(product => product.category_id === categoryId)
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading categories...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">Categories</h1>
          <p className="text-gray-600">Browse products by category</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {categories.map(category => {
            const categoryProducts = getProductsByCategory(category.id)
            return (
              <div key={category.id} className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow">
                <div className="p-6">
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">{category.name}</h3>
                  <p className="text-gray-600 text-sm mb-4">{category.description}</p>
                  
                  <div className="mb-4">
                    <span className="text-sm text-gray-500">
                      {categoryProducts.length} product{categoryProducts.length !== 1 ? 's' : ''}
                    </span>
                  </div>

                  {categoryProducts.length > 0 && (
                    <div className="mb-4">
                      <div className="grid grid-cols-2 gap-2">
                        {categoryProducts.slice(0, 4).map(product => (
                          <div key={product.id} className="aspect-square bg-gray-200 rounded-lg flex items-center justify-center">
                            {product.images && product.images.length > 0 ? (
                              <img
                                src={product.images[0].startsWith('http') ? product.images[0] : `/api/uploads/${product.images[0]}`}
                                alt={product.name}
                                className="w-full h-full object-cover rounded-lg"
                                onError={(e) => {
                                  e.currentTarget.src = `https://via.placeholder.com/150x150/f3f4f6/9ca3af?text=${encodeURIComponent(product.name)}`
                                }}
                              />
                            ) : (
                              <img
                                src={`https://via.placeholder.com/150x150/f3f4f6/9ca3af?text=${encodeURIComponent(product.name)}`}
                                alt={product.name}
                                className="w-full h-full object-cover rounded-lg"
                              />
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  <Link 
                    href={`/products?category=${category.id}`}
                    className="block w-full bg-blue-600 text-white text-center py-2 px-4 rounded-md hover:bg-blue-700 transition-colors"
                  >
                    View Products
                  </Link>
                </div>
              </div>
            )
          })}
        </div>

        {categories.length === 0 && (
          <div className="text-center py-12">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">No categories found</h2>
            <p className="text-gray-600">Categories will appear here once they are added to the system.</p>
          </div>
        )}
      </div>
    </div>
  )
}
