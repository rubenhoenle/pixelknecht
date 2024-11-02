const Endpoint = Object.freeze({
  MODE_API_URL: "/mode",
  SERVER_API_URL: "/api/server",
});

function loadData(endpoint) {
  fetch(endpoint, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Failed to update configuration");
    })
    .then((data) => {
      console.log("Response from server:", data);

      if (endpoint == Endpoint.MODE_API_URL) {
        document.getElementById("enable").checked = data.enabled;
        document.getElementById("imageUrl").value = data.imageUrl;
        document.getElementById("x").value = data.posX;
        document.getElementById("y").value = data.posY;

        document.getElementById("configResult").style.display = "block";
        document.getElementById("configResult").textContent =
          "Mode configuration loaded successfully!";

        const resultImage = document.getElementById("resultImage");
        resultImage.src = data.imageUrl;

        document.getElementById("configResult").style.display = "block";
      } else if (endpoint == Endpoint.SERVER_API_URL) {
        document.getElementById("host").value = data.host;
        document.getElementById("port").value = data.port;
      }
    })
    .catch((error) => {
      console.error("Error:", error);
      document.getElementById("configResult").style.display = "block";
      document.getElementById("configResult").textContent =
        "Failed to load configuration. Please try again.";
    });
}

function updateConfig(endpoint, body) {
  fetch(endpoint, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  })
    .then((response) => {
      if (response.ok) {
        return response.json();
      }
      throw new Error("Failed to update configuration");
    })
    .then((data) => {
      console.log("Response from server:", data);
      if (endpoint == Endpoint.MODE_API_URL) {
        document.getElementById("configResult").style.display = "block";
        document.getElementById("configResult").textContent =
          "Configuration sent successfully!";

        const resultImage = document.getElementById("resultImage");
        resultImage.src = imageUrl;
        resultImage.style.objectPosition = `${x} ${y}`;

        document.getElementById("configResult").style.display = "block";
      }
      console.error("Error:", error);
      document.getElementById("configResult").style.display = "block";
      document.getElementById("configResult").textContent =
        "Failed to send configuration. Please try again.";
    });
}

function submitServerConfig() {
  const host = document.getElementById("host").value;
  const port = parseInt(document.getElementById("port").value, 10);

  const serverData = {
    host: host,
    port: port,
  };

  updateConfig(Endpoint.SERVER_API_URL, serverData);
}

function submitModeConfig() {
  const enable = document.getElementById("enable").checked;
  const imageUrl = document.getElementById("imageUrl").value;
  const x = parseInt(document.getElementById("x").value, 10);
  const y = parseInt(document.getElementById("y").value, 10);

  const modeData = {
    enabled: enable,
    posY: y,
    posX: x,
    imageUrl: imageUrl,
  };

  updateConfig(Endpoint.MODE_API_URL, modeData);
}

window.onload = function () {
  loadData(Endpoint.MODE_API_URL);
  loadData(Endpoint.SERVER_API_URL);
};
