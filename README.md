# CS2 Server Manager

**Automated setup for multiple CS2 dedicated servers with competitive plugins and tournament support.**

Deploy 3-5 servers in minutes. Perfect for LAN parties, tournaments, or community servers.

🏆 **Tournament Ready:** Integrates with [MatchZy Auto Tournament](https://mat.sivert.io/) for complete automation.

---

## 🚀 Quick Start

### Prerequisites
```bash
sudo apt-get update
sudo apt-get install -y lib32gcc-s1 lib32stdc++6 steamcmd tmux curl jq unzip tar rsync
```

### Installation

```bash
# Clone and enter directory
git clone https://github.com/sivert-io/cs2-server-manager.git
cd cs2-server-manager

# Interactive mode (recommended for first-time setup)
./manage.sh
# Choose Option 1, press Enter for defaults (3 servers)
# Wait 15-30 minutes for CS2 download (~60GB)

# Non-interactive mode (automated with defaults)
./manage.sh install
# Auto-confirms with defaults: 3 servers, downloads plugins, updates master
```

### Start Servers
```bash
./manage.sh start          # Non-interactive
# OR
./manage.sh                # Interactive menu → Option 3
```

---

## 🔌 Included Plugins

Automatically installed and configured:

- **[Metamod:Source](https://www.metamodsource.net/)** - Plugin framework
- **[CounterStrikeSharp](https://github.com/roflmuffin/CounterStrikeSharp)** - Modern C# plugin framework
- **[MatchZy Enhanced Fork](https://github.com/sivert-io/MatchZy)** - Match management with tournament automation
  - Additional events for [MatchZy Auto Tournament](https://mat.sivert.io/) integration
  - Real-time tracking, auto-progression, advanced tournament features
  - Falls back to official release if unavailable
- **[CS2-AutoUpdater](https://github.com/dran1x/CS2-AutoUpdater)** - Automatic server updates

---

## 🎮 Usage

### Interactive Menu
```bash
./manage.sh
```

Main operations:
- **1** - Install/redeploy servers
- **2** - Server status
- **3/4/5** - Start/stop/restart all servers
- **12** - Update CS2 game files
- **13** - Update plugins
- **15** - Repair servers (no re-download)
- **16** - Debug mode (foreground)

### Non-Interactive Commands
```bash
./manage.sh install           # Install with defaults
./manage.sh start             # Start all servers
./manage.sh stop              # Stop all servers
./manage.sh status            # Show status
./manage.sh update-game       # Update CS2
./manage.sh update-plugins    # Update plugins
./manage.sh repair            # Repair servers
./manage.sh help              # Show all commands
```

### Direct Server Management
```bash
sudo ./scripts/cs2_tmux.sh start          # Start all
sudo ./scripts/cs2_tmux.sh start 1        # Start server 1
sudo ./scripts/cs2_tmux.sh stop           # Stop all
sudo ./scripts/cs2_tmux.sh status         # Check status
sudo ./scripts/cs2_tmux.sh attach 1       # Console (Ctrl+B,D to detach)
sudo ./scripts/cs2_tmux.sh debug 1        # Debug mode (Ctrl+C to stop)
sudo ./scripts/cs2_tmux.sh logs 1 100     # Show last 100 log lines
```

---

## ⚙️ Configuration

### Server Ports (auto-increment by 10)
| Server | Game Port | GOTV Port |
|--------|-----------|-----------|
| 1      | 27015     | 27020     |
| 2      | 27025     | 27030     |
| 3      | 27035     | 27040     |

**RCON Password:** `ntlan2025` (change in interactive install or configs)

### Custom Configs
Place your custom configs in `overrides/game/csgo/`:
```
overrides/game/csgo/
├── cfg/MatchZy/         # MatchZy configs (config.cfg, admins.json, etc.)
└── addons/              # Plugin configs
```

Files in `overrides/` are **preserved during all updates**.

### Environment Variables
```bash
NUM_SERVERS=5 RCON_PASSWORD=mypass ./manage.sh install
```

---

## 🐛 Troubleshooting

**Server won't start:**
```bash
sudo ./scripts/cs2_tmux.sh debug 1    # See all errors in real-time
```

**Plugins not loading / steamclient.so errors:**
```bash
./manage.sh repair                     # Validates and fixes without re-download
```

**View logs:**
```bash
sudo ./scripts/cs2_tmux.sh logs 1 100
```

---

## 📂 File Structure

```
/home/cs2/
├── master-install/    # Master CS2 (60GB)
├── server-1/          # Server instances
├── server-2/
└── server-3/

cs2-server-manager/
├── scripts/           # Management scripts
├── game_files/        # Downloaded plugins
├── overrides/         # Your custom configs (preserved)
└── manage.sh          # Main interface
```

---

## 🔗 Links

- **[MatchZy Enhanced Fork](https://github.com/sivert-io/MatchZy)** - Tournament automation version
- **[MatchZy Auto Tournament](https://mat.sivert.io/)** - Complete tournament platform ([GitHub](https://github.com/sivert-io/matchzy-auto-tournament))
- [MatchZy Original](https://github.com/shobhit-pathak/MatchZy) - Official version
- [CounterStrikeSharp Docs](https://docs.cssharp.dev/)
- [CS2 Dedicated Server Guide](https://developer.valvesoftware.com/wiki/Counter-Strike_2/Dedicated_Servers)

---

**Made for the CS2 Community** 🎮
