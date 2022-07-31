package main

import "github.com/alexflint/go-arg"

type CliArguments struct {
	RootDir                         string   `arg:"--root-dir,required,env:ROOT_DIR" help:"Root-Directory to serve- and transcode files from"`
	TranscodeExtensions             []string `arg:"env:TRANSCODE_EXTENSIONS" help:"List of file-extensions for which a transcoding is offered"`
	PlayerExtensions                []string `arg:"env:PLAYER_EXTENSIONS" help:"List of file-extensions for which a player is offered"`
	TempDir                         string   `arg:"--temp-dir,env:TEMP_DIR" help:"Temporary directory where the transcoding-results will be stored"`
	HttpPort                        string   `arg:"--port,env:PORT" help:"Port to bind to"`
	HttpListen                      string   `arg:"--listen,env:LISTEN" help:"Address to bind to (is. '::' or '127.0.0.1')"`
	LifetimeMinutes                 uint32   `arg:"--lifetime,env:LIFETIME" help:"Number of minutes after which the transcoding-results will be deleted, counted from the last visit"`
	MinimalTranscodeDurationSeconds uint64   `arg:"--minimal-transcode-duration,env:MINIMAL_TRANSCODE_DURATION" help:"Number of seconds after which the transcoding is considered ready"`
	Acceleration                    string   `arg:"--acceleration,env:ACCELERATION" help:"Use Hardware-Acceleration (Supported Values: none (Uses libx264 on the CPU) and h264_v4l2m2m (Uses the v4l2 m2m Accelerator on newer Raspberry-Pis)"`
}

func (*CliArguments) Version() string {
	return "live-hls-transcode 1.1.1"
}

func (arguments *CliArguments) HttpBind() string {
	return arguments.HttpListen + ":" + arguments.HttpPort
}

type CliArgumentsParser struct {
	cliArguments CliArguments
	parsed       bool
}

func NewCliArgumentsParser() CliArgumentsParser {
	return CliArgumentsParser{}
}

func (c CliArgumentsParser) GetCliArguments() CliArguments {
	if !c.parsed {

		c.cliArguments.TranscodeExtensions = []string{
			"avi", "ts", "m2ts", "mp2", "mpeg", "mpg", "wmv",
		}
		c.cliArguments.PlayerExtensions = []string{
			"mp4", "m4v",
		}
		c.cliArguments.TempDir = "/tmp/live-hls-transcode"
		c.cliArguments.HttpListen = ""
		c.cliArguments.HttpPort = "8048"
		c.cliArguments.LifetimeMinutes = 1440
		c.cliArguments.MinimalTranscodeDurationSeconds = 60
		c.cliArguments.Acceleration = "none"

		arg.MustParse(&c.cliArguments)
		c.parsed = true
	}
	return c.cliArguments
}
