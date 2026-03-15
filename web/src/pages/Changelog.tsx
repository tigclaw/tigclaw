import SEO from '../components/SEO'

const releases = [
  {
    version: 'v0.1.0',
    date: '2026-03-15',
    tag: 'Latest',
    tagColor: 'bg-green-500/20 text-green-400',
    changes: [
      { type: 'feat', text: 'Zero-Trust Key Substitution Engine — AES-256-GCM encrypted vault with machine-bound key derivation' },
      { type: 'feat', text: 'Reverse proxy gateway with SSE streaming support (FlushInterval: -1)' },
      { type: 'feat', text: 'Token Bucket rate limiter with per-IP tracking and auto-cleanup' },
      { type: 'feat', text: 'Strict Mode — blocks real API keys from leaking through OpenClaw' },
      { type: 'feat', text: 'CLI commands: serve, init, keys (add/list/update/remove), status, version' },
      { type: 'feat', text: 'Auto-migration: tigclaw init scans OpenClaw configs and replaces plaintext keys' },
      { type: 'feat', text: 'tigclaw-scan.sh vulnerability scanner for market validation' },
      { type: 'feat', text: 'Official website with SEO, blog, docs, and pricing pages' },
    ],
  },
]

export default function ChangelogPage() {
  return (
    <>
      <SEO
        title="Changelog"
        description="Release notes and version history for Tigclaw. Track new features, improvements, and security updates."
        path="/changelog"
      />

      <div className="max-w-3xl mx-auto px-6 py-20">
        <div className="mb-16">
          <h1 className="text-4xl md:text-5xl font-black text-white mb-4">Changelog</h1>
          <p className="text-lg text-zinc-400">
            All notable changes to Tigclaw are documented here.
          </p>
        </div>

        <div className="space-y-12">
          {releases.map((release, i) => (
            <div key={i} className="relative">
              {/* Version Header */}
              <div className="flex items-center gap-3 mb-6">
                <h2 className="text-2xl font-bold text-white">{release.version}</h2>
                <span className={`text-xs px-2 py-1 rounded-full font-medium ${release.tagColor}`}>
                  {release.tag}
                </span>
                <span className="text-sm text-zinc-500">{release.date}</span>
              </div>

              {/* Changes List */}
              <div className="space-y-3 pl-4 border-l-2 border-zinc-800">
                {release.changes.map((change, j) => (
                  <div key={j} className="flex items-start gap-3">
                    <span className={`shrink-0 text-xs px-2 py-0.5 rounded font-mono font-medium mt-0.5 ${
                      change.type === 'feat' ? 'bg-green-500/20 text-green-400' :
                      change.type === 'fix' ? 'bg-blue-500/20 text-blue-400' :
                      'bg-yellow-500/20 text-yellow-400'
                    }`}>
                      {change.type}
                    </span>
                    <span className="text-sm text-zinc-300">{change.text}</span>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  )
}
