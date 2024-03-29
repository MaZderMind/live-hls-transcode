package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type Transcoder struct {
	minimalTranscodeDurationSeconds uint64
	acceleration                    string
}

func NewTranscoder(minimalTranscodeDurationMilliseconds uint64, acceleration string) Transcoder {
	return Transcoder{
		minimalTranscodeDurationMilliseconds,
		acceleration,
	}
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

	totalDuration     time.Duration
	processedDuration time.Duration

	transcoder *Transcoder
}

func newTranscoderHandle(
	sourceFile string,
	destinationFolder string,
	cmd *exec.Cmd,
	cancel context.CancelFunc,
	stdOut io.ReadCloser,

	totalDuration time.Duration,
	transcoder *Transcoder,
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

		totalDuration,
		time.Duration(0),

		transcoder,
	}
}

func (transcoder Transcoder) StartTranscoder(sourceFile string, destinationFolder string) *TranscoderHandle {
	log.Printf("%s: Probing Duration", sourceFile)
	duration := transcoder.ProbeDuration(sourceFile)

	cmdHead := []string{
		"-v", "error",
		"-hide_banner",
		"-y",
		"-nostdin",
		"-progress", "pipe:1",
		"-i", sourceFile,
		"-analyzeduration", "20000000",
	}

	cmdAccNone := []string{
		"-threads", "16",
		"-c:v:0", "libx264",
		"-pix_fmt", "yuv420p",
		"-bufsize", "8192k",
		"-crf", "20",
		"-minrate", "100k",
		"-maxrate", "6000k",
		"-profile", "main",
		"-level", "4.0",
	}

	cmdAccH264V4l2m2m := []string{
		"-threads", "8",
		"-vf", "scale=iw*sar:ih,setsar=1:1",
		"-c:v:0", "h264_v4l2m2m",
		"-pix_fmt", "yuv420p",
		"-bufsize", "8192k",
		"-crf", "20",
		"-b:v:0", "4M",
		"-level", "4.0",
	}

	cmdAudio := []string{
		"-threads", "8",
		"-c:a", "aac",
		"-b:a", "192k",
		"-ar:a", "48000",
	}

	cmdTail := []string{
		"-flags", "+cgop",
		"-g", "60",
		"-hls_playlist_type", "event",
		"-hls_time", "10",
		path.Join(destinationFolder, "index.m3u8"),
	}

	cmdArgs := make([]string, 0)
	cmdArgs = append(cmdArgs, cmdHead...)
	if transcoder.acceleration == "h264_v4l2m2m" {
		cmdArgs = append(cmdArgs, cmdAccH264V4l2m2m...)
	} else {
		cmdArgs = append(cmdArgs, cmdAccNone...)
	}
	cmdArgs = append(cmdArgs, cmdAudio...)
	cmdArgs = append(cmdArgs, cmdTail...)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "ffmpeg", cmdArgs...)

	log.Printf("%s: Starting Transcoder-Command: %v", sourceFile, cmd)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Unable to open StdoutPipe", err)
	}

	handle := newTranscoderHandle(
		sourceFile,
		destinationFolder,
		cmd,
		cancel,
		stdOut,
		duration,
		&transcoder,
	)

	err = cmd.Start()
	if err != nil {
		log.Print("Unable to start Transcoding-Process", err)
		return &handle
	}

	handle.isRunning = true
	go handle.readStdOut()
	return &handle
}

func (transcoder Transcoder) ProbeDuration(source string) time.Duration {
	log.Printf("%s: Probing Duration", source)
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		source,
	)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("%s: Unable to probe Duration", source)
		return time.Duration(0)
	}

	ffprobeOutput := struct {
		Format struct {
			Duration string
		}
	}{}

	err = json.Unmarshal(output, &ffprobeOutput)
	if err != nil {
		log.Printf("%s: Unable to read stdout-Format: %s\n%s", source, err, output)
		return time.Duration(0)
	}

	numSeconds, err := strconv.ParseFloat(ffprobeOutput.Format.Duration, 64)
	if err != nil {
		log.Printf("%s: Unable to parse as float64: %s\n%s", source, err, ffprobeOutput.Format.Duration)
		return time.Duration(0)
	}

	duration := time.Duration(numSeconds) * time.Second
	log.Printf("%s: Probed duration to be %s", source, duration)
	return duration
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

func (handle TranscoderHandle) TotalDuration() time.Duration {
	return handle.totalDuration
}

func (handle TranscoderHandle) ProcessedDuration() time.Duration {
	return handle.processedDuration
}

func (handle TranscoderHandle) ProcessedPercent() float64 {
	return handle.processedDuration.Seconds() / handle.totalDuration.Seconds() * 100
}

func (handle *TranscoderHandle) Stop() {
	if !handle.isRunning {
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

	for {
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

		if k == "out_time_ms" {
			microseconds, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				log.Printf("Unable to Parse ffmpeg out_time_ms-Value %v as Int64: %v", v, err)
				microseconds = 0
			}
			if microseconds > handle.transcoder.minimalTranscodeDurationSeconds*1_000_000 {
				handle.isRunning = true
				handle.isMinimalWindowDone = true
			}

			processedDurationExact := time.Duration(microseconds) * time.Microsecond
			handle.processedDuration = processedDurationExact.Round(time.Second)
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
