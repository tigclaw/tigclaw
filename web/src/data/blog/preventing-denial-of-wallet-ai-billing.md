## What is a Denial of Wallet (DoW) Attack?

Historically, hackers launched Distributed Denial of Service (DDoS) attacks to crash your servers and take your website offline. The goal was simple: overwhelm the target with traffic until it buckled under the load.

Today, cloud hyperscalers like AWS, Google Cloud, and Azure have fundamentally changed this equation. Their auto-scaling infrastructure automatically absorbs traffic spikes by spinning up more servers, more containers, and more compute capacity. The result? Your website stays online. Your users never notice a thing.

But you receive a **massive, unexpected bill** at the end of the month.

This is known as a **Denial of Wallet (DoW) attack** — also called an Economic Denial of Sustainability (EDoS) attack. Instead of crashing your service, the attacker drains your bank account.

---

## Why LLMs Are the Perfect DoW Target

When it comes to Generative AI and Large Language Models (LLMs), DoW attacks are uniquely devastating for several reasons:

### 1. Extreme Computational Cost

A single API call to GPT-4 can cost between $0.01 and $0.10 depending on prompt length and response tokens. That sounds trivial — until an attacker automates thousands of such calls per hour. A sustained attack generating 10,000 requests per hour at $0.05 each costs the victim **$500/hour** or **$12,000/day**.

### 2. The Pay-Per-Token Model

Unlike traditional web hosting where you pay a fixed monthly fee, LLM APIs charge strictly per token consumed. Every single character in the prompt and response has a price tag. Attackers exploit this by crafting artificially long prompts:

- "Summarize this 10,000-word document in extreme detail"
- "Generate a comprehensive 5,000-word essay on quantum mechanics"
- "Write a complete Python application with full documentation and test suite"

Each of these prompts forces the LLM to generate thousands of expensive output tokens.

### 3. Streaming Connections Hold Resources

Modern AI chatbots use Server-Sent Events (SSE) for streaming responses. Each streaming connection holds a server thread open for 10–60 seconds while the LLM generates its response token by token. An attacker can open hundreds of concurrent SSE connections, each generating a long response, effectively consuming all available server resources while maximizing your API bill.

---

## Anatomy of an API Cost Abuse Attack

A typical AI API cost abuse attack doesn't look like a massive traffic spike — it looks like "low and slow" legitimate usage. This is what makes it so dangerous and difficult to detect:

**Phase 1: Reconnaissance**
The attacker discovers an exposed, authenticated endpoint on your self-hosted AI chatbot. This could be an OpenClaw instance, a custom LangChain deployment, or any API endpoint that accepts natural language prompts and forwards them to a paid LLM provider.

**Phase 2: Distributed Low-Rate Abuse**
Instead of flooding your service from a single IP, the attacker uses a distributed botnet to send requests from hundreds of different IP addresses. Each individual IP sends only a few requests per minute — well below any reasonable rate limit threshold.

**Phase 3: Prompt Optimization**
The attacker crafts prompts specifically designed to maximize token consumption. Common techniques include:

- Requesting extremely detailed, verbose responses
- Asking the LLM to "expand on every point" or "provide comprehensive examples"
- Using system prompt overrides to force maximum output length
- Sending encoded or compressed payloads that expand into massive prompts after preprocessing

**Phase 4: Sustained Drain**
The attack runs continuously for days or weeks. Because the request rate per IP is low and the requests look like legitimate user queries, standard web application firewalls (WAFs) and basic rate limiters fail to trigger.

By the time you notice the unusual spike in your OpenAI or Anthropic billing dashboard, the damage is already done.

---

## Why Standard Defenses Fall Short

Most developers attempt to stop DoW attacks by adding basic API throttling at the application layer. While this seems logical, it's fundamentally inadequate for several reasons:

### Application-Layer Rate Limiting Is Too Late

By the time your application code (Node.js, Python, etc.) processes the request, parses the JSON body, validates headers, and decides to reject it, you have already consumed significant server resources:

- TCP handshake completed
- TLS/SSL negotiation finished
- HTTP headers fully parsed
- Request body read into memory
- Authentication middleware executed
- JSON body deserialized

