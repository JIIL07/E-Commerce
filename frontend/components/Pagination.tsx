'use client'

import { ChevronLeft, ChevronRight, MoreHorizontal } from 'lucide-react'

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
  showFirstLast?: boolean
  maxVisiblePages?: number
  className?: string
}

export default function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  showFirstLast = true,
  maxVisiblePages = 5,
  className = ''
}: PaginationProps) {
  if (totalPages <= 1) return null

  const getVisiblePages = () => {
    const pages: (number | string)[] = []
    
    if (totalPages <= maxVisiblePages) {
      // Show all pages if total is less than max visible
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i)
      }
    } else {
      // Calculate start and end pages
      let startPage = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2))
      let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1)
      
      // Adjust start page if we're near the end
      if (endPage - startPage + 1 < maxVisiblePages) {
        startPage = Math.max(1, endPage - maxVisiblePages + 1)
      }
      
      // Add first page and ellipsis if needed
      if (startPage > 1) {
        pages.push(1)
        if (startPage > 2) {
          pages.push('...')
        }
      }
      
      // Add visible pages
      for (let i = startPage; i <= endPage; i++) {
        pages.push(i)
      }
      
      // Add ellipsis and last page if needed
      if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
          pages.push('...')
        }
        pages.push(totalPages)
      }
    }
    
    return pages
  }

  const visiblePages = getVisiblePages()

  const handlePageClick = (page: number | string) => {
    if (typeof page === 'number') {
      onPageChange(page)
    }
  }

  const handlePrevious = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1)
    }
  }

  const handleNext = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1)
    }
  }

  return (
    <nav className={`flex items-center justify-center space-x-2 ${className}`} aria-label="Pagination">
      {/* Previous Button */}
      <button
        onClick={handlePrevious}
        disabled={currentPage === 1}
        className={`
          flex items-center px-3 py-2 text-sm font-medium rounded-lg
          ${currentPage === 1
            ? 'text-gray-400 cursor-not-allowed'
            : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
          }
          transition-colors
        `}
        aria-label="Previous page"
      >
        <ChevronLeft className="w-4 h-4 mr-1" />
        Previous
      </button>

      {/* First Page Button */}
      {showFirstLast && currentPage > 2 && (
        <>
          <button
            onClick={() => handlePageClick(1)}
            className="px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900 rounded-lg transition-colors"
          >
            1
          </button>
          {currentPage > 3 && (
            <span className="px-2 py-2 text-gray-400">
              <MoreHorizontal className="w-4 h-4" />
            </span>
          )}
        </>
      )}

      {/* Page Numbers */}
      {visiblePages.map((page, index) => (
        <div key={index}>
          {typeof page === 'number' ? (
            <button
              onClick={() => handlePageClick(page)}
              className={`
                px-3 py-2 text-sm font-medium rounded-lg transition-colors
                ${page === currentPage
                  ? 'bg-blue-600 text-white'
                  : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                }
              `}
              aria-label={`Page ${page}`}
              aria-current={page === currentPage ? 'page' : undefined}
            >
              {page}
            </button>
          ) : (
            <span className="px-2 py-2 text-gray-400">
              <MoreHorizontal className="w-4 h-4" />
            </span>
          )}
        </div>
      ))}

      {/* Last Page Button */}
      {showFirstLast && currentPage < totalPages - 1 && (
        <>
          {currentPage < totalPages - 2 && (
            <span className="px-2 py-2 text-gray-400">
              <MoreHorizontal className="w-4 h-4" />
            </span>
          )}
          <button
            onClick={() => handlePageClick(totalPages)}
            className="px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 hover:text-gray-900 rounded-lg transition-colors"
          >
            {totalPages}
          </button>
        </>
      )}

      {/* Next Button */}
      <button
        onClick={handleNext}
        disabled={currentPage === totalPages}
        className={`
          flex items-center px-3 py-2 text-sm font-medium rounded-lg
          ${currentPage === totalPages
            ? 'text-gray-400 cursor-not-allowed'
            : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
          }
          transition-colors
        `}
        aria-label="Next page"
      >
        Next
        <ChevronRight className="w-4 h-4 ml-1" />
      </button>
    </nav>
  )
}
