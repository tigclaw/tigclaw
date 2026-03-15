import { Routes, Route } from 'react-router-dom'
import Layout from './components/Layout'
import SEO from './components/SEO'
import BlogPage from './pages/Blog'
import DocsPage from './pages/Docs'
import PricingPage from './pages/Pricing'
import ChangelogPage from './pages/Changelog'
import './App.css'

function HomePage() {
  return (
    <>
      <SEO />

      {/* Hero Section */}
      <header className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-orange-500/5 via-transparent to-transparent" />
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[600px] bg-orange-500/10 rounded-full blur-[120px] -translate-y-1/2" />

        <div className="relative z-10 max-w-4xl mx-auto px-6 pt-20 pb-32 text-center">
          <div className="inline-flex items-center gap-2 mb-6 px-4 py-1.5 rounded-full border border-zinc-800 bg-zinc-900/50 text-sm text-zinc-400">
            <span className="w-2 h-2 rounded-full bg-green-400 animate-pulse" />
            Open Source · 100% Local · Zero Cloud
          </div>

          <h1 className="text-5xl md:text-7xl font-black text-white leading-[1.1] tracking-tight mb-6">
            Your AI Gateway
            <br />
            <span className="bg-gradient-to-r from-orange-400 via-amber-300 to-orange-500 bg-clip-text text-transparent">
              is Running Naked
            </span>
          </h1>

          <p className="text-lg md:text-xl text-zinc-400 max-w-2xl mx-auto mb-10 leading-relaxed">
            Tigclaw is the zero-trust security gateway for OpenClaw.
            It shields your API keys, rate-limits abuse, and blocks prompt injection attacks
            — <span className="text-white font-medium">without sending a single byte to the cloud</span>.
          </p>

          <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-8">
            <div className="group relative">
              <div className="absolute -inset-1 bg-gradient-to-r from-orange-500 to-amber-500 rounded-xl blur opacity-25 group-hover:opacity-50 transition-opacity" />
              <code className="relative block bg-zinc-900 border border-zinc-700 rounded-xl px-6 py-3 text-sm md:text-base text-orange-300 font-mono cursor-pointer hover:border-orange-500/50 transition-colors"
                    onClick={() => navigator.clipboard.writeText('curl -sSL https://tigclaw.com/install.sh | bash')}>
                $ curl -sSL https://tigclaw.com/install.sh | bash
              </code>
            </div>
          </div>
          <p className="text-xs text-zinc-500">Click to copy · Works on Linux, macOS & Windows (WSL)</p>
        </div>
      </header>

      {/* Threat Section */}
      <section className="max-w-6xl mx-auto px-6 py-24">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">Real Threats. Real Damage.</h2>
          <p className="text-zinc-400 max-w-xl mx-auto">These aren't hypothetical risks — they've already happened to real users.</p>
        </div>
        <div className="grid md:grid-cols-2 gap-4">
          {[
            { icon: '💳', title: 'Credit Card Drained', desc: 'Port 3001 exposed to the internet. Bot scripts burned through $2,000 of OpenAI credits overnight.', severity: 'CRITICAL' },
            { icon: '🔑', title: 'API Keys Stolen', desc: '1.5M tokens leaked from plaintext config files. Your sk- key is one CVE away from theft.', severity: 'CRITICAL' },
            { icon: '🧠', title: 'Prompt Injection RCE', desc: 'DAN jailbreak exploits can write to SOUL.md, creating persistent backdoors on your server.', severity: 'HIGH' },
            { icon: '📋', title: 'Data Exfiltration', desc: 'Employees paste ID numbers and financial data into AI prompts. Compliance nightmare.', severity: 'HIGH' },
          ].map((threat, i) => (
            <div key={i} className="group bg-zinc-900/50 border border-zinc-800 rounded-2xl p-6 hover:border-red-500/30 transition-all hover:bg-zinc-900">
              <div className="flex items-start gap-4">
                <span className="text-3xl">{threat.icon}</span>
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <h3 className="text-lg font-semibold text-white">{threat.title}</h3>
                    <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${threat.severity === 'CRITICAL' ? 'bg-red-500/20 text-red-400' : 'bg-orange-500/20 text-orange-400'}`}>{threat.severity}</span>
                  </div>
                  <p className="text-zinc-400 text-sm leading-relaxed">{threat.desc}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Features Section */}
      <section className="max-w-6xl mx-auto px-6 py-24">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">Four Layers of Defense</h2>
          <p className="text-zinc-400 max-w-xl mx-auto">Tigclaw sits between the internet and your OpenClaw instance, filtering every request through military-grade protection.</p>
        </div>
        <div className="grid md:grid-cols-2 gap-6">
          {[
            { icon: '🔐', title: 'Zero-Trust Key Vault', desc: 'Real API keys never touch OpenClaw. Stored in AES-256-GCM encrypted vault with machine-bound keys. OpenClaw only sees disposable fake keys.', color: 'from-orange-500 to-amber-500' },
            { icon: '🚦', title: 'Anti-DoW Rate Limiter', desc: 'Token bucket algorithm cuts connections at the TCP level when abuse is detected. Your credit card survives the night.', color: 'from-red-500 to-pink-500' },
            { icon: '🧠', title: 'Local SLM Firewall', desc: 'Sub-50MB ONNX model runs locally to detect and block prompt injection, jailbreaks, and social engineering — in under 50ms.', color: 'from-purple-500 to-violet-500' },
            { icon: '🔏', title: 'Streaming DLP', desc: 'Real-time sensitive data masking in both directions. ID numbers and credit cards are tokenized before reaching the AI.', color: 'from-blue-500 to-cyan-500' },
          ].map((feature, i) => (
            <div key={i} className="relative group bg-zinc-900/50 border border-zinc-800 rounded-2xl p-8 hover:border-zinc-700 transition-all">
              <div className={`absolute top-0 left-0 right-0 h-px bg-gradient-to-r ${feature.color} opacity-0 group-hover:opacity-100 transition-opacity`} />
              <span className="text-4xl block mb-4">{feature.icon}</span>
              <h3 className="text-xl font-bold text-white mb-3">{feature.title}</h3>
              <p className="text-zinc-400 text-sm leading-relaxed">{feature.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Architecture */}
      <section className="max-w-4xl mx-auto px-6 py-24">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">How It Works</h2>
        </div>
        <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-8 font-mono text-sm text-zinc-400 overflow-x-auto">
          <pre className="whitespace-pre leading-relaxed">{`
  🌐 Internet Traffic (WhatsApp / Telegram / API)
         │
         ▼
  ┌────────────────────────────────────────┐
  │  🐯 Tigclaw Gateway  (:443 / :9000)   │
  │                                        │
  │  1. 🚦 Rate Limit (Token Bucket)       │
  │  2. 🔐 Key Swap (sk-tigclaw → sk-real) │
  │  3. 🔏 DLP Masking                     │
  │  4. 🧠 Prompt Firewall                 │
  └──────────────┬─────────────────────────┘
                 │  127.0.0.1 only
                 ▼
  ┌────────────────────────────────────────┐
  │  🐙 OpenClaw  (hidden from internet)   │
  │     config: "sk-tigclaw-a1b2c3d4"      │
  │     (fake key — useless if stolen)     │
  └────────────────────────────────────────┘
          `}</pre>
        </div>
      </section>

      {/* CTA */}
      <section className="max-w-4xl mx-auto px-6 py-24 text-center">
        <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">Stop Running Naked.</h2>
        <p className="text-zinc-400 mb-8 max-w-lg mx-auto">Your API keys deserve better than plaintext config files. Install Tigclaw in 30 seconds.</p>
        <a href="https://github.com/tigclaw/tigclaw" target="_blank" rel="noreferrer"
           className="inline-flex items-center gap-2 bg-white text-black font-bold px-8 py-4 rounded-xl text-lg hover:bg-zinc-200 transition-colors hover:shadow-lg hover:shadow-white/10">
          ⭐ Star on GitHub
        </a>
      </section>
    </>
  )
}

function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/blog" element={<BlogPage />} />
        <Route path="/docs" element={<DocsPage />} />
        <Route path="/pricing" element={<PricingPage />} />
        <Route path="/changelog" element={<ChangelogPage />} />
      </Routes>
    </Layout>
  )
}

export default App
