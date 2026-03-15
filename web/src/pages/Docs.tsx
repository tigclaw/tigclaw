import SEO from '../components/SEO'

export default function DocsPage() {
  return (
    <>
      <SEO
        title="Documentation"
        description="Complete guide to installing, configuring, and managing Tigclaw — the zero-trust security gateway for OpenClaw."
        path="/docs"
      />

      <div className="max-w-5xl mx-auto px-6 py-20">
        <div className="mb-16">
          <h1 className="text-4xl md:text-5xl font-black text-white mb-4">Documentation</h1>
          <p className="text-lg text-zinc-400">
            Everything you need to secure your AI gateway in under 5 minutes.
          </p>
        </div>

        {/* Quick Start */}
        <section className="mb-16">
          <h2 className="text-2xl font-bold text-white mb-6 flex items-center gap-2">
            <span className="text-orange-400">01</span> Quick Start
          </h2>
          <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6 space-y-6">
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Install Tigclaw</h3>
              <pre className="bg-zinc-950 border border-zinc-800 rounded-xl p-4 text-sm overflow-x-auto">
                <code className="text-orange-300">{'# Linux / macOS\ncurl -sSL https://tigclaw.com/install.sh | bash\n\n# Windows (PowerShell)\nirm https://tigclaw.com/install.ps1 | iex'}</code>
              </pre>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Initialize & Auto-Migrate Keys</h3>
              <pre className="bg-zinc-950 border border-zinc-800 rounded-xl p-4 text-sm overflow-x-auto">
                <code className="text-green-300">{'$ tigclaw init\n🔍 Found OpenClaw config: ~/.openclaw/config.json\n🔑 Found plaintext key: sk-proj-abc...xyz\n✅ Migrated → sk-tigclaw-a1b2c3d4 (encrypted in vault)'}</code>
              </pre>
            </div>
            <div>
              <h3 className="text-lg font-semibold text-white mb-2">Start the Gateway</h3>
              <pre className="bg-zinc-950 border border-zinc-800 rounded-xl p-4 text-sm overflow-x-auto">
                <code className="text-cyan-300">{'$ tigclaw serve\n🐯 Tigclaw Security Gateway\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n  Listening on    : :9000\n  Upstream (OC)   : http://127.0.0.1:3001\n  Strict Mode     : true\n  Protected Keys  : 2'}</code>
              </pre>
            </div>
          </div>
        </section>

        {/* CLI Reference */}
        <section className="mb-16">
          <h2 className="text-2xl font-bold text-white mb-6 flex items-center gap-2">
            <span className="text-orange-400">02</span> CLI Reference
          </h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-800">
                  <th className="text-left py-3 pr-4 text-zinc-400 font-medium">Command</th>
                  <th className="text-left py-3 pl-4 text-zinc-400 font-medium">Description</th>
                </tr>
              </thead>
              <tbody className="text-zinc-300">
                {[
                  ['tigclaw serve', 'Start the security gateway'],
                  ['tigclaw init', 'Auto-migrate plaintext keys from OpenClaw'],
                  ['tigclaw keys add <key>', 'Add a new API key to the encrypted vault'],
                  ['tigclaw keys list', 'List all protected fake keys'],
                  ['tigclaw keys update <fake> <new>', 'Rotate the real key (OpenClaw needs no change)'],
                  ['tigclaw keys remove <fake>', 'Delete a key from vault'],
                  ['tigclaw status', 'Show security score and configuration'],
                  ['tigclaw version', 'Print version number'],
                ].map(([cmd, desc], i) => (
                  <tr key={i} className="border-b border-zinc-800/50 hover:bg-zinc-900/50">
                    <td className="py-3 pr-4"><code className="text-orange-300 bg-zinc-900 px-2 py-1 rounded text-xs">{cmd}</code></td>
                    <td className="py-3 pl-4 text-zinc-400">{desc}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* Configuration */}
        <section className="mb-16">
          <h2 className="text-2xl font-bold text-white mb-6 flex items-center gap-2">
            <span className="text-orange-400">03</span> Configuration
          </h2>
          <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6">
            <p className="text-zinc-400 mb-4">
              Tigclaw stores its config in <code className="text-orange-300 bg-zinc-950 px-2 py-0.5 rounded text-xs">~/.tigclaw/config.json</code>
            </p>
            <pre className="bg-zinc-950 border border-zinc-800 rounded-xl p-4 text-sm overflow-x-auto">
              <code className="text-zinc-300">{JSON.stringify({
                listen_addr: ":9000",
                upstream_addr: "http://127.0.0.1:3001",
                data_dir: "~/.tigclaw",
                strict_mode: true,
                rate_limit: 60
              }, null, 2)}</code>
            </pre>
            <div className="mt-6 space-y-3 text-sm">
              {[
                ['listen_addr', 'Gateway listening address (default :9000)'],
                ['upstream_addr', 'OpenClaw backend URL'],
                ['strict_mode', 'Block non-tigclaw keys from passing through'],
                ['rate_limit', 'Max requests per second per IP (0 = unlimited)'],
              ].map(([key, desc], i) => (
                <div key={i} className="flex gap-2">
                  <code className="text-orange-300 bg-zinc-950 px-2 py-0.5 rounded text-xs shrink-0">{key}</code>
                  <span className="text-zinc-400">{desc}</span>
                </div>
              ))}
            </div>
          </div>
        </section>
      </div>
    </>
  )
}
