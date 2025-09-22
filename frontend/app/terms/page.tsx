'use client'

import Navbar from '../../components/Navbar'

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="container mx-auto px-4 py-12">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white rounded-2xl shadow-lg p-8">
            <h1 className="text-4xl font-bold mb-8 text-center">Terms of Service</h1>
            <p className="text-gray-600 text-center mb-12">
              Last updated: {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
            </p>

            <div className="prose prose-lg max-w-none">
              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">1. Acceptance of Terms</h2>
                <p className="mb-4">
                  By accessing and using our website and services, you accept and agree to be bound by 
                  the terms and provision of this agreement. If you do not agree to abide by the above, 
                  please do not use this service.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">2. Use License</h2>
                <p className="mb-4">
                  Permission is granted to temporarily download one copy of the materials on our website 
                  for personal, non-commercial transitory viewing only. This is the grant of a license, 
                  not a transfer of title, and under this license you may not:
                </p>
                <ul className="list-disc pl-6 mb-4">
                  <li>Modify or copy the materials</li>
                  <li>Use the materials for any commercial purpose or for any public display</li>
                  <li>Attempt to reverse engineer any software contained on the website</li>
                  <li>Remove any copyright or other proprietary notations from the materials</li>
                </ul>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">3. User Accounts</h2>
                <p className="mb-4">
                  When you create an account with us, you must provide information that is accurate, 
                  complete, and current at all times. You are responsible for safeguarding the password 
                  and for all activities that occur under your account.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">4. Prohibited Uses</h2>
                <p className="mb-4">You may not use our service:</p>
                <ul className="list-disc pl-6 mb-4">
                  <li>For any unlawful purpose or to solicit others to perform unlawful acts</li>
                  <li>To violate any international, federal, provincial, or state regulations, rules, laws, or local ordinances</li>
                  <li>To infringe upon or violate our intellectual property rights or the intellectual property rights of others</li>
                  <li>To harass, abuse, insult, harm, defame, slander, disparage, intimidate, or discriminate</li>
                  <li>To submit false or misleading information</li>
                  <li>To upload or transmit viruses or any other type of malicious code</li>
                </ul>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">5. Product Information</h2>
                <p className="mb-4">
                  We strive to provide accurate product information, but we do not warrant that product 
                  descriptions or other content is accurate, complete, reliable, current, or error-free. 
                  If a product offered by us is not as described, your sole remedy is to return it.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">6. Pricing and Payment</h2>
                <p className="mb-4">
                  All prices are subject to change without notice. We reserve the right to modify or 
                  discontinue the service (or any part thereof) temporarily or permanently with or 
                  without notice. You agree to pay all charges incurred by you or any users of your 
                  account and credit card (or other applicable payment mechanism).
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">7. Returns and Refunds</h2>
                <p className="mb-4">
                  We offer a 30-day return policy for most items. Items must be in original condition 
                  with tags attached. Refunds will be processed within 5-7 business days after we 
                  receive the returned item. Some items may be excluded from our return policy.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">8. Disclaimer</h2>
                <p className="mb-4">
                  The information on this website is provided on an "as is" basis. To the fullest extent 
                  permitted by law, this Company excludes all representations, warranties, conditions 
                  and terms relating to our website and the use of this website.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">9. Limitations</h2>
                <p className="mb-4">
                  In no event shall our company or its suppliers be liable for any damages (including, 
                  without limitation, damages for loss of data or profit, or due to business interruption) 
                  arising out of the use or inability to use the materials on our website.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">10. Governing Law</h2>
                <p className="mb-4">
                  These terms and conditions are governed by and construed in accordance with the laws 
                  and you irrevocably submit to the exclusive jurisdiction of the courts in that state or location.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">11. Changes to Terms</h2>
                <p className="mb-4">
                  We reserve the right, at our sole discretion, to modify or replace these Terms at any time. 
                  If a revision is material, we will try to provide at least 30 days notice prior to any new 
                  terms taking effect.
                </p>
              </section>

              <section className="mb-8">
                <h2 className="text-2xl font-semibold mb-4">12. Contact Information</h2>
                <p className="mb-4">
                  If you have any questions about these Terms of Service, please contact us at:
                </p>
                <div className="bg-gray-50 p-4 rounded-lg">
                  <p><strong>Email:</strong> legal@ecommerce.com</p>
                  <p><strong>Phone:</strong> +1 (555) 123-4567</p>
                  <p><strong>Address:</strong> 123 Commerce Street, Business City, BC 12345</p>
                </div>
              </section>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
