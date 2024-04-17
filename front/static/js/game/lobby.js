// Unused now

class Options {
    FPS
    FrameDuration
    SnakeLen
    SnakesAmount
    Duration
    Field
}

class OptionsManager {
    constructor() {
        this.options = new Options()
        this.optionsDiv = document.getElementById("options")
    }

    setOptions(options) {
        this.options = options
        this.optionsDiv.innerText = JSON.stringify(this.options)
    }
}

manager = new OptionsManager()

const HTTP = "http"
const WS = "ws"

class Lobby {
    run(id) {
        this.socket = new WebSocket(`${WS}://${window.location.host}/ws/play/${id}`)
        this.socket.onmessage = function (event) {
            const e = JSON.parse(event.data)
            if (e.event === "update_options") {
                manager.setOptions(e.options)
            }
        }
    }
}
