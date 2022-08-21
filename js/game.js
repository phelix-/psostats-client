import * as THREE from '/js/three.module.js';
import { OrbitControls } from '/js/OrbitControls.js';

var scene, camera, renderer, controls, draughts, board;
let frame = 0;
let frameFraction = 0;
const playbackSpeed = 30;
var dataFrames;
var players = {};
let visibleMonsters = {};
const playerColors = ["red", "blue", "green", "yellow"];
const exampleSphere = new THREE.SphereGeometry(5,8,8)
const canvas = document.getElementById("map-canvas")
let paused = true;
let cameraUnset = true;
const pauseButton = document.getElementById("pause-button");
const playbackPositionSlider = document.getElementById("playback-position");
const playbackTimer = document.getElementById("playback-timer");
const monsterMeshes = {
    "5": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#e8eca6", wireframe: true}), "heightOffset": 3},
    "9": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#755138", wireframe: true}), "heightOffset": 3},
    "10": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#668149", wireframe: true}), "heightOffset": 3},
    "11": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#01a417", wireframe: true}), "heightOffset": 3},
    "7": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#3f5f72", wireframe: true}), "heightOffset": 3},
    "8": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#853530", wireframe: true}), "heightOffset": 3},
    "3": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#ff0000"}), "heightOffset": 3},
    "4": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#003f9d", wireframe: true}), "heightOffset": 3},
    "1": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#818181", wireframe: true}), "heightOffset": 3},

    "16": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#43a401", wireframe: true}), "heightOffset": 3},
    "17": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "18": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a48301", wireframe: true}), "heightOffset": 3},
    "13": { "geometry": new THREE.CylinderGeometry(2,2,12), "material": new THREE.MeshBasicMaterial( {color: "#ffd501", wireframe: true}), "heightOffset": 8},
    "19": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#01d0ff"}), "heightOffset": 3},
    "12": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#84ff01", wireframe: true}), "heightOffset": 3},
    "15": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a16dff", wireframe: true}), "heightOffset": 3},

    "50": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#707070", wireframe: true}), "heightOffset": 3},
    "24": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#702a00", wireframe: true}), "heightOffset": 3},
    "28": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#6b5353"}), "heightOffset": 3},
    "29": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#b96a07"}), "heightOffset": 3},
    "26": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#326dff", wireframe: true}), "heightOffset": 3},
    "27": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#c00000", wireframe: true}), "heightOffset": 3},
    "25": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#0002b2", wireframe: true}), "heightOffset": 3},

    "41": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#43a401", wireframe: true}), "heightOffset": 3},
    "42": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "43": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a48301", wireframe: true}), "heightOffset": 3},
    "38": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#886060"}), "heightOffset": 3},
    "40": {"geometry": new THREE.CylinderGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#886060"}), "heightOffset": 3},
    "Delsaber": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a16dff", wireframe: true}), "heightOffset": 3},
    "31": { "geometry": new THREE.SphereGeometry(4,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "32": { "geometry": new THREE.SphereGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "33": { "geometry": new THREE.SphereGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "37": { "geometry": new THREE.CylinderGeometry(5,5,8), "material": new THREE.MeshBasicMaterial( {color: "#00ffe0", wireframe: true}), "heightOffset": 3},
    "36": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#d59d1c", wireframe: true}), "heightOffset": 3},

    "52": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a41f01", wireframe: true}), "heightOffset": 3},
    "53": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#6ea401", wireframe: true}), "heightOffset": 3},
    "59": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#566015", wireframe: true}), "heightOffset": 3},
    "60": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#540d0d", wireframe: true}), "heightOffset": 3},
    "54": { "geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#ff4f00"}), "heightOffset": 3},
    "61": { "geometry": new THREE.CylinderGeometry(12,12,24), "material": new THREE.MeshBasicMaterial( {color: "#e3d806", wireframe: true}), "heightOffset": 3},
    "55": { "geometry": new THREE.CylinderGeometry(12,12,24), "material": new THREE.MeshBasicMaterial( {color: "#e34e0f", wireframe: true}), "heightOffset": 3},
    "56": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ff4400", wireframe: true}), "heightOffset": 3},
    "57": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#00ceff", wireframe: true}), "heightOffset": 3},
    "58": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#e6ff00", wireframe: true}), "heightOffset": 3},
    "63": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#987171", wireframe: true}), "heightOffset": 3},
    "62": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#159a42", wireframe: true}), "heightOffset": 3},

    "64": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3}, // Dolmolm
    "65": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3}, // Dolmdarl
    "67": { "geometry": new THREE.BoxGeometry(6,6,6), "material": new THREE.MeshBasicMaterial( {color: "#6c6c6c"}), "heightOffset": 3}, // Recobox
    "68": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#777777"}), "heightOffset": 3}, // Recon
    "66": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3}, // Morfos
    "69": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3}, // Sinow Zoa
    "70": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3}, // Sinow Zele
    "72": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3}, // Delbiter
    "71": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#ffffff"}), "heightOffset": 3}, // Deldepth

    "82": { "geometry": new THREE.SphereGeometry(6,8, 8), "material": new THREE.MeshBasicMaterial( {color: "#d07eff"}), "heightOffset": 3},
    "83": { "geometry": new THREE.CylinderGeometry(4,4,20), "material": new THREE.MeshBasicMaterial( {color: "#45147a"}), "heightOffset": 3},
    "84": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#862f00"}), "heightOffset": 3},
    "87": { "geometry": new THREE.SphereGeometry(4,8, 8), "material": new THREE.MeshBasicMaterial( {color: "#862f00"}), "heightOffset": 3},

    "96": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "97": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "98": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "88": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#00c70a", wireframe: true}), "heightOffset": 3},
    "99": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},

    "104": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ff8000", wireframe: true}), "heightOffset": 3},
    "101": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "103": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#01a417", wireframe: true}), "heightOffset": 3},
    "102": { "geometry": new THREE.CylinderGeometry(12,12,24), "material": new THREE.MeshBasicMaterial( {color: "#c900b5", wireframe: true}), "heightOffset": 3},
    "90": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#49811f", wireframe: true}), "heightOffset": 3},
    "89": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#b0512a", wireframe: true}), "heightOffset": 3},
    "91": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#ff7e95"}), "heightOffset": 3},
    "94": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#fa53b3", wireframe: true}), "heightOffset": 3},
    "93": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ab5bff", wireframe: true}), "heightOffset": 3},
};

