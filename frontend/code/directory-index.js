$(function () {
	const $search = $('#search');
	const $list = $('#filter-list');

	$search.on('keyup', function () {
		const search = $search.val().trim().toLowerCase();
		const searchEmpty = search.length === 0;

		$list.find('> li').each(function (_, li) {
			const $li = $(li);
			const name = $li.find('a').text();
			const visible = searchEmpty || name.trim().toLowerCase().indexOf(search) !== -1;
			$li.css('display', visible ? 'block' : 'none');
		});
	});

	$list.find('> li').each(function (_, li) {
		const $li = $(li);
		const url = $li.data('url');

		const playbackInfo = loadPlaybackInfo(url);
		if (playbackInfo) {
			let progress = 0.1; // default progrerss for live transcodings which do not have a duration
			if (playbackInfo.startOffset && playbackInfo.duration) {
				progress = playbackInfo.startOffset / playbackInfo.duration
			}
			$li
				.find('.playback-progress')
				.css('width', (progress * 100) + '%');
		}
	});
});