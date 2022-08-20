import * as THREE from '/js/three.module.js';
import { OrbitControls } from '/js/OrbitControls.js';
import {Vector3} from "/js/three.module.js";

var scene, camera, renderer, controls, draughts, board;
var frame = 0;
var dataFrames;
var players = {};
let visibleMonsters = {};
const playerColors = ["red", "blue", "green", "yellow"];
const exampleSphere = new THREE.SphereGeometry(5,8,8)
const canvas = document.getElementById("map-canvas")
let paused = true;
let cameraUnset = true;
const pauseButton = document.getElementById("pause-button")
const monsterMeshes = {
    "El Rappy": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#e8eca6", wireframe: true}), "heightOffset": 3},
    "Bartle": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#755138", wireframe: true}), "heightOffset": 3},
    "Barble": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#668149", wireframe: true}), "heightOffset": 3},
    "Tollaw": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#01a417", wireframe: true}), "heightOffset": 3},
    "Gulgus": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#3f5f72", wireframe: true}), "heightOffset": 3},
    "Gulgus-gue": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#853530", wireframe: true}), "heightOffset": 3},
    "Mothvert": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#ff0000"}), "heightOffset": 3},
    "Mothvist": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#9d007d", wireframe: true}), "heightOffset": 3},
    "Hildelt": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#818181", wireframe: true}), "heightOffset": 3},

    "Vulmer": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#43a401", wireframe: true}), "heightOffset": 3},
    "Govulmer": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "Melqueek": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a48301", wireframe: true}), "heightOffset": 3},
    "Ob Lily": { "geometry": new THREE.CylinderGeometry(2,2,12), "material": new THREE.MeshBasicMaterial( {color: "#ffd501", wireframe: true}), "heightOffset": 8},
    "Pofuilly Slime": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#01d0ff"}), "heightOffset": 3},
    "Crimson Assassin": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#84ff01", wireframe: true}), "heightOffset": 3},
    "Nano Dragon": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a16dff", wireframe: true}), "heightOffset": 3},
    "Gillchich": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#707070", wireframe: true}), "heightOffset": 3},
    "Dubchich": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#702a00", wireframe: true}), "heightOffset": 3},
    "Canabin": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#6b5353"}), "heightOffset": 3},
    "Canune": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#b96a07"}), "heightOffset": 3},
    "Sinow Blue": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#326dff", wireframe: true}), "heightOffset": 3},
    "Sinow Red": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#c00000", wireframe: true}), "heightOffset": 3},
    "Baranz": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#0002b2", wireframe: true}), "heightOffset": 3},
};

