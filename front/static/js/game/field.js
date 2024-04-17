class Field {
    constructor(canvasElement) {
        this.canvas = canvasElement;
        this.ctx = canvasElement.getContext('2d');
        this.colors = {
            default: '#ffffff' // Default color is white
        };
        this.cellSize = 13; // Adjust as needed
    }

    visualize(data) {
        this.canvas.width = data[0].length * this.cellSize;
        this.canvas.height = data.length * this.cellSize;

        for (let i = 0; i < data.length; i++) {
            for (let j = 0; j < data[i].length; j++) {
                this.ctx.fillStyle = this.getColor(data[i][j]);
                this.ctx.fillRect(j * this.cellSize, i * this.cellSize, this.cellSize, this.cellSize);
            }
        }
    }

    setColor(value, color) {
        this.colors[value] = color;
    }

    getColor(value) {
        return this.colors[value] || this.colors.default;
    }
}

const field = new Field(canvasElement);

field.setColor('E', '#075264');
field.setColor('S', '#79d4fd');
field.setColor('F', '#cab100');

function defaultField() {
    const fieldData = [];
    const width = 40;
    const height = 40;

    for (let x = 0; x < height; x++) {
        fieldData.push([]);
        for (let y = 0; y < width; y++) {
            fieldData[x].push("E");
        }
    }
    field.visualize(fieldData);
}

defaultField();
