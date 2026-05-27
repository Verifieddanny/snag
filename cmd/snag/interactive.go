package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/Verifieddanny/snag/internal/downloader"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type step int

type downloadMsg struct {
	filename string
	err      error
}

type downloadProgressMsg struct {
	line string
}

const (
	stepURL step = iota
	stepMediaType
	stepQuality
	stepConfirm
	stepDownloading
	stepDone
)

type model struct {
	step        step
	url         string
	platform    string
	mediaType   string
	quality     string
	outputDir   string
	filename    string
	err         error
	textInput   textinput.Model
	spinner     spinner.Model
	choices     []string
	cursor      int
	downloadLog []string
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Paste your URL here"
	ti.Focus()
	ti.Width = 60

	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		step:      stepURL,
		textInput: ti,
		spinner:   s,
		outputDir: "./downloads",
		quality:   "best",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) View() string {
	switch m.step {
	case stepURL:
		title := titleStyle.Render("🎬 Snag — grab media from anywhere")
		hint := dimStyle.Render("(press Enter to continue)")
		content := fmt.Sprintf("\n%s\n\n%s\n\n%s", title, m.textInput.View(), hint)
		return boxStyle.Render(content) + "\n"

	case stepMediaType:
		return m.renderChoices("What do you want to download?", []string{"🎥 Video", "🎵 Audio Only"})

	case stepQuality:
		return m.renderChoices("Select quality:", []string{"🔥 Best", "📺 720p", "📱 480p", "📟 360p"})

	case stepConfirm:
		title := titleStyle.Render("📋 Summary")
		summary := fmt.Sprintf(
			"  %s  %s\n  %s  %s\n  %s  %s\n  %s  %s\n  %s  %s",
			dimStyle.Render("URL:"), urlStyle.Render(m.url),
			dimStyle.Render("Platform:"), highlightStyle.Render(m.platform),
			dimStyle.Render("Type:"), m.mediaType,
			dimStyle.Render("Quality:"), m.quality,
			dimStyle.Render("Output:"), m.outputDir,
		)
		hint := dimStyle.Render("Press Enter to start • q to cancel")
		content := fmt.Sprintf("\n%s\n\n%s\n\n%s", title, summary, hint)
		return boxStyle.Render(content) + "\n"

	case stepDownloading:
		title := titleStyle.Render("📥 Downloading...")
		var logs strings.Builder
		for _, l := range m.downloadLog {
			logs.WriteString("  ")
			logs.WriteString(dimStyle.Render(l))
			logs.WriteString("\n")
		}
		content := fmt.Sprintf("\n%s %s\n\n%s", m.spinner.View(), title, logs.String())
		return boxStyle.Render(content) + "\n"

	case stepDone:
		if m.err != nil {
			title := errorStyle.Render("❌ Download failed")
			content := fmt.Sprintf("\n%s\n\n  %v\n\n%s", title, m.err, dimStyle.Render("Press q to quit"))
			return boxStyle.Render(content) + "\n"
		}
		title := successStyle.Render("✅ Download complete!")
		content := fmt.Sprintf("\n%s\n\n  %s %s\n\n%s", title, dimStyle.Render("Saved to:"), urlStyle.Render(m.filename), dimStyle.Render("Press q to quit"))
		return boxStyle.Render(content) + "\n"
	}
	return ""
}

