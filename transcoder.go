package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Transcoder struct {
}

func NewTranscoder() Transcoder {
	return Transcoder{}
}

type TranscoderHandle struct {
	sourceFile        string
	destinationFolder string
	cancel            context.CancelFunc
	stdOut            io.ReadCloser

	isRunning  bool
	isReady    bool
	isFinished bool
}

func newTranscoderHandle(sourceFile string, destinationFolder string, cancel context.CancelFunc, stdOut io.ReadCloser) TranscoderHandle {
	return TranscoderHandle{
		sourceFile,
		destinationFolder,
		cancel,
		stdOut,

		false,
		false,
		false,
	}

}

func (transcoder Transcoder) StartTranscoder(sourceFile string, destinationFolder string) TranscoderHandle {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		"-v", "error",
		"-hide_banner",
		"-y",
		"-nostdin",
		"-progress", "pipe:1",
		"-threads", "16",
		"-i", sourceFile,
		"-analyzeduration", "20000000",
		"-c:v:0", "libx264",
		"-pix_fmt", "yuv420p",
		"-bufsize", "8192k",
		"-crf", "20",
		"-minrate", "100k",
		"-maxrate", "6000k",
		"-profile", "main",
		"-level", "4.0",
		"-threads", "8",
		"-c:a", "aac",
		"-b:a", "192k",
		"-ar:a", "48000",
		"-flags", "+cgop",
		"-g", "60",
		"-hls_playlist_type", "event",
		"-hls_time", "5",
		path.Join(destinationFolder, "index.m3u8"),
	)

	log.Printf("%s: Starting Transcoder-Command: %v", sourceFile, cmd)
	readCloser, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Unable to open StdoutPipe", err)
	}

	handle := newTranscoderHandle(sourceFile, destinationFolder, cancel, readCloser)

	err = cmd.Start()
	if err != nil {
		log.Print("Unable to start Transcoding-Process", err)
		return handle
	}

	handle.isRunning = true
	go handle.readStdOut()
	return handle
}

func (handle TranscoderHandle) IsRunning() bool {
	return handle.isRunning
}

func (handle TranscoderHandle) IsReady() bool {
	return handle.isReady
}

func (handle TranscoderHandle) IsFinished() bool {
	return handle.isFinished
}

func (handle TranscoderHandle) Stop() {
	handle.cancel()
	log.Printf("%s: Deleting Temp-Dir: %s", handle.sourceFile, handle.destinationFolder)
	err := os.RemoveAll(handle.destinationFolder)
	if err != nil {
		log.Printf("%s: Unable to delete Temp-Dir: %s", handle.sourceFile, handle.destinationFolder)
	}
}

func (handle TranscoderHandle) readStdOut() {
	reader := bufio.NewReader(handle.stdOut)

	for ; ; {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Printf("%s: Error Reading from Transcoder-StdOut, Stopping Process", handle.sourceFile)
			handle.Stop()
			handle.isRunning = false
			return
		}

		log.Printf("%s: Read Line: >%s<", handle.sourceFile, strings.Trim(line, "\n"))
	}
}
