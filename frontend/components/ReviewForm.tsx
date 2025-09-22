'use client'

import { useState } from 'react'
import { Star, Send, ThumbsUp, ThumbsDown, AlertCircle, CheckCircle } from 'lucide-react'

interface ReviewFormProps {
  productId: string
  onReviewSubmitted?: () => void
  userReview?: {
    id: string
    rating: number
    comment: string
    helpful?: boolean
  }
  onUpdate?: (review: { rating: number; comment: string; helpful?: boolean }) => void
}

export default function ReviewForm({ 
  productId, 
  onReviewSubmitted, 
  userReview,
  onUpdate 
}: ReviewFormProps) {
  const [rating, setRating] = useState(userReview?.rating || 0)
  const [comment, setComment] = useState(userReview?.comment || '')
  const [helpful, setHelpful] = useState<boolean | undefined>(userReview?.helpful)
  const [hoveredRating, setHoveredRating] = useState(0)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (rating === 0) {
      setError('Please select a rating')
      return
    }

    if (!comment.trim()) {
      setError('Please write a comment')
      return
    }

    try {
      setSubmitting(true)
      setError('')
      setSuccess('')

      const token = localStorage.getItem('token')
      if (!token) {
        setError('Please log in to write a review')
        return
      }

      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'
      const isEditing = !!userReview
      
      const response = await fetch(`${apiUrl}/api/reviews${isEditing ? `/${userReview.id}` : ''}`, {
        method: isEditing ? 'PUT' : 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          productId,
          rating,
          comment: comment.trim(),
          helpful
        })
      })

      if (response.ok) {
        setSuccess(isEditing ? 'Review updated successfully!' : 'Review submitted successfully!')
        if (!isEditing) {
          setRating(0)
          setComment('')
          setHelpful(undefined)
        }
        onReviewSubmitted?.()
        if (onUpdate) {
          onUpdate({ rating, comment: comment.trim(), helpful })
        }
      } else {
        const errorData = await response.json()
        setError(errorData.message || `Failed to ${isEditing ? 'update' : 'submit'} review`)
      }
    } catch (error) {
      setError('An error occurred. Please try again.')
    } finally {
      setSubmitting(false)
    }
  }

  const isEditing = !!userReview

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h3 className="text-lg font-semibold mb-4 flex items-center">
        {isEditing ? 'Edit Your Review' : 'Write a Review'}
        {isEditing && <span className="ml-2 text-sm bg-blue-100 text-blue-800 px-2 py-1 rounded-full">Editing</span>}
      </h3>
      
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4 flex items-center">
          <AlertCircle className="w-5 h-5 mr-2" />
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-4 flex items-center">
          <CheckCircle className="w-5 h-5 mr-2" />
          {success}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Rating */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-3">
            Rating *
          </label>
          <div className="flex space-x-1">
            {[1, 2, 3, 4, 5].map((star) => (
              <button
                key={star}
                type="button"
                onClick={() => setRating(star)}
                onMouseEnter={() => setHoveredRating(star)}
                onMouseLeave={() => setHoveredRating(0)}
                className="focus:outline-none transform hover:scale-110 transition-transform"
              >
                <Star
                  size={28}
                  className={`${
                    star <= (hoveredRating || rating)
                      ? 'text-yellow-400 fill-current'
                      : 'text-gray-300 hover:text-yellow-200'
                  } transition-colors`}
                />
              </button>
            ))}
          </div>
          {rating > 0 && (
            <p className="text-sm text-gray-600 mt-2">
              {rating === 1 && 'Poor'}
              {rating === 2 && 'Fair'}
              {rating === 3 && 'Good'}
              {rating === 4 && 'Very Good'}
              {rating === 5 && 'Excellent'}
            </p>
          )}
        </div>

        {/* Review Text */}
        <div>
          <label htmlFor="comment" className="block text-sm font-medium text-gray-700 mb-2">
            Your Review *
          </label>
          <textarea
            id="comment"
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            rows={5}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            placeholder="Share your detailed thoughts about this product..."
            maxLength={500}
            required
          />
          <p className="text-sm text-gray-500 mt-1">
            {comment.length}/500 characters
          </p>
        </div>

        {/* Helpful Rating */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-3">
            Was this product helpful for you?
          </label>
          <div className="flex space-x-4">
            <button
              type="button"
              onClick={() => setHelpful(true)}
              className={`flex items-center space-x-2 px-4 py-2 rounded-lg border transition-colors ${
                helpful === true
                  ? 'bg-green-50 border-green-200 text-green-700'
                  : 'bg-gray-50 border-gray-200 text-gray-600 hover:bg-green-50 hover:border-green-200'
              }`}
            >
              <ThumbsUp size={16} />
              <span>Yes</span>
            </button>
            <button
              type="button"
              onClick={() => setHelpful(false)}
              className={`flex items-center space-x-2 px-4 py-2 rounded-lg border transition-colors ${
                helpful === false
                  ? 'bg-red-50 border-red-200 text-red-700'
                  : 'bg-gray-50 border-gray-200 text-gray-600 hover:bg-red-50 hover:border-red-200'
              }`}
            >
              <ThumbsDown size={16} />
              <span>No</span>
            </button>
          </div>
        </div>

        {/* Submit Button */}
        <div className="flex justify-end space-x-3">
          {isEditing && (
            <button
              type="button"
              onClick={() => {
                setRating(userReview?.rating || 0)
                setComment(userReview?.comment || '')
                setHelpful(userReview?.helpful || undefined)
              }}
              className="px-6 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
          )}
          <button
            type="submit"
            disabled={submitting || rating === 0 || !comment.trim()}
            className="flex items-center space-x-2 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            <Send size={16} />
            <span>{submitting ? 'Submitting...' : isEditing ? 'Update Review' : 'Submit Review'}</span>
          </button>
        </div>
      </form>
    </div>
  )
}
