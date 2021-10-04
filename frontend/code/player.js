$(function () {
    const $video = $('video');
    const videoEl = $video.get(0);
    const $dialog = $('#resume-dialog');

    const originalUrl = $video.data('url');
    const playbackInfo = loadPlaybackInfo(originalUrl);

    let actualUrl = $video.data('src');
    if (playbackInfo) {
        $dialog
            .find('#position')
            .text(
                dayjs.duration({ seconds: playbackInfo.startOffset }).humanize()
            )
            .end()
            .find('#date')
            .text(
                dayjs.duration(dayjs(playbackInfo.created).diff(now)).humanize(true)
            )
            .end()
            .css('display', 'block');

        actualUrl += '#t=' + playbackInfo.startOffset;
    }
    $video.attr('src', actualUrl);

    $video.on('timeupdate', function () {
        storePlaybackInfo(originalUrl, videoEl.currentTime, videoEl.duration);
    });

    const $resumeButton = $dialog.find('#resume-button');
    const $resetButton = $dialog.find('#reset-button')
    $resumeButton.on('click', function (e) {
        e.preventDefault();

        videoEl.play();
        $dialog.css('display', 'none');
    });

    $resetButton.on('click', function (e) {
        e.preventDefault();

        videoEl.currentTime = 0;
        videoEl.play();
        $dialog.css('display', 'none');
    });

    $video.on('play', function (e) {
        $dialog.css('display', 'none');
    })
})
