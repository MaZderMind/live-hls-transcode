const maxAge = dayjs.duration({ months: 3 });
const minDuration = dayjs.duration({ seconds: 30 });
const now = dayjs();

function loadPlaybackInfo(url) {
    const playbackInfoStr = localStorage.getItem(url)
    if (!playbackInfoStr) {
        console.info('No stored playbackInfo for', url);
        return;
    }

    const playbackInfo = JSON.parse(playbackInfoStr);

    const created = dayjs(playbackInfo.created);
    if (created.add(maxAge).isBefore(now)) {
        console.log('playbackInfo for ', url, 'is from', created.format(), 'and thus too old');
        return;
    }

    if (playbackInfo.startOffset < minDuration.asSeconds()) {
        console.log('startOffset for ', url, 'is only', playbackInfo.startOffset, 'and thus too small');
        return;
    }

    console.log('playbackInfo for ', url, 'is', playbackInfo);
    return playbackInfo;
}

function storePlaybackInfo(url, currentTime, duration) {
    const playbackInfo = {
        startOffset: currentTime,
        duration: duration,
        created: dayjs(),
    };

    const playbackInfoStr = JSON.stringify(playbackInfo);
    localStorage.setItem(url, playbackInfoStr);
}