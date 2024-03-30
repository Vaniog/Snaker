const logs = document.getElementById("logs")

function log(str) {
    const log = document.createElement("span")
    log.innerHTML = str
    logs.prepend(log)
}