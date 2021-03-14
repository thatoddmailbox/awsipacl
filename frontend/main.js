window.addEventListener("load", function() {
	document.getElementById("test").addEventListener("click", function() {
		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/login");
		xhr.onreadystatechange = function() {
			if (xhr.readyState == 4) {
				var response = xhr.responseText;
				console.log(response);
				console.log(JSON.parse(response));
			}
		};
		xhr.send();
	});
});