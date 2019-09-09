package main

import "github.com/alexflint/go-arg"

type CliArguments struct {
	RootDir    string `arg:"--root-dir,required"`
	Extensions []string
	TempDir    string `arg:"--temp-dir"`
}

func (CliArguments) Version() string {
	return "live-hls-transcode 1.0.0"
}

func ParseCliArguments() CliArguments {
	var args CliArguments
	arg.MustParse(&args)
	return args
}
