export default function AdminHome() {
  return (
    <main className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-4xl font-bold text-center text-gray-900 mb-8">
          🔧 Admin Panel
        </h1>
        <div className="text-center">
          <p className="text-lg text-gray-600 mb-4">
            E-Commerce Platform Administration
          </p>
          <div className="bg-white rounded-lg shadow-md p-6 max-w-md mx-auto">
            <h2 className="text-xl font-semibold mb-2">Admin Dashboard</h2>
            <p className="text-gray-600">
              Manage products, orders, and users from here.
            </p>
          </div>
        </div>
      </div>
    </main>
  )
}
