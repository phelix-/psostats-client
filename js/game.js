import * as THREE from '/js/three.module.js';
import { OrbitControls } from '/js/OrbitControls.js';

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

function createCanabinGeometry() {
    const canabinGeometry = new THREE.Group()
    const canabinMesh = new THREE.Mesh(new THREE.CylinderGeometry(4,4,2),  new THREE.MeshToonMaterial( {color: "#6b5353"}));
    const canabinFaceMesh = new THREE.Mesh(new THREE.CylinderGeometry(1,1,2),  new THREE.MeshToonMaterial( {color: "#1f359c"}))
    canabinFaceMesh.position.set(0, 1, 0);
    canabinGeometry.add(canabinMesh);
    canabinGeometry.add(canabinFaceMesh);
    canabinGeometry.rotateZ(1.57)
    return canabinGeometry;
}
function createCanuneGeometry() {
    const canabinGeometry = new THREE.Group()
    const canabinMesh = new THREE.Mesh(new THREE.CylinderGeometry(4,4,2),  new THREE.MeshPhongMaterial( {color: "#b96a07"}));
    const canabinFaceMesh = new THREE.Mesh(new THREE.CylinderGeometry(1,1,2),  new THREE.MeshPhongMaterial( {color: "#1f359c"}))
    canabinFaceMesh.position.set(0, 1, 0);
    canabinGeometry.add(canabinMesh);
    canabinGeometry.add(canabinFaceMesh);
    canabinGeometry.rotateZ(1.57)
    return canabinGeometry;
}

function createBaranzGeometry() {
    const body = new THREE.Mesh(new THREE.CylinderGeometry(14,10,20),  new THREE.MeshToonMaterial( {color: "#0002b2"}));
    const canabinFaceMesh = new THREE.Mesh(new THREE.CylinderGeometry(1,1,2),  new THREE.MeshPhongMaterial( {color: "#1f359c"}))
    canabinFaceMesh.position.set(0, 1, 0);
    const group = new THREE.Group()
    group.add(body);
    group.add(canabinFaceMesh);
    return group;
}

function createClawGeometry() {
    const combined = new THREE.Group()
    const body = new THREE.Mesh(new THREE.CylinderGeometry(
    4, 4, 2, 8, 2,
        false, Math.PI * 0.25, Math.PI * 1.5),  new THREE.MeshToonMaterial( {color: "#807373"}));
    const tail = new THREE.Mesh(new THREE.ConeGeometry(1,8,8),  new THREE.MeshToonMaterial( {color: "#807373"}))
    body.rotateY(-Math.PI * 0.5);
    tail.rotateZ(-Math.PI * 0.5);
    tail.position.set(5, 0, 0);
    // canabinFaceMesh.position.set(0, 1, 0);
    combined.add(body);
    combined.add(tail);
    body.castShadow = true;
    tail.castShadow = true;
    // combined.rotateZ(1.57)
    return combined;
}

