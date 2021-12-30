function initThreeScene() {
    const wInnerHeight = window.innerHeight;
    const wInnerWidth = window.innerWidth;
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, wInnerWidth / wInnerHeight, 0.1, 1000);

    const renderer = new THREE.WebGLRenderer({
        antialias: true,
    });
    renderer.outputEncoding = THREE.sRGBEncoding;
    renderer.setSize(wInnerWidth, wInnerHeight);
    document.body.appendChild(renderer.domElement);
    return {
        camera,
        renderer,
        scene,
    }
}

function makeNode(meshOptions) {
    const geometry = new THREE.SphereGeometry(2);
    const material = new THREE.MeshLambertMaterial(meshOptions);
    return new THREE.Mesh(geometry, material);
}

const NODE_COLORS = [0x99ee22, 0xee9922, 0x9922ee, 0x22ee99, 0x2299ee];
function animate(state) {
    if (!state.running) {
        return;
    }

    requestAnimationFrame(() => animate(state));

    const {
        camera,
        nodes,
        frameQueue,
        renderer,
        scene,
        traceLines,
        traceMeshes,
    } = state;

    const frame = frameQueue.shift();
    if (!frame) {
        state.running = false;
        return;
    }

    while (traceMeshes.length) {
        const traceMesh = traceMeshes.pop();
        traceMesh.geometry.dispose();
        traceMesh.material.dispose();
        scene.remove(traceMesh);
    }

    const frameNodeIds = new Set();
    let i = 0;
    for (let nodeId of frame.ids) {
        frameNodeIds.add(nodeId);
        // New node, add to scene.
        if (!nodes.has(nodeId)) {
            nodes.set(nodeId, makeNode({ color: NODE_COLORS[nodes.size % NODE_COLORS.length] }));
            traceLines.set(nodeId, []);
            scene.add(nodes.get(nodeId));
        }

        // Update the node's position.
        nodes.get(nodeId).position.set(
            frame.positions[i],
            frame.positions[i + 1],
            frame.positions[i + 2],
        );

        if (traceLines.get(nodeId).length >= 120) {
            traceLines.get(nodeId).shift();
        }
        traceLines.get(nodeId).push(new THREE.Vector3(
            frame.positions[i],
            frame.positions[i + 1],
            frame.positions[i + 2],
        ));

        const lineMaterial = new THREE.LineBasicMaterial({
            color: 0xffffff,
            opacity: 0.4,
            transparent: true,
        });
        const geometry = new THREE.BufferGeometry().setFromPoints(traceLines.get(nodeId));
        const newTraceMesh = new THREE.Line(geometry, lineMaterial);
        traceMeshes.push(newTraceMesh);
        scene.add(newTraceMesh);

        i += 3;
    }

    // Remove any nodes no longer being rendered.
    nodes.forEach((node, key) => {
        if (!frameNodeIds.has(key)) {
            traceLines.delete(key);
            node.geometry.dispose();
            node.material.dispose();
            scene.remove(node);
            nodes.delete(key);
        }
    });

    renderer.render(scene, camera);
};

function initDataLink(animationState) {
    const socket = new WebSocket("wss://localhost:8822/sim");
    let requestingFrames = false;
    function requestFrames(interval) {
        if (requestingFrames) {
            console.log("Frame request in progress... returning early.")
            return;
        }

        if (socket.readyState !== socket.OPEN) {
            console.log("Websocket closed... exiting.");
            window.clearInterval(interval);
            return;
        }

        if (animationState.frameQueue.length > 10) {
            return;
        }

        requestingFrames = true;
        console.log(`Requesting more frames. Frames remaining: ${animationState.frameQueue.length}`);
        socket.send("get_frames");
    }

    socket.addEventListener('message', (event) => {
        requestingFrames = false;
        console.log(`Received frames. Frames remaining: ${animationState.frameQueue.length}`);
        animationState.frameQueue.push(...JSON.parse(event.data));
        if (!animationState.running) {
            animationState.running = true;
            requestAnimationFrame(() => animate(animationState));
        }
    });
    socket.addEventListener('open', () => {
        console.log('Websocket connected. Initializing polling.');
        const interval = window.setInterval(() => requestFrames(interval), 100);
        requestFrames(interval);
        setTimeout(() => {
            if (socket.readyState === socket.OPEN) {
                socket.close()
            }
        }, 250000);
    });
}

function initApp() {
    const {
        camera,
        renderer,
        scene,
    } = initThreeScene();

    const animationState = {
        camera,
        frameQueue: [],
        nodes: new Map(),
        renderer,
        running: false,
        scene,
        traceLines: new Map(),
        traceMeshes: [],
    };

    scene.fog = new THREE.Fog(0xffffff, 10, 1000);

    const directionalLight = new THREE.DirectionalLight(0x404040, 0.8);
    directionalLight.position.set(5, 2, 0);
    scene.add(directionalLight);
    
    const directionalLight2 = new THREE.DirectionalLight(0x404040, 0.2);
    directionalLight2.position.set(-5, -2, 0);
    scene.add(directionalLight2);

    const directionalLight3 = new THREE.DirectionalLight(0x202020, 0.9);
    directionalLight3.position.set(0, 0, 20);
    scene.add(directionalLight3);

    camera.position.set(0, -20, 50);
    camera.lookAt(0, 0, 0);
    window.camera = camera;

    initDataLink(animationState);
}

window.addEventListener('DOMContentLoaded', initApp);