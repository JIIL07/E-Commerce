'use client'

import { useState } from 'react'
import { ChevronDown, ChevronUp, X, Filter } from 'lucide-react'

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

interface FilterPanelProps {
  filters: FilterGroup[]
  selectedFilters: Record<string, string[]>
  onFilterChange: (filterId: string, values: string[]) => void
  onClearAll: () => void
  className?: string
}

export default function FilterPanel({
  filters,
  selectedFilters,
  onFilterChange,
  onClearAll,
  className = ''
}: FilterPanelProps) {
  const [expandedFilters, setExpandedFilters] = useState<Set<string>>(new Set())

  const toggleFilter = (filterId: string) => {
    const newExpanded = new Set(expandedFilters)
    if (newExpanded.has(filterId)) {
      newExpanded.delete(filterId)
    } else {
      newExpanded.add(filterId)
    }
    setExpandedFilters(newExpanded)
  }

  const handleOptionChange = (filterId: string, optionValue: string, checked: boolean) => {
    const currentValues = selectedFilters[filterId] || []
    let newValues: string[]

    if (checked) {
      newValues = [...currentValues, optionValue]
    } else {
      newValues = currentValues.filter(value => value !== optionValue)
    }

    onFilterChange(filterId, newValues)
  }

  const handleRangeChange = (filterId: string, value: string) => {
    onFilterChange(filterId, [value])
  }

  const clearFilter = (filterId: string) => {
    onFilterChange(filterId, [])
  }

  const hasActiveFilters = Object.values(selectedFilters).some(values => values.length > 0)

  return (
    <div className={`bg-white rounded-lg shadow-lg p-6 ${className}`}>
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-gray-900 flex items-center">
          <Filter className="w-5 h-5 mr-2" />
          Filters
        </h3>
        {hasActiveFilters && (
          <button
            onClick={onClearAll}
            className="text-sm text-blue-600 hover:text-blue-700 font-medium"
          >
            Clear All
          </button>
        )}
      </div>

      <div className="space-y-4">
        {filters.map((filter) => {
          const isExpanded = expandedFilters.has(filter.id)
          const hasSelection = selectedFilters[filter.id]?.length > 0

          return (
            <div key={filter.id} className="border-b border-gray-200 pb-4 last:border-b-0">
              <button
                onClick={() => toggleFilter(filter.id)}
                className="w-full flex items-center justify-between text-left"
              >
                <div className="flex items-center">
                  <span className="font-medium text-gray-900">{filter.label}</span>
                  {hasSelection && (
                    <span className="ml-2 bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full">
                      {selectedFilters[filter.id].length}
                    </span>
                  )}
                </div>
                {isExpanded ? (
                  <ChevronUp className="w-4 h-4 text-gray-500" />
                ) : (
                  <ChevronDown className="w-4 h-4 text-gray-500" />
                )}
              </button>

              {isExpanded && (
                <div className="mt-3 space-y-2">
                  {filter.type === 'range' ? (
                    <div className="space-y-3">
                      <input
                        type="range"
                        min={filter.min || 0}
                        max={filter.max || 100}
                        value={selectedFilters[filter.id]?.[0] || filter.min || 0}
                        onChange={(e) => handleRangeChange(filter.id, e.target.value)}
                        className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
                      />
                      <div className="flex justify-between text-sm text-gray-600">
                        <span>${filter.min || 0}</span>
                        <span>${selectedFilters[filter.id]?.[0] || filter.min || 0}</span>
                        <span>${filter.max || 100}</span>
                      </div>
                    </div>
                  ) : (
                    <div className="space-y-2 max-h-48 overflow-y-auto">
                      {filter.options.map((option) => {
                        const isSelected = selectedFilters[filter.id]?.includes(option.value) || false

                        return (
                          <label
                            key={option.value}
                            className="flex items-center space-x-3 cursor-pointer hover:bg-gray-50 p-2 rounded"
                          >
                            <input
                              type={filter.type}
                              name={filter.id}
                              value={option.value}
                              checked={isSelected}
                              onChange={(e) => handleOptionChange(
                                filter.id,
                                option.value,
                                e.target.checked
                              )}
                              className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                            />
                            <span className="text-sm text-gray-700 flex-1">
                              {option.label}
                            </span>
                            {option.count !== undefined && (
                              <span className="text-xs text-gray-500">
                                ({option.count})
                              </span>
                            )}
                          </label>
                        )
                      })}
                    </div>
                  )}

                  {hasSelection && (
                    <button
                      onClick={() => clearFilter(filter.id)}
                      className="flex items-center text-sm text-red-600 hover:text-red-700 mt-2"
                    >
                      <X className="w-3 h-3 mr-1" />
                      Clear
                    </button>
                  )}
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}
