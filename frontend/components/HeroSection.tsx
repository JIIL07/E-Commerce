'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { ArrowRight, Play, Star, Users, Award, Truck } from 'lucide-react'

export default function HeroSection() {
  const [currentSlide, setCurrentSlide] = useState(0)

  const slides = [
    {
      id: 1,
      title: "Discover Amazing Products",
      subtitle: "Shop the latest trends with unbeatable prices",
      description: "From electronics to fashion, find everything you need in one place with fast shipping and excellent customer service.",
      image: "https://images.unsplash.com/photo-1441986300917-64674bd600d8?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80",
      cta: "Shop Now",
      ctaLink: "/products",
      stats: [
        { label: "Happy Customers", value: "50K+", icon: Users },
        { label: "Products", value: "10K+", icon: Award },
        { label: "Fast Delivery", value: "24h", icon: Truck }
      ]
    },
    {
      id: 2,
      title: "Premium Quality",
      subtitle: "Curated selection of top brands",
      description: "We partner with the best brands to bring you premium quality products at competitive prices with guaranteed satisfaction.",
      image: "https://images.unsplash.com/photo-1472851294608-062f824d29cc?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80",
      cta: "Explore Brands",
      ctaLink: "/categories",
      stats: [
        { label: "Brands", value: "500+", icon: Award },
        { label: "Reviews", value: "25K+", icon: Star },
        { label: "Satisfaction", value: "99%", icon: Users }
      ]
    },
    {
      id: 3,
      title: "Smart Shopping",
      subtitle: "AI-powered recommendations",
      description: "Get personalized product recommendations based on your preferences and shopping history for the ultimate shopping experience.",
      image: "https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?ixlib=rb-4.0.3&auto=format&fit=crop&w=2070&q=80",
      cta: "Get Started",
      ctaLink: "/register",
      stats: [
        { label: "AI Accuracy", value: "95%", icon: Star },
        { label: "Time Saved", value: "2h", icon: Truck },
        { label: "Satisfaction", value: "98%", icon: Users }
      ]
    }
  ]

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentSlide((prev) => (prev + 1) % slides.length)
    }, 8000)

    return () => clearInterval(timer)
  }, [slides.length])

  const currentSlideData = slides[currentSlide]

  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100">
      {/* Animated Background */}
      <div className="absolute inset-0">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-600/10 via-purple-600/10 to-pink-600/10"></div>
        <div className="absolute top-0 left-0 w-full h-full">
          <div className="absolute top-20 left-20 w-72 h-72 bg-blue-400/20 rounded-full blur-3xl animate-float"></div>
          <div className="absolute top-40 right-20 w-96 h-96 bg-purple-400/20 rounded-full blur-3xl animate-float" style={{ animationDelay: '1s' }}></div>
          <div className="absolute bottom-20 left-1/3 w-80 h-80 bg-pink-400/20 rounded-full blur-3xl animate-float" style={{ animationDelay: '2s' }}></div>
        </div>
      </div>

      <div className="relative z-10 container mx-auto px-4 py-20">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          {/* Content */}
          <div className="text-center lg:text-left space-y-8">
            <div className="space-y-4">
              <div className="inline-flex items-center px-4 py-2 bg-white/80 backdrop-blur-sm rounded-full text-sm font-medium text-blue-600 border border-blue-200">
                <Star className="w-4 h-4 mr-2 fill-current" />
                Trusted by 50,000+ customers
              </div>
              
              <h1 className="text-responsive-4xl font-bold text-gray-900 leading-tight">
                {currentSlideData.title}
              </h1>
              
              <h2 className="text-responsive-xl text-blue-600 font-semibold">
                {currentSlideData.subtitle}
              </h2>
              
              <p className="text-responsive-lg text-gray-600 max-w-2xl mx-auto lg:mx-0">
                {currentSlideData.description}
              </p>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-3 gap-6 py-8">
              {currentSlideData.stats.map((stat, index) => (
                <div key={index} className="text-center lg:text-left">
                  <div className="flex items-center justify-center lg:justify-start mb-2">
                    <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-500 rounded-lg flex items-center justify-center mr-2">
                      <stat.icon className="w-4 h-4 text-white" />
                    </div>
                    <span className="text-2xl font-bold text-gray-900">{stat.value}</span>
                  </div>
                  <p className="text-sm text-gray-600">{stat.label}</p>
                </div>
              ))}
            </div>

            {/* CTA Buttons */}
            <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
              <Link
                href={currentSlideData.ctaLink}
                className="group btn-primary text-lg px-8 py-4 inline-flex items-center justify-center"
              >
                {currentSlideData.cta}
                <ArrowRight className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
              </Link>
              
              <button className="group btn-secondary text-lg px-8 py-4 inline-flex items-center justify-center">
                <Play className="w-5 h-5 mr-2 group-hover:scale-110 transition-transform" />
                Watch Demo
              </button>
            </div>
          </div>

          {/* Image */}
          <div className="relative">
            <div className="relative z-10">
              <div className="aspect-square rounded-3xl overflow-hidden shadow-2xl">
                <img
                  src={currentSlideData.image}
                  alt={currentSlideData.title}
                  className="w-full h-full object-cover transition-transform duration-1000 hover:scale-105"
                />
              </div>
              
              {/* Floating Cards */}
              <div className="absolute -top-6 -left-6 bg-white/90 backdrop-blur-sm rounded-2xl p-4 shadow-xl border border-white/20">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 bg-gradient-to-r from-green-400 to-green-500 rounded-xl flex items-center justify-center">
                    <Truck className="w-6 h-6 text-white" />
                  </div>
                  <div>
                    <p className="font-semibold text-gray-900">Free Shipping</p>
                    <p className="text-sm text-gray-600">On orders over $50</p>
                  </div>
                </div>
              </div>

              <div className="absolute -bottom-6 -right-6 bg-white/90 backdrop-blur-sm rounded-2xl p-4 shadow-xl border border-white/20">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 bg-gradient-to-r from-blue-400 to-blue-500 rounded-xl flex items-center justify-center">
                    <Award className="w-6 h-6 text-white" />
                  </div>
                  <div>
                    <p className="font-semibold text-gray-900">Best Price</p>
                    <p className="text-sm text-gray-600">Guaranteed</p>
                  </div>
                </div>
              </div>
            </div>

            {/* Background Decoration */}
            <div className="absolute inset-0 bg-gradient-to-r from-blue-400/20 to-purple-400/20 rounded-3xl transform rotate-3 scale-105 -z-10"></div>
          </div>
        </div>

        {/* Slide Indicators */}
        <div className="flex justify-center mt-12 space-x-3">
          {slides.map((_, index) => (
            <button
              key={index}
              onClick={() => setCurrentSlide(index)}
              className={`w-3 h-3 rounded-full transition-all duration-300 ${
                index === currentSlide
                  ? 'bg-blue-600 scale-125'
                  : 'bg-gray-300 hover:bg-gray-400'
              }`}
              aria-label={`Go to slide ${index + 1}`}
            />
          ))}
        </div>
      </div>
    </section>
  )
}
