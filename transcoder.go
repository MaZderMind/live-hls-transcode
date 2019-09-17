package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// 2 Minutes
const minimalWindowMs = 2 * 60000000

type Transcoder struct {
}

func NewTranscoder() Transcoder {
	return Transcoder{}
}

type TranscoderHandle struct {
	sourceFile        string
	destinationFolder string
	cmd               *exec.Cmd
	cancel            context.CancelFunc
	stdOut            io.ReadCloser

	isRunning           bool
	isMinimalWindowDone bool
	isFinished          bool
}

func newTranscoderHandle(
	sourceFile string,
	destinationFolder string,
	cmd *exec.Cmd,
	cancel context.CancelFunc,
	stdOut io.ReadCloser,
) TranscoderHandle {
	return TranscoderHandle{
		sourceFile,
		destinationFolder,
		cmd,
		cancel,
		stdOut,

		false,
		false,
		false,
	}
}

func (transcoder Transcoder) StartTranscoder(sourceFile string, destinationFolder string) *TranscoderHandle {
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

	handle := newTranscoderHandle(sourceFile, destinationFolder, cmd, cancel, readCloser)

	err = cmd.Start()
	if err != nil {
		log.Print("Unable to start Transcoding-Process", err)
		return &handle
	}

	handle.isRunning = true
	go handle.readStdOut()
	return &handle
}

func (handle TranscoderHandle) IsRunning() bool {
	return handle.isRunning
}

func (handle TranscoderHandle) IsReady() bool {
	return handle.isMinimalWindowDone || handle.isFinished
}

func (handle TranscoderHandle) IsFinished() bool {
	return handle.isFinished
}

func (handle *TranscoderHandle) Stop() {
	if ! handle.isRunning {
		return
	}

	handle.cancel()
	log.Printf("%s: Deleting Temp-Dir: %s", handle.sourceFile, handle.destinationFolder)
	err := os.RemoveAll(handle.destinationFolder)

	handle.isRunning = false

	if err != nil {
		log.Printf("%s: Unable to delete Temp-Dir: %s", handle.sourceFile, handle.destinationFolder)
	}
}

func (handle *TranscoderHandle) disarm() {
	handle.cmd = nil
	handle.cancel = nil
}

func (handle *TranscoderHandle) readStdOut() {
	reader := bufio.NewReader(handle.stdOut)

	for ; ; {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			handle.checkExitCode()
			return
		} else if err != nil {
			log.Printf("%s: Error Reading from Transcoder-StdOut, Stopping Process: %s", handle.sourceFile, err)
			handle.Stop()
			handle.disarm()
			return
		}

		k, v := splitKeyValue(line)

		// when the transcoder has processed minimalWindow Miliseconds, we take the Stream as "ready"
		if k == "out_time_ms" {
			ms, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				log.Fatalf("Unable to Parse ffmpeg out_time_ms-Value %v as Int64: %v", v, err)
			}
			if ms > minimalWindowMs {
				handle.isRunning = true
				handle.isMinimalWindowDone = true
			}
		} else if k == "out_time" {
			log.Printf("%s: Transcoding... %s", handle.sourceFile, v)
		}
	}
}

func splitKeyValue(line string) (string, string) {
	pieces := strings.SplitN(strings.Trim(line, "\n"), "=", 2)
	return pieces[0], pieces[1]
}

func (handle *TranscoderHandle) checkExitCode() {
	err := handle.cmd.Wait()

	if err != nil {
		log.Printf("%s: Error while waiting for Child-Process: %v", handle.sourceFile, err)
		handle.Stop()
		handle.disarm()
		return
	}

	if handle.cmd.ProcessState.ExitCode() != 0 {
		log.Printf("%s: Child-Process failed with Exit-Code: %v", handle.sourceFile, handle.cmd.ProcessState.ExitCode())
		handle.Stop()
		handle.disarm()
		return
	}

	log.Printf("%s: Successfully finished Transcoding", handle.sourceFile)
	handle.isRunning = false
	handle.isFinished = true
	handle.disarm()
}
