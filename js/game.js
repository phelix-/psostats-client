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
const pauseButton = document.getElementById("pause-button")
const playbackPositionSlider = document.getElementById("playback-position")
const monsterMeshes = {
    "El Rappy": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#e8eca6", wireframe: true}), "heightOffset": 3},
    "Bartle": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#755138", wireframe: true}), "heightOffset": 3},
    "Barble": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#668149", wireframe: true}), "heightOffset": 3},
    "Tollaw": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#01a417", wireframe: true}), "heightOffset": 3},
    "Gulgus": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#3f5f72", wireframe: true}), "heightOffset": 3},
    "Gulgus-gue": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#853530", wireframe: true}), "heightOffset": 3},
    "Mothvert": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#ff0000"}), "heightOffset": 3},
    "Mothvist": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#003f9d", wireframe: true}), "heightOffset": 3},
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

    "Arlan": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#43a401", wireframe: true}), "heightOffset": 3},
    "Merlan": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "Del-D": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a48301", wireframe: true}), "heightOffset": 3},
    "Claw": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#886060"}), "heightOffset": 3},
    "Bulclaw": {"geometry": new THREE.CylinderGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#886060"}), "heightOffset": 3},
    "Delsaber": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a16dff", wireframe: true}), "heightOffset": 3},
    "Gran Sorcerer": { "geometry": new THREE.SphereGeometry(4,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "Gee R": { "geometry": new THREE.SphereGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "Gee L": { "geometry": new THREE.SphereGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "Indi Belra": { "geometry": new THREE.CylinderGeometry(5,5,8), "material": new THREE.MeshBasicMaterial( {color: "#00ffe0", wireframe: true}), "heightOffset": 3},
    "Dark Bringer": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#d59d1c", wireframe: true}), "heightOffset": 3},

    "Merillia": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a41f01", wireframe: true}), "heightOffset": 3},
    "Meriltas": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#6ea401", wireframe: true}), "heightOffset": 3},
    "Ul Gibbon": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#566015", wireframe: true}), "heightOffset": 3},
    "Zol Gibbon": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#540d0d", wireframe: true}), "heightOffset": 3},
    "Gee": { "geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#ff4f00"}), "heightOffset": 3},
    "Gibbles": { "geometry": new THREE.CylinderGeometry(12,12,24), "material": new THREE.MeshBasicMaterial( {color: "#e3d806", wireframe: true}), "heightOffset": 3},
    "Gi Gue": { "geometry": new THREE.CylinderGeometry(12,12,24), "material": new THREE.MeshBasicMaterial( {color: "#e34e0f", wireframe: true}), "heightOffset": 3},
    "Mericarol": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ff4400", wireframe: true}), "heightOffset": 3},
    "Merikle": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#00ceff", wireframe: true}), "heightOffset": 3},
    "Mericus": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#e6ff00", wireframe: true}), "heightOffset": 3},
    "Sinow Spigell": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#987171", wireframe: true}), "heightOffset": 3},
    "Sinow Berill": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#159a42", wireframe: true}), "heightOffset": 3},

    "Dolmolm": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "Dolmdarl": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "Recobox": { "geometry": new THREE.BoxGeometry(5,5,5), "material": new THREE.MeshBasicMaterial( {color: "#6c6c6c"}), "heightOffset": 3},
    "Recon": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#777777"}), "heightOffset": 3},
    "Morfos": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "Sinow Zoa": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "Sinow Zele": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "Delbiter": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "Deldepth": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#ffffff"}), "heightOffset": 3},

    "Ill Gill": { "geometry": new THREE.SphereGeometry(6,8, 8), "material": new THREE.MeshBasicMaterial( {color: "#d07eff"}), "heightOffset": 3},
    "Del Lily": { "geometry": new THREE.CylinderGeometry(4,4,20), "material": new THREE.MeshBasicMaterial( {color: "#45147a"}), "heightOffset": 3},
    "Epsilon": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#862f00"}), "heightOffset": 3},
    "Epsigard": { "geometry": new THREE.SphereGeometry(4,8, 8), "material": new THREE.MeshBasicMaterial( {color: "#862f00"}), "heightOffset": 3},

    "Boota": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "Ze Boota": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "Ba Boota": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    "Astark": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#00c70a", wireframe: true}), "heightOffset": 3},
    "Dorphon": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},

    "Sand Rappy": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ff8000", wireframe: true}), "heightOffset": 3},
    "Goran": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "Pyro Goran": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#01a417", wireframe: true}), "heightOffset": 3},
    "Goran Detonator": { "geometry": new THREE.CylinderGeometry(12,12,24), "material": new THREE.MeshBasicMaterial( {color: "#c900b5", wireframe: true}), "heightOffset": 3},
    "Satellite Lizard": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#49811f", wireframe: true}), "heightOffset": 3},
    "Yowie": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#b0512a", wireframe: true}), "heightOffset": 3},
    "Merissa A": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#ff7e95"}), "heightOffset": 3},
    "Zu": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#fa53b3", wireframe: true}), "heightOffset": 3},
    "Girtablulu": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ab5bff", wireframe: true}), "heightOffset": 3},
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
            playbackPositionSlider.value = frame;
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