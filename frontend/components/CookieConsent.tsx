'use client'

import { useState, useEffect } from 'react'
import { Cookie, X, Settings, Check, AlertCircle } from 'lucide-react'

export default function CookieConsent() {
  const [isVisible, setIsVisible] = useState(false)
  const [showSettings, setShowSettings] = useState(false)
  const [preferences, setPreferences] = useState({
    necessary: true,
    analytics: false,
    marketing: false,
    preferences: false
  })

  useEffect(() => {
    const consent = localStorage.getItem('cookie-consent')
    if (!consent) {
      setIsVisible(true)
    }
  }, [])

  const handleAcceptAll = () => {
    const allPreferences = {
      necessary: true,
      analytics: true,
      marketing: true,
      preferences: true
    }
    setPreferences(allPreferences)
    localStorage.setItem('cookie-consent', JSON.stringify(allPreferences))
    setIsVisible(false)
  }

  const handleAcceptSelected = () => {
    localStorage.setItem('cookie-consent', JSON.stringify(preferences))
    setIsVisible(false)
    setShowSettings(false)
  }

  const handleRejectAll = () => {
    const minimalPreferences = {
      necessary: true,
      analytics: false,
      marketing: false,
      preferences: false
    }
    setPreferences(minimalPreferences)
    localStorage.setItem('cookie-consent', JSON.stringify(minimalPreferences))
    setIsVisible(false)
  }

  const togglePreference = (key: keyof typeof preferences) => {
    if (key === 'necessary') return // Can't disable necessary cookies
    setPreferences(prev => ({
      ...prev,
      [key]: !prev[key]
    }))
  }

  if (!isVisible) return null

  return (
    <div className="fixed bottom-0 left-0 right-0 z-50 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-2xl shadow-2xl border border-gray-200 overflow-hidden">
          {!showSettings ? (
            /* Main Consent Banner */
            <div className="p-6">
              <div className="flex items-start space-x-4">
                <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center flex-shrink-0">
                  <Cookie className="w-6 h-6 text-blue-600" />
                </div>
                
                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    We use cookies to enhance your experience
                  </h3>
                  <p className="text-gray-600 text-sm leading-relaxed">
                    We use cookies and similar technologies to provide, protect, and improve our services. 
                    Some cookies are necessary for our website to function, while others help us understand 
                    how you use our site so we can improve it.
                  </p>
                  
                  <div className="flex flex-col sm:flex-row gap-3 mt-4">
                    <button
                      onClick={handleAcceptAll}
                      className="bg-blue-600 text-white px-6 py-3 rounded-xl font-medium hover:bg-blue-700 transition-colors"
                    >
                      Accept All
                    </button>
                    <button
                      onClick={() => setShowSettings(true)}
                      className="bg-gray-100 text-gray-700 px-6 py-3 rounded-xl font-medium hover:bg-gray-200 transition-colors flex items-center"
                    >
                      <Settings className="w-4 h-4 mr-2" />
                      Customize
                    </button>
                    <button
                      onClick={handleRejectAll}
                      className="text-gray-600 px-6 py-3 rounded-xl font-medium hover:bg-gray-100 transition-colors"
                    >
                      Reject All
                    </button>
                  </div>
                </div>
                
                <button
                  onClick={handleRejectAll}
                  className="p-2 rounded-lg hover:bg-gray-100 transition-colors flex-shrink-0"
                >
                  <X className="w-5 h-5 text-gray-500" />
                </button>
              </div>
            </div>
          ) : (
            /* Settings Panel */
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold text-gray-900">
                  Cookie Preferences
                </h3>
                <button
                  onClick={() => setShowSettings(false)}
                  className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <X className="w-5 h-5 text-gray-500" />
                </button>
              </div>

              <div className="space-y-6">
                {/* Necessary Cookies */}
                <div className="flex items-start space-x-4">
                  <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center flex-shrink-0">
                    <Check className="w-5 h-5 text-green-600" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-semibold text-gray-900">Necessary Cookies</h4>
                      <div className="bg-green-100 text-green-800 text-xs px-2 py-1 rounded-full">
                        Always Active
                      </div>
                    </div>
                    <p className="text-gray-600 text-sm">
                      These cookies are essential for the website to function properly. 
                      They cannot be disabled and don't store any personal information.
                    </p>
                  </div>
                </div>

                {/* Analytics Cookies */}
                <div className="flex items-start space-x-4">
                  <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center flex-shrink-0">
                    <AlertCircle className="w-5 h-5 text-blue-600" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-semibold text-gray-900">Analytics Cookies</h4>
                      <button
                        onClick={() => togglePreference('analytics')}
                        className={`w-12 h-6 rounded-full transition-colors ${
                          preferences.analytics ? 'bg-blue-600' : 'bg-gray-300'
                        }`}
                      >
                        <div
                          className={`w-5 h-5 bg-white rounded-full transition-transform ${
                            preferences.analytics ? 'translate-x-6' : 'translate-x-0.5'
                          }`}
                        />
                      </button>
                    </div>
                    <p className="text-gray-600 text-sm">
                      These cookies help us understand how visitors interact with our website 
                      by collecting and reporting information anonymously.
                    </p>
                  </div>
                </div>

                {/* Marketing Cookies */}
                <div className="flex items-start space-x-4">
                  <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center flex-shrink-0">
                    <AlertCircle className="w-5 h-5 text-purple-600" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-semibold text-gray-900">Marketing Cookies</h4>
                      <button
                        onClick={() => togglePreference('marketing')}
                        className={`w-12 h-6 rounded-full transition-colors ${
                          preferences.marketing ? 'bg-blue-600' : 'bg-gray-300'
                        }`}
                      >
                        <div
                          className={`w-5 h-5 bg-white rounded-full transition-transform ${
                            preferences.marketing ? 'translate-x-6' : 'translate-x-0.5'
                          }`}
                        />
                      </button>
                    </div>
                    <p className="text-gray-600 text-sm">
                      These cookies are used to deliver advertisements more relevant to you 
                      and your interests across different websites.
                    </p>
                  </div>
                </div>

                {/* Preference Cookies */}
                <div className="flex items-start space-x-4">
                  <div className="w-10 h-10 bg-orange-100 rounded-lg flex items-center justify-center flex-shrink-0">
                    <Settings className="w-5 h-5 text-orange-600" />
                  </div>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <h4 className="font-semibold text-gray-900">Preference Cookies</h4>
                      <button
                        onClick={() => togglePreference('preferences')}
                        className={`w-12 h-6 rounded-full transition-colors ${
                          preferences.preferences ? 'bg-blue-600' : 'bg-gray-300'
                        }`}
                      >
                        <div
                          className={`w-5 h-5 bg-white rounded-full transition-transform ${
                            preferences.preferences ? 'translate-x-6' : 'translate-x-0.5'
                          }`}
                        />
                      </button>
                    </div>
                    <p className="text-gray-600 text-sm">
                      These cookies remember your choices and preferences to provide 
                      a more personalized experience.
                    </p>
                  </div>
                </div>
              </div>

              <div className="flex flex-col sm:flex-row gap-3 mt-8">
                <button
                  onClick={handleAcceptSelected}
                  className="bg-blue-600 text-white px-6 py-3 rounded-xl font-medium hover:bg-blue-700 transition-colors"
                >
                  Save Preferences
                </button>
                <button
                  onClick={handleAcceptAll}
                  className="bg-gray-100 text-gray-700 px-6 py-3 rounded-xl font-medium hover:bg-gray-200 transition-colors"
                >
                  Accept All
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
