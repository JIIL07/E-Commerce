'use client'

import { useState } from 'react'
import { 
  Truck, 
  Shield, 
  CreditCard, 
  Headphones, 
  Award, 
  Clock,
  ArrowRight,
  CheckCircle
} from 'lucide-react'

interface Feature {
  id: string
  title: string
  description: string
  icon: React.ComponentType<any>
  color: string
  benefits: string[]
}

interface FeaturesSectionProps {
  title?: string
  subtitle?: string
  features?: Feature[]
}

export default function FeaturesSection({ 
  title = "Why Choose Us",
  subtitle = "We provide exceptional service and value to our customers",
  features
}: FeaturesSectionProps) {
  const [activeFeature, setActiveFeature] = useState(0)

  const defaultFeatures: Feature[] = [
    {
      id: 'shipping',
      title: 'Fast & Free Shipping',
      description: 'Get your orders delivered quickly with our free shipping on orders over $50. We partner with reliable carriers to ensure your packages arrive safely and on time.',
      icon: Truck,
      color: 'from-blue-500 to-blue-600',
      benefits: [
        'Free shipping on orders over $50',
        'Express delivery options available',
        'Real-time tracking updates',
        'Safe and secure packaging'
      ]
    },
    {
      id: 'security',
      title: 'Secure Payments',
      description: 'Your payment information is protected with bank-level security. We use industry-standard encryption and never store your payment details on our servers.',
      icon: Shield,
      color: 'from-green-500 to-green-600',
      benefits: [
        'SSL encryption for all transactions',
        'PCI DSS compliant payment processing',
        'Multiple payment methods accepted',
        'Fraud protection guarantee'
      ]
    },
    {
      id: 'support',
      title: '24/7 Customer Support',
      description: 'Our dedicated support team is available around the clock to help you with any questions or concerns. We pride ourselves on providing exceptional customer service.',
      icon: Headphones,
      color: 'from-purple-500 to-purple-600',
      benefits: [
        '24/7 live chat support',
        'Phone support during business hours',
        'Email support with quick response',
        'Comprehensive FAQ section'
      ]
    },
    {
      id: 'quality',
      title: 'Premium Quality',
      description: 'We carefully curate our product selection to ensure you receive only the highest quality items. Every product is tested and verified before being offered to our customers.',
      icon: Award,
      color: 'from-orange-500 to-orange-600',
      benefits: [
        'Rigorous quality control process',
        'Premium brand partnerships',
        'Quality guarantee on all products',
        'Regular product testing and reviews'
      ]
    },
    {
      id: 'returns',
      title: 'Easy Returns',
      description: 'Not satisfied with your purchase? No problem! We offer a hassle-free 30-day return policy with free return shipping for your convenience.',
      icon: Clock,
      color: 'from-pink-500 to-pink-600',
      benefits: [
        '30-day return policy',
        'Free return shipping',
        'Easy online return process',
        'Quick refund processing'
      ]
    },
    {
      id: 'rewards',
      title: 'Loyalty Rewards',
      description: 'Join our loyalty program and earn points with every purchase. Redeem your points for discounts, exclusive products, and special offers.',
      icon: CreditCard,
      color: 'from-indigo-500 to-indigo-600',
      benefits: [
        'Earn points on every purchase',
        'Exclusive member discounts',
        'Early access to sales',
        'Birthday and anniversary rewards'
      ]
    }
  ]

  const displayFeatures = features || defaultFeatures
  const currentFeature = displayFeatures[activeFeature]

  return (
    <section className="py-20 bg-white">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-16">
          <div className="inline-flex items-center px-4 py-2 bg-blue-100 rounded-full text-sm font-medium text-blue-600 mb-4">
            <Award className="w-4 h-4 mr-2" />
            Our Advantages
          </div>
          
          <h2 className="text-responsive-3xl font-bold text-gray-900 mb-4">
            {title}
          </h2>
          
          <p className="text-responsive-lg text-gray-600 max-w-2xl mx-auto">
            {subtitle}
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          {/* Features List */}
          <div className="space-y-4">
            {displayFeatures.map((feature, index) => (
              <div
                key={feature.id}
                onClick={() => setActiveFeature(index)}
                className={`p-6 rounded-2xl cursor-pointer transition-all duration-300 ${
                  activeFeature === index
                    ? 'bg-gradient-to-r from-blue-50 to-purple-50 border-2 border-blue-200 shadow-lg'
                    : 'bg-gray-50 hover:bg-gray-100 border-2 border-transparent hover:border-gray-200'
                }`}
              >
                <div className="flex items-start space-x-4">
                  <div className={`w-12 h-12 bg-gradient-to-r ${feature.color} rounded-xl flex items-center justify-center flex-shrink-0`}>
                    <feature.icon className="w-6 h-6 text-white" />
                  </div>
                  
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">
                      {feature.title}
                    </h3>
                    <p className="text-gray-600 text-sm line-clamp-2">
                      {feature.description}
                    </p>
                  </div>
                  
                  <ArrowRight className={`w-5 h-5 text-gray-400 transition-transform ${
                    activeFeature === index ? 'rotate-90' : ''
                  }`} />
                </div>
              </div>
            ))}
          </div>

          {/* Feature Detail */}
          <div className="relative">
            <div className="bg-gradient-to-br from-white to-gray-50 rounded-3xl p-8 shadow-2xl border border-gray-100">
              {/* Icon */}
              <div className={`w-20 h-20 bg-gradient-to-r ${currentFeature.color} rounded-2xl flex items-center justify-center mb-6`}>
                <currentFeature.icon className="w-10 h-10 text-white" />
              </div>

              {/* Title */}
              <h3 className="text-2xl font-bold text-gray-900 mb-4">
                {currentFeature.title}
              </h3>

              {/* Description */}
              <p className="text-gray-600 mb-6 leading-relaxed">
                {currentFeature.description}
              </p>

              {/* Benefits */}
              <div className="space-y-3">
                <h4 className="font-semibold text-gray-900 mb-3">Key Benefits:</h4>
                {currentFeature.benefits.map((benefit, index) => (
                  <div key={index} className="flex items-center space-x-3">
                    <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0" />
                    <span className="text-gray-700">{benefit}</span>
                  </div>
                ))}
              </div>

              {/* CTA */}
              <div className="mt-8">
                <button className="btn-primary">
                  Learn More
                  <ArrowRight className="w-4 h-4 ml-2" />
                </button>
              </div>
            </div>

            {/* Background Decoration */}
            <div className="absolute inset-0 bg-gradient-to-r from-blue-400/10 to-purple-400/10 rounded-3xl transform rotate-2 scale-105 -z-10"></div>
          </div>
        </div>

        {/* Bottom Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mt-16 pt-16 border-t border-gray-200">
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <Truck className="w-8 h-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-gray-900 mb-2">24h</h3>
            <p className="text-gray-600">Fast Delivery</p>
          </div>
          
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-green-500 to-green-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <Shield className="w-8 h-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-gray-900 mb-2">100%</h3>
            <p className="text-gray-600">Secure Payments</p>
          </div>
          
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-purple-500 to-purple-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <Headphones className="w-8 h-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-gray-900 mb-2">24/7</h3>
            <p className="text-gray-600">Support</p>
          </div>
          
          <div className="text-center">
            <div className="w-16 h-16 bg-gradient-to-r from-orange-500 to-orange-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <Award className="w-8 h-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-gray-900 mb-2">5★</h3>
            <p className="text-gray-600">Quality Rating</p>
          </div>
        </div>
      </div>
    </section>
  )
}
