let socket

const HTTP = "http"
const WS = "ws"

function runGame(url) {
    // TODO убрать
    try {
        socket.close()
    } catch (e) {
    }
    socket = new WebSocket(`${WS}://${window.location.host}/${url}`)
    loadingAnimation.start()
    socket.onmessage = (event) => {
        gameUpdate(JSON.parse(event.data))
    }
}

function rotate(direction) {
    socket.send(JSON.stringify({
        "event": "rotate", "direction": direction
    }))
}

document.addEventListener('keydown', function (event) {
    const key = event.key;
    switch (key) {
        case "ArrowLeft":
            rotate("Left")
            break;
        case "ArrowRight":
            rotate("Right")
            break;
        case "ArrowUp":
            rotate("Up")
            break;
        case "ArrowDown":
            rotate("Down")
            break;
        case "w":
            rotate("Up")
            break
        case "a":
            rotate("Left")
            break
        case "s":
            rotate("Down")
            break
        case "d":
            rotate("Right")
            break
    }
});

function createGame() {
    fetch(`${HTTP}://${window.location.host}/find-hub/`)
        .then(r => r.json())
        .then(data => {
            runGame(`ws/play/${data.id}`)
        })
}

document.getElementById("playBtn")
    .addEventListener("click", () => createGame())


function gameUpdate(data) {
    loadingAnimation.stop()
    const fieldData = [[]];
    const width = 40
    const height = 40

    for (let x = 0; x < width; x++) {
        fieldData.push([])
        for (let y = 0; y < height; y++) {
            fieldData[x].push("E")
        }
    }


    for (let snake in data["snakes"]) {
        for (let i in data["snakes"][snake]["body"]) {
            const point = data["snakes"][snake]["body"][i]
            fieldData[point["Y"]][point["X"]] = "S"
        }
    }
    let fPoint = data["food"]["point"]
    fieldData[fPoint["Y"]][fPoint["X"]] = "F"

    field.visualize(fieldData)
}