function getGeometry(unitxtId) {
    switch (unitxtId) {
        case 25: return createBaranzGeometry();
        case 28: return createCanabinGeometry();
        case 29: return createCanuneGeometry();
        case 38: return createClawGeometry();
        default:
            const meshInfo = monsterMeshes[unitxtId]
            return new THREE.Mesh(meshInfo.geometry, meshInfo.material);
    }
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

    /* Gillchic */ "50": {"geometry": new THREE.ConeGeometry(5,8), "material": new THREE.MeshToonMaterial( {color: "#707070"}), "heightOffset": 3},
    "24": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshToonMaterial( {color: "#702a00"}), "heightOffset": 3},
    /* Canabin */ "28": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshPhongMaterial( {color: "#6b5353"}), "heightOffset": 3},
    /* Canune */ "29": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshToonMaterial( {color: "#b96a07"}), "heightOffset": 3},
    "26": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshToonMaterial( {color: "#326dff", wireframe: true}), "heightOffset": 3},
    /* Sinow Red */ "27": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshToonMaterial( {color: "#c00000", wireframe: true}), "heightOffset": 3},
    "25": { "geometry": new THREE.SphereGeometry(12,8,8), "material": new THREE.MeshToonMaterial( {color: "#0002b2", wireframe: true}), "heightOffset": 3},

    "41": {"geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#43a401", wireframe: true}), "heightOffset": 3},
    "42": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "43": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a48301", wireframe: true}), "heightOffset": 3},
    /* Claw */ "38": {"geometry": new THREE.CylinderGeometry(4,4,2), "material": new THREE.MeshBasicMaterial( {color: "#886060"}), "heightOffset": 3},
    "40": {"geometry": new THREE.CylinderGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#886060"}), "heightOffset": 3},
    /* Delsaber */ "30": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#a16dff", wireframe: true}), "heightOffset": 3},
    /* Sorc */ "31": { "geometry": new THREE.SphereGeometry(4,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "32": { "geometry": new THREE.SphereGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    "33": { "geometry": new THREE.SphereGeometry(2,8,8), "material": new THREE.MeshBasicMaterial( {color: "#5900ff", wireframe: true}), "heightOffset": 3},
    /* Indie Belra */ "37": { "geometry": new THREE.CylinderGeometry(5,5,8), "material": new THREE.MeshBasicMaterial( {color: "#00ffe0", wireframe: true}), "heightOffset": 3},
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

function init() {
    pauseButton.onclick = togglePlaybackPause;
    mapFullscreenButton.onclick = toggleFullscreenMap;
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
            if (monsterLocation.HP > 0) {
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
        if (cameraUnset ||
            Math.abs(camera.position.x - primaryPlayerState.X) +
            Math.abs(camera.position.y - primaryPlayerState.Y) +
            Math.abs(camera.position.z - primaryPlayerState.Z) > 1000
        ) {
            camera.position.x = primaryPlayerState.X;
            camera.position.y = primaryPlayerState.Y + 300;
            camera.position.z = primaryPlayerState.Z - 300;
            controls.target.set(primaryPlayerState.X, primaryPlayerState.Y, primaryPlayerState.Z);
            cameraUnset = false;
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
                    let monsterMesh = getGeometry(monsterInfo.UnitxtId);
                    monsterMesh["monsterId"] = monsterId;
                    if (monsterMesh["children"] && monsterMesh.children.length > 0) {
                        // for groups
                        monsterMesh.children[0]["monsterId"] = monsterId;
                    }
                    monsterMesh.castShadow = true;
                    monster = new THREE.Group();
                    monster.add(monsterMesh);
                } else {
                    console.warn("Missing mesh", monsterInfo.UnitxtId, monsterInfo.Name)
                    monster = new THREE.Mesh(new THREE.SphereGeometry(3,8,8), new THREE.MeshBasicMaterial( {color: "white", wireframe: true}))
                }
                visibleMonsters[monsterId] = monster;
                monster.castShadow = true;
                scene.add(monster)
            }
            let monsterLocationElement = currentDataFrame.MonsterLocation[monsterId];
            monster.position.x = monsterLocationElement.X ;
            monster.position.y = monsterLocationElement.Y + 3 ;
            monster.position.z = monsterLocationElement.Z ;
            let monsterAngle = hasFacing ? ((monsterLocationElement.Facing * (Math.PI * 2)) / 0xFFFF) + (Math.PI / 2) : 0;
            // 16383 -> 3.14
            // 0 ->
            // 49151 -> 0
            monster.rotation.set(0, monsterAngle, 0);
            if (monsterLocationElement.Frozen && monster["iceMesh"] === undefined) {
                const iceMesh = new THREE.Mesh(new THREE.OctahedronGeometry(10), new THREE.MeshBasicMaterial({color: iceColor, opacity: 0.3, transparent: true}));
                // iceMesh["monsterId"] = monsterId
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

window.onload = init;