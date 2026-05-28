package downloader


type Options struct {
	URL       string
	OutputDir string
	AudioOnly bool
	Quality   string
}


func (o Options) BuildArgs() []string {
	args := []string{}

	if o.AudioOnly {
		args = append(args, "-x", "--audio-format", "mp3")
	} else {
		var format string
		switch o.Quality {
		case "720":
			format = "bestvideo[height<=720]+bestaudio/best[height<=720]"
		case "480":
			format = "bestvideo[height<=480]+bestaudio/best[height<=480]"
		case "360":
			format = "bestvideo[height<=360]+bestaudio/best[height<=360]"
		default:
			format = "bestvideo+bestaudio/best"
		}
		args = append(args, "-f", format, "--merge-output-format", "mp4", "--postprocessor-args", "ffmpeg:-c:v libx264 -c:a aac")
	}

	args = append(args, "--no-warnings", "--newline", "--output", o.OutputDir, o.URL)
	return args
}
