import { useParams, Navigate, Link } from 'react-router-dom'
import { Helmet } from 'react-helmet-async'
import ReactMarkdown from 'react-markdown'
import { blogPosts } from '../data/blogData'
import SEO from '../components/SEO'

export default function BlogDetail() {
  const { slug } = useParams<{ slug: string }>()
  const post = blogPosts.find(p => p.slug === slug)

  if (!post) {
    return <Navigate to="/blog" replace />
  }

  return (
    <>
      <SEO 
        title={post.title}
        description={post.excerpt}
        path={`/blog/${post.slug}`}
        type="article"
      />
      {/* Article Schema Data for Google Rich Snippets */}
      <Helmet>
        <script type="application/ld+json">
          {JSON.stringify({
            "@context": "https://schema.org",
            "@type": "Article",
            "headline": post.title,
            "description": post.excerpt,
            "datePublished": post.date,
            "author": {
              "@type": "Organization",
              "name": "Tigclaw Security Labs"
            }
          })}
        </script>
      </Helmet>

      <div className="max-w-3xl mx-auto px-6 py-20">
        <Link to="/blog" className="inline-flex items-center text-sm text-zinc-400 hover:text-white mb-8 transition-colors">
          <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Back to Blog
        </Link>
        
        <header className="mb-12 pb-8 border-b border-zinc-800">
          <div className="flex items-center gap-3 mb-6">
            <span className={`text-xs px-2 py-1 rounded-full font-medium ${post.tagColor}`}>
              {post.tag}
            </span>
            <span className="text-sm text-zinc-500">{post.date}</span>
            <span className="text-sm text-zinc-600">·</span>
            <span className="text-sm text-zinc-500">{post.readTime} read</span>
          </div>
          <h1 className="text-3xl md:text-5xl font-black text-white leading-tight mb-6">
            {post.title}
          </h1>
          <p className="text-xl text-zinc-400 leading-relaxed">
            {post.excerpt}
          </p>
        </header>

        <article className="prose prose-invert prose-orange max-w-none 
            prose-headings:font-bold prose-h2:text-2xl prose-h2:mt-12 prose-h2:mb-6
            prose-h3:text-xl prose-h3:mt-8 prose-h3:mb-4
            prose-p:text-zinc-300 prose-p:leading-relaxed prose-p:mb-6
            prose-ul:text-zinc-300 prose-ul:mb-6 prose-li:mb-2
            prose-strong:text-white prose-strong:font-semibold
            prose-code:text-orange-300 prose-code:bg-zinc-900 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-code:before:content-none prose-code:after:content-none">
          <ReactMarkdown>{post.content}</ReactMarkdown>
        </article>

        <div className="mt-16 pt-8 border-t border-zinc-800 text-center">
          <h3 className="text-2xl font-bold text-white mb-4">Secure your AI Infrastructure Today</h3>
          <p className="text-zinc-400 mb-8 max-w-lg mx-auto">
            Don't wait for your API keys to leak. Install the Tigclaw zero-trust gateway in under 60 seconds.
          </p>
          <a href="https://github.com/tigclaw/tigclaw" target="_blank" rel="noreferrer"
             className="inline-block bg-orange-500 hover:bg-orange-400 text-black font-semibold px-8 py-3 rounded-xl transition-all hover:shadow-lg hover:shadow-orange-500/25">
            Get Tigclaw on GitHub
          </a>
        </div>
      </div>
    </>
  )
}
