import * as THREE from '/js/three.module.js';
import { OrbitControls } from '/js/OrbitControls.js';

var scene, camera, renderer, controls, draughts, board;
let frame = 0;
let frameFraction = 0;
const playbackSpeed = 30;
var players = {};
let visibleMonsters = {};
const playerColors = ["red", "blue", "green", "yellow"];
const exampleSphere = new THREE.SphereGeometry(5,8,8)
const coneGeom = new THREE.ConeGeometry(5, 7);
const combinedGeom = new THREE.BufferGeometry();

combinedGeom.merge(exampleSphere);
// combinedGeom.merge(coneGeom);
const canvas = document.getElementById("map-canvas")
let paused = true;
let currentMap = -1;
let cameraUnset = true;
const pauseButton = document.getElementById("pause-button");
const mapRow = document.getElementById("map-row");
const mapFullscreenButton = document.getElementById("map-fullscreen-toggle");
const playbackPositionSlider = document.getElementById("playback-position");
const playbackTimer = document.getElementById("playback-timer");
const mapEnemy = document.getElementById("map-enemy");
const mapHp = document.getElementById("map-hp");
const mapTp = document.getElementById("map-tp");
const equippedWeapon = document.getElementById("map-equipped-weapon");
const floorName = document.getElementById("playback-floor-name")
const darksquare = new THREE.MeshBasicMaterial( { color: "#303030" });
let hoveredEnemy = null;
let hoveredEnemyMesh = null;
let storedMaterial = null;
const showUnitxtId = true;
const showPlayerCoordinates = true;
const showMonsterCoordinates = true;
const floorNames = {
    "0":"Pioneer II",
    "1":"Forest 1",
    "2":"Forest 2",
    "3":"Cave 1",
    "4":"Cave 2",
    "5":"Cave 3",
    "6":"Mine 1",
    "7":"Mine 2",
    "8":"Ruins 1",
    "9":"Ruins 2",
    "10":"Ruins 3",
    "11":"Under the Dome",
    "12":"Underground Channel",
    "13":"Control Room",
    "14":"????",
    "15":"Lobby",
    "16":"BA Spaceship",
    "17":"BA Temple",
    "18":"Lab",
    "19":"Temple Alpha",
    "20":"Temple Beta 2",
    "21":"Spaceship Alpha",
    "22":"Spaceship Beta",
    "23":"CCA",
    "24":"Jungle North",
    "25":"Jungle East",
    "26":"Mountain",
    "27":"Seaside",
    "28":"Seabed Upper",
    "29":"Seabed Lower",
    "30":"Cliffs of Gal Da Val",
    "31":"Test Subject Disposal Area",
    "32":"Temple Final",
    "33":"Spaceship Final",
    "34":"Seaside at Night",
    "35":"Control Tower",
    "36":"Crater East",
    "37":"Crater West",
    "38":"Crater South",
    "39":"Crater North",
    "40":"Crater Interior",
    "41":"Desert 1",
    "42":"Desert 2",
    "43":"Desert 3",
    "44":"Meteor Impact Site",
    "45":"Pioneer II"
}
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
    /*Dragon*/"44": { "geometry": new THREE.SphereGeometry(30,8,8), "material": new THREE.MeshBasicMaterial( {color: "#0002b2", wireframe: true}), "heightOffset": 3},

    "16": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#43a401", wireframe: true}), "heightOffset": 3},
    "17": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "18": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a48301", wireframe: true}), "heightOffset": 3},
    "13": { "geometry": new THREE.CylinderGeometry(2,2,12), "material": new THREE.MeshBasicMaterial( {color: "#ffd501", wireframe: true}), "heightOffset": 8},
    "19": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#01d0ff"}), "heightOffset": 3},
    "12": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#84ff01", wireframe: true}), "heightOffset": 3},
    "15": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a16dff", wireframe: true}), "heightOffset": 3},
    /* Pan Arms */ "21": { "geometry": new THREE.SphereGeometry(8,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ff996d", wireframe: true}), "heightOffset": 3},
    /* Migium */ "22": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a76dff", wireframe: true}), "heightOffset": 3},
    /* Hidoom */ "23": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ff996d", wireframe: true}), "heightOffset": 3},

    /* Gillchic */"50": {"geometry": new THREE.ConeGeometry(5,8), "material": new THREE.MeshBasicMaterial( {color: "#707070", wireframe: true}), "heightOffset": 3},
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
    /* Delsaber */ "30": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a16dff", wireframe: true}), "heightOffset": 3},
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
    /* Dragon */ "76": { "geometry": new THREE.SphereGeometry(30), "material": new THREE.MeshBasicMaterial( {color: "#159a42", wireframe: true}), "heightOffset": 3},
    /* Barba Ray */ "73": { "geometry": new THREE.SphereGeometry(8,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ff6d00", wireframe: true}), "heightOffset": 3},
    /* Pig Ray */ "74": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7679e8"}), "heightOffset": 3},
    /* Gal Gryphon */ "77": { "geometry": new THREE.SphereGeometry(30), "material": new THREE.MeshBasicMaterial( {color: "#159a42", wireframe: true}), "heightOffset": 3},

    /*Dolmolm */ "64": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    /*Dolmdarl */ "65": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    /*Recobox */ "67": { "geometry": new THREE.BoxGeometry(6,6,6), "material": new THREE.MeshBasicMaterial( {color: "#6c6c6c"}), "heightOffset": 3},
    /*Recon */ "68": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#777777"}), "heightOffset": 3},
    /*Morfos */ "66": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    /*Sinow Zoa */ "69": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    /*Sinow Zele */ "70": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    /*Delbiter */ "72": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ffffff", wireframe: true}), "heightOffset": 3},
    /*Deldepth */ "71": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#ffffff"}), "heightOffset": 3},
    /* Olga Flow */ "78": { "geometry": new THREE.SphereGeometry(30), "material": new THREE.MeshBasicMaterial( {color: "#a200ff", wireframe: true}), "heightOffset": 3},
    /* gael */ "85": { "geometry": new THREE.SphereGeometry(15), "material": new THREE.MeshBasicMaterial( {color: "#ff4d00", wireframe: true}), "heightOffset": 3},
    /* Giel */ "86": { "geometry": new THREE.SphereGeometry(15), "material": new THREE.MeshBasicMaterial( {color: "#0051ff", wireframe: true}), "heightOffset": 3},

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
        if (paused) {
            pauseButton.innerText = "play_arrow";
        } else {
            pauseButton.innerText = "pause";
        }
    }
    mapFullscreenButton.onclick = function () {
        if (mapRow.classList.contains("fullscreen-map")) {
            mapRow.classList.remove("fullscreen-map");
            mapFullscreenButton.innerText = "fullscreen";
        } else {
            mapRow.classList.add("fullscreen-map");
            mapFullscreenButton.innerText = "fullscreen_exit";
        }
        onWindowResize();
    }
    playbackPositionSlider.setAttribute("min", "0")
    playbackPositionSlider.setAttribute("max", "" + (dataFrames.length - 1))
    playbackPositionSlider.value = 0
    playbackPositionSlider.oninput = (e) => {
        paused = true;
        frame = e.target.value;
        playbackTimer.innerText = Math.floor(frame / 60) + ":" + String(frame % 60).padStart(2, '0');
    }

    scene = new THREE.Scene();
    camera = new THREE.PerspectiveCamera(75, canvas.parentNode.clientWidth / canvas.parentNode.clientHeight, 0.1, 5000);

    renderer = new THREE.WebGLRenderer({canvas: canvas});
    renderer.setSize(canvas.parentNode.clientWidth, canvas.parentNode.clientHeight);
    controls = new OrbitControls(camera, renderer.domElement);

    controls.maxPolarAngle = Math.PI / 2;

    controls.enableDamping = true;
    board = new THREE.Group();


    window.requestAnimationFrame(animate);
}

function animate() {
    controls.update();
    if (dataFrames && dataFrames.length > 0) {
        if (!paused && frameFraction % playbackSpeed === 0) {
            frame++;
            frame = frame % dataFrames.length;
            playbackPositionSlider.value = frame;
            playbackTimer.innerText = Math.floor(frame / 60) + ":" + String(frame % 60).padStart(2, '0');
            // gameTimelineChart.options.plugins.annotation.annotations.cursor.xMin = frame;
            // gameTimelineChart.options.plugins.annotation.annotations.cursor.xMax = frame;
            // gameTimelineChart.update();
        }
        let currentDataFrame = dataFrames[frame % dataFrames.length]
        equippedWeapon.innerText = weapons[currentDataFrame.Weapon];
        mapHp.innerText = currentDataFrame.HP;
        mapHp.style.width = ((Number(currentDataFrame.HP) * 100) / maxHp) + "%";
        if (maxTp > 0) {
            mapTp.innerText = currentDataFrame.TP;
            mapTp.style.width = ((Number(currentDataFrame.TP) * 100) / maxTp) + "%";
        }
        if (hoveredEnemy != null) {
            let monsterLocation = currentDataFrame.MonsterLocation[hoveredEnemy];
            if (monsterLocation.HP > 0) {
                let unitxtId = showUnitxtId ? monsters[hoveredEnemy].UnitxtId + " " : "";
                let coordinates = showMonsterCoordinates
                    ? `(${monsterLocation.X}, ${monsterLocation.Y}, ${monsterLocation.Z}) ${monsterLocation.Facing}`
                    : "";
                mapEnemy.innerText = `${monsters[hoveredEnemy].Name} ${unitxtId} ${coordinates} ${monsterLocation.HP}`;
            } else {
                mapEnemy.innerText = monsters[hoveredEnemy].Name;
            }

        } else {
            mapEnemy.innerText = "";
        }
        if (currentDataFrame.Map !== currentMap) {
            currentMap = currentDataFrame.Map;
            floorName.innerText = floorNames[currentMap];
            cameraUnset = true;
            scene.remove(board)
            board.clear();
            const meshes = meshesByFloor[currentMap].meshes
            const normals = meshesByFloor[currentMap].normals
            for (let i = 0; i < meshes.length; i++) {
                let geom = new THREE.BufferGeometry()
                geom.setFromPoints(meshes[i])
                geom.setAttribute('normal', new THREE.BufferAttribute(new Float32Array(normals[i]), 3));
                let cube = new THREE.Mesh(geom, darksquare)
                board.add(cube)
            }
            scene.add(board)
        }
        let playerIndex = 0;
        for (let playerId in currentDataFrame.PlayerByGcLocation) {
            let player = players[playerId]
            if (!player) {
                let playerGeometry = hasFacing ? coneGeom : exampleSphere;
                player = new THREE.Mesh(playerGeometry, new THREE.MeshBasicMaterial( {color: playerColors[playerIndex]}))
                players[playerId] = player;
                scene.add(player);
            }
            player.position.x = currentDataFrame.PlayerByGcLocation[playerId].X;
            player.position.y = currentDataFrame.PlayerByGcLocation[playerId].Y + 5;
            player.position.z = currentDataFrame.PlayerByGcLocation[playerId].Z;

            let playerYAngle = hasFacing ? ((currentDataFrame.PlayerByGcLocation[playerId].Facing * 6.28) / 0xFFFF) + 1.57 : 0;

            let rotationX = 0;
            let rotationY = playerYAngle;
            let rotationZ = 1.57;
            player.rotation.set(rotationX, rotationY, rotationZ);
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
            let monsterAngle = hasFacing ? ((currentDataFrame.MonsterLocation[monsterId].Facing * 6.28) / 0xFFFF) + 1.57 : 0;
            // 16383 -> 3.14
            // 0 ->
            // 49151 -> 0
            monster.rotation.set(0, monsterAngle, 1.57);
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

    camera.aspect = canvas.parentNode.clientWidth / canvas.parentNode.clientHeight;
    camera.updateProjectionMatrix();

    renderer.setSize( canvas.parentNode.clientWidth, canvas.parentNode.clientHeight );

}

function onDocumentMouseMove(event) {
    let mouse = new THREE.Vector2();
    const parentNode = canvas.parentNode;
    mouse.x = ((event.clientX - parentNode.offsetLeft) / parentNode.clientWidth) * 2 - 1;
    mouse.y = - ((event.clientY - parentNode.offsetTop) / parentNode.clientHeight) * 2 + 1;
    let raycaster = new THREE.Raycaster();
    raycaster.setFromCamera( mouse, camera );
    let intersects = raycaster.intersectObject( scene , true);

    let currentEnemy = null;
    if(intersects.length > 0) {
        for (let id in visibleMonsters) {
            if (intersects[0].object === visibleMonsters[id]) {
                currentEnemy = id;
            }
        }
    }
    if (hoveredEnemy !== currentEnemy) {
        if (hoveredEnemyMesh != null) {
            hoveredEnemyMesh.material = storedMaterial
        }
        hoveredEnemy = currentEnemy;
        if (hoveredEnemy != null) {
            hoveredEnemyMesh = intersects[0].object;
            storedMaterial = hoveredEnemyMesh.material
            hoveredEnemyMesh.material = new THREE.MeshBasicMaterial({color: "0xff9900"});
        } else {
            hoveredEnemyMesh = null;
        }
    }
}


window.addEventListener('resize', onWindowResize);
document.addEventListener( 'mousemove', onDocumentMouseMove, false );

window.onload = init;