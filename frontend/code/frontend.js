$(function() {
	var replaceIntervalMs = 5 * 1000;
	if ($('[data-replace]').length > 0) {
		console.log('initialize live-update');

		setInterval(function() {
			$.ajax({
				url: window.location.href,
				dataType: 'html',
				success: function(html) {
					var $newDom = $('<div>').html(html);

					var $newReplaceables = $newDom.find('[data-replace]');
					$newReplaceables.each(function(_, newReplaceable) {
						var key = $(newReplaceable).data('replace');
						console.log('updating replaceable block', key);
						$('[data-replace=' + key + ']').replaceWith(newReplaceable);
					});

					if ($('[data-isready]').data('isready') && window.location.href.indexOf('autoplay') !== -1) {
						console.log('stream is ready, redirecting to playlist');
						window.location.href = '?stream&playlist'
					}
				}
			})
		}, replaceIntervalMs)
	}
});