All of this consumes CPU, memory, and I/O. In a severe distributed attack, this application-layer parsing causes the server to crash from memory exhaustion long before it reaches the "rate limit check" code.

### Simple Request Counting Is Insufficient

Standard rate limiters count requests per second (RPS). But in an AI context, not all requests are equal. A prompt that generates 50 tokens costs 100x less than a prompt that generates 5,000 tokens. A simple "How are you?" should not be counted the same as "Write me a complete business plan."

### IP-Based Blocking Creates False Positives

Aggressive IP blocking can inadvertently block legitimate users — especially those behind corporate NATs, VPNs, or carrier-grade NATs where many users share a single public IP address.

---

## The Tigclaw Defense: Token Bucket at TCP Level

To effectively kill a DoW attack before it burns your AI budget, you must sever the connection at the **lowest possible layer** — before the HTTP request is even fully parsed.

**Tigclaw** implements a highly optimized **Token Bucket algorithm** at the TCP layer. It acts as a strict reverse proxy sitting between the internet and your AI application:

### How the Token Bucket Works

Imagine a bucket that holds tokens:

- The bucket has a maximum capacity (e.g., 20 tokens) — this is the **burst limit**.
- Tokens are added to the bucket at a constant rate (e.g., 1 token every 3 seconds) — this is the **sustained rate**.
- When a user makes a request, the system checks if there are enough tokens in the bucket.
- If tokens are available, the request is processed and tokens are consumed.
- If the bucket is empty, the connection is **immediately dropped**.

This algorithm beautifully accommodates real human traffic. A legitimate user might send 4 quick messages in 10 seconds (draining tokens fast), but then they pause to read the AI's response, giving the bucket time to refill. A malicious bot, however, sends continuous requests without pausing. It quickly empties the bucket and is blocked from incurring further costs.

### Tigclaw's Advanced Capabilities

Tigclaw goes beyond basic token bucket counting:

- **Per-IP Connection Tracking:** Tigclaw tracks connection metadata before the HTTP payload is even fully parsed. It operates at the TCP connection level, not the HTTP request level.
- **Cost-Aware Limiting:** Tigclaw can be configured to penalize connections that hold open long SSE streaming sessions unnecessarily. A connection streaming for 60 seconds consumes more "cost tokens" than a quick 2-second response.
- **Instant Drop with `429`:** When a threshold is breached, Tigclaw drops the connection with a stark `429 Too Many Requests` response. It does not waste CPU cycles generating friendly JSON error messages or HTML error pages for bots.
- **Automatic IP Cleanup:** Expired rate limit entries are automatically pruned from memory, keeping Tigclaw's footprint minimal even under sustained attack.
- **Compiled Go Performance:** Written in pure Go, Tigclaw handles rate limiting in highly concurrent goroutines. It can process tens of thousands of concurrent connections with minimal memory overhead — far exceeding what a Python or Node.js middleware can achieve.

---

## Implementing Tigclaw for DoW Protection

Setting up Tigclaw's rate limiting takes less than 60 seconds:

```bash
# Install Tigclaw
curl -fsSL https://get.tigclaw.dev | sh

# Initialize with your AI provider
tigclaw init --endpoint https://api.openai.com

# Configure rate limits
tigclaw serve --rate-limit 20 --rate-window 60s
```

With this configuration:
- Each IP gets 20 tokens per 60-second window
- Burst capacity allows legitimate users to send quick successive messages
- Bots that exceed the threshold are instantly blocked at TCP level
- Your upstream API costs are protected by the hardened perimeter

---

## Conclusion

Denial of Wallet attacks are the new DDoS. As AI inference costs continue to dominate cloud budgets, protecting your LLM endpoints from automated cost abuse is no longer optional — it is existential.

By placing Tigclaw in front of your open-source AI deployments, you create a hardened perimeter that ensures your credit card survives the night. The Token Bucket algorithm, implemented at the TCP level in compiled Go, provides the performance and precision needed to stop sophisticated DoW attacks dead in their tracks — without degrading the experience for legitimate users.

Don't wait for the $10,000 bill. Install Tigclaw today.
