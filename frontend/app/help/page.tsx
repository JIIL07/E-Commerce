'use client'

import { useState } from 'react'
import Navbar from '../../components/Navbar'
import {
  Search,
  ChevronDown,
  ChevronRight,
  MessageCircle,
  Phone,
  Mail,
  Clock,
  HelpCircle,
  ShoppingBag,
  CreditCard,
  Truck,
  User,
  Shield,
  Star
} from 'lucide-react'

export default function HelpPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [expandedSection, setExpandedSection] = useState<string | null>(null)

  const faqSections = [
    {
      id: 'shopping',
      title: 'Shopping & Orders',
      icon: ShoppingBag,
      questions: [
        {
          question: 'How do I place an order?',
          answer: 'To place an order, simply browse our products, add items to your cart, and proceed to checkout. You can pay using various methods including credit cards, PayPal, and other secure payment options.'
        },
        {
          question: 'Can I modify or cancel my order?',
          answer: 'You can modify or cancel your order within 30 minutes of placing it. After that, the order will be processed and cannot be changed. Contact our support team for assistance.'
        },
        {
          question: 'How do I track my order?',
          answer: 'Once your order ships, you will receive a tracking number via email. You can also track your order by logging into your account and visiting the "Orders" section.'
        }
      ]
    },
    {
      id: 'payment',
      title: 'Payment & Billing',
      icon: CreditCard,
      questions: [
        {
          question: 'What payment methods do you accept?',
          answer: 'We accept all major credit cards (Visa, MasterCard, American Express), PayPal, Apple Pay, Google Pay, and bank transfers.'
        },
        {
          question: 'Is my payment information secure?',
          answer: 'Yes, we use industry-standard SSL encryption and never store your payment information on our servers. All transactions are processed securely through trusted payment processors.'
        },
        {
          question: 'Can I get a refund?',
          answer: 'We offer a 30-day return policy for most items. Refunds are processed within 5-7 business days after we receive the returned item.'
        }
      ]
    },
    {
      id: 'shipping',
      title: 'Shipping & Delivery',
      icon: Truck,
      questions: [
        {
          question: 'How long does shipping take?',
          answer: 'Standard shipping takes 3-5 business days, while express shipping takes 1-2 business days. International shipping may take 7-14 business days depending on the destination.'
        },
        {
          question: 'Do you ship internationally?',
          answer: 'Yes, we ship to most countries worldwide. Shipping costs and delivery times vary by location. Check our shipping calculator for specific details.'
        },
        {
          question: 'What if my package is damaged or lost?',
          answer: 'If your package arrives damaged or is lost in transit, please contact our support team immediately. We will work with the shipping carrier to resolve the issue and provide a replacement or refund.'
        }
      ]
    },
    {
      id: 'account',
      title: 'Account & Profile',
      icon: User,
      questions: [
        {
          question: 'How do I create an account?',
          answer: 'Click the "Sign Up" button in the top right corner, fill in your details, and verify your email address. You can also sign up using your Google or Facebook account.'
        },
        {
          question: 'I forgot my password. How do I reset it?',
          answer: 'Click "Forgot Password" on the login page, enter your email address, and follow the instructions sent to your email to reset your password.'
        },
        {
          question: 'How do I update my profile information?',
          answer: 'Log into your account, go to "Profile" or "Settings", and update your information. Make sure to save your changes before leaving the page.'
        }
      ]
    },
    {
      id: 'security',
      title: 'Security & Privacy',
      icon: Shield,
      questions: [
        {
          question: 'How do you protect my personal information?',
          answer: 'We use industry-standard security measures including SSL encryption, secure servers, and regular security audits. We never share your personal information with third parties without your consent.'
        },
        {
          question: 'Can I delete my account?',
          answer: 'Yes, you can delete your account at any time by going to Settings > Account > Delete Account. Please note that this action cannot be undone.'
        }
      ]
    }
  ]

  const toggleSection = (sectionId: string) => {
    setExpandedSection(expandedSection === sectionId ? null : sectionId)
  }

  const filteredSections = faqSections.filter(section =>
    section.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    section.questions.some(q =>
      q.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
      q.answer.toLowerCase().includes(searchQuery.toLowerCase())
    )
  )

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50">
      <Navbar />

      <main className="container mx-auto px-4 py-12">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              Help Center
            </h1>
            <p className="text-xl text-gray-600 mb-8">
              Find answers to common questions and get the support you need
            </p>

            {/* Search Bar */}
            <div className="relative max-w-2xl mx-auto">
              <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
              <input
                type="text"
                placeholder="Search for help..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full pl-12 pr-4 py-4 border border-gray-300 rounded-2xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent shadow-lg"
              />
            </div>
          </div>

          {/* Quick Contact */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
            <div className="bg-white/80 backdrop-blur-sm rounded-2xl p-6 shadow-lg border border-white/20 text-center">
              <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center mx-auto mb-4">
                <MessageCircle className="w-6 h-6 text-blue-600" />
              </div>
              <h3 className="font-bold text-lg mb-2">Live Chat</h3>
              <p className="text-gray-600 text-sm mb-4">Get instant help from our support team</p>
              <button className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors">
                Start Chat
              </button>
            </div>

            <div className="bg-white/80 backdrop-blur-sm rounded-2xl p-6 shadow-lg border border-white/20 text-center">
              <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center mx-auto mb-4">
                <Phone className="w-6 h-6 text-green-600" />
              </div>
              <h3 className="font-bold text-lg mb-2">Phone Support</h3>
              <p className="text-gray-600 text-sm mb-4">Call us at +1 (555) 123-4567</p>
              <p className="text-xs text-gray-500">Mon-Fri 9AM-6PM EST</p>
            </div>

            <div className="bg-white/80 backdrop-blur-sm rounded-2xl p-6 shadow-lg border border-white/20 text-center">
              <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center mx-auto mb-4">
                <Mail className="w-6 h-6 text-purple-600" />
              </div>
              <h3 className="font-bold text-lg mb-2">Email Support</h3>
              <p className="text-gray-600 text-sm mb-4">Send us an email anytime</p>
              <p className="text-xs text-gray-500">support@ecommerce.com</p>
            </div>
          </div>

          {/* FAQ Sections */}
          <div className="space-y-6">
            {filteredSections.map((section) => (
              <div key={section.id} className="bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg border border-white/20 overflow-hidden">
                <button
                  onClick={() => toggleSection(section.id)}
                  className="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-center">
                    <div className="w-10 h-10 bg-gradient-to-r from-blue-500 to-purple-500 rounded-xl flex items-center justify-center mr-4">
                      <section.icon className="w-5 h-5 text-white" />
                    </div>
                    <h2 className="text-xl font-semibold text-gray-900">{section.title}</h2>
                  </div>
                  {expandedSection === section.id ? (
                    <ChevronDown className="w-6 h-6 text-gray-500" />
                  ) : (
                    <ChevronRight className="w-6 h-6 text-gray-500" />
                  )}
                </button>

                {expandedSection === section.id && (
                  <div className="px-6 pb-6">
                    <div className="space-y-4">
                      {section.questions.map((faq, index) => (
                        <div key={index} className="border-l-4 border-blue-500 pl-4">
                          <h3 className="font-semibold text-gray-900 mb-2">{faq.question}</h3>
                          <p className="text-gray-600 leading-relaxed">{faq.answer}</p>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Contact Form */}
          <div className="mt-16 bg-white/80 backdrop-blur-sm rounded-2xl shadow-lg border border-white/20 p-8">
            <h2 className="text-2xl font-bold mb-6 text-center">Still Need Help?</h2>
            <p className="text-gray-600 text-center mb-8">
              Can't find what you're looking for? Send us a message and we'll get back to you as soon as possible.
            </p>

            <form className="max-w-2xl mx-auto space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
                    Name
                  </label>
                  <input
                    type="text"
                    id="name"
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Your name"
                  />
                </div>
                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
                    Email
                  </label>
                  <input
                    type="email"
                    id="email"
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="your@email.com"
                  />
                </div>
              </div>

              <div>
                <label htmlFor="subject" className="block text-sm font-medium text-gray-700 mb-2">
                  Subject
                </label>
                <input
                  type="text"
                  id="subject"
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="What can we help you with?"
                />
              </div>

              <div>
                <label htmlFor="message" className="block text-sm font-medium text-gray-700 mb-2">
                  Message
                </label>
                <textarea
                  id="message"
                  rows={5}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                  placeholder="Please describe your issue or question in detail..."
                ></textarea>
              </div>

              <div className="text-center">
                <button
                  type="submit"
                  className="bg-gradient-to-r from-blue-600 to-purple-600 text-white px-8 py-3 rounded-xl hover:scale-105 transition-all duration-300 shadow-lg hover:shadow-xl font-medium"
                >
                  Send Message
                </button>
              </div>
            </form>
          </div>
        </div>
      </main>
    </div>
  )
}
