<div align="center">

# 🎬 Snag

**Grab media from anywhere. Paste a link. Get your video.**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=for-the-badge)]()

A beautiful, interactive CLI tool to download videos and audio from **any platform** — without watermarks, at the highest quality.

</div>

---

## ✨ Features

- 🌍 **Works everywhere** — YouTube, TikTok, Instagram, X (Twitter), Facebook, Reddit, and [1000+ more sites](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md)
- 🚫 **No watermarks** — downloads clean TikTok videos without the watermark
- 🎥 **Best quality** — grabs the highest resolution video + audio and merges them
- 🎵 **Audio extraction** — download just the audio as MP3
- 📺 **Quality selection** — choose between Best, 720p, 480p, or 360p
- 🎨 **Beautiful interactive UI** — styled terminal interface with step-by-step prompts
- ⚡ **Fast** — streams directly from source, no middleman servers
- 🔧 **Flag mode** — skip the prompts with CLI flags for scripting

---

## 📸 Preview

```
╭──────────────────────────────────────────────────────────╮
│                                                          │
│  🎬 Snag — grab media from anywhere                     │
│                                                          │
│  Paste your URL here                                     │
│                                                          │
│  (press Enter to continue)                               │
│                                                          │
╰──────────────────────────────────────────────────────────╯
```

```
╭──────────────────────────────────────────────────────────╮
│                                                          │
│  ✅ Download complete!                                  │
│                                                          │
│    Saved to: ./downloads/My Video Title.mp4              │
│                                                          │
│    Press q to quit                                       │
│                                                          │
╰──────────────────────────────────────────────────────────╯
```

---

## 🚀 Installation

### Prerequisites

You need `yt-dlp` and `ffmpeg` installed on your system:

```bash
# macOS (Homebrew)
brew install yt-dlp ffmpeg

# Linux (apt)
sudo apt install yt-dlp ffmpeg

# Windows (Chocolatey)
choco install yt-dlp ffmpeg
```

### Install Snag

```bash
# Using Go
go install github.com/Verifieddanny/snag/cmd/snag@latest
```

### Build from source

```bash
git clone https://github.com/Verifieddanny/snag.git
cd snag/snag-cli
go build -o bin/snag ./cmd/snag/
```

---

## 📖 Usage

### Interactive Mode (Recommended)

Just run `snag` with no arguments and follow the prompts:

```bash
snag
```

The interactive UI will guide you through:

1. **Paste your URL** — from any supported platform
2. **Choose media type** — Video or Audio Only
3. **Select quality** — Best, 720p, 480p, or 360p
4. **Confirm & download** — watch the progress in real time

### Flag Mode (For scripting)

```bash
# Download video at best quality
snag -url="https://www.youtube.com/watch?v=dQw4w9WgXcQ"

# Download audio only as MP3
snag -url="https://youtu.be/example" -audio-only

# Download at 720p to a custom directory
snag -url="https://youtu.be/example" -quality=720 -o=~/Videos

# Download TikTok without watermark
snag -url="https://vt.tiktok.com/ZSxV8G3Lp/"
```

### Available Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-url` | — | URL to download |
| `-o` | `./downloads` | Output directory |
| `-audio-only` | `false` | Extract audio as MP3 |
| `-quality` | `best` | Video quality: `best`, `720`, `480`, `360` |

---

## 🌍 Supported Platforms

| Platform | Video | Audio | No Watermark |
|----------|-------|-------|--------------|
| YouTube | ✅ | ✅ | ✅ |
| TikTok | ✅ | ✅ | ✅ |
| Instagram | ✅ | ✅ | ✅ |
| X (Twitter) | ✅ | ✅ | ✅ |
| Facebook | ✅ | ✅ | ✅ |
| Reddit | ✅ | ✅ | ✅ |
| [1000+ more](https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md) | ✅ | ✅ | varies |

---

## 🏗️ Project Structure

```
snag-cli/
├── cmd/
│   └── snag/
│       ├── main.go            # Entry point & flag mode
│       └── interactive.go     # Bubbletea interactive UI
├── internal/
│   └── downloader/
│       └── downloader.go      # yt-dlp wrapper & download logic
├── go.mod
└── Makefile
```

---

## 🗺️ Roadmap

- [x] Interactive CLI with Bubbletea
- [x] Multi-platform support (YouTube, TikTok, Instagram, X, etc.)
- [x] Quality selection
- [x] Audio extraction
- [x] No-watermark TikTok downloads
- [ ] Web app + API
- [ ] Desktop app (macOS)
- [ ] Homebrew formula
- [ ] Playlist support
- [ ] Batch downloads from file

---

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit PRs.

```bash
# Clone the repo
git clone https://github.com/Verifieddanny/snag.git
cd snag/snag-cli

# Install dependencies
go mod tidy

# Run in development
go run ./cmd/snag/
```

---

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

**Built by [Danny](https://github.com/Verifieddanny)** 🚀

If Snag saved you time, give it a ⭐

</div>