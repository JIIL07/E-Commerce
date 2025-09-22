'use client'

import { useState, useEffect, useCallback } from 'react'
import { useSearchParams } from 'next/navigation'
import Navbar from '../../components/Navbar'
import Breadcrumbs from '../../components/Breadcrumbs'
import ProductCard from '../../components/ProductCard'
import Pagination from '../../components/Pagination'
import FilterPanel from '../../components/FilterPanel'
import LoadingSpinner from '../../components/LoadingSpinner'
import { Search, Filter, SortAsc, SortDesc } from 'lucide-react'

interface Product {
  id: string
  name: string
  price: number
  image: string
  category: {
    name: string
  }
  stock: number
  averageRating: number
  reviewCount: number
}

interface FilterOption {
  value: string
  label: string
  count?: number
}

interface FilterGroup {
  id: string
  label: string
  options: FilterOption[]
  type: 'checkbox' | 'radio' | 'range'
  min?: number
  max?: number
}

export default function SearchPage() {
  const searchParams = useSearchParams()
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [totalResults, setTotalResults] = useState(0)
  const [searchQuery, setSearchQuery] = useState('')
  const [sortBy, setSortBy] = useState('relevance')
  const [showFilters, setShowFilters] = useState(false)
  const [selectedFilters, setSelectedFilters] = useState<Record<string, string[]>>({})

  const filters: FilterGroup[] = [
    {
      id: 'category',
      label: 'Category',
      type: 'checkbox',
      options: [
        { value: 'electronics', label: 'Electronics', count: 45 },
        { value: 'clothing', label: 'Clothing', count: 32 },
        { value: 'home', label: 'Home & Garden', count: 28 },
        { value: 'sports', label: 'Sports', count: 19 },
        { value: 'books', label: 'Books', count: 15 }
      ]
    },
    {
      id: 'price',
      label: 'Price Range',
      type: 'range',
      options: [],
      min: 0,
      max: 1000
    },
    {
      id: 'rating',
      label: 'Rating',
      type: 'checkbox',
      options: [
        { value: '5', label: '5 Stars', count: 23 },
        { value: '4', label: '4+ Stars', count: 45 },
        { value: '3', label: '3+ Stars', count: 67 },
        { value: '2', label: '2+ Stars', count: 89 }
      ]
    },
    {
      id: 'availability',
      label: 'Availability',
      type: 'checkbox',
      options: [
        { value: 'in-stock', label: 'In Stock', count: 120 },
        { value: 'low-stock', label: 'Low Stock', count: 15 },
        { value: 'out-of-stock', label: 'Out of Stock', count: 8 }
      ]
    }
  ]

  const sortOptions = [
    { value: 'relevance', label: 'Relevance' },
    { value: 'price-asc', label: 'Price: Low to High' },
    { value: 'price-desc', label: 'Price: High to Low' },
    { value: 'rating', label: 'Customer Rating' },
    { value: 'newest', label: 'Newest First' },
    { value: 'name', label: 'Name A-Z' }
  ]

  const fetchProducts = useCallback(async () => {
    setLoading(true)
    try {
      const query = searchParams.get('q') || ''
      setSearchQuery(query)

      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Mock data
      const mockProducts: Product[] = Array.from({ length: 20 }, (_, i) => ({
        id: `product-${i + 1}`,
        name: `Product ${i + 1} ${query ? `- ${query}` : ''}`,
        price: Math.random() * 500 + 10,
        image: `https://picsum.photos/400/400?random=${i + 1}`,
        category: {
          name: ['Electronics', 'Clothing', 'Home', 'Sports', 'Books'][i % 5]
        },
        stock: Math.floor(Math.random() * 100),
        averageRating: Math.random() * 2 + 3,
        reviewCount: Math.floor(Math.random() * 100)
      }))

      setProducts(mockProducts)
      setTotalPages(5)
      setTotalResults(100)
    } catch (error) {
      console.error('Error fetching products:', error)
    } finally {
      setLoading(false)
    }
  }, [searchParams])

  useEffect(() => {
    fetchProducts()
  }, [fetchProducts])

  const handleFilterChange = (filterId: string, values: string[]) => {
    setSelectedFilters(prev => ({
      ...prev,
      [filterId]: values
    }))
    setCurrentPage(1)
  }

  const handleClearAllFilters = () => {
    setSelectedFilters({})
    setCurrentPage(1)
  }

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const handleSortChange = (newSortBy: string) => {
    setSortBy(newSortBy)
    setCurrentPage(1)
  }

  const breadcrumbItems = [
    { label: 'Search', href: '/search' },
    { label: searchQuery || 'All Products' }
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="container mx-auto px-4 py-8">
        {/* Breadcrumbs */}
        <Breadcrumbs items={breadcrumbItems} className="mb-6" />

        {/* Search Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            {searchQuery ? `Search results for "${searchQuery}"` : 'All Products'}
          </h1>
          <p className="text-gray-600">
            {totalResults} products found
          </p>
        </div>

        <div className="flex flex-col lg:flex-row gap-8">
          {/* Filters Sidebar */}
          <div className="lg:w-80 flex-shrink-0">
            <div className="lg:sticky lg:top-8">
              <FilterPanel
                filters={filters}
                selectedFilters={selectedFilters}
                onFilterChange={handleFilterChange}
                onClearAll={handleClearAllFilters}
              />
            </div>
          </div>

          {/* Main Content */}
          <div className="flex-1">
            {/* Sort and View Options */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4">
              <div className="flex items-center space-x-4">
                <button
                  onClick={() => setShowFilters(!showFilters)}
                  className="lg:hidden flex items-center space-x-2 px-4 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  <Filter className="w-4 h-4" />
                  <span>Filters</span>
                </button>
              </div>

              <div className="flex items-center space-x-4">
                <span className="text-sm text-gray-600">Sort by:</span>
                <select
                  value={sortBy}
                  onChange={(e) => handleSortChange(e.target.value)}
                  className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  {sortOptions.map(option => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            {/* Mobile Filters */}
            {showFilters && (
              <div className="lg:hidden mb-6">
                <FilterPanel
                  filters={filters}
                  selectedFilters={selectedFilters}
                  onFilterChange={handleFilterChange}
                  onClearAll={handleClearAllFilters}
                />
              </div>
            )}

            {/* Products Grid */}
            {loading ? (
              <div className="flex justify-center items-center py-16">
                <LoadingSpinner size="lg" text="Loading products..." />
              </div>
            ) : products.length > 0 ? (
              <>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 mb-8">
                  {products.map(product => (
                    <ProductCard
                      key={product.id}
                      product={product}
                      onAddToCart={(productId) => {
                        console.log('Add to cart:', productId)
                      }}
                    />
                  ))}
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <Pagination
                    currentPage={currentPage}
                    totalPages={totalPages}
                    onPageChange={handlePageChange}
                  />
                )}
              </>
            ) : (
              <div className="text-center py-16">
                <Search className="w-16 h-16 text-gray-400 mx-auto mb-4" />
                <h3 className="text-xl font-semibold text-gray-900 mb-2">No products found</h3>
                <p className="text-gray-600 mb-6">
                  Try adjusting your search or filter criteria
                </p>
                <button
                  onClick={handleClearAllFilters}
                  className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  Clear All Filters
                </button>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
