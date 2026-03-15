import { Link, useLocation } from 'react-router-dom'
import { useState } from 'react'

const navLinks = [
  { label: 'Docs', href: '/docs' },
  { label: 'Blog', href: '/blog' },
  { label: 'Changelog', href: '/changelog' },
]

export default function Layout({ children }: { children: React.ReactNode }) {
  const location = useLocation()
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <div className="min-h-screen flex flex-col">
      {/* Navbar */}
      <nav className="sticky top-0 z-50 border-b border-zinc-800/50 bg-zinc-950/80 backdrop-blur-xl">
        <div className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
          <Link to="/" className="flex items-center gap-2 group">
            <span className="text-2xl group-hover:scale-110 transition-transform">🐯</span>
            <span className="text-xl font-bold text-white tracking-tight">Tigclaw</span>
          </Link>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-6">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                to={link.href}
                className={`text-sm transition-colors ${
                  location.pathname === link.href
                    ? 'text-orange-400 font-medium'
                    : 'text-zinc-400 hover:text-white'
                }`}
              >
                {link.label}
              </Link>
            ))}
            <a
              href="https://github.com/tigclaw/tigclaw"
              target="_blank"
              rel="noreferrer"
              className="text-sm text-zinc-400 hover:text-white transition-colors"
            >
              GitHub
            </a>
            <Link
              to="/docs"
              className="text-sm bg-orange-500 hover:bg-orange-400 text-black font-semibold px-4 py-2 rounded-lg transition-all hover:shadow-lg hover:shadow-orange-500/25"
            >
              Get Started
            </Link>
          </div>

          {/* Mobile Hamburger */}
          <button
            className="md:hidden text-zinc-400 hover:text-white"
            onClick={() => setMobileOpen(!mobileOpen)}
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {mobileOpen ? (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              ) : (
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              )}
            </svg>
          </button>
        </div>

        {/* Mobile Menu */}
        {mobileOpen && (
          <div className="md:hidden border-t border-zinc-800 bg-zinc-950 px-6 py-4 space-y-3">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                to={link.href}
                onClick={() => setMobileOpen(false)}
                className="block text-sm text-zinc-400 hover:text-white"
              >
                {link.label}
              </Link>
            ))}
            <a href="https://github.com/tigclaw/tigclaw" target="_blank" rel="noreferrer"
               className="block text-sm text-zinc-400 hover:text-white">
              GitHub
            </a>
          </div>
        )}
      </nav>

      {/* Main Content */}
      <main className="flex-1">{children}</main>

      {/* Footer */}
      <footer className="border-t border-zinc-800 bg-zinc-950">
        <div className="max-w-6xl mx-auto px-6 py-12">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8 mb-8">
            <div>
              <h4 className="text-white font-semibold text-sm mb-3">Product</h4>
              <div className="space-y-2">
                <Link to="/docs" className="block text-sm text-zinc-500 hover:text-zinc-300 transition-colors">Documentation</Link>

                <Link to="/changelog" className="block text-sm text-zinc-500 hover:text-zinc-300 transition-colors">Changelog</Link>
              </div>
            </div>
            <div>
              <h4 className="text-white font-semibold text-sm mb-3">Resources</h4>
              <div className="space-y-2">
                <Link to="/blog" className="block text-sm text-zinc-500 hover:text-zinc-300 transition-colors">Blog</Link>
                <a href="https://github.com/tigclaw/tigclaw" target="_blank" rel="noreferrer" className="block text-sm text-zinc-500 hover:text-zinc-300 transition-colors">GitHub</a>
                <a href="https://github.com/tigclaw/tigclaw/discussions" target="_blank" rel="noreferrer" className="block text-sm text-zinc-500 hover:text-zinc-300 transition-colors">Community</a>
              </div>
            </div>
            <div>
              <h4 className="text-white font-semibold text-sm mb-3">Security</h4>
              <div className="space-y-2">
                <Link to="/blog/openclaw-vulnerability-scan" className="block text-sm text-zinc-500 hover:text-zinc-300 transition-colors">Vulnerability Report</Link>
                <a href="https://github.com/tigclaw/tigclaw/security" target="_blank" rel="noreferrer" className="block text-sm text-zinc-500 hover:text-zinc-300 transition-colors">Security Policy</a>
              </div>
            </div>
            <div>
              <h4 className="text-white font-semibold text-sm mb-3">Legal</h4>
              <div className="space-y-2">
                <span className="block text-sm text-zinc-500">MIT License</span>
              </div>
            </div>
          </div>
          <div className="border-t border-zinc-800 pt-8 flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <span>🐯</span>
              <span className="text-sm text-zinc-500">Tigclaw — Zero-Trust AI Security Gateway</span>
            </div>
            <p className="text-xs text-zinc-600">100% Open Source · 100% Local · Zero Cloud</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
