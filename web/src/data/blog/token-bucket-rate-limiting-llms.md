## The Unique Problem of LLM Rate Limiting

Standard web applications measure rate limits in "Requests Per Second" (RPS). If a user loads an image, that's one request. If they load a web page with 20 assets, that's 20 requests. The computational cost per request is negligible — a few milliseconds of CPU time and a few kilobytes of bandwidth.

**Generative AI turns this paradigm completely upside down.**

A single request to an LLM — asking it to summarize a massive PDF, write a complex multi-file application, or analyze a dataset — can tie up a GPU on the backend for **10 to 60 seconds** and cost anywhere from **$0.01 to $1.00** depending on the model and token count. A GPT-4 request processing 128K tokens of context can cost over $5 for a single invocation.

Therefore, a simplistic RPS rate limit is entirely inadequate for protecting an AI deployment. You need a mechanism that accounts for:

- **Burst capacity:** Real users send rapid-fire messages in short bursts, then pause to read responses
- **Variable cost:** Not all requests are equal; a 50-token prompt costs 100x less than a 50,000-token prompt
- **Long-running connections:** Server-Sent Events (SSE) streaming ties up server resources for extended durations
- **Distributed attacks:** Sophisticated attackers spread requests across hundreds of IP addresses

The **Token Bucket algorithm** solves all of these problems — and when implemented at the TCP level rather than the application level, it provides an impenetrable shield for your AI infrastructure.

---

## The Token Bucket Algorithm Explained

The Token Bucket is one of the most elegant algorithms in computer science for traffic shaping. Its brilliance lies in its simplicity and its natural accommodation of bursty traffic patterns.

### The Mental Model

Imagine a physical bucket that holds tokens (think of them as permission slips to make API requests):

1. **Capacity:** The bucket has a maximum capacity — say, 20 tokens. This is the **burst limit**. A legitimate user can send up to 20 quick requests before any throttling kicks in.

2. **Refill Rate:** New tokens are added to the bucket at a constant rate — say, 1 token every 3 seconds. This is the **sustained rate**. Over time, a user can make ~20 requests per minute.

3. **Consumption:** When a request arrives, the system checks the bucket:
   - **Tokens available?** The request is processed. One token is consumed.
   - **Bucket empty?** The request is immediately **rejected** with a `429 Too Many Requests` response.

4. **Accumulation:** If the user is quiet for a while, tokens accumulate up to the maximum capacity. This means a returning user can immediately send a burst of messages — just like a real human would after reading a long AI response.

### Why Token Bucket Beats Fixed Window

The simpler approach — counting requests per fixed time window (e.g., "max 60 requests per minute") — has a critical flaw: the **boundary problem**.

A user can send 60 requests at 11:59:59 and another 60 requests at 12:00:01 — a burst of 120 requests in 2 seconds — because the counter resets at the minute boundary. This defeats the entire purpose of rate limiting.

The Token Bucket has no boundaries. The refill is continuous. The consumption is continuous. There is no moment where the counter "resets" and allows a spike.

### Why Token Bucket Beats Sliding Window

Sliding window algorithms track request timestamps in a sorted list and compute the request count over the trailing N seconds. While more accurate than fixed windows, they have significant memory and computational overhead:

- Each request requires inserting a timestamp and pruning expired entries
- Memory usage grows linearly with request volume
- Under high concurrency, the sorted list becomes a contention point

The Token Bucket requires only **two numbers per IP address**: the current token count and the timestamp of the last refill. Total memory per tracked IP: **16 bytes**. Even with 100,000 concurrent IPs, the total memory footprint is under 2MB.

---

## Why TCP-Level Implementation Matters

Here is where most developers make a critical architectural mistake: they implement rate limiting as **application-layer middleware**.

### The Typical (Wrong) Approach

```javascript
// Express.js middleware - DON'T DO THIS for AI protection
const rateLimit = require('express-rate-limit');
app.use('/api/chat', rateLimit({
  windowMs: 60 * 1000,
  max: 20,
  message: { error: "Too many requests" }
}));
```

This looks reasonable. But by the time this middleware executes, an enormous amount of work has already happened:

### The Hidden Cost of Application-Layer Processing

For every single incoming request — including malicious ones — your server has already:

1. **TCP Handshake** — 3 packets exchanged (SYN, SYN-ACK, ACK). Server allocates a socket buffer.
2. **TLS/SSL Negotiation** — CPU-intensive cryptographic operations: key exchange, cipher suite negotiation, certificate validation. This alone can consume 1-5ms of CPU time per connection.
3. **HTTP Parsing** — The full HTTP request (headers + body) is read into memory. For AI requests with long prompts, this could be kilobytes or megabytes of JSON data.
4. **Middleware Chain** — Express.js parses cookies, sessions, CORS headers, authentication tokens, and deserializes the JSON body before the rate limiter middleware even sees the request.
5. **Rate Limit Check** — Finally, the rate limiter checks the counter and rejects the request.
6. **Response Generation** — The server generates a JSON response body (`{ "error": "Too many requests" }`), serializes it, adds HTTP headers, and sends it back through the TLS layer.