function init() {
    draughts = new Draughts();
    pauseButton.onclick = function() {
        paused = !paused;
    }
    playbackPositionSlider.setAttribute("min", "0")
    playbackPositionSlider.setAttribute("max", "" + (jsonMeshes.dataFrames.length - 1))
    playbackPositionSlider.value = 0
    playbackPositionSlider.oninput = (e) => {
        paused = true;
        frame = e.target.value;
        playbackTimer.innerText = (frame / 60).toFixed(0) + ":" + String(frame % 60).padStart(2, '0');
    }

    scene = new THREE.Scene();
    camera = new THREE.PerspectiveCamera(75, canvas.clientWidth / canvas.clientHeight, 0.1, 5000);

    renderer = new THREE.WebGLRenderer({canvas: canvas});
    renderer.setSize(canvas.clientWidth, canvas.clientHeight);
    controls = new OrbitControls(camera, renderer.domElement);

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
        if (!paused && frameFraction % playbackSpeed === 0) {
            frame++;
            frame = frame % dataFrames.length;
            playbackPositionSlider.value = frame;
            playbackTimer.innerText = (frame / 60).toFixed(0) + ":" + String(frame % 60).padStart(2, '0');
            // gameTimelineChart.options.plugins.annotation.annotations.cursor.xMin = frame;
            // gameTimelineChart.options.plugins.annotation.annotations.cursor.xMax = frame;
            // gameTimelineChart.update();
        }
        let currentDataFrame = dataFrames[frame % dataFrames.length]
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
                let monsterMeshInfo = monsterMeshes[monsterInfo.UnitxtId]
                if (monsterMeshInfo) {
                    monster = new THREE.Mesh(monsterMeshInfo.geometry, monsterMeshInfo.material);
                } else {
                    console.warn("Missing mesh", monsterInfo.UnitxtId, monsterInfo.Name)
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
        frameFraction++;
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