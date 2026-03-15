import { Helmet } from 'react-helmet-async'

interface SEOProps {
  title?: string
  description?: string
  path?: string
  image?: string
  type?: string
}

const SITE = {
  name: 'Tigclaw',
  url: 'https://tigclaw.com',
  defaultTitle: 'Tigclaw — Zero-Trust AI Security Gateway',
  defaultDescription: 'Open-source security gateway for OpenClaw. Protect your API keys, rate-limit abuse, and block prompt injection attacks — 100% locally.',
  defaultImage: '/og-image.png',
}

export default function SEO({ title, description, path = '/', image, type = 'website' }: SEOProps) {
  const pageTitle = title ? `${title} | ${SITE.name}` : SITE.defaultTitle
  const pageDescription = description || SITE.defaultDescription
  const pageUrl = `${SITE.url}${path}`
  const pageImage = `${SITE.url}${image || SITE.defaultImage}`

  return (
    <Helmet>
      <title>{pageTitle}</title>
      <meta name="description" content={pageDescription} />
      <link rel="canonical" href={pageUrl} />

      {/* Open Graph */}
      <meta property="og:type" content={type} />
      <meta property="og:title" content={pageTitle} />
      <meta property="og:description" content={pageDescription} />
      <meta property="og:url" content={pageUrl} />
      <meta property="og:image" content={pageImage} />
      <meta property="og:site_name" content={SITE.name} />

      {/* Twitter */}
      <meta name="twitter:card" content="summary_large_image" />
      <meta name="twitter:title" content={pageTitle} />
      <meta name="twitter:description" content={pageDescription} />
      <meta name="twitter:image" content={pageImage} />

      {/* Schema.org */}
      <script type="application/ld+json">{JSON.stringify({
        "@context": "https://schema.org",
        "@type": "SoftwareApplication",
        "name": "Tigclaw",
        "applicationCategory": "SecurityApplication",
        "operatingSystem": "Linux, macOS, Windows",
        "offers": { "@type": "Offer", "price": "0", "priceCurrency": "USD" },
        "description": pageDescription,
        "url": SITE.url,
      })}</script>
    </Helmet>
  )
}
