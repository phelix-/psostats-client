import * as THREE from '/js/three.module.js';
import { OrbitControls } from '/js/OrbitControls.js';
import { getGeometry } from  '/js/monster_geometry.js';

let scene, camera, renderer, controls, floorGroup;
let frame = 0;
let frameFraction = 0;
const playbackSpeed = 30;
let players = {};
let visibleMonsters = {};
const playerColors = ["red", "blue", "green", "yellow"];
const exampleSphere = new THREE.SphereGeometry(5,8,8)
const coneGeom = new THREE.ConeGeometry(5, 7);
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
const MAP_HP_BAR = document.getElementById("map-hp");
const MAP_TP_BAR = document.getElementById("map-tp");
const MAP_EQUIPPED_WEAPON = document.getElementById("map-equipped-weapon");
const MAP_FLOOR_NAME = document.getElementById("playback-floor-name");
const floorMaterial = new THREE.MeshLambertMaterial( { color: "#303030" });
let hoveredPlayerId = null;
let hoveredPlayerMesh = null;
let storedPlayerMaterial = null;
let hoveredEnemyId = null;
let hoveredEnemyMesh = null;
let storedEnemyMaterial = null;
// const config = getConfig();
const showUnitxtId = Boolean(window.localStorage.getItem("showUnitxtId"));
const showFacing = Boolean(window.localStorage.getItem("showFacing"));
const showPlayerCoordinates = Boolean(window.localStorage.getItem("showPlayerCoordinates"));
const showMonsterCoordinates = Boolean(window.localStorage.getItem("showMonsterCoordinates"));
let followingGc = null;
let priorPlayerLocation = null;
let cameraPositionDiff = null;

const iceColor = '#4fb1db'
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

function getConfig() {
    let gameConfig = window.localStorage.getItem("gameConfig");
    if (!gameConfig) {
        gameConfig = {
            "showUnitxtId": false,
            "showFacing": false,
            "showPlayerCoordinates": true,
            "showMonsterCoordinates": true,
        };

        window.localStorage.setItem("gameConfig", gameConfig);
    }
    return gameConfig;
}

function toggleFullscreenMap() {
    if (mapRow.classList.contains("fullscreen-map")) {
        mapRow.classList.remove("fullscreen-map");
        mapFullscreenButton.innerText = "fullscreen";
    } else {
        mapRow.classList.add("fullscreen-map");
        mapFullscreenButton.innerText = "fullscreen_exit";
    }
    onWindowResize();
}

function togglePlaybackPause() {
    paused = !paused;
    if (paused) {
        pauseButton.innerText = "play_arrow";
    } else {
        pauseButton.innerText = "pause";
    }
}

function watchPlayer(playerIndex) {
    priorPlayerLocation = null;
    cameraPositionDiff = null;
    if (playerIndex === -1 || gcsByGem.length <= playerIndex) {
        followingGc = null;
    } else {
        followingGc = gcsByGem[playerIndex];
        cameraUnset = true;
    }
}

function init() {
    followingGc = viewingGc;
    pauseButton.onclick = togglePlaybackPause;
    mapFullscreenButton.onclick = toggleFullscreenMap;
    playbackPositionSlider.setAttribute("min", "0")
    playbackPositionSlider.setAttribute("max", "" + (dataFrames.length - 1))
    playbackPositionSlider.value = 0
    playbackPositionSlider.oninput = (e) => {
        if (!paused) {
            togglePlaybackPause()
        }
        frame = e.target.value;
        playbackTimer.innerText = Math.floor(frame / 60) + ":" + String(frame % 60).padStart(2, '0');
    }

    scene = new THREE.Scene();
    camera = new THREE.PerspectiveCamera(75, canvas.parentNode.clientWidth / canvas.parentNode.clientHeight, 0.1, 5000);
    addAmbientLight(scene);
    addDirectLight(scene);

    renderer = new THREE.WebGLRenderer({canvas: canvas});
    renderer.setSize(canvas.parentNode.clientWidth, canvas.parentNode.clientHeight);
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFSoftShadowMap;
    controls = new OrbitControls(camera, renderer.domElement);

    controls.maxPolarAngle = Math.PI / 2;

    controls.enableDamping = true;
    floorGroup = new THREE.Group();

    window.requestAnimationFrame(animate);
}

