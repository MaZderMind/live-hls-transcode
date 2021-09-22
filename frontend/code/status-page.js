$(function () {
	const replaceIntervalMs = 5 * 1000;
	if ($('[data-replace]').length > 0) {
		console.log('initialize live-update');

		setInterval(function () {
			if (window['STOP_UPDATE'] || (window.localStorage['STOP_UPDATE'] === 'true')) {
				return
			}

			$.ajax({
				url: window.location.href,
				dataType: 'html',
				success: function (html) {
					const $newDom = $('<div>').html(html);

					let autoplayEnabled = $('.autoplay').length > 0;
					let isReady = $newDom.find('[data-isready]').data('isready');
					if (autoplayEnabled && isReady) {
						console.log('stream is ready, redirecting to playlist');
						window.location.href = '?stream&playlist'
					}

					const $newReplaceables = $newDom.find('[data-replace]');
					$newReplaceables.each(function (_, newReplaceable) {
						const key = $(newReplaceable).data('replace');
						console.log('updating replaceable block', key);
						$('[data-replace=' + key + ']').replaceWith(newReplaceable);
					});
				}
			})
		}, replaceIntervalMs)
	}
});
