export default function Footer() {
  return (
    <footer className="bg-gray-800 text-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <p className="text-gray-400">
            © 2024 E-Commerce Platform. Built with Next.js, Go, and PostgreSQL.
          </p>
          <div className="mt-4 flex justify-center space-x-6">
            <a href="/api/docs" className="text-gray-400 hover:text-white">
              API Docs
            </a>
            <a href="http://localhost:9090" className="text-gray-400 hover:text-white">
              Prometheus
            </a>
            <a href="http://localhost:3001" className="text-gray-400 hover:text-white">
              Grafana
            </a>
          </div>
        </div>
      </div>
    </footer>
  )
}