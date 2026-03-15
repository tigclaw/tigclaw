#!/bin/bash

# ==============================================================================
# 🐯 Tigclaw Security Scanner for OpenClaw
# ==============================================================================
# version: 0.1.0
# description: A lightweight, read-only script to scan your local OpenClaw 
#              installation for catastrophic security misconfigurations.
# ==============================================================================

# --- Color Definitions ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# --- Score System ---
SCORE=100

echo -e "${CYAN}${BOLD}"
echo "  _______ _           _               "
echo " |__   __(_)         | |              "
echo "    | |   _  __ _ ___| | __ ___      "
echo "    | |  | |/ _\` / __| |/ _\` \ \ /\ / /"
echo "    | |  | | (_| \__ \ | (_| |\ V  V / "
echo "    |_|  |_|\__, |___/_|\__,_| \_/\_/  "
echo "             __/ |                    "
echo "            |___/    Security Scanner "
echo -e "${NC}"
echo -e "Initializing read-only security check for OpenClaw environment...\n"

sleep 1

# ---------------------------------------------------------
# 1. ROOT PRIVILEGE CHECK
# ---------------------------------------------------------
echo -n "Checking execution privileges... "
if [ "$EUID" -eq 0 ]; then
  echo -e "[${RED}FAIL${NC}]"
  echo -e "  ${RED}>> DANGER:${NC} You are running this scan as root. If OpenClaw also runs as root, a Prompt Injection (DAN) could give a global attacker full control over your entire server.\n"
  ((SCORE-=20))
else
  echo -e "[${GREEN}PASS${NC}]"
  echo -e "  Non-root user detected. Good.\n"
fi

sleep 0.5

# ---------------------------------------------------------
# 2. DEFAULT PORT EXPOSURE (0.0.0.0) CHECK
# ---------------------------------------------------------
echo -n "Checking network exposure on port 3001 (OpenClaw Default)... "
if command -v ss >/dev/null 2>&1; then
  PORT_CHECK=$(ss -tlnp 2>/dev/null | grep ":3001" | grep -E "0\.0\.0\.0|\*|\:\:")
elif command -v netstat >/dev/null 2>&1; then
  PORT_CHECK=$(netstat -tlnp 2>/dev/null | grep ":3001" | grep -E "0\.0\.0\.0|\*|\:\:")
else
  PORT_CHECK="Command not found, skipping."
fi

if [[ -n "$PORT_CHECK" && "$PORT_CHECK" != "Command not found, skipping." ]]; then
  echo -e "[${RED}FAIL${NC}]"
  echo -e "  ${RED}>> DANGER:${NC} Port 3001 is bound to 0.0.0.0/::. Your OpenClaw instance is naked on the public internet."
  echo -e "  ${RED}>> RISK:${NC} Scanners on Shodan can easily discover your instance and abuse your LLM quotas (Anti-DoW risk).\n"
  ((SCORE-=30))
elif [ "$PORT_CHECK" == "Command not found, skipping." ]; then
  echo -e "[${YELLOW}WARN${NC}] (ss/netstat not found)\n"
else
  echo -e "[${GREEN}PASS${NC}]"
  echo -e "  No public exposure detected on port 3001.\n"
fi

sleep 0.5

# ---------------------------------------------------------
# 3. PLAINTEXT API KEY CHECK (Leaked Bottoms)
# ---------------------------------------------------------
echo -n "Scanning for plaintext API Keys in common OpenClaw config files... "
KEY_FOUND=0
# Scanning common directories
CONFIG_PATHS=(
  "$HOME/.openclaw/config.json"
  "$HOME/.openclaw/openclaw.json"
  "/etc/openclaw/config.json"
  "./config.json"
  "./openclaw.json"
)

FOUND_PATH=""
for path in "${CONFIG_PATHS[@]}"; do
  if [ -f "$path" ]; then
    if grep -q "sk-" "$path" 2>/dev/null; then
      KEY_FOUND=1
      FOUND_PATH="$path"
      break
    fi
  fi
done

