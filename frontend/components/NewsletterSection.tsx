'use client'

import { useState } from 'react'
import { Mail, Send, CheckCircle, Gift, Star, Users } from 'lucide-react'

interface NewsletterSectionProps {
  title?: string
  subtitle?: string
  showBenefits?: boolean
  showSocialProof?: boolean
}

export default function NewsletterSection({ 
  title = "Stay in the Loop",
  subtitle = "Get the latest updates, exclusive offers, and insider tips delivered to your inbox",
  showBenefits = true,
  showSocialProof = true
}: NewsletterSectionProps) {
  const [email, setEmail] = useState('')
  const [isSubscribed, setIsSubscribed] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

  const benefits = [
    {
      icon: Gift,
      title: 'Exclusive Offers',
      description: 'Get access to special discounts and early bird deals'
    },
    {
      icon: Star,
      title: 'Product Updates',
      description: 'Be the first to know about new products and features'
    },
    {
      icon: Users,
      title: 'Community Access',
      description: 'Join our community of like-minded shoppers'
    }
  ]

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email.trim()) return

    setIsLoading(true)
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    setIsSubscribed(true)
    setIsLoading(false)
    setEmail('')
  }

  if (isSubscribed) {
    return (
      <section className="py-20 bg-gradient-to-br from-green-500 via-emerald-500 to-teal-500">
        <div className="container mx-auto px-4">
          <div className="max-w-2xl mx-auto text-center">
            <div className="bg-white/10 backdrop-blur-md rounded-3xl p-12 border border-white/20">
              <div className="w-20 h-20 bg-white rounded-full flex items-center justify-center mx-auto mb-6">
                <CheckCircle className="w-10 h-10 text-green-600" />
              </div>
              
              <h2 className="text-3xl font-bold text-white mb-4">
                Welcome to Our Community!
              </h2>
              
              <p className="text-green-100 text-lg mb-8">
                Thank you for subscribing! You'll receive your first newsletter within the next few minutes.
              </p>
              
              <button
                onClick={() => setIsSubscribed(false)}
                className="bg-white text-green-600 px-8 py-3 rounded-xl font-semibold hover:bg-gray-100 transition-colors"
              >
                Subscribe Another Email
              </button>
            </div>
          </div>
        </div>
      </section>
    )
  }

  return (
    <section className="py-20 bg-gradient-to-br from-blue-600 via-purple-600 to-pink-600 relative overflow-hidden">
      {/* Background Effects */}
      <div className="absolute inset-0">
        <div className="absolute top-0 left-0 w-full h-full">
          <div className="absolute top-20 left-20 w-72 h-72 bg-white/10 rounded-full blur-3xl animate-float"></div>
          <div className="absolute top-40 right-20 w-96 h-96 bg-white/10 rounded-full blur-3xl animate-float" style={{ animationDelay: '1s' }}></div>
          <div className="absolute bottom-20 left-1/3 w-80 h-80 bg-white/10 rounded-full blur-3xl animate-float" style={{ animationDelay: '2s' }}></div>
        </div>
      </div>

      <div className="container mx-auto px-4 relative z-10">
        <div className="max-w-4xl mx-auto">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            {/* Content */}
            <div className="text-center lg:text-left">
              <div className="inline-flex items-center px-4 py-2 bg-white/20 backdrop-blur-sm rounded-full text-sm font-medium text-white mb-6">
                <Mail className="w-4 h-4 mr-2" />
                Newsletter
              </div>
              
              <h2 className="text-responsive-3xl font-bold text-white mb-6">
                {title}
              </h2>
              
              <p className="text-responsive-lg text-blue-100 mb-8">
                {subtitle}
              </p>

              {/* Benefits */}
              {showBenefits && (
                <div className="space-y-4 mb-8">
                  {benefits.map((benefit, index) => (
                    <div key={index} className="flex items-center space-x-3">
                      <div className="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center">
                        <benefit.icon className="w-5 h-5 text-white" />
                      </div>
                      <div>
                        <h4 className="font-semibold text-white">{benefit.title}</h4>
                        <p className="text-blue-100 text-sm">{benefit.description}</p>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Social Proof */}
              {showSocialProof && (
                <div className="flex items-center space-x-6 text-blue-100">
                  <div className="text-center">
                    <div className="text-2xl font-bold text-white">50K+</div>
                    <div className="text-sm">Subscribers</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-white">4.9★</div>
                    <div className="text-sm">Rating</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-white">Weekly</div>
                    <div className="text-sm">Updates</div>
                  </div>
                </div>
              )}
            </div>

            {/* Newsletter Form */}
            <div className="bg-white/10 backdrop-blur-md rounded-3xl p-8 border border-white/20">
              <form onSubmit={handleSubmit} className="space-y-6">
                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-white mb-2">
                    Email Address
                  </label>
                  <input
                    type="email"
                    id="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="Enter your email address"
                    className="w-full px-4 py-4 text-gray-900 rounded-xl border-0 focus:outline-none focus:ring-2 focus:ring-white/50 placeholder-gray-500"
                    required
                  />
                </div>

                <button
                  type="submit"
                  disabled={isLoading || !email.trim()}
                  className="w-full bg-white text-blue-600 py-4 px-6 rounded-xl font-semibold hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 flex items-center justify-center"
                >
                  {isLoading ? (
                    <div className="spinner w-5 h-5 mr-2"></div>
                  ) : (
                    <Send className="w-5 h-5 mr-2" />
                  )}
                  {isLoading ? 'Subscribing...' : 'Subscribe Now'}
                </button>

                <p className="text-blue-100 text-sm text-center">
                  We respect your privacy. Unsubscribe at any time.
                </p>
              </form>

              {/* Trust Indicators */}
              <div className="mt-6 pt-6 border-t border-white/20">
                <div className="flex items-center justify-center space-x-6 text-blue-100 text-sm">
                  <div className="flex items-center">
                    <CheckCircle className="w-4 h-4 mr-1" />
                    <span>No Spam</span>
                  </div>
                  <div className="flex items-center">
                    <CheckCircle className="w-4 h-4 mr-1" />
                    <span>Easy Unsubscribe</span>
                  </div>
                  <div className="flex items-center">
                    <CheckCircle className="w-4 h-4 mr-1" />
                    <span>Secure</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
