import { Link } from 'react-router-dom'
import SEO from '../components/SEO'

const posts = [
  {
    slug: 'openclaw-vulnerability-scan',
    title: 'Your OpenClaw Instance is Probably Hacked Right Now — Here\'s How to Check',
    excerpt: 'We scanned thousands of self-hosted OpenClaw instances. 73% had exposed API keys in plaintext config files. Is yours one of them?',
    date: '2026-03-15',
    readTime: '8 min',
    tag: 'Security',
    tagColor: 'bg-red-500/20 text-red-400',
  },
  {
    slug: 'zero-trust-key-substitution',
    title: 'Zero-Trust Key Substitution: How Tigclaw Protects Your $10,000 API Budget',
    excerpt: 'A deep-dive into AES-256-GCM encrypted vaults, machine-bound key derivation, and in-memory-only real key substitution.',
    date: '2026-03-14',
    readTime: '12 min',
    tag: 'Architecture',
    tagColor: 'bg-blue-500/20 text-blue-400',
  },
  {
    slug: 'prompt-injection-defense',
    title: 'Blocking DAN Attacks with a 50MB Local AI Model — No Cloud Required',
    excerpt: 'How Tigclaw\'s SLM Firewall detects and neutralizes prompt injection, jailbreak, and social engineering attacks in under 50ms.',
    date: '2026-03-12',
    readTime: '10 min',
    tag: 'AI Security',
    tagColor: 'bg-purple-500/20 text-purple-400',
  },
  {
    slug: 'anti-dow-rate-limiting',
    title: 'Anti-DoW: Stop Botnets From Draining Your OpenAI Credits Overnight',
    excerpt: 'Token bucket rate limiting at the TCP level — how Tigclaw saves your credit card from Denial-of-Wallet attacks.',
    date: '2026-03-10',
    readTime: '6 min',
    tag: 'Defense',
    tagColor: 'bg-orange-500/20 text-orange-400',
  },
]

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
          {posts.map((post) => (
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
