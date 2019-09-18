package main

import "github.com/alexflint/go-arg"

type CliArguments struct {
	RootDir         string   `arg:"--root-dir,required" help:"Root-Directory to serve- and transcode files from"`
	Extensions      []string `help:"List of file-extensions for which a stream is offered"`
	TempDir         string   `arg:"--temp-dir" help:"Temporary directory where the transcoding-results will be stored"`
	HttpBind        string   `arg:"--bind" help:"IP-Address and Port to bind to (ie. '::8042' or '127.0.0.1:8000')"`
	LifetimeMinutes int32    `arg:"--lifetime" help:"Number of minutes after which the transcoding-results will be deleted, counted from the last visit"`
}

func (CliArguments) Version() string {
	return "live-hls-transcode 1.0.0"
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
		c.cliArguments.HttpBind = ":8042"
		c.cliArguments.LifetimeMinutes = 1440

		arg.MustParse(&c.cliArguments)
		c.parsed = true
	}
	return c.cliArguments
}
