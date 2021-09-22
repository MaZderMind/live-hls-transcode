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
			$li
				.find('.playback-progress')
				.css('width', (playbackInfo.startOffset / playbackInfo.duration * 100) + '%');
		}
	});
});