func (m model) renderChoices(title string, choices []string) string {
	var s strings.Builder
	t := titleStyle.Render(title)
	fmt.Fprintf(&s, "\n%s\n\n", t)
	for i, choice := range choices {
		if m.cursor == i {
			fmt.Fprintf(&s, "  %s %s\n", highlightStyle.Render("▸"), highlightStyle.Render(choice))
		} else {
			fmt.Fprintf(&s, "    %s\n", dimStyle.Render(choice))
		}
	}
	fmt.Fprintf(&s, "\n%s\n", dimStyle.Render("(↑/↓ to move • Enter to select)"))
	content := s.String()
	return boxStyle.Render(content) + "\n"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step != stepDownloading {
				return m, tea.Quit
			}
		}

		switch m.step {
		case stepURL:
			switch msg.String() {
			case "enter":
				m.url = m.textInput.Value()
				if m.url == "" {
					return m, nil
				}
				m.platform = detectPlatform(m.url)
				m.step = stepMediaType
				m.choices = []string{"video", "audio"}
				m.cursor = 0
				return m, nil
			}
			m.textInput, _ = m.textInput.Update(msg)
			return m, nil

		case stepMediaType:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				m.mediaType = m.choices[m.cursor]
				if m.mediaType == "audio" {
					m.quality = "best"
					m.step = stepConfirm
				} else {
					m.step = stepQuality
					m.choices = []string{"best", "720", "480", "360"}
					m.cursor = 0
				}
				return m, nil
			}

		case stepQuality:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				m.quality = m.choices[m.cursor]
				m.step = stepConfirm
				return m, nil
			}

		case stepConfirm:
			switch msg.String() {
			case "enter":
				m.step = stepDownloading
				return m, tea.Batch(m.spinner.Tick, m.startDownload())
			}

		case stepDone:
			switch msg.String() {
			case "enter", "q":
				return m, tea.Quit
			}

			return m, nil
		}

	case downloadProgressMsg:
		m.downloadLog = append(m.downloadLog, msg.line)
		if len(m.downloadLog) > 5 {
			m.downloadLog = m.downloadLog[len(m.downloadLog)-5:]
		}

	case spinner.TickMsg:
		if m.step == stepDownloading {
			m.spinner, _ = m.spinner.Update(msg)
			return m, m.spinner.Tick
		}

	case downloadMsg:
		m.filename = msg.filename
		m.err = msg.err
		m.step = stepDone
		return m, nil
	}

	return m, nil
}

func (m model) startDownload() tea.Cmd {
	return func() tea.Msg {
		opts := downloader.Options{
			URL:       m.url,
			OutputDir: m.outputDir + "/%(title)s.%(ext)s",
			AudioOnly: m.mediaType == "audio",
			Quality:   m.quality,
		}

		args := opts.BuildArgs()
		cmd := exec.Command("yt-dlp", args...)

		stdout, _ := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout

		if err := cmd.Start(); err != nil {
			return downloadMsg{err: err}
		}

		var filename string
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			p.Send(downloadProgressMsg{line: line})
			if strings.Contains(line, "Destination:") {
				parts := strings.SplitN(line, "Destination: ", 2)
				if len(parts) == 2 {
					filename = parts[1]
				}
			}
			if strings.Contains(line, "Merging formats into") {
				start := strings.Index(line, "\"")
				end := strings.LastIndex(line, "\"")
				if start != -1 && end != start {
					filename = line[start+1 : end]
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return downloadMsg{err: err}
		}

		err := cmd.Wait()
		if filename == "" {
			filename = "unknown"
		}
		return downloadMsg{filename: filename, err: err}
	}
}
func detectPlatform(urlStr string) string {
	switch {
	case strings.Contains(urlStr, "youtube.com") || strings.Contains(urlStr, "youtu.be"):
		return "YouTube"
	case strings.Contains(urlStr, "tiktok.com"):
		return "TikTok"
	case strings.Contains(urlStr, "instagram.com"):
		return "Instagram"
	case strings.Contains(urlStr, "x.com") || strings.Contains(urlStr, "twitter.com"):
		return "X (Twitter)"
	case strings.Contains(urlStr, "facebook.com") || strings.Contains(urlStr, "fb.watch"):
		return "Facebook"
	case strings.Contains(urlStr, "reddit.com"):
		return "Reddit"
	default:
		u, err := url.Parse(urlStr)
		if err != nil || u.Host == "" {
			return "Unknown"
		}

		host := strings.TrimPrefix(u.Host, "www.")
		parts := strings.Split(host, ".")

		if len(parts) > 0 {
			caser := cases.Title(language.English)
			return caser.String(parts[0])
		}

		return "Unknown"
	}
}
