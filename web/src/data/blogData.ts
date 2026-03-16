import post1 from './blog/openclaw-api-leak-cve-2026-25253.md?raw'
import post2 from './blog/preventing-denial-of-wallet-ai-billing.md?raw'
import post3 from './blog/prompt-injection-defense-llm-jailbreak.md?raw'
import post4 from './blog/zero-trust-ai-gateway-architecture.md?raw'
import post5 from './blog/open-source-ai-gateway-comparison-2026.md?raw'
import post6 from './blog/clawhub-skills-supply-chain-attacks.md?raw'
import post7 from './blog/securing-self-hosted-ai-agents.md?raw'
import post8 from './blog/token-bucket-rate-limiting-llms.md?raw'

export interface BlogPost {
  slug: string
  title: string
  excerpt: string
  date: string
  readTime: string
  tag: string
  tagColor: string
  content: string
}

export const blogPosts: BlogPost[] = [
  {
    slug: 'openclaw-api-leak-cve-2026-25253',
    title: 'OpenClaw API Key Leaks: Are You Vulnerable to CVE-2026-25253?',
    excerpt: 'We scanned thousands of self-hosted OpenClaw instances. 73% had exposed API keys in plaintext config files. A deep dive into the CVE-2026-25253 vulnerability and how to fix it.',
    date: '2026-03-15',
    readTime: '12 min',
    tag: 'Vulnerability',
    tagColor: 'bg-red-500/20 text-red-400',
    content: post1
  },
  {
    slug: 'preventing-denial-of-wallet-ai-billing',
    title: 'How to Prevent Denial of Wallet (DoW) Attacks on Your AI Infrastructure',
    excerpt: 'API cost abuse is the new DDoS. Learn how botnets are draining LLM credits and how to implement TCP-level rate limiting to survive Denial of Wallet attacks.',
    date: '2026-03-14',
    readTime: '11 min',
    tag: 'Cloud Billing',
    tagColor: 'bg-orange-500/20 text-orange-400',
    content: post2
  },
  {
    slug: 'prompt-injection-defense-llm-jailbreak',
    title: 'Stopping LLM Jailbreaks: Prompt Injection Defense Strategies for 2026',
    excerpt: 'Prompt injection remains the OWASP #1 threat for LLMs. Discover how to block DAN attacks, agent hijacking, and implement a local SLM firewall.',
    date: '2026-03-12',
    readTime: '14 min',
    tag: 'AI Security',
    tagColor: 'bg-purple-500/20 text-purple-400',
    content: post3
  },
  {
    slug: 'zero-trust-ai-gateway-architecture',
    title: 'Why Your AI Gateway Needs a Zero-Trust Architecture',
    excerpt: 'Stop relying on perimeter defense. Learn why modern self-hosted AI requires Zero-Trust principles, encrypted vaults, and machine-bound authentication.',
    date: '2026-03-10',
    readTime: '13 min',
    tag: 'Architecture',
    tagColor: 'bg-blue-500/20 text-blue-400',
    content: post4
  },
  {
    slug: 'open-source-ai-gateway-comparison-2026',
    title: 'Top Open Source AI Gateways in 2026: Kong vs LiteLLM vs Tigclaw',
    excerpt: 'Looking for a reverse proxy for your LLM deployments? We compare LiteLLM, Kong AI Gateway, and Tigclaw across security, performance, and features.',
    date: '2026-03-08',
    readTime: '15 min',
    tag: 'Comparison',
    tagColor: 'bg-indigo-500/20 text-indigo-400',
    content: post5
  },
  {
    slug: 'clawhub-skills-supply-chain-attacks',
    title: '7.1% of ClawHub Skills Leak Keys: Preventing AI Supply Chain Attacks',
    excerpt: 'An analysis of 4,000 third-party AI skills reveals a massive supply chain vulnerability in the self-hosted ecosystem. How to audit and protect your data.',
    date: '2026-03-05',
    readTime: '12 min',
    tag: 'Vulnerability',
    tagColor: 'bg-red-500/20 text-red-400',
    content: post6
  },
  {
    slug: 'securing-self-hosted-ai-agents',
    title: '5 Security Best Practices for Self-Hosted AI Agents in 2026',
    excerpt: 'Deploying OpenClaw or AutoGPT? Follow this definitive checklist to secure your server, protect your API keys, and isolate model execution.',
    date: '2026-03-01',
    readTime: '14 min',
    tag: 'Guide',
    tagColor: 'bg-teal-500/20 text-teal-400',
    content: post7
  },
  {
    slug: 'token-bucket-rate-limiting-llms',
    title: 'Why TCP-Level Token Bucket Rate Limiting is Crucial for LLM APIs',
    excerpt: 'An engineering deep-dive into how Tigclaw uses the Token Bucket algorithm to protect generative AI applications from sophisticated traffic spikes.',
    date: '2026-02-28',
    readTime: '13 min',
    tag: 'Defense',
    tagColor: 'bg-emerald-500/20 text-emerald-400',
    content: post8
  }
]