const rooms = {
    101: [1175.000000,0.000000,-2545.000000],
    102: [1460.000000,0.000000,-2770.000000],
    110: [1245.000000,0.000000,-2915.000000],
    111: [350.000000,0.000000,-1850.000000],
    112: [645.000000,0.000000,-2465.000000],
    201: [1460.000000,0.000000,-2915.000000],
    202: [430.000000,0.000000,-2465.000000],
    203: [270.000000,0.000000,-2465.000000],
    1: [-300.000000,0.000000,-50.000000],
    2: [600.000000,0.000000,-950.000000],
    3: [-200.000000,0.000000,-2400.000000],
    4: [2060.000000,0.000000,-1995.000000],
    5: [395.000000,0.000000,-2665.000000],
    6: [55.000000,0.000000,-2465.000000],
    7: [-500.000000,0.000000,-950.000000],
    8: [-250.000000,0.000000,-1150.000000],
    30: [1410.000000,0.000000,-2545.000000],
    35: [945.000000,0.000000,-2865.000000],
    40: [945.000000,0.000000,-2465.000000],
    45: [350.000000,0.000000,-2150.000000],
    20: [0.000000,0.000000,-300.000000],
    21: [350.000000,0.000000,-650.000000],
    22: [-500.000000,0.000000,-650.000000],
    23: [1460.000000,0.000000,-1995.000000],
    10: [0.000000,0.000000,-650.000000],
    11: [350.000000,0.000000,-950.000000],
    12: [395.000000,0.000000,-2915.000000],
    13: [-200.000000,0.000000,-2150.000000],
    14: [1810.000000,0.000000,-1995.000000],
    15: [-300.000000,0.000000,-300.000000],
    16: [0.000000,0.000000,-1150.000000],
    120: [-250.000000,0.000000,-650.000000],
    121: [645.000000,0.000000,-2915.000000],
    122: [350.000000,0.000000,-1200.000000],
    123: [1460.000000,0.000000,-2245.000000],
    124: [50.000000,0.000000,-2150.000000],
    125: [0.000000,0.000000,-900.000000],
    60: [350.000000,0.000000,-1525.000000],
}

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
    controls = new OrbitControls(camera, renderer.domElement);

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
                player = new THREE.Mesh(exampleSphere, new THREE.MeshBasicMaterial( {color: playerColors[playerIndex]}))
                players[playerId] = player;
                scene.add(player);
            }
            player.position.x = currentDataFrame.PlayerByGcLocation[playerId].X;
            player.position.y = currentDataFrame.PlayerByGcLocation[playerId].Y + 5;
            player.position.z = currentDataFrame.PlayerByGcLocation[playerId].Z;
            if (playerIndex === 0) {
                let room = currentDataFrame.PlayerByGcLocation[playerId].Room
                // if (rooms[room]) {
                //     controls.target.set(rooms[room][0], rooms[room][1], rooms[room][2])
                // } else {
                    if (cameraUnset ||
                        Math.abs(camera.position.x - player.position.x) > 1000 ||
                        Math.abs(camera.position.y - player.position.y) > 1000 ||
                        Math.abs(camera.position.x - player.position.x) > 1000
                    ) {
                        camera.position.x = player.position.x;
                        camera.position.y = player.position.y + 300;
                        camera.position.z = player.position.z - 300;
                        controls.target.set(player.position.x,player.position.y,player.position.z);
                        cameraUnset = false;
                    }

                    // camera.position.set(player.position.x, player.position.y + 300, player.position.z - 300)
                    // controls.target.set(player.position.x,player.position.y,player.position.z);
                // }

                // this.camera.position.lerp(this.cameraTargetPosition, 0.4);
                // camera.position.z = player.position.z;
                // camera.position.x = player.position.x - 500;
                // camera.lookAt(player.position);
            }
            playerIndex++;
        }

        for (let monsterId in currentDataFrame.MonsterLocation) {
            let monster = visibleMonsters[monsterId]
            if (!monster) {
                let monsterInfo = monsters["" + monsterId]
                if (!monsterInfo) {
                    console.warn("Missing", monsterId)
                }
                let monsterMeshInfo = monsterMeshes[monsterInfo.Name]
                if (monsterMeshInfo) {
                    monster = new THREE.Mesh(monsterMeshInfo.geometry, monsterMeshInfo.material);
                } else {
                    console.warn("Missing mesh", monsterInfo.Name)
                    monster = new THREE.Mesh(new THREE.SphereGeometry(3,8,8), new THREE.MeshBasicMaterial( {color: "white", wireframe: true}))
                }
                visibleMonsters[monsterId] = monster
                scene.add(monster)
            }
            monster.position.x = currentDataFrame.MonsterLocation[monsterId].X ;
            monster.position.y = currentDataFrame.MonsterLocation[monsterId].Y + 3 ;
            monster.position.z = currentDataFrame.MonsterLocation[monsterId].Z ;
        }
        for (let monsterId in visibleMonsters) {
            let monster = currentDataFrame.MonsterLocation[monsterId]
            if (!monster && visibleMonsters[monsterId]) {
                scene.remove(visibleMonsters[monsterId]);
                visibleMonsters[monsterId] = null;
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