if [ $KEY_FOUND -eq 1 ]; then
  echo -e "[${RED}FAIL${NC}]"
  echo -e "  ${RED}>> LETHAL DANGER:${NC} Found plaintext OpenAI/Anthropic keys ('sk-...') in: $FOUND_PATH"
  echo -e "  ${RED}>> RISK:${NC} Any zero-day Local File Inclusion (LFI) in OpenClaw or its plugins will immediately steal this key. Your credit card limit is at severe risk.\n"
  ((SCORE-=40))
else
  echo -e "[${GREEN}PASS${NC}]"
  echo -e "  No plaintext 'sk-' API keys found in standard config paths.\n"
fi

sleep 0.5

# ---------------------------------------------------------
# 4. SOUL.MD PROMPT INJECTION CHECK
# ---------------------------------------------------------
echo -n "Scanning SOUL.md for dangerous Shell execution prompt instructions... "
SOUL_FOUND=0
SOUL_PATHS=(
  "$HOME/.openclaw/SOUL.md"
  "./SOUL.md"
)

FOUND_SOUL_PATH=""
for path in "${SOUL_PATHS[@]}"; do
  if [ -f "$path" ]; then
    # Look for instructions that might allow shell command execution
    if grep -Eiq "shell|bash|execute|command line|终端|执行命令" "$path" 2>/dev/null; then
      SOUL_FOUND=1
      FOUND_SOUL_PATH="$path"
      break
    fi
  fi
done

if [ $SOUL_FOUND -eq 1 ]; then
  echo -e "[${YELLOW}WARN${NC}]"
  echo -e "  ${YELLOW}>> WARNING:${NC} Found words like 'shell/bash/execute' in your $FOUND_SOUL_PATH"
  echo -e "  ${YELLOW}>> RISK:${NC} Giving an LLM raw shell execution access without a robust semantic firewall (SLM) is highly susceptible to DAN Prompt Injections.\n"
  ((SCORE-=10))
else
  echo -e "[${GREEN}PASS${NC}]"
  echo -e "  No immediate dangerous macro instructions found in SOUL.md.\n"
fi

sleep 0.5

# ==============================================================================
# REPORTING
# ==============================================================================
echo -e "---------------------------------------------------------"
echo -ne "Final tigclaw-scan.sh Security Score: "

if [ $SCORE -eq 100 ]; then
  echo -e "${GREEN}${BOLD}$SCORE / 100${NC}"
  echo -e "${GREEN}Excellent! Your instance is well secured.${NC}\n"
elif [ $SCORE -ge 70 ]; then
  echo -e "${YELLOW}${BOLD}$SCORE / 100${NC}"
  echo -e "${YELLOW}Warning: Minor security issues represent surface-level risks.${NC}\n"
else
  echo -e "${RED}${BOLD}$SCORE / 100${NC}"
  echo -e "${RED}catastrophic failure: YOUR INSTANCE IS EXTREMELY VULNERABLE.${NC}"
  echo -e "${RED}Your API Keys and your server are at high risk of compromise!${NC}\n"
fi

# Marketing hook for Tigclaw
echo -e "${CYAN}================================================================================${NC}"
if [ $SCORE -lt 100 ]; then
  echo -e "${BOLD}Don't want your credit card stolen? Stop running naked on the internet.${NC}"
  echo -e "Protect your OpenClaw instance with ${CYAN}Tigclaw${NC} — The Zero-Trust AI Security Gateway."
  echo ""
  echo -e "🔒 ${GREEN}Zero-Trust Key Substitution${NC} (Never store real API keys in config again)"
  echo -e "🚦 ${GREEN}Anti-DoW Rate Limiting${NC} (Stop botnet billing attacks in their tracks)"
  echo -e "🧠 ${GREEN}Local SLM Firewall${NC} (Block Prompt Injections seamlessly)"
  echo ""
  echo -e "Get early access and star us at: ${BOLD}https://github.com/tigclaw/tigclaw${NC}"
  echo -e "Official site: ${BOLD}https://tigclaw.com${NC}"
fi
echo -e "${CYAN}================================================================================${NC}"
