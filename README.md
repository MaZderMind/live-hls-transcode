# File to Live-HLS-Transcoder
Deeply buried in the basements of our NAS-systems and storage-servers, we all have one: a more or less large collection of video-files in
peculiar formats from the late 90s. AVIs, DivX'es, some WMVs or even some old and dusted FLVs? Not important enough or too many, to
manually convert them to a more modern format, but still, on a rainy day, while browsing your NAS with your Phone/Tablet, you wonder
– what is this? Might it be funny? A blast fro mthe past? Remind me of the good old days?

That's where this very specialized piece of modern software art got you covered:

You can browse through your files with a mobile-friendly Web-UI (better than the apache/nginx directory index!), access any files that
your device can natively read and start a ffmpeg-based transcoder process to h264/AAC in HLS for the files it can't read. After only some
seconds you can start to stream this funny little clip you laughed so hard over when you were a teen to our device.

The HLS-Stream is of type *EVENT*, so for compatible Players* it starts as a Live-Stream, that soon becomes navigatable as the
transcoding continues to convert more and more of the source material. Once the transcoder is done completely, the HLS-Playlist gets
closed and behaves just like a normal file.

After a configurable lifetime after the last access (default: 24 hours), the transcoded files are deleted.

*) Tested with kind'a recent iOS Devices

## Screenshots
![](doc/1.png | width=150)
![](doc/2.png | width=150)
![](doc/1.png | width=150)
![](doc/1.png | width=150)


## Build & Run
```
go get github.com/MaZderMind/live-hls-transcode
cd $GOPATH/src/MaZderMind/live-hls-transcode

make dependencies
make
./live-hls-transcode --root-dir /video/
-> http://localhost:8048
```

## Run Dev-Server (with hot-reload)
```
go get github.com/MaZderMind/live-hls-transcode
cd $GOPATH/src/MaZderMind/live-hls-transcode

make dependencies
make run
```
