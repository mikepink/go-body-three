function initThreeScene() {
    const wInnerHeight = window.innerHeight;
    const wInnerWidth = window.innerWidth;
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, wInnerWidth / wInnerHeight, 0.1, 1000);

    const renderer = new THREE.WebGLRenderer();
    renderer.setSize(wInnerWidth, wInnerHeight);
    document.body.appendChild(renderer.domElement);
    return {
        camera,
        renderer,
        scene,
    }
}

function initApp() {
    const {
        camera,
        renderer,
        scene,
    } = initThreeScene();

    const geometry = new THREE.BoxGeometry();
    const material = new THREE.MeshBasicMaterial({ color: 0xddffaa });
    const cube = new THREE.Mesh(geometry, material);
    scene.add(cube);

    camera.position.z = 5;
    function animate() {
        requestAnimationFrame(animate);

        cube.rotation.x += 0.01;
        cube.rotation.y += 0.01;

        renderer.render(scene, camera);
    };
    animate();
}

window.addEventListener('DOMContentLoaded', initApp);