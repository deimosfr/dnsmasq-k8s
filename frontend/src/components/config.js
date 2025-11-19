async function getConfig() {
    const response = await fetch('/api/v1/config');
    const data = await response.text();
    return data;
}

async function displayConfig() {
    const config = await getConfig();
    const configDiv = document.getElementById('config');
    configDiv.textContent = config;
    if (window.hljs) {
        window.hljs.highlightAll();
    }
}

displayConfig();
