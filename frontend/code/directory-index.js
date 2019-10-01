var $search = $('#search');
var $list = $('#filter-list');
$search.on('keyup', function() {
	let search = $search.val().trim().toLowerCase();
	var searchEmpty = search.length === 0;

	$list.find('> li').each(function(_, li) {
		var $li = $(li);
		var name = $li.find('a').text();
		var visible = searchEmpty || name.trim().toLowerCase().indexOf(search) !== -1;
		$li.css('display', visible ? 'block' : 'none');
	});
});
