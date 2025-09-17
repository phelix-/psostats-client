import * as THREE from "./three.module.js";

export function getGeometry(unitxtId) {
    switch (unitxtId) {
        case 5:
            return sandRappy("#fff200")
        case 9: // Bartle
            return grunt("#755138");
        case 10: // Barble
            return grunt("#668149");
        case 11: // Tollaw
            return grunt("#01A417");
        case 25:
            return baranz();
        case 28:
            return canabin();
        case 29:
            return canune();
        case 38:
            return claw();
        case 93:
            return girta();
        case 94:
            return zu();
        case 95:
            return pazuzu();
        case 96: // Boota
            return grunt("#d69292");
        case 97: // Ze Boota
            return grunt("#9d36ac");
        case 98: // Ba Boota
            return grunt("#572325");
        case 99:
            return dorphon();
        case 104:
            return sandRappy();
        case 105: // Del Rappy
            return sandRappy("#5b4b67");
        case 106:
            return saintMilion();
        default:
            let meshInfo = monsterMeshes[unitxtId]
            if (meshInfo === undefined) {
                console.log(`Undefined mesh for ${unitxtId}`);
                meshInfo = {
                    geometry: new THREE.SphereGeometry(3, 8, 8),
                    material: new THREE.MeshBasicMaterial({color: "white", wireframe: true})
                };
            }
            return new THREE.Mesh(meshInfo.geometry, meshInfo.material);
    }
}

function canabin() {
    const canabinGeometry = new THREE.Group()
    const canabinMesh = new THREE.Mesh(new THREE.CylinderGeometry(4,4,2),  new THREE.MeshToonMaterial( {color: "#6b5353"}));
    canabinMesh.castShadow = true;
    canabinGeometry.add(canabinMesh);
    const canabinFaceMesh = new THREE.Mesh(new THREE.CylinderGeometry(1,1,2),  new THREE.MeshToonMaterial( {color: "#1f359c"}))
    canabinFaceMesh.position.set(0, 1, 0);
    canabinFaceMesh.castShadow = true;
    canabinGeometry.add(canabinFaceMesh);
    canabinGeometry.rotateZ(1.57)
    return canabinGeometry;
}
function canune() {
    const canabinGeometry = new THREE.Group()
    const canabinMesh = new THREE.Mesh(new THREE.CylinderGeometry(4,4,2),  new THREE.MeshPhongMaterial( {color: "#b96a07"}));
    canabinMesh.castShadow = true;
    canabinGeometry.add(canabinMesh);
    const canabinFaceMesh = new THREE.Mesh(new THREE.CylinderGeometry(1,1,2),  new THREE.MeshPhongMaterial( {color: "#1f359c"}))
    canabinFaceMesh.position.set(0, 1, 0);
    canabinFaceMesh.castShadow = true;
    canabinGeometry.add(canabinFaceMesh);
    canabinGeometry.rotateZ(1.57)
    return canabinGeometry;
}

function baranz() {
    const group = new THREE.Group()
    const body = new THREE.Mesh(new THREE.CylinderGeometry(14,10,20),  new THREE.MeshToonMaterial( {color: "#0002b2"}));
    group.add(body);
    body.castShadow = true;
    const face = new THREE.Mesh(new THREE.CylinderGeometry(1,1,2),  new THREE.MeshPhongMaterial( {color: "#1f359c"}))
    face.position.set(0, 1, 0);
    face.castShadow = true;
    group.add(face);
    return group;
}

function claw() {
    const combined = new THREE.Group()
    const body = new THREE.Mesh(new THREE.CylinderGeometry(
        4, 4, 2, 8, 2,
        false, Math.PI * 0.25, Math.PI * 1.5),  new THREE.MeshToonMaterial( {color: "#807373"}));
    const tail = new THREE.Mesh(new THREE.ConeGeometry(1,8,8),  new THREE.MeshToonMaterial( {color: "#807373"}))
    body.rotateY(-Math.PI * 0.5);
    tail.rotateZ(-Math.PI * 0.5);
    tail.position.set(5, 0, 0);
    combined.add(body);
    combined.add(tail);
    body.castShadow = true;
    tail.castShadow = true;
    return combined;
}

