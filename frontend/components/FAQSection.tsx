'use client'

import { useState } from 'react'
import { ChevronDown, ChevronUp, HelpCircle, Search } from 'lucide-react'

interface FAQ {
  id: string
  question: string
  answer: string
  category: string
  isPopular?: boolean
}

interface FAQSectionProps {
  title?: string
  subtitle?: string
  faqs?: FAQ[]
  showSearch?: boolean
  showCategories?: boolean
}

export default function FAQSection({ 
  title = "Frequently Asked Questions",
  subtitle = "Find answers to common questions about our services",
  faqs,
  showSearch = true,
  showCategories = true
}: FAQSectionProps) {
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set())
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCategory, setSelectedCategory] = useState('all')

  const defaultFAQs: FAQ[] = [
    {
      id: '1',
      question: 'How do I place an order?',
      answer: 'Placing an order is simple! Browse our products, add items to your cart, and proceed to checkout. You can pay using various methods including credit cards, PayPal, and other secure payment options. You\'ll receive a confirmation email once your order is placed.',
      category: 'Orders',
      isPopular: true
    },
    {
      id: '2',
      question: 'What payment methods do you accept?',
      answer: 'We accept all major credit cards (Visa, MasterCard, American Express), PayPal, Apple Pay, Google Pay, and bank transfers. All payments are processed securely through trusted payment processors.',
      category: 'Payment',
      isPopular: true
    },
    {
      id: '3',
      question: 'How long does shipping take?',
      answer: 'Standard shipping takes 3-5 business days, while express shipping takes 1-2 business days. International shipping may take 7-14 business days depending on the destination. You can track your order using the tracking number provided in your confirmation email.',
      category: 'Shipping',
      isPopular: true
    },
    {
      id: '4',
      question: 'Can I return or exchange items?',
      answer: 'Yes! We offer a 30-day return policy for most items. Items must be in original condition with tags attached. Refunds will be processed within 5-7 business days after we receive the returned item. Some items may be excluded from our return policy.',
      category: 'Returns',
      isPopular: true
    },
    {
      id: '5',
      question: 'Is my personal information secure?',
      answer: 'Absolutely! We use industry-standard security measures including SSL encryption, secure servers, and regular security audits. We never store your payment information on our servers and never share your personal information with third parties without your consent.',
      category: 'Security',
      isPopular: false
    },
    {
      id: '6',
      question: 'Do you ship internationally?',
      answer: 'Yes, we ship to most countries worldwide. Shipping costs and delivery times vary by location. You can check shipping options and costs during checkout. International orders may be subject to customs duties and taxes.',
      category: 'Shipping',
      isPopular: false
    },
    {
      id: '7',
      question: 'How can I track my order?',
      answer: 'Once your order ships, you will receive a tracking number via email. You can also track your order by logging into your account and visiting the "Orders" section. Click on your order to see detailed tracking information.',
      category: 'Orders',
      isPopular: false
    },
    {
      id: '8',
      question: 'What if my package is damaged or lost?',
      answer: 'If your package arrives damaged or is lost in transit, please contact our support team immediately. We will work with the shipping carrier to resolve the issue and provide a replacement or refund. We take full responsibility for items damaged during shipping.',
      category: 'Shipping',
      isPopular: false
    },
    {
      id: '9',
      question: 'Can I cancel my order?',
      answer: 'You can cancel your order within 30 minutes of placing it. After that, the order will be processed and cannot be changed. If you need to cancel after this time, please contact our support team as soon as possible.',
      category: 'Orders',
      isPopular: false
    },
    {
      id: '10',
      question: 'Do you offer customer support?',
      answer: 'Yes! We provide 24/7 customer support through live chat, email, and phone. Our support team is available to help with any questions or concerns you may have. You can also check our help center for self-service options.',
      category: 'Support',
      isPopular: true
    }
  ]

  const displayFAQs = faqs || defaultFAQs

  const categories = ['all', ...Array.from(new Set(displayFAQs.map(faq => faq.category)))]

  const filteredFAQs = displayFAQs.filter(faq => {
    const matchesSearch = faq.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         faq.answer.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesCategory = selectedCategory === 'all' || faq.category === selectedCategory
    return matchesSearch && matchesCategory
  })

  const toggleExpanded = (id: string) => {
    const newExpanded = new Set(expandedItems)
    if (newExpanded.has(id)) {
      newExpanded.delete(id)
    } else {
      newExpanded.add(id)
    }
    setExpandedItems(newExpanded)
  }

  const popularFAQs = displayFAQs.filter(faq => faq.isPopular)

  return (
    <section className="py-20 bg-gradient-to-br from-gray-50 to-blue-50">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-16">
          <div className="inline-flex items-center px-4 py-2 bg-blue-100 rounded-full text-sm font-medium text-blue-600 mb-4">
            <HelpCircle className="w-4 h-4 mr-2" />
            Help Center
          </div>
          
          <h2 className="text-responsive-3xl font-bold text-gray-900 mb-4">
            {title}
          </h2>
          
          <p className="text-responsive-lg text-gray-600 max-w-2xl mx-auto">
            {subtitle}
          </p>
        </div>

        {/* Search and Filters */}
        {showSearch && (
          <div className="max-w-2xl mx-auto mb-12">
            <div className="relative">
              <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Search FAQs..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-12 pr-4 py-4 border border-gray-300 rounded-2xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white shadow-lg"
              />
            </div>
          </div>
        )}

        {/* Categories */}
        {showCategories && (
          <div className="flex flex-wrap justify-center gap-3 mb-12">
            {categories.map((category) => (
              <button
                key={category}
                onClick={() => setSelectedCategory(category)}
                className={`px-6 py-3 rounded-full font-medium transition-all duration-200 ${
                  selectedCategory === category
                    ? 'bg-blue-600 text-white shadow-lg'
                    : 'bg-white text-gray-700 hover:bg-gray-100 shadow-md'
                }`}
              >
                {category === 'all' ? 'All Questions' : category}
              </button>
            ))}
          </div>
        )}

        {/* Popular FAQs */}
        {selectedCategory === 'all' && searchQuery === '' && (
          <div className="mb-12">
            <h3 className="text-xl font-bold text-gray-900 mb-6 text-center">Popular Questions</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {popularFAQs.slice(0, 4).map((faq) => (
                <div
                  key={faq.id}
                  className="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-300 cursor-pointer"
                  onClick={() => toggleExpanded(faq.id)}
                >
                  <div className="flex items-start justify-between">
                    <h4 className="font-semibold text-gray-900 mb-2 flex-1 pr-4">
                      {faq.question}
                    </h4>
                    {expandedItems.has(faq.id) ? (
                      <ChevronUp className="w-5 h-5 text-gray-500 flex-shrink-0" />
                    ) : (
                      <ChevronDown className="w-5 h-5 text-gray-500 flex-shrink-0" />
                    )}
                  </div>
                  {expandedItems.has(faq.id) && (
                    <p className="text-gray-600 leading-relaxed">
                      {faq.answer}
                    </p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* All FAQs */}
        <div className="max-w-4xl mx-auto">
          <div className="space-y-4">
            {filteredFAQs.map((faq) => (
              <div
                key={faq.id}
                className="bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-all duration-300"
              >
                <button
                  onClick={() => toggleExpanded(faq.id)}
                  className="w-full p-6 text-left flex items-start justify-between hover:bg-gray-50 transition-colors"
                >
                  <div className="flex-1 pr-4">
                    <div className="flex items-center mb-2">
                      <span className="px-3 py-1 bg-blue-100 text-blue-600 text-sm font-medium rounded-full mr-3">
                        {faq.category}
                      </span>
                      {faq.isPopular && (
                        <span className="px-3 py-1 bg-orange-100 text-orange-600 text-sm font-medium rounded-full">
                          Popular
                        </span>
                      )}
                    </div>
                    <h3 className="text-lg font-semibold text-gray-900">
                      {faq.question}
                    </h3>
                  </div>
                  {expandedItems.has(faq.id) ? (
                    <ChevronUp className="w-6 h-6 text-gray-500 flex-shrink-0" />
                  ) : (
                    <ChevronDown className="w-6 h-6 text-gray-500 flex-shrink-0" />
                  )}
                </button>
                
                {expandedItems.has(faq.id) && (
                  <div className="px-6 pb-6">
                    <div className="border-t border-gray-200 pt-4">
                      <p className="text-gray-600 leading-relaxed">
                        {faq.answer}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* No Results */}
        {filteredFAQs.length === 0 && (
          <div className="text-center py-12">
            <HelpCircle className="w-16 h-16 text-gray-400 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-900 mb-2">No FAQs Found</h3>
            <p className="text-gray-600 mb-6">
              Try adjusting your search terms or category filter
            </p>
            <button
              onClick={() => {
                setSearchQuery('')
                setSelectedCategory('all')
              }}
              className="btn-primary"
            >
              Clear Filters
            </button>
          </div>
        )}

        {/* Contact Support */}
        <div className="text-center mt-16">
          <div className="bg-white rounded-3xl shadow-xl p-8 max-w-2xl mx-auto">
            <h3 className="text-2xl font-bold text-gray-900 mb-4">
              Still Need Help?
            </h3>
            <p className="text-gray-600 mb-6">
              Can't find what you're looking for? Our support team is here to help!
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <button className="btn-primary">
                Contact Support
              </button>
              <button className="btn-secondary">
                Live Chat
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
