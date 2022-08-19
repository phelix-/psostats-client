import * as THREE from '/js/three.module.js';
import { OrbitControls } from '/js/OrbitControls.js';

var scene, camera, renderer, controls, draughts, board;
var frame = 0;
var dataFrames;
var players = {};
var monsters = {};
const playerColors = ["red", "blue", "green", "yellow"];
const exampleSphere = new THREE.SphereGeometry(5,8,8)
const canvas = document.getElementById("map-canvas")
let paused = true;
const pauseButton = document.getElementById("pause-button")

function init() {
    draughts = new Draughts();
    pauseButton.onclick = function() {
        paused = !paused;
    }

    scene = new THREE.Scene();
    camera = new THREE.PerspectiveCamera(75, canvas.clientWidth / canvas.clientHeight, 0.1, 5000);

    renderer = new THREE.WebGLRenderer({canvas: canvas});
    renderer.setSize(canvas.clientWidth, canvas.clientHeight);
    // document.body.appendChild(renderer.domElement);

    camera.position.y = 500;
    camera.position.z = 0;

    controls = new OrbitControls(camera, renderer.domElement);

    controls.target.set(4.5, 0, 4.5);

    // controls.enablePan = false;
    controls.maxPolarAngle = Math.PI / 2;

    controls.enableDamping = true;
    dataFrames = jsonMeshes.dataFrames
    const meshes = jsonMeshes.meshes
    const normals = jsonMeshes.normals
    const darksquare = new THREE.MeshBasicMaterial( { color: 0x101010 });
    board = new THREE.Group();
    for (let i = 0; i < meshes.length; i++) {
        let geom = new THREE.BufferGeometry()
        geom.setFromPoints(meshes[i])
        geom.setAttribute('normal', new THREE.BufferAttribute(new Float32Array(normals[i]), 3));
        let cube = new THREE.Mesh(geom, darksquare)
        board.add(cube)
    }
    scene.add(board);


    window.requestAnimationFrame(animate);
}

function animate() {
    controls.update();
    if (dataFrames && dataFrames.length > 0) {
        let slowedFrame = (frame / 30).toFixed(0)
        let currentDataFrame = dataFrames[slowedFrame % dataFrames.length]
        let playerIndex = 0;
        for (let playerId in currentDataFrame.PlayerByGcLocation) {

            let player = players[playerId]
            if (!player) {
                player = new THREE.Mesh(exampleSphere, new THREE.MeshBasicMaterial( {color: playerColors[playerIndex], wireframe: true}))
                players[playerId] = player;
                scene.add(player);
            }
            player.position.x = currentDataFrame.PlayerByGcLocation[playerId].X;
            player.position.y = currentDataFrame.PlayerByGcLocation[playerId].Y + 5;
            player.position.z = currentDataFrame.PlayerByGcLocation[playerId].Z;
            if (playerIndex === 0) {
                controls.target.set(player.position.x,player.position.y,player.position.z);
                // this.camera.position.lerp(this.cameraTargetPosition, 0.4);
                // camera.position.z = player.position.z;
                // camera.position.x = player.position.x - 500;
                // camera.lookAt(player.position);
            }
            playerIndex++;
        }

        for (let monsterId in currentDataFrame.MonsterLocation) {
            let monster = monsters[monsterId]
            if (!monster) {
                monster = new THREE.Mesh(new THREE.SphereGeometry(3,8,8), new THREE.MeshBasicMaterial( {color: "white", wireframe: true}))
                monsters[monsterId] = monster
                scene.add(monster)
            }
            monster.position.x = currentDataFrame.MonsterLocation[monsterId].X ;
            monster.position.y = currentDataFrame.MonsterLocation[monsterId].Y + 3 ;
            monster.position.z = currentDataFrame.MonsterLocation[monsterId].Z ;
        }
        for (let monsterId in monsters) {
            let monster = currentDataFrame.MonsterLocation[monsterId]
            if (!monster && monsters[monsterId]) {
                scene.remove(monsters[monsterId]);
                monsters[monsterId] = null;
            }
        }
    }
    if (!paused) {
        frame++;
    }

    renderer.render(scene, camera);

    window.requestAnimationFrame(animate);

}

function onWindowResize() {

    camera.aspect = canvas.clientWidth / canvas.clientHeight;
    camera.updateProjectionMatrix();

    renderer.setSize( canvas.clientWidth, canvas.clientHeight );

}


window.addEventListener('resize', onWindowResize);

window.onload = init;