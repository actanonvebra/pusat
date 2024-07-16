var input = document.getElementById("input");
var apiUrl = null;

input.addEventListener("focus", function() {
    if (input.value === "Example: 127.0.0.1") {
        input.value = "";
    }
});

input.addEventListener("blur", function() {
    if (input.value === "") {
        input.value = "Example: 127.0.0.1";
    }
});

function checkIpAddress(ip) {
    const ipv4Pattern = 
        /^(\d{1,3}\.){3}\d{1,3}$/;
    const ipv6Pattern = 
        /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/;
    return ipv4Pattern.test(ip) || ipv6Pattern.test(ip);
}

input.addEventListener("keydown", async function(event) {
    if (event.key === "Enter") {
        if (!checkIpAddress(input.value)) {
            alert("Invalid IP address.");
            input.value = "Example: 127.0.0.1";
            return;
        }




        document.getElementById('input').disabled = true;
        var ip = input.value;
        input.value = "PROCESSING...";
        apiUrl = `http://localhost:8000/data?ioc=${ip}`
        data = await getJSON(ip);
        addInfoBox(data);
        input.value = "Example: 127.0.0.1";
        

        document.getElementById('input').disabled = false;
    }
});

async function getJSON(ip) {
    return   fetch(apiUrl)
    .then(response => {
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .catch(error => {
      console.error('Error:', error);
    });
}


function addInfoBox(data) {

    data = JSON.parse(data);
    var infoContainer = document.getElementById("info-container");

    var infoBox = document.createElement("div");
    infoBox.className = "info-box";

    var collapsible = document.createElement("div");
    collapsible.className = "collapsible";

    var checkboxId = "collapsible-head-" + Math.random().toString(36).substring(7);
    var inputCheckbox = document.createElement("input");
    inputCheckbox.type = "checkbox";
    inputCheckbox.id = checkboxId;

    var label = document.createElement("label");
    label.htmlFor = checkboxId;
    label.innerText = data.ip;

    var collapsibleText = document.createElement("div");
    collapsibleText.className = "collapsible-text";

    var openPortsTitle = document.createElement("h2");
    openPortsTitle.innerText = "Open Ports / Protocol";
    collapsibleText.appendChild(openPortsTitle);

    var portsDiv = document.createElement("div");
    portsDiv.className = "ports";

    if (Array.isArray(data.port_data)) {
        data.port_data.forEach(function(port) {
            var p = document.createElement("p");
            p.innerText = `PORT: ${port.port} / PROTOCOL: ${port.protocol.toUpperCase()}`;
            portsDiv.appendChild(p);
        });
    } else {
        var noPortsMessage = document.createElement("p");
        noPortsMessage.innerText = "No port data available.";
        portsDiv.appendChild(noPortsMessage);
    }

    collapsibleText.appendChild(portsDiv);

    var cveCount = document.createElement("p");
    cveCount.className = "info-item";
    cveCount.innerText = "CVE Count: " + (data.cve_count !== undefined ? data.cve_count : "Unknown");
    collapsibleText.appendChild(cveCount);

    var operatingSystem = document.createElement("p");
    operatingSystem.className = "info-item";
    operatingSystem.innerText = "Operating System: " + (data.os !== 0 ? data.os : "Unknown");
    collapsibleText.appendChild(operatingSystem);

    var asnNumber = document.createElement("p");
    asnNumber.className = "info-item";
    asnNumber.innerText = "ASN Number: " + (data.asn !== 0 ? data.asn : "Unknown");
    collapsibleText.appendChild(asnNumber);

    var country = document.createElement("p");
    country.className = "info-item";
    country.innerText = "Country: " + (data.country_code !== 0 ? data.country_code : "Unknown");
    collapsibleText.appendChild(country);

    var status = data.status === 0 ? "Safe" : data.status === 1 ? "Malicious" : "Unknown";
    var statusClass = data.status === 0 ? "status-safe" : data.status === 1 ? "status-malicious" : "";

    var ipStatus = document.createElement("p");
    ipStatus.className = "info-item";
    ipStatus.innerText = `IP Status: ${status}`;
    collapsibleText.appendChild(ipStatus);

    var statusDisplay = document.createElement("div");
    statusDisplay.className = "status " + statusClass;
    statusDisplay.innerText = status;
    collapsibleText.appendChild(statusDisplay);

    collapsible.appendChild(inputCheckbox);
    collapsible.appendChild(label);
    collapsible.appendChild(collapsibleText);
    infoBox.appendChild(collapsible);

    // Yeni bilgi kutusunu en üste ekleyin
    infoContainer.prepend(infoBox);
}