All of this work is **wasted** for a request that was going to be rejected anyway. Under a distributed Denial of Wallet attack sending 10,000 requests per second from 1,000 different IPs, this wasted processing can crash your Node.js or Python server from **memory exhaustion** long before the AI API is ever called.

The rate limiter successfully blocks the AI API calls, but the application server itself goes down — converting a DoW attack into a plain old DoS attack.

### The TCP-Level Alternative

TCP-level rate limiting operates **before** any of the above processing occurs:

1. A new TCP connection arrives
2. The kernel notifies the rate limiting process
3. The process checks the source IP against the token bucket (16 bytes of state, one memory lookup)
4. **If blocked:** The TCP connection is immediately reset (`RST`) or the socket is closed. Zero TLS negotiation. Zero HTTP parsing. Zero JSON deserialization. Total CPU cost: **microseconds**.
5. **If allowed:** The connection proceeds to the normal application stack

The difference in resource consumption between application-layer and TCP-layer rate limiting under attack is not 2x or 5x — it is **100x to 1000x**.

---

## Tigclaw's High-Performance Implementation

**Tigclaw** pushes rate limiting to the absolute edge of your local infrastructure. Written in pure Go, it leverages Go's goroutine-based concurrency model to handle rate limiting with minimal overhead.

### Architecture

```
[Internet] → [Tigclaw TCP Listener] → [Token Bucket Engine] → [TLS Termination] → [HTTP Proxy] → [Your AI App]
                                              ↓
                                        [429 / RST if blocked]
```

The Token Bucket engine sits at the **first stage** of the connection pipeline. Before TLS negotiation, before HTTP parsing, before any application logic — the source IP is checked against the bucket.

### Key Technical Details

**Per-IP State:**
- Each tracked IP uses exactly 16 bytes: `{ tokens: float64, lastRefill: int64 }`
- 100,000 concurrent IPs = 1.6MB of memory
- State is stored in a lock-free concurrent hash map for maximum throughput

**Goroutine-Per-Connection:**
- Each accepted connection spawns a lightweight goroutine (~4KB stack)
- Go's runtime multiplexes goroutines across OS threads efficiently
- 10,000 concurrent connections = ~40MB of goroutine stacks

**Automatic Cleanup:**
- A background goroutine periodically scans the IP table and removes entries that haven't been seen in the last 10 minutes
- This prevents memory leaks during sustained attacks from botnets with rotating IPs

**Configurable Behavior:**
- `--rate-limit 20` — Burst capacity (tokens in bucket)
- `--rate-window 60s` — Time to refill the bucket completely
- `--rate-cost-sse 3` — SSE streaming connections consume 3x tokens (because they hold resources longer)

### Performance Benchmarks

Under synthetic load testing with 50,000 concurrent connections from 10,000 unique IPs:

| Metric | Application-Layer (Express.js) | TCP-Layer (Tigclaw) |
|---|---|---|
| Memory under attack | 2.1 GB (before crash) | 48 MB |
| Requests rejected/sec | ~3,000 (then OOM) | ~200,000 |
| Reject latency (p99) | 45ms | 0.3ms |
| Server stability | ❌ Crashed at 8,000 conn | ✅ Stable at 50,000 conn |

The numbers speak for themselves. Application-layer rate limiting is a speed bump. TCP-layer rate limiting is a concrete wall.

---

## Implementation Guide

### Basic Setup

```bash
# Install Tigclaw
curl -fsSL https://get.tigclaw.dev | sh

# Add your API key with fake key substitution
tigclaw keys add --provider openai --real-key "sk-proj-..."

# Start the gateway with rate limiting
tigclaw serve --rate-limit 20 --rate-window 60s --port 9090
```

### Point Your Application at Tigclaw

```bash
# In your AI application's .env:
OPENAI_API_KEY=sk-tigclaw-f7a2b91c
OPENAI_BASE_URL=http://localhost:9090/v1
```

### Monitor Rate Limiting in Real Time

```bash
tigclaw status
# Rate Limit Stats:
# Active IPs tracked: 147
# Requests allowed (1h): 2,341
# Requests blocked (1h): 12
# Memory usage: 2.4 MB
```

---

## Conclusion

Rate limiting for AI applications is fundamentally different from rate limiting for traditional web applications. The extreme cost per request, the long-running streaming connections, and the sophisticated distributed attack patterns demand a solution that operates at the lowest possible network layer.

**Tigclaw** pushes the Token Bucket algorithm to the TCP connection level, implemented in compiled Go for maximum performance. It protects not just your wallet from upstream OpenAI costs, but also protects your server's RAM and CPU from crashing under load.

When an abuser crosses the threshold, Tigclaw severs the connection in microseconds. No wasted TLS handshakes. No wasted HTTP parsing. No wasted JSON serialization. Just an instant `RST` packet and a log entry.

It is the ultimate defensive perimeter for self-hosted AI. Install it in 60 seconds and sleep soundly.
