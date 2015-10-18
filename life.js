window.addEventListener('load', function() {
	var b = document.getElementById('board');
	var running;
	
	var q = document.getElementById('quit');
	q.addEventListener('click', function() {
		replace(b, "/quit");
		clearInterval(running);
	});
	
	var r = document.getElementById('random');
	r.addEventListener('click', function() {
		replace(b, "/random");
	});
	
	var s = document.getElementById('step');
	s.addEventListener('click', function() {
		replace(b, "/step");
	});
	
	s = document.getElementById('start');
	s.addEventListener('click', function() {
			running = setInterval(function(){
			replace(b, "/step");
		}, 200);
		replace(b, "/step");
	});
	
	s = document.getElementById('stop');
	s.addEventListener('click', function() {
		clearInterval(running);
	});
});

function replace(node, target) {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', encodeURI(target));
	xhr.onload = function() {
		if (xhr.status === 200) {
			node.innerHTML = xhr.responseText;
		} else {
			alert('Request failed.  Returned status of ' + xhr.status);
		}
	};
	xhr.send();
}