function addAmbientLight(scene) {
    let light = new THREE.DirectionalLight( "#ffffff", 0.55);
    light.castShadow = false;
    scene.add(light);
}

function addDirectLight(scene) {
    let light = new THREE.DirectionalLight("#ffffff", 1);
    let shadowmapSize = 4096;
    light.castShadow = true;
    light.shadow.normalBias = 0.01;

    light.shadow.mapSize.width = shadowmapSize;
    light.shadow.mapSize.height = shadowmapSize;
    light.shadow.camera.near = -750;
    light.shadow.camera.far = 280;

    light.shadow.camera.left = -1000;
    light.shadow.camera.bottom = -1000;
    light.shadow.camera.right = 1000;
    light.shadow.camera.top = 1000;

    light.shadow.blurSamples = 12;
    scene.add(light);
}

function animate() {
    controls.update();
    if (dataFrames && dataFrames.length > 0) {
        if (!paused && frameFraction % playbackSpeed === 0) {
            frame++;
            frame = frame % dataFrames.length;
            playbackPositionSlider.value = frame;
            playbackTimer.innerText = Math.floor(frame / 60) + ":" + String(frame % 60).padStart(2, '0');
        }
        let currentDataFrame = dataFrames[frame % dataFrames.length];
        let primaryPlayerState = currentDataFrame.PlayerByGcLocation[viewingGc];
        updatePlayerInfo(currentDataFrame);

        if (hoveredEnemyId != null) {
            let monsterLocation = currentDataFrame.MonsterLocation[hoveredEnemyId];
            if (monsterLocation && monsterLocation.HP > 0) {
                let unitxtId = showUnitxtId ? monsters[hoveredEnemyId].UnitxtId + " " : "";
                let facing = showFacing ? " " + monsterLocation.Facing : "";
                let coordinates = showMonsterCoordinates
                    ? `(${monsterLocation.X}, ${monsterLocation.Y}, ${monsterLocation.Z})${facing}`
                    : "";
                mapEnemy.innerText = `${monsters[hoveredEnemyId].Name} ${unitxtId} ${coordinates} ${monsterLocation.HP}`;
            } else {
                mapEnemy.innerText = monsters[hoveredEnemyId].Name;
            }
        } else {
            mapEnemy.innerText = "";
        }
        if (currentDataFrame.Map !== currentMap) {
            currentMap = currentDataFrame.Map;
            MAP_FLOOR_NAME.innerText = floorNames[currentMap];
            cameraUnset = true;
            scene.remove(floorGroup)
            floorGroup.clear();
            const meshes = meshesByFloor[currentMap].meshes
            for (let i = 0; i < meshes.length; i++) {
                let geom = new THREE.BufferGeometry()
                geom.setFromPoints(meshes[i]);
                geom.computeVertexNormals();
                geom.normalizeNormals();
                const floorMesh = new THREE.Mesh(geom, floorMaterial);
                floorMesh.receiveShadow = true;
                floorMesh.castShadow = true;
                floorGroup.add(floorMesh);
            }
            floorGroup.receiveShadow = true;
            floorGroup.castShadow = true;
            scene.add(floorGroup)
        }
        let playerIndex = -1;
        for (let playerId of gcsByGem) {
            playerIndex++;
            let player = players[playerId];
            let playerState = currentDataFrame.PlayerByGcLocation[playerId];
            if (player && (playerState.Floor !== primaryPlayerState.Floor || playerState.Warping)) {
                scene.remove(player);
                players[playerId] = null;
            }
            if (playerState.Floor !== primaryPlayerState.Floor || playerState.Warping) {
                continue;
            }
            if (!player) {
                let playerGeometry = hasFacing ? coneGeom : exampleSphere;
                player = new THREE.Mesh(playerGeometry, new THREE.MeshToonMaterial( {color: playerColors[playerIndex], /*emissive: new THREE.Color(0xAA00FF), emissiveIntensity: 100*/}))
                player["playerId"] = playerId;
                players[playerId] = player;
                player.castShadow = true
                scene.add(player);
            }
            player.position.x = playerState.X;
            player.position.y = playerState.Y + 5;
            player.position.z = playerState.Z;

            let playerYAngle = hasFacing ? ((playerState.Facing * (Math.PI * 2)) / 0xFFFF) + (Math.PI / 2) : 0;
            player.rotation.set(0, playerYAngle, Math.PI / 2);
        }
        if (cameraUnset) {
            let locationToUse = followingGc ? currentDataFrame.PlayerByGcLocation[followingGc] : primaryPlayerState;
            camera.position.x = locationToUse.X;
            camera.position.y = locationToUse.Y + 100;
            camera.position.z = locationToUse.Z - 100;
            controls.target.set(locationToUse.X, locationToUse.Y, locationToUse.Z);
            cameraUnset = false;
            priorPlayerLocation = null;
            cameraPositionDiff = null;
        }

        if (followingGc) {
            let followedPlayerLocation = currentDataFrame.PlayerByGcLocation[followingGc];
            if (priorPlayerLocation && frameFraction % playbackSpeed === 0) {
                cameraPositionDiff = new THREE.Vector3(
                (followedPlayerLocation.X - priorPlayerLocation.X) / playbackSpeed,
                (followedPlayerLocation.Y - priorPlayerLocation.Y) / playbackSpeed,
                (followedPlayerLocation.Z - priorPlayerLocation.Z) / playbackSpeed
                );
            }
            priorPlayerLocation = followedPlayerLocation;
            if (cameraPositionDiff && !paused) {
                camera.position.x += cameraPositionDiff.x;
                camera.position.y += cameraPositionDiff.y;
                camera.position.z += cameraPositionDiff.z;
                controls.target.x += cameraPositionDiff.x;
                controls.target.y += cameraPositionDiff.y;
                controls.target.z += cameraPositionDiff.z;
            }
        }

        for (let monsterId in currentDataFrame.MonsterLocation) {
            let monster = visibleMonsters[monsterId]
            if (!monster) {
                let monsterInfo = monsters["" + monsterId]
                if (!monsterInfo) {
                    console.warn("Missing", monsterId)
                }
                let monsterMesh = getGeometry(monsterInfo.UnitxtId);
                monsterMesh["monsterId"] = monsterId;
                if (monsterMesh["children"] && monsterMesh.children.length > 0) {
                    // for groups
                    monsterMesh.children[0]["monsterId"] = monsterId;
                }
                monsterMesh.castShadow = true;
                monster = new THREE.Group();
                monster.add(monsterMesh);

                visibleMonsters[monsterId] = monster;
                monster.castShadow = true;
                scene.add(monster)
            }
            let monsterLocationElement = currentDataFrame.MonsterLocation[monsterId];
            monster.position.x = monsterLocationElement.X ;
            monster.position.y = monsterLocationElement.Y + 3 ;
            monster.position.z = monsterLocationElement.Z ;
            let monsterAngle = hasFacing ? ((monsterLocationElement.Facing * (Math.PI * 2)) / 0xFFFF) + (Math.PI / 2) : 0;

            monster.rotation.set(0, monsterAngle, 0);
            if (monsterLocationElement.Frozen && monster["iceMesh"] === undefined) {
                const iceMesh = new THREE.Mesh(new THREE.OctahedronGeometry(10), new THREE.MeshBasicMaterial({color: iceColor, opacity: 0.3, transparent: true}));
                monster.add(iceMesh);
                monster["iceMesh"] = iceMesh;
            } else if (!monsterLocationElement.Frozen && monster["iceMesh"] !== undefined) {
                monster.remove(monster["iceMesh"]);
                monster["iceMesh"] = undefined;
            }
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

function updatePlayerInfo(currentDataFrame) {
    let gcToDisplay = hoveredPlayerId !== null ? hoveredPlayerId : viewingGc;
    let playerLocation = currentDataFrame.PlayerByGcLocation[gcToDisplay];
    let facing = showFacing ? " " + playerLocation.Facing : "";
    let coordinates = showPlayerCoordinates ? `(${playerLocation.X}, ${playerLocation.Y}, ${playerLocation.Z})${facing}` : "";
    if (hoveredPlayerId !== null && hoveredPlayerId !== viewingGc) {
        let player = charactersByGc[hoveredPlayerId];
        MAP_EQUIPPED_WEAPON.innerText = `${player.name} (${player.class})${coordinates}`;
        MAP_HP_BAR.innerText = "";
        MAP_HP_BAR.style.width = "0";
        MAP_TP_BAR.innerText = "";
        MAP_TP_BAR.style.width = "0";
    } else {
        let player = charactersByGc[viewingGc];
        MAP_EQUIPPED_WEAPON.innerText = `${player.name} (${player.class})${coordinates}\n
        ${weapons[currentDataFrame.Weapon]}`;
        MAP_HP_BAR.innerText = currentDataFrame.HP;
        MAP_HP_BAR.style.width = ((Number(currentDataFrame.HP) * 100) / maxHp) + "%";
        if (maxTp > 0) {
            MAP_TP_BAR.innerText = currentDataFrame.TP;
            MAP_TP_BAR.style.width = ((Number(currentDataFrame.TP) * 100) / maxTp) + "%";
        }
    }
}

function onWindowResize() {
    camera.aspect = canvas.parentNode.clientWidth / canvas.parentNode.clientHeight;
    camera.updateProjectionMatrix();
    renderer.setSize( canvas.parentNode.clientWidth, canvas.parentNode.clientHeight );
}

function onDocumentMouseMove(event) {
    if (cameraUnset) {
        return;
    }
    let mouse = new THREE.Vector2();
    const parentNode = canvas.parentNode;
    mouse.x = ((event.pageX - parentNode.offsetLeft) / parentNode.clientWidth) * 2 - 1;
    mouse.y = - ((event.pageY - parentNode.offsetTop) / parentNode.clientHeight) * 2 + 1;
    let raycaster = new THREE.Raycaster();
    raycaster.setFromCamera( mouse, camera );
    let intersects = raycaster.intersectObject( scene , true);

    let intersectedEnemy = null
    let intersectedEnemyId = null;
    let intersectedPlayer = null;
    let intersectedPlayerId = null;
    for (let intersectedMesh of intersects) {
        if (intersectedMesh.object["playerId"] in players) {
            intersectedPlayer = intersectedMesh.object;
            intersectedPlayerId = intersectedMesh.object["playerId"];
        }
        if (intersectedMesh.object["monsterId"] in visibleMonsters) {
            intersectedEnemy = intersectedMesh.object;
            intersectedEnemyId = intersectedMesh.object["monsterId"]
        }
    }

    if (hoveredPlayerId !== intersectedPlayerId) {
        if (hoveredPlayerMesh != null) {
            hoveredPlayerMesh.material = storedPlayerMaterial;
        }
        hoveredPlayerId = intersectedPlayerId;
        if (hoveredPlayerId != null) {
            hoveredPlayerMesh = intersectedPlayer;
            storedPlayerMaterial = hoveredPlayerMesh.material;
            hoveredPlayerMesh.material = new THREE.MeshBasicMaterial({color: "#ffffff"});
        } else {
            hoveredPlayerMesh = null;
        }
    }

    if (hoveredEnemyId !== intersectedEnemyId) {
        if (hoveredEnemyMesh != null) {
            hoveredEnemyMesh.material = storedEnemyMaterial;
        }
        hoveredEnemyId = intersectedEnemyId;
        if (hoveredEnemyId != null) {
            hoveredEnemyMesh = intersectedEnemy;
            storedEnemyMaterial = hoveredEnemyMesh.material
            hoveredEnemyMesh.material = new THREE.MeshBasicMaterial({color: "#ff9900"});
        } else {
            hoveredEnemyMesh = null;
        }
    }
}

window.addEventListener('resize', onWindowResize);
document.addEventListener( 'mousemove', onDocumentMouseMove, false );
onkeydown = (event) => {
    switch (event.code) {
        case "Space":
            event.preventDefault();
            togglePlaybackPause();
            break;
        case "Digit0":
            watchPlayer(-1);
            break;
        case "Digit1":
            watchPlayer(0);
            break
        case "Digit2":
            watchPlayer(1);
            break;
        case "Digit3":
            watchPlayer(2);
            break;
        case "Digit4":
            watchPlayer(3);
            break;
    }
};

window.onload = init;