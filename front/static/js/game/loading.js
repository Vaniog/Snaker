function drawLoading(canvasElement) {
    const ctx = canvasElement.getContext('2d');
    ctx.fillStyle = "#075264";
    const radius = 40;
    let angle = 0;
    let animationId = null;
    const numLines = 12;
    let isAnimating = false;

    function drawFrame() {
        ctx.clearRect(0, 0, canvasElement.width, canvasElement.height);
        ctx.fillRect(0, 0, canvasElement.width, canvasElement.height);
        ctx.save();
        ctx.translate(canvasElement.width / 2, canvasElement.height / 2);
        ctx.rotate(angle);

        for (let i = 0; i < numLines; i++) {
            ctx.beginPath();
            ctx.rotate(Math.PI / (numLines / 2));
            ctx.moveTo(radius, 0);
            ctx.lineTo(radius * 0.8, 0);
            ctx.stroke();
        }

        ctx.restore();
        angle += Math.PI / 100; // Adjust the speed of rotation
        animationId = requestAnimationFrame(drawFrame);
    }

    function startLoading() {
        if (!isAnimating) {
            isAnimating = true;
            drawFrame();
        }
    }

    function stopLoading() {
        if (isAnimating) {
            isAnimating = false;
            cancelAnimationFrame(animationId);
            ctx.clearRect(0, 0, canvasElement.width, canvasElement.height);
        }
    }

    return {
        start: startLoading, stop: stopLoading
    };
}

const canvasElement = document.getElementById('canvas');
const loadingAnimation = drawLoading(canvasElement);