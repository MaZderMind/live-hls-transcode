package main

import "github.com/alexflint/go-arg"

type CliArguments struct {
	RootDir                         string   `arg:"--root-dir,required,env:ROOT_DIR" help:"Root-Directory to serve- and transcode files from"`
	Extensions                      []string `arg:"env:EXTENSIONS",help:"List of file-extensions for which a stream is offered"`
	TempDir                         string   `arg:"--temp-dir,env:TEMP_DIR" help:"Temporary directory where the transcoding-results will be stored"`
	HttpPort                        string   `arg:"--port,env:PORT" help:"Port to bind to"`
	HttpListen                      string   `arg:"--listen,env:LISTEN" help:"Address to bind to (is. '::' or '127.0.0.1')"`
	LifetimeMinutes                 uint32   `arg:"--lifetime,env:LIFETIME" help:"Number of minutes after which the transcoding-results will be deleted, counted from the last visit"`
	MinimalTranscodeDurationSeconds uint64   `arg:"--minimal-transcode-duration,env:MINIMAL_TRANSCODE_DURATION" help:"Number of seconds after which the transcoding is considered ready"`
}

func (*CliArguments) Version() string {
	return "live-hls-transcode 1.1"
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

		c.cliArguments.Extensions = []string{
			"avi", "ts", "m2ts", "mp2", "mpeg", "mpg", "wmv",
		}
		c.cliArguments.TempDir = "/tmp/live-hls-transcode"
		c.cliArguments.HttpListen = ""
		c.cliArguments.HttpPort = "8048"
		c.cliArguments.LifetimeMinutes = 1440
		c.cliArguments.MinimalTranscodeDurationSeconds = 60

		arg.MustParse(&c.cliArguments)
		c.parsed = true
	}
	return c.cliArguments
}
