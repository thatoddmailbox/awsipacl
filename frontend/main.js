window.addEventListener("load", function() {
	document.getElementById("screenPasswordSubmit").addEventListener("click", function() {
		var password = document.getElementById("screenPasswordInput").value;
		var params = new URLSearchParams();
		params.set("password", password);

		var xhr = new XMLHttpRequest();
		xhr.open("POST", "/login");
		xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
		xhr.onreadystatechange = function() {
			if (xhr.readyState == 4) {
				var response = xhr.responseText;
				console.log(JSON.parse(response));
			}
		};
		xhr.send(params);
	});
});