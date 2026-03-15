import { Link } from 'react-router-dom'
import SEO from '../components/SEO'
import { blogPosts } from '../data/blogData'

export default function BlogPage() {
  return (
    <>
      <SEO
        title="Blog"
        description="Security research, vulnerability analysis, and technical deep-dives on protecting self-hosted AI platforms."
        path="/blog"
      />

      <div className="max-w-4xl mx-auto px-6 py-20">
        <div className="mb-16">
          <h1 className="text-4xl md:text-5xl font-black text-white mb-4">Blog</h1>
          <p className="text-lg text-zinc-400">
            Security research, vulnerability analysis, and technical deep-dives.
          </p>
        </div>

        <div className="space-y-6">
          {blogPosts.map((post) => (
            <Link
              key={post.slug}
              to={`/blog/${post.slug}`}
              className="group block bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6 md:p-8 hover:border-zinc-700 transition-all"
            >
              <div className="flex items-center gap-3 mb-3">
                <span className={`text-xs px-2 py-1 rounded-full font-medium ${post.tagColor}`}>
                  {post.tag}
                </span>
                <span className="text-xs text-zinc-500">{post.date}</span>
                <span className="text-xs text-zinc-600">·</span>
                <span className="text-xs text-zinc-500">{post.readTime} read</span>
              </div>
              <h2 className="text-xl md:text-2xl font-bold text-white group-hover:text-orange-400 transition-colors mb-2">
                {post.title}
              </h2>
              <p className="text-zinc-400 text-sm leading-relaxed">{post.excerpt}</p>
            </Link>
          ))}
        </div>
      </div>
    </>
  )
}
