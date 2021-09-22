$(function () {
    const maxAge = dayjs.duration({ months: 3 });
    const minDuration = dayjs.duration({ seconds: 30 });
    const now = dayjs();

    const $video = $('video');
    const videoEl = $video.get(0);
    const $dialog = $('#resume-dialog');

    const originalUrl = $video.data('src');
    const playbackInfo = loadPlaybackInfo(originalUrl);

    let actualUrl = originalUrl;
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
        storePlaybackInfo(originalUrl, videoEl.currentTime);
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


    function loadPlaybackInfo(url) {
        const playbackInfoStr = localStorage.getItem(url)
        if (!playbackInfoStr) {
            console.info('No stored playbackInfo');
            return;
        }

        const playbackInfo = JSON.parse(playbackInfoStr);

        const created = dayjs(playbackInfo.created);
        if (created.add(maxAge).isBefore(now)) {
            console.log('playbackInfo is from', created.format(), 'and thus too old');
            return;
        }

        if (playbackInfo.startOffset < minDuration.asSeconds()) {
            console.log('startOffset is only', playbackInfo.startOffset, 'and thus too small');
            return;
        }

        return playbackInfo;
    }

    function storePlaybackInfo(url, currentTime) {
        const playbackInfo = {
            startOffset: currentTime,
            created: dayjs(),
        };
        const playbackInfoStr = JSON.stringify(playbackInfo);
        localStorage.setItem(url, playbackInfoStr);
    }
})
