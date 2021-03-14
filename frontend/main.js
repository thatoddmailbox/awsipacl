var password;

function api(path, data, callback) {
	var params = new URLSearchParams();
	for (var key in data) {
		params.set(key, data[key]);
	}

	var xhr = new XMLHttpRequest();
	xhr.open("POST", "/" + path);
	xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
	xhr.onreadystatechange = function() {
		if (xhr.readyState == 4) {
			callback(JSON.parse(xhr.responseText));
		}
	};
	xhr.send(params);
}

function screen(name) {
	document.querySelector(".screen:not(.hidden)").classList.add("hidden");
	document.getElementById("screen" + name).classList.remove("hidden");
}

window.addEventListener("load", function() {
	document.getElementById("screenPasswordSubmit").addEventListener("click", function() {
		password = document.getElementById("screenPasswordInput").value;

		api("login", {
			password: password
		}, function(data) {
			if (data.status == "error") {
				alert(data.error);
				return;
			}

			document.getElementById("screenMainTitle").innerText = data.title;
			document.getElementById("screenMainDescription").innerText = data.description;
			document.getElementById("screenMainCurrentIP").innerText = data.clientIP;

			var ipsElement = document.getElementById("screenMainIPs");
			for (var i = 0; i < data.ips.length; i++) {
				var ip = data.ips[i];

				var ipRow = document.createElement("tr");

				var ipRowIP = document.createElement("td");
				ipRowIP.innerText = ip.ip;
				ipRow.appendChild(ipRowIP);

				var ipRowDescription = document.createElement("td");
				ipRowDescription.innerText = ip.description;
				ipRow.appendChild(ipRowDescription);

				var ipRowActions = document.createElement("td");
				ipRow.appendChild(ipRowActions);

				ipsElement.appendChild(ipRow);
			}

			screen("Main");
		});
	});
});