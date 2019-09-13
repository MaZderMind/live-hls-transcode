package main

import "github.com/alexflint/go-arg"

type CliArguments struct {
	RootDir    string `arg:"--root-dir,required"`
	Extensions []string
	TempDir    string `arg:"--temp-dir"`
	HttpBind   string `arg:"--bind"`
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
			"avi", "ts", "m2ts", "mp2", "mpeg", "mpg",
		}
		c.cliArguments.TempDir = "/tmp/live-hls-transcode"
		c.cliArguments.HttpBind = ":8042"

		arg.MustParse(&c.cliArguments)
		c.parsed = true
	}
	return c.cliArguments
}
