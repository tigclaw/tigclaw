# 🐯 Tigclaw — Zero-Trust AI Security Gateway

> **Your self-hosted AI is running naked on the internet. Fix it in 30 seconds.**

Tigclaw is an open-source security gateway for [OpenClaw](https://github.com/openclaw/openclaw) and other self-hosted AI platforms. It sits between the internet and your AI instance, protecting your API keys, rate-limiting abuse, and blocking prompt injection attacks — **100% locally, without sending a single byte to the cloud.**

---

## 🔥 The Problem

| Incident | Severity |
|----------|----------|
| **CVE-2026-25253**: WebSocket RCE — any website can hijack your instance | 🔴 Critical |
| **1.5M API Tokens Leaked** from plaintext config files | 🔴 Critical |
| **41.7%** of popular plugins contain vulnerabilities, **7.1%** steal keys | 🔴 Critical |
| API keys stored in **plaintext** in `config.json` | 🟠 High |
| Prompt injection can write persistent backdoors to `SOUL.md` | 🟠 High |

**If your OpenClaw port is exposed to the internet, your credit card is at risk. Right now.**

---

## 🛡️ The Solution

```
🌐 Internet Traffic
       │
       ▼
┌──────────────────────────────────────┐
│  🐯 Tigclaw Gateway  (:9000)        │
│                                      │
│  1. 🚦 Rate Limit (Anti-DoW)        │
│  2. 🔐 Key Swap (fake → real)       │
│  3. 🔏 DLP Masking                  │
│  4. 🧠 Prompt Firewall              │
└──────────────┬───────────────────────┘
               │  127.0.0.1 only
               ▼
┌──────────────────────────────────────┐
│  🐙 OpenClaw  (hidden from internet)│
│     config: "sk-tigclaw-a1b2c3d4"   │
│     (useless if stolen)             │
└──────────────────────────────────────┘
```

**Tigclaw replaces your real API keys with disposable fake keys.** Even if an attacker completely compromises OpenClaw, they get nothing.

---

## ⚡ Quick Start

### Install
```bash
# Linux / macOS
curl -sSL https://tigclaw.com/install.sh | bash

# Windows (PowerShell)
irm https://tigclaw.com/install.ps1 | iex
```

### Initialize (auto-migrate existing keys)
```bash
tigclaw init
# 🔍 Found OpenClaw config: ~/.openclaw/config.json
# 🔑 Found plaintext key: sk-proj-abc...xyz
# ✅ Migrated → sk-tigclaw-a1b2c3d4 (encrypted in vault)
```

### Start the Gateway
```bash
tigclaw serve
# 🐯 Tigclaw Security Gateway
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
#   Listening on    : :9000
#   Upstream (OC)   : http://127.0.0.1:3001
#   Strict Mode     : true
#   Protected Keys  : 2
```

### Manage Keys
```bash
tigclaw keys add sk-your-real-key-here       # Add a new key
tigclaw keys list                             # List all fake keys
tigclaw keys update sk-tigclaw-xxx new-key    # Rotate real key (seamless)
tigclaw keys remove sk-tigclaw-xxx            # Delete a key
tigclaw status                                # Security score
```

---

## 🔐 How Zero-Trust Key Substitution Works

1. You register your real API key with `tigclaw keys add`
2. Tigclaw encrypts it with **AES-256-GCM** (machine-bound key derivation)
3. A fake key `sk-tigclaw-xxxx` is generated and written to OpenClaw's config
4. When OpenClaw makes a request, Tigclaw intercepts it **in memory**
5. The fake key is swapped for the real key → forwarded to OpenAI/Anthropic
6. Response streams back untouched → real key vanishes from memory

**Result:** Your real API key exists only in encrypted form on disk and for milliseconds in memory. It never appears in any config file, log, or network trace.

---

## 🏗️ Tech Stack

| Component | Technology | Why |
|-----------|-----------|-----|
| Gateway (Data Plane) | **Go** | `httputil.ReverseProxy` + SSE streaming, <5ms latency |
| Encryption | **AES-256-GCM** | Standard library, hardware-accelerated |
| Database | **SQLite** (pure Go) | Zero external dependencies |
| CLI | **Cobra** | Industry-standard Go CLI framework |

---

## 🧪 Run Vulnerability Scan (No Install Required)

Don't trust us? Run the read-only scanner first:

```bash
curl -sSL https://raw.githubusercontent.com/tigclaw/tigclaw/main/tigclaw-scan.sh | bash
```

This scans your OpenClaw instance for:
- 🌐 Public port exposure
- 🔑 Plaintext API keys in config
- 👤 Root privilege risks
- 🧠 Dangerous SOUL.md instructions

---

## 📜 License

MIT — Use it, fork it, sell it. Just keep your AI safe.

---

<p align="center">
  <b>🐯 Stop running naked. Start running Tigclaw.</b>
  <br><br>
  <a href="https://tigclaw.com">Website</a> · 
  <a href="https://github.com/tigclaw/tigclaw/issues">Issues</a> · 
  <a href="https://github.com/tigclaw/tigclaw/discussions">Discussions</a>
</p>
