'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Calendar, User, ArrowRight, Tag, Clock, Eye } from 'lucide-react'

interface BlogPost {
  id: string
  title: string
  excerpt: string
  content: string
  image: string
  author: {
    name: string
    avatar: string
  }
  publishedAt: string
  readTime: number
  views: number
  tags: string[]
  category: string
}

interface BlogSectionProps {
  title?: string
  subtitle?: string
  posts?: BlogPost[]
  showViewAll?: boolean
  maxPosts?: number
}

export default function BlogSection({ 
  title = "Latest News & Updates",
  subtitle = "Stay informed with our latest insights and industry news",
  posts,
  showViewAll = true,
  maxPosts = 3
}: BlogSectionProps) {
  const [activePost, setActivePost] = useState(0)

  const defaultPosts: BlogPost[] = [
    {
      id: '1',
      title: 'The Future of E-Commerce: Trends to Watch in 2024',
      excerpt: 'Discover the latest trends shaping the e-commerce landscape and how they will impact your business.',
      content: 'E-commerce is evolving rapidly with new technologies and consumer behaviors...',
      image: 'https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      author: {
        name: 'Sarah Johnson',
        avatar: 'https://images.unsplash.com/photo-1494790108755-2616b612b786?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80'
      },
      publishedAt: '2024-01-15',
      readTime: 5,
      views: 1250,
      tags: ['E-commerce', 'Technology', 'Trends'],
      category: 'Business'
    },
    {
      id: '2',
      title: 'Sustainable Shopping: How to Make Eco-Friendly Choices',
      excerpt: 'Learn how to make more sustainable shopping decisions and reduce your environmental impact.',
      content: 'Sustainability is becoming increasingly important in consumer decision-making...',
      image: 'https://images.unsplash.com/photo-1542601906990-b4d3fb778b09?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      author: {
        name: 'Michael Chen',
        avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80'
      },
      publishedAt: '2024-01-12',
      readTime: 7,
      views: 890,
      tags: ['Sustainability', 'Environment', 'Shopping'],
      category: 'Lifestyle'
    },
    {
      id: '3',
      title: 'Product Review: Top 10 Gadgets of 2024',
      excerpt: 'Our comprehensive review of the most innovative and useful gadgets released this year.',
      content: 'We\'ve tested hundreds of gadgets to bring you the best of 2024...',
      image: 'https://images.unsplash.com/photo-1498049794561-7780e7231661?ixlib=rb-4.0.3&auto=format&fit=crop&w=1000&q=80',
      author: {
        name: 'Emily Rodriguez',
        avatar: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixlib=rb-4.0.3&auto=format&fit=crop&w=150&q=80'
      },
      publishedAt: '2024-01-10',
      readTime: 12,
      views: 2100,
      tags: ['Gadgets', 'Reviews', 'Technology'],
      category: 'Reviews'
    }
  ]

  const displayPosts = posts || defaultPosts
  const featuredPost = displayPosts[activePost]
  const otherPosts = displayPosts.filter((_, index) => index !== activePost).slice(0, maxPosts - 1)

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }

  return (
    <section className="py-20 bg-white">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-16">
          <div className="inline-flex items-center px-4 py-2 bg-blue-100 rounded-full text-sm font-medium text-blue-600 mb-4">
            <Tag className="w-4 h-4 mr-2" />
            Latest Articles
          </div>
          
          <h2 className="text-responsive-3xl font-bold text-gray-900 mb-4">
            {title}
          </h2>
          
          <p className="text-responsive-lg text-gray-600 max-w-2xl mx-auto">
            {subtitle}
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Featured Post */}
          <div className="lg:col-span-2">
            <div className="group bg-white rounded-3xl shadow-2xl overflow-hidden hover:shadow-3xl transition-all duration-300">
              <div className="relative aspect-video overflow-hidden">
                <img
                  src={featuredPost.image}
                  alt={featuredPost.title}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
                
                {/* Category Badge */}
                <div className="absolute top-6 left-6 px-3 py-1 bg-blue-600 text-white text-sm font-medium rounded-full">
                  {featuredPost.category}
                </div>
                
                {/* Read Time */}
                <div className="absolute top-6 right-6 px-3 py-1 bg-white/90 backdrop-blur-sm text-gray-700 text-sm font-medium rounded-full flex items-center">
                  <Clock className="w-4 h-4 mr-1" />
                  {featuredPost.readTime} min read
                </div>
              </div>
              
              <div className="p-8">
                {/* Tags */}
                <div className="flex flex-wrap gap-2 mb-4">
                  {featuredPost.tags.map((tag) => (
                    <span
                      key={tag}
                      className="px-3 py-1 bg-gray-100 text-gray-600 text-sm rounded-full"
                    >
                      {tag}
                    </span>
                  ))}
                </div>
                
                {/* Title */}
                <h3 className="text-2xl font-bold text-gray-900 mb-4 group-hover:text-blue-600 transition-colors">
                  {featuredPost.title}
                </h3>
                
                {/* Excerpt */}
                <p className="text-gray-600 mb-6 leading-relaxed">
                  {featuredPost.excerpt}
                </p>
                
                {/* Author & Meta */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <img
                      src={featuredPost.author.avatar}
                      alt={featuredPost.author.name}
                      className="w-10 h-10 rounded-full object-cover mr-3"
                    />
                    <div>
                      <p className="font-medium text-gray-900">{featuredPost.author.name}</p>
                      <p className="text-sm text-gray-500">{formatDate(featuredPost.publishedAt)}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center text-gray-500">
                    <Eye className="w-4 h-4 mr-1" />
                    <span className="text-sm">{featuredPost.views}</span>
                  </div>
                </div>
                
                {/* CTA */}
                <div className="mt-6">
                  <Link
                    href={`/blog/${featuredPost.id}`}
                    className="group inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
                  >
                    Read More
                    <ArrowRight className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform" />
                  </Link>
                </div>
              </div>
            </div>
          </div>

          {/* Other Posts */}
          <div className="space-y-6">
            {otherPosts.map((post, index) => (
              <div
                key={post.id}
                className="group bg-white rounded-2xl shadow-lg overflow-hidden hover:shadow-xl transition-all duration-300 cursor-pointer"
                onClick={() => setActivePost(displayPosts.findIndex(p => p.id === post.id))}
              >
                <div className="relative aspect-video overflow-hidden">
                  <img
                    src={post.image}
                    alt={post.title}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  />
                  <div className="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"></div>
                  
                  {/* Category Badge */}
                  <div className="absolute top-3 left-3 px-2 py-1 bg-white/90 backdrop-blur-sm text-gray-700 text-xs font-medium rounded-full">
                    {post.category}
                  </div>
                </div>
                
                <div className="p-6">
                  {/* Title */}
                  <h4 className="font-bold text-gray-900 mb-2 line-clamp-2 group-hover:text-blue-600 transition-colors">
                    {post.title}
                  </h4>
                  
                  {/* Excerpt */}
                  <p className="text-gray-600 text-sm mb-4 line-clamp-2">
                    {post.excerpt}
                  </p>
                  
                  {/* Meta */}
                  <div className="flex items-center justify-between text-sm text-gray-500">
                    <div className="flex items-center">
                      <User className="w-4 h-4 mr-1" />
                      <span>{post.author.name}</span>
                    </div>
                    <div className="flex items-center">
                      <Clock className="w-4 h-4 mr-1" />
                      <span>{post.readTime}m</span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* View All Button */}
        {showViewAll && (
          <div className="text-center mt-12">
            <Link
              href="/blog"
              className="group inline-flex items-center px-8 py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-semibold rounded-2xl hover:from-blue-700 hover:to-purple-700 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-1"
            >
              View All Articles
              <ArrowRight className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform" />
            </Link>
          </div>
        )}
      </div>
    </section>
  )
}
