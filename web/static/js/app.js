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

    const directionalLight = new THREE.DirectionalLight(0x404040, 0.8);
    directionalLight.position.set(20, 10, 0);
    scene.add(directionalLight);

    const geometry = new THREE.SphereGeometry(2);
    const material = new THREE.MeshLambertMaterial({ color: 0x99ee22 });
    const cube = new THREE.Mesh(geometry, material);
    scene.add(cube);

    camera.position.z = 20;
    function animate() {
        requestAnimationFrame(animate);

        cube.rotation.x += 0.01;
        cube.rotation.y += 0.01;

        renderer.render(scene, camera);
    };
    animate();
}

window.addEventListener('DOMContentLoaded', initApp);