function zu() {
    const body = new THREE.Mesh(
        new THREE.SphereGeometry(12, 8, 8),
        new THREE.MeshToonMaterial({color: "#fa53b3"})
    );
    body.castShadow = true;
    body.position.set(0, 15, 0);
    return body;
}

function pazuzu() {
    const body = new THREE.Mesh(
        new THREE.SphereGeometry(12, 8, 8),
        new THREE.MeshToonMaterial({color: "#00b522"})
    );
    body.castShadow = true;
    body.position.set(0, 15, 0);
    return body;
}


function girta() {
    const body = new THREE.Mesh(
        new THREE.SphereGeometry(18, 8, 8),
        new THREE.MeshToonMaterial({color: "#780078"})
    );
    body.castShadow = true;
    return body;
}

function dorphon() {
    const body = new THREE.Mesh(
        new THREE.CapsuleGeometry(24, 24, 8, ),
        new THREE.MeshToonMaterial({color: "#685147"})
    );
    body.castShadow = true;
    body.rotateZ(Math.PI * 0.5);
    return body;
}

function saintMilion() {
    const body = new THREE.Mesh(
        new THREE.OctahedronGeometry(8),
        new THREE.MeshToonMaterial({color: "#78000a"})
    );
    body.castShadow = true;
    return body;
}

function grunt(color) {
    const body = new THREE.Mesh(
        new THREE.ConeGeometry(8, 8, 8, 8),
        new THREE.MeshToonMaterial({color: color})
    );
    body.castShadow = true;
    body.rotateZ(Math.PI * 0.5);
    return body;
}

function sandRappy(color = "#ff8000") {
    const combined = new THREE.Group();
    const body = new THREE.Mesh(
        new THREE.SphereGeometry(5, 8, 8),
        new THREE.MeshToonMaterial({color: color})
    );
    body.castShadow = true;
    combined.add(body);
    const nose = new THREE.Mesh(
        new THREE.ConeGeometry(2, 4, 8, 8),
        new THREE.MeshToonMaterial({color: "#4c4959"})
    )
    nose.position.set(-5, 2.5, 0);
    nose.rotateZ(Math.PI * 0.5);
    nose.castShadow = true;
    combined.add(nose);

    return combined;
}

const monsterMeshes = {
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

    "101": { "geometry": new THREE.SphereGeometry(3,8,8), "material": new THREE.MeshBasicMaterial( {color: "#7601a4", wireframe: true}), "heightOffset": 3},
    "103": { "geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#01a417", wireframe: true}), "heightOffset": 3},
    "102": { "geometry": new THREE.CylinderGeometry(12,12,24), "material": new THREE.MeshBasicMaterial( {color: "#c900b5", wireframe: true}), "heightOffset": 3},
    "90": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#49811f", wireframe: true}), "heightOffset": 3},
    "89": {"geometry": new THREE.SphereGeometry(5,8,8), "material": new THREE.MeshBasicMaterial( {color: "#b0512a", wireframe: true}), "heightOffset": 3},
    "91": { "geometry": new THREE.CylinderGeometry(6,6,1), "material": new THREE.MeshBasicMaterial( {color: "#ff7e95"}), "heightOffset": 3},
    "94": { "geometry": new THREE.SphereGeometry(6,8,8), "material": new THREE.MeshBasicMaterial( {color: "#fa53b3", wireframe: true}), "heightOffset": 3},
    "93": { "geometry": new THREE.SphereGeometry(20,8,8), "material": new THREE.MeshBasicMaterial( {color: "#ab5bff", wireframe: true}), "heightOffset": 3},
};
