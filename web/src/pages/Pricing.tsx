import SEO from '../components/SEO'

const plans = [
  {
    name: 'Community',
    price: 'Free',
    period: 'forever',
    description: 'For individual developers and small teams',
    features: [
      'Zero-Trust Key Substitution',
      'AES-256-GCM Encrypted Vault',
      'Token Bucket Rate Limiting',
      'Strict Mode Enforcement',
      'SSE Streaming Support',
      'CLI Management',
      'SQLite Local Database',
      'Community Support',
    ],
    cta: 'Get Started',
    ctaLink: '/docs',
    highlighted: false,
  },
  {
    name: 'Pro',
    price: '$19',
    period: '/month',
    description: 'For production deployments and teams',
    features: [
      'Everything in Community, plus:',
      '🧠 Local SLM Prompt Firewall',
      '🔏 Real-time DLP Masking',
      'Web Dashboard (SOC)',
      'Multi-Key Load Balancing',
      'Alerting & Notifications',
      'Token Usage Analytics',
      'Priority Email Support',
    ],
    cta: 'Coming Soon',
    ctaLink: '#',
    highlighted: true,
  },
  {
    name: 'Enterprise',
    price: 'Custom',
    period: '',
    description: 'For organizations with custom security needs',
    features: [
      'Everything in Pro, plus:',
      'SSO / LDAP Integration',
      'Compliance Audit Logs',
      'Custom SLM Model Training',
      'Multi-Gateway Clustering',
      'SLA & Dedicated Support',
      'On-premise Deployment',
      'White-label Option',
    ],
    cta: 'Contact Us',
    ctaLink: 'mailto:eelegshe@gmail.com',
    highlighted: false,
  },
]

export default function PricingPage() {
  return (
    <>
      <SEO
        title="Pricing"
        description="Tigclaw is free and open-source. Pro features available for production deployments starting at $19/month."
        path="/pricing"
      />

      <div className="max-w-6xl mx-auto px-6 py-20">
        <div className="text-center mb-16">
          <h1 className="text-4xl md:text-5xl font-black text-white mb-4">
            Simple, Transparent Pricing
          </h1>
          <p className="text-lg text-zinc-400 max-w-xl mx-auto">
            The core gateway is free forever. Pay only for advanced AI security features.
          </p>
        </div>

        <div className="grid md:grid-cols-3 gap-6 max-w-5xl mx-auto">
          {plans.map((plan, i) => (
            <div
              key={i}
              className={`relative rounded-2xl p-8 flex flex-col ${
                plan.highlighted
                  ? 'bg-gradient-to-b from-orange-500/10 to-zinc-900 border-2 border-orange-500/50'
                  : 'bg-zinc-900/50 border border-zinc-800'
              }`}
            >
              {plan.highlighted && (
                <div className="absolute -top-3 left-1/2 -translate-x-1/2 bg-orange-500 text-black text-xs font-bold px-3 py-1 rounded-full">
                  MOST POPULAR
                </div>
              )}
              <h3 className="text-lg font-bold text-white mb-1">{plan.name}</h3>
              <div className="mb-2">
                <span className="text-4xl font-black text-white">{plan.price}</span>
                {plan.period && <span className="text-zinc-400 text-sm">{plan.period}</span>}
              </div>
              <p className="text-sm text-zinc-400 mb-6">{plan.description}</p>

              <ul className="space-y-3 mb-8 flex-1">
                {plan.features.map((f, j) => (
                  <li key={j} className="flex items-start gap-2 text-sm text-zinc-300">
                    <span className="text-green-400 mt-0.5 shrink-0">✓</span>
                    {f}
                  </li>
                ))}
              </ul>

              <a
                href={plan.ctaLink}
                className={`block text-center py-3 rounded-xl font-semibold text-sm transition-all ${
                  plan.highlighted
                    ? 'bg-orange-500 hover:bg-orange-400 text-black hover:shadow-lg hover:shadow-orange-500/25'
                    : 'bg-zinc-800 hover:bg-zinc-700 text-white'
                }`}
              >
                {plan.cta}
              </a>
            </div>
          ))}
        </div>

        <div className="text-center mt-12 text-sm text-zinc-500">
          All plans include unlimited requests, unlimited keys, and full source code access.
        </div>
      </div>
    </>
  )
}
