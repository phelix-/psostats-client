'use strict';

let sortColumn = "";
let sortAscending = null;

const barriers = {
    "None": {atp: 0, ata: 0},
    "Red Ring": {atp: 20, ata: 20},
    "Ranger Wall": {atp: 0, ata: 20},
    "Kasami Bracer": {atp: 35, ata: 0},
    "Combat Gear": {atp: 35, ata: 0},
    "Safety Heart": {atp: 0, ata: 0},
    "S-Parts ver2.01": {atp: 0, ata: 15},
    "Black Ring (1)": {atp: 50, ata: 0},
    "Black Ring (2)": {atp: 100, ata: 0},
    "Black Ring (3)": {atp: 150, ata: 0},
}

const possWeapons = [
    "Ancient Saber",
    "Delsaber's Buster",
    "Durandal",
    "Excalibur",
    "Flamberge",
    "Galatine",
    "Kaladbolg",
    "Kusanagi",
    "Lavis Cannon",
    "Red Saber",
    "Two Kamui",
    "Guren",
    "Kamui",
    "Orotiagito",
    "Raikiri",
    "Sange",
    "Shichishito",
    "Shouren",
    "Yamigarasu",
    "Yasha",
    "Asuka",
    "Jizai",
    "Musashi",
    "Sange & Yasha",
    "Yamato"
]
const frameData = {
    "Saber": {n1: 29, n1c: 13, n2: 24, n2c: 10, n3: 31, h1: 37, h1c: 21, h2: 29, h2c: 15, h3: 34},
    "Sword": {n1: 39, n1c: 17, n2: 31, n2c: 15, n3: 43, h1: 46, h1c: 24, h2: 35, h2c: 19, h3: 43},
    "Dagger": {n1: 40, n1c: 21, n2: 30, n2c: 15, n3: 50, h1: 46, h1c: 27, h2: 35, h2c: 20, h3: 49},
    "Partisan": {n1: 39, n1c: 17, n2: 30, n2c: 14, n3: 33, h1: 46, h1c: 24, h2: 35, h2c: 19, h3: 35},
    "Slicer": {n1: 40, n1c: 21, n2: 32, n2c: 12, n3: 42, h1: 47, h1c: 28, h2: 37, h2c: 17, h3: 43},
    "Double Saber": {n1: 39, n1c: 22, n2: 27, n2c: 11, n3: 51, h1: 46, h1c: 29, h2: 32, h2c: 16, h3: 49},
    "Claw": {n1: 27, n1c: 16, n2: 22, n2c: 11, n3: 37, h1: 35, h1c: 24, h2: 27, h2c: 16, h3: 39},
    "Katana": {n1: 30, n1c: 14, n2: 29, n2c: 15, n3: 44, h1: 38, h1c: 22, h2: 33, h2c: 19, h3: 44},
    "Twin Sword": {n1: 37, n1c: 18, n2: 34, n2c: 19, n3: 51, h1: 44, h1c: 25, h2: 37, h2c: 22, h3: 49},
    "Fist": {n1: 26, n1c: 19, n2: 26, n2c: 19, n3: 35, h1: 33, h1c: 26, h2: 29, h2c: 22, h3: 36},
    "Master Raven": {n1: 26, h1: 36},
    "L&K38 Combat": {n1: 46, h1: 55},
    "Handgun": {n1: 27, n1c: 14, n2: 25, n2c: 11, n3: 19, h1: 34, h1c: 22, h2: 30, h2c: 16, h3: 25},
    "Rifle": {n1: 29, n1c: 15, n2: 25, n2c: 12, n3: 20, h1: 37, h1c: 23, h2: 30, h2c: 17, h3: 26},
    "Mechgun": {n1: 49, n1c: 12, n2: 45, n2c: 10, n3: 42, h1: 58, h1c: 21, h2: 50, h2c: 15, h3: 48},
    "Shot": {n1: 50, n1c: 25, n2: 43, n2c: 21, n3: 34, h1: 56, h1c: 31, h2: 46, h2c: 24, h3: 38},
    "Launcher": {n1: 46, n1c: 21, n2: 41, n2c: 19, n3: 36, h1: 52, h1c: 27, h2: 44, h2c: 22, h3: 39},
    "Cane": {n1: 29, n1c: 13, n2: 27, n2c: 13, n3: 39, h1: 37, h1c: 21, h2: 32, h2c: 18, h3: 40},
    "Rod": {n1: 29, n1c: 14, n2: 27, n2c: 14, n3: 40, h1: 37, h1c: 22, h2: 32, h2c: 19, h3: 41},
    "Wand": {n1: 30, n1c: 14, n2: 29, n2c: 15, n3: 40, h1: 37, h1c: 21, h2: 33, h2c: 19, h3: 41},
    "Card": {n1: 33, n1c: 18, n2: 30, n2c: 17, n3: 47, h1: 40, h1c: 25, h2: 34, h2c: 21, h3: 47},
}

const femaleFrameData = {
    "Saber": {n1: 29, n1c: 13, n2: 26, n2c: 12, n3: 35, h1: 37, h1c: 21, h2: 31, h2c: 17, h3: 37},
    "Double Saber": {n1: 35, n1c: 18, n2: 26, n2c: 10, n3: 45, h1: 42, h1c: 25, h2: 31, h2c: 16, h3: 45},
    "Claw": {n1: 27, n1c: 16, n2: 22, n2c: 11, n3: 34, h1: 35, h1c: 24, h2: 27, h2c: 16, h3: 36},
    "Katana": {n1: 30, n1c: 14, n2: 29, n2c: 15, n3: 41, h1: 38, h1c: 22, h2: 33, h2c: 19, h3: 42},
    "Fist": {n1: 25, n1c: 18, n2: 23, n2c: 16, n3: 39, h1: 32, h1c: 25, h2: 27, h2c: 20, h3: 39},
    "Last Swan": {n1: 26, n1c: 14, n2: 25, n2c: 12, n3: 21, h1: 36, h1c: 24, h2: 32, h2c: 19, h3: 28},
    "Cane": {n1: 30, n1c: 14, n2: 27, n2c: 13, n3: 39, h1: 38, h1c: 22, h2: 32, h2c: 18, h3: 40},
    "Rod": {n1: 33, n1c: 18, n2: 30, n2c: 17, n3: 40, h1: 40, h1c: 25, h2: 34, h2c: 21, h3: 41},
}

const classSpecificFrameData = {
    HUcaseal: {
        "Dagger": {n1: 37, n1c: 18, n2: 29, n2c: 14, n3: 41, h1: 44, h1c: 25, h2: 33, h2c: 18, h3: 42},
        "Double Saber": {n1: 35, n1c: 18, n2: 27, n2c: 11, n3: 45, h1: 42, h1c: 25, h2: 32, h2c: 16, h3: 45},
        "Claw": {n1: 32, n1c: 21, n2: 25, n2c: 14, n3: 34, h1: 38, h1c: 27, h2: 29, h2c: 18, h3: 36},
        "Fist": {n1: 22, n1c: 15, n2: 21, n2c: 14, n3: 36, h1: 30, h1c: 23, h2: 26, h2c: 19, h3: 38},
    },
    RAmarl: {
        "Claw": {n1: 32, n1c: 21, n2: 25, n2c: 14, n3: 34, h1: 38, h1c: 27, h2: 29, h2c: 18, h3: 36},
        "Handgun": {n1: 26, n1c: 13, n2: 24, n2c: 10, n3: 19, h1: 34, h1c: 21, h2: 29, h2c: 15, h3: 25},
    },
    FOmar: {
        "Claw": {n1: 32, n1c: 21, n2: 22, n2c: 11, n3: 41, h1: 38, h1c: 27, h2: 27, h2c: 16, h3: 42},
        "Fist": {n1: 24, n1c: 17, n2: 23, n2c: 16, n3: 45, h1: 31, h1c: 24, h2: 27, h2c: 20, h3: 45},
        "Rod": {n1: 31, n1c: 16, n2: 27, n2c: 14, n3: 40, h1: 38, h1c: 23, h2: 32, h2c: 19, h3: 41},
        "Wand": {n1: 30, n1c: 14, n2: 30, n2c: 16, n3: 42, h1: 37, h1c: 21, h2: 34, h2c: 20, h3: 42},
    },
    FOmarl: {
        "Saber": {n1: 31, n1c: 15, n2: 26, n2c: 12, n3: 37, h1: 39, h1c: 23, h2: 31, h2c: 17, h3: 39},
        "Sword": {n1: 38, n1c: 16, n2: 31, n2c: 15, n3: 53, h1: 45, h1c: 23, h2: 35, h2c: 19, h3: 51},
        "Dagger": {n1: 42, n1c: 23, n2: 34, n2c: 19, n3: 53, h1: 48, h1c: 29, h2: 37, h2c: 22, h3: 51},
        "Partisan": {n1: 38, n1c: 16, n2: 32, n2c: 16, n3: 42, h1: 46, h1c: 24, h2: 36, h2c: 20, h3: 42},
        "Slicer": {n1: 39, n1c: 20, n2: 33, n2c: 13, n3: 38, h1: 45, h1c: 26, h2: 38, h2c: 18, h3: 40},
        "Double Saber": {n1: 38, n1c: 21, n2: 27, n2c: 11, n3: 48, h1: 44, h1c: 27, h2: 32, h2c: 16, h3: 47},
        "Claw": {n1: 27, n1c: 16, n2: 22, n2c: 11, n3: 34, h1: 35, h1c: 24, h2: 27, h2c: 16, h3: 36},
        "Katana": {n1: 30, n1c: 14, n2: 29, n2c: 15, n3: 41, h1: 38, h1c: 22, h2: 33, h2c: 19, h3: 42},
        "Fist": {n1: 24, n1c: 17, n2: 21, n2c: 14, n3: 33, h1: 31, h1c: 24, h2: 25, h2c: 18, h3: 35},
        "Handgun": {n1: 30, n1c: 17, n2: 28, n2c: 14, n3: 23, h1: 37, h1c: 24, h2: 32, h2c: 18, h3: 28},
        "Last Swan": {n1: 27, n1c: 15, n2: 25, n2c: 12, n3: 21, h1: 36, h1c: 24, h2: 32, h2c: 19, h3: 28},
        "Rifle": {n1: 31, n1c: 17, n2: 28, n2c: 15, n3: 26, h1: 38, h1c: 24, h2: 32, h2c: 19, h3: 30},
        "Shot": {n1: 43, n1c: 18, n2: 32, n2c: 10, n3: 27, h1: 50, h1c: 25, h2: 37, h2c: 15, h3: 32},
        "L&K38 Combat": {n1: 45, h1: 54},
        "Cane": {n1: 29, n1c: 13, n2: 27, n2c: 13, n3: 39, h1: 37, h1c: 21, h2: 32, h2c: 18, h3: 40},
        "Rod": {n1: 31, n1c: 16, n2: 27, n2c: 14, n3: 36, h1: 38, h1c: 23, h2: 32, h2c: 19, h3: 38},
        "Wand": {n1: 30, n1c: 14, n2: 30, n2c: 16, n3: 40, h1: 37, h1c: 21, h2: 34, h2c: 20, h3: 41},
    }
}

function accuracyModifierForAttackType(attackType, special) {
    if (attackType === 'NORMAL') {
        return 1.0;
    } else if (attackType === 'HEAVY') {
        return 0.7;
    } else if (attackType === 'SPECIAL') {
        return getSpecialAccuracyModifier(special);
    }
}

function getSpecialAccuracyModifier(special) {
    if (special === 'Vjaya' || special === "Dark Flow" || special === "Frozen Shooter") {
        return 0.7;
    } else {
        return 0.5;
    }
}

function getDamageModifierForAttackType(attackType, special) {
    if (attackType === 'NORMAL') {
        return 0.9;
    } else if (attackType === 'HEAVY') {
        return 1.7;
    } else if (attackType === 'SPECIAL') {
        return getSpecialDamageModifier(special);
    } else if (attackType === 'NONE') {
        return 0;
    }
}

function getSpecialDamageModifier(special) {
    if (special === 'Arrest') {
        return 0.5;
    } else if (special === 'Raikiri') {
        return 0.875;
    } else if (special === 'Lavis Cannon') {
        return 0.5;
    } else if (special === 'Lavis Blade') {
        return 0.583;
    } else if (special === 'Dark Flow' || special === 'TJS' || special === "Frozen Shooter") {
        return 1.7;
    } else if (special === 'Orotiagito') {
        return 1.75;
    } else if (special === 'Charge' || special === 'Spirit' || special === 'Berserk') {
        return 3.0;
    } else if (special === 'Vjaya') {
        return 5.1;
    } else {
        return 0;
    }
}

function getEvpModifier(frozen, paralyzed) {
    let modifier = 1.0;
    if (frozen) {
        modifier -= 0.3;
    }
    if (paralyzed) {
        modifier -= 0.15;
    }
    return modifier;
}

function createMonsterRow(
    special, autoCombo, weapon, enemy,
    evpModifier, base_ata, snGlitch, atpInput, comboInput, range
) {
    let modified_evp = enemy.evp * evpModifier;

    let baseDamage = calculateBaseDamage(atpInput, enemy);
    let damageToUse = atpInput.useMaxDamageRoll ? baseDamage.nMax : baseDamage.nMin;
    if (autoCombo) {
        return generateAutoCombo(special, weapon, enemy, modified_evp, base_ata, snGlitch, damageToUse, atpInput, comboInput, range)
    }

    let accuracyResult = getAccuracyForCombo(base_ata, comboInput, special, modified_evp, snGlitch, range);
    let comboDamage = getDamageForCombo(enemy.hp, atpInput, comboInput, special, damageToUse);
    let percentDamage = 100 * (comboDamage.total / enemy.hp);
    if (percentDamage > 100) {
        percentDamage = 100;
    }

    return {
        name: enemy.name,
        hp: enemy.hp,
        percentDamage: percentDamage,
        comboDamage: comboDamage.total,
        overallAccuracy: accuracyResult.overallAccuracy,
        overallMinAccuracy: accuracyResult.overallMinAccuracy,
        a1Damage: comboDamage.a1Damage,
        a1Type: comboInput.a1Type,
        a1Accuracy: accuracyResult.a1Accuracy,
        a1MinAccuracy: accuracyResult.a1MinAccuracy,
        a2Type: comboInput.a2Type,
        a2Damage: comboDamage.a2Damage,
        a2Accuracy: accuracyResult.a2Accuracy,
        a2MinAccuracy: accuracyResult.a2MinAccuracy,
        a3Type: comboInput.a3Type,
        a3Damage: comboDamage.a3Damage,
        a3Accuracy: accuracyResult.a3Accuracy,
        a3MinAccuracy: accuracyResult.a3MinAccuracy,
    }
}

// 0 if range does not apply (smartlink checked, melee weapon, class is RA)
function getRange() {
    const smartlink = $('#smartlinkInput').is(":checked");
    if (smartlink) {
        return 0;
    }
    let className = $('#class-select').val();
    if (className.startsWith("RA")) {
        return 0;
    }
    return selectedWeapon.horizontalDistance;
}

function getAccuracyForCombo(base_ata, comboInput, special, modified_evp, snGlitch, range) {
    let a1Accuracy = calculateAccuracy(base_ata, comboInput.a1Type, special, 1.0, modified_evp);
    let a2Accuracy = calculateAccuracy(base_ata, comboInput.a2Type, special, 1.3, modified_evp);
    let a3Accuracy = calculateAccuracy(base_ata, comboInput.a3Type, special, 1.69, modified_evp);
    let a1MinAccuracy = a1Accuracy;
    let a2MinAccuracy = a2Accuracy;
    let a3MinAccuracy = a3Accuracy;
    if (range > 0) {
        a1MinAccuracy = calculateAccuracy(base_ata, comboInput.a1Type, special, 1.0, modified_evp, range);
        a2MinAccuracy = calculateAccuracy(base_ata, comboInput.a2Type, special, 1.3, modified_evp, range);
        a3MinAccuracy = calculateAccuracy(base_ata, comboInput.a3Type, special, 1.69, modified_evp, range);
    }

    // Account for SN glitch - I'm assuming optimistic case where they're able to glitch
    // if the accuracy is better but not if it's worse
    let glitchedA1Accuracy = a1Accuracy;
    let glitchedA1MinAccuracy = a1MinAccuracy;
    if (snGlitch && a2Accuracy > a1Accuracy && comboInput.a2Type !== 'NONE') {
        glitchedA1Accuracy = a2Accuracy;
        glitchedA1MinAccuracy = a2MinAccuracy;
    }

    let glitchedA2Accuracy = a2Accuracy;
    let glitchedA2MinAccuracy = a2MinAccuracy;
    if (snGlitch && a3Accuracy > a2Accuracy && comboInput.a3Type !== 'NONE') {
        glitchedA2Accuracy = a3Accuracy;
        glitchedA2MinAccuracy = a3MinAccuracy;
    }
    let overallAccuracy = Math.pow((glitchedA1Accuracy * 0.01), comboInput.a1Hits)
        * Math.pow((glitchedA2Accuracy * 0.01), comboInput.a2Hits)
        * Math.pow((a3Accuracy * 0.01), comboInput.a3Hits);
    let minOverallAccuracy = (Math.pow((glitchedA1MinAccuracy * 0.01), comboInput.a1Hits)
        * Math.pow((glitchedA2MinAccuracy * 0.01), comboInput.a2Hits)
        * Math.pow((a3MinAccuracy * 0.01), comboInput.a3Hits)) * 100;
    overallAccuracy *= 100;
    return {overallAccuracy, a1Accuracy, a2Accuracy, a3Accuracy,
        overallMinAccuracy: minOverallAccuracy, a1MinAccuracy, a2MinAccuracy, a3MinAccuracy};
}

function getDamageForCombo(enemyHp, atpInput, comboInput, special, baseDamage) {
    let total = 0;
    let a1Damage = 0;
    let a2Damage = 0;
    let a3Damage = 0;
    let hpCutModifier = special === "Devil's" ? getDevilsModifier(atpInput.playerClass) : getDemonsModifier(atpInput.playerClass);
    if (comboInput.a1Type === "SPECIAL" && ["Demon's", "Devil's"].includes(special)) {
        for (let i = 0; i < comboInput.a1Hits; i++) {
            a1Damage = (enemyHp - total) * hpCutModifier;
            total += a1Damage;
        }
    } else {
        a1Damage = getDamageModifierForAttackType(comboInput.a1Type, special) * baseDamage;
        total += a1Damage * comboInput.a1Hits;
    }
    if (comboInput.a2Type === "SPECIAL" && ["Demon's", "Devil's"].includes(special)) {
        for (let i = 0; i < comboInput.a2Hits; i++) {
            a2Damage = (enemyHp - total) * hpCutModifier;
            total += a2Damage;
        }
    } else {
        a2Damage = getDamageModifierForAttackType(comboInput.a2Type, special) * baseDamage;
        total += a2Damage * comboInput.a2Hits;
    }
    if (comboInput.a3Type === "SPECIAL" && ["Demon's", "Devil's"].includes(special)) {
        for (let i = 0; i < comboInput.a3Hits; i++) {
            a3Damage = (enemyHp - total) * hpCutModifier;
            total += a3Damage;
        }
    } else {
        a3Damage = getDamageModifierForAttackType(comboInput.a3Type, special) * baseDamage;
        total += a3Damage * comboInput.a3Hits;
    }

    return {total, a1Damage, a2Damage, a3Damage};
}

function getDemonsModifier(playerClass) {
    switch (playerClass) {
        case "HUcast":
        case "HUcaseal":
        case "RAcast":
        case "RAcaseal":
            return 0.45;
        default:
            return 0.75;
    }
}

function getDevilsModifier(playerClass) {
    switch (playerClass) {
        case "HUcast":
        case "HUcaseal":
        case "RAcast":
        case "RAcaseal":
            return 0.20;
        default:
            return 0.50;
    }
}

function generateAutoCombo(
    special, weapon, enemy, modified_evp, base_ata, snGlitch, baseDamage, atpInput, comboInput, range
) {
    let className = $('#class-select').val();
    let frameData = getFrameDataForWeapon(weapon, className);

    let bestCombo = null;
    let bestComboFrames = null;
    let bestComboDamage = null;
    let bestComboAccuracy = null;
    const attacks = ["NONE", "NORMAL", "HEAVY", "SPECIAL"]
    for (let a1 in attacks) {
        for (let a2 in attacks) {
            if (attacks[a2] !== "NONE" && !!weapon.combo && weapon.combo.attack2 === "NONE") {
                continue;
            }
            for (let a3 in attacks) {
                if (attacks[a3] !== "NONE" && !!weapon.combo && weapon.combo.attack3 === "NONE") {
                    continue;
                }
                comboInput.a1Type = attacks[a1];
                comboInput.a2Type = attacks[a2];
                comboInput.a3Type = attacks[a3];
                let frames = getFramesForCombo(attacks[a1], attacks[a2], attacks[a3], frameData.animationFrameData)
                let accuracyResult = getAccuracyForCombo(base_ata, comboInput, special, modified_evp, snGlitch, range);
                let comboDamage = getDamageForCombo(enemy.hp, atpInput, comboInput, special, baseDamage);
                if (accuracyResult.overallAccuracy < 100 && bestCombo != null) {
                    continue;
                }
                if (comboInput.a1Type === "NONE" ||
                    comboInput.a2Type === "NONE" && comboInput.a3Type !== "NONE") {
                    continue;
                }
                if (bestCombo == null ||
                    (bestComboDamage.total < enemy.hp && comboDamage.total > bestComboDamage.total) ||
                    (bestComboDamage.total >= enemy.hp && comboDamage.total >= enemy.hp && frames < bestComboFrames)
                ) {
                    bestCombo = [attacks[a1], attacks[a2], attacks[a3]];
                    bestComboDamage = comboDamage;
                    bestComboFrames = frames;
                    bestComboAccuracy = accuracyResult;
                }
            }
        }
    }

    let percentDamage = 100 * (bestComboDamage.total / enemy.hp);
    if (percentDamage > 100) {
        percentDamage = 100;
    }

    let comboName = "";
    for (let i = 0; i < 3; i++) {
        comboName += bestCombo[i] === "NONE" ? "." : bestCombo[i][0]
    }
    return {
        name: enemy.name + " (" + comboName + " " + bestComboFrames + "f)",
        hp: enemy.hp,
        percentDamage: percentDamage,
        comboDamage: bestComboDamage.total,
        overallAccuracy: bestComboAccuracy.overallAccuracy,
        a1Damage: bestComboDamage.a1Damage,
        a1Type: bestCombo[0],
        a1Accuracy: bestComboAccuracy.a1Accuracy,
        a2Type: bestCombo[1],
        a2Damage: bestComboDamage.a2Damage,
        a2Accuracy: bestComboAccuracy.a2Accuracy,
        a3Type: bestCombo[2],
        a3Damage: bestComboDamage.a3Damage,
        a3Accuracy: bestComboAccuracy.a3Accuracy,
    }
}

function formatAccuracyText(accuracy, minAccuracy, showAccuracyRange) {
    if (showAccuracyRange) {
        return `${formatAccuracy(minAccuracy)}% - ${formatAccuracy(accuracy)}%`;
    } else {
        return `${formatAccuracy(accuracy)}%`;
    }
}

function appendMonsterRow(rowEntry, showAccuracyRange) {
    let comboKill = rowEntry.comboDamage > rowEntry.hp;
    let damageBgColor = comboKill ? 'rgb(61,73,61)' : 'rgb(73,73,61)';

    let a1Text = `${rowEntry.a1Damage.toFixed(0)} (${formatAccuracyText(rowEntry.a1Accuracy, rowEntry.a1MinAccuracy, showAccuracyRange)})`;
    let a2Text = `${rowEntry.a2Damage.toFixed(0)} (${formatAccuracyText(rowEntry.a2Accuracy, rowEntry.a2MinAccuracy, showAccuracyRange)})`;
    let a3Text = `${rowEntry.a3Damage.toFixed(0)} (${formatAccuracyText(rowEntry.a3Accuracy, rowEntry.a3MinAccuracy, showAccuracyRange)})`;

    return $('<tr/>')
        .append($('<th/>', {
            'colspan': 2,
            'data-label': 'monster',
            'text': rowEntry.name
        }))
        .append($('<td/>', {
            'style': 'padding: 0'
        }).append($('<div>', {
            'style': 'background: rgba(255,150,150,0.1)'
        }).append($('<div>', {
            'style': `overflow: visible; white-space: nowrap; background: ${damageBgColor}; padding: 0.78571429em 0.78571429em; width: ${rowEntry.percentDamage}%`,
            'text': `${rowEntry.comboDamage.toFixed(0)} (${rowEntry.percentDamage.toFixed(0)}%)`,
            'title': rowEntry.comboDamage.toFixed(0) + '/' + rowEntry.hp
        }))))
        .append($('<td/>', {
            'data-label': 'accuracy',
            'text': formatAccuracyText(rowEntry.overallAccuracy, rowEntry.overallMinAccuracy, showAccuracyRange),
            'style': rowEntry.overallAccuracy >= 100.0 ? 'background: rgba(150,255,150,0.1)' : 'background: rgba(255,150,150,0.1)'
        }))
        .append($('<td/>', {
            'data-label': 'a1-accuracy',
            'text': a1Text,
            'style': rowEntry.a1Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
        }))
        .append($('<td/>', {
            'data-label': 'a2-accuracy',
            'text': a2Text,
            'style': rowEntry.a2Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
        }))
        .append($('<td/>', {
            'data-label': 'a3-accuracy',
            'text': a3Text,
            'style': rowEntry.a3Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
        }));
}

function calculateBaseDamage(atpInput, enemy) {
    let areaPercent = enemy.ccaMiniboss ? 0 : atpInput.areaPercent;
    let minWeaponAtp = atpInput.minAtp * ((areaPercent * 0.01) + 1);
    let maxWeaponAtp = minWeaponAtp + (atpInput.maxAtp - atpInput.minAtp);
    let shiftaModifier = 0;
    if (atpInput.shifta > 0) {
        shiftaModifier = ((1.3 * (atpInput.shifta - 1)) + 10) * 0.01;
    }
    let zalureModifier = 0;
    if (atpInput.zalure > 0) {
        zalureModifier = ((1.3 * (atpInput.zalure - 1)) + 10) * 0.01;
    }
    let minShiftaAtp = shiftaModifier * atpInput.classMinAtp;
    let maxShiftaAtp = (shiftaModifier * atpInput.classMaxAtp) + (((atpInput.maxAtp - atpInput.minAtp) + (atpInput.classMaxAtp - atpInput.classMinAtp)) * shiftaModifier);

    let effectiveMinAtp = atpInput.classMinAtp + minWeaponAtp + minShiftaAtp;
    let effectiveMaxAtp = atpInput.classMaxAtp + maxWeaponAtp + maxShiftaAtp;

    let effectiveDfp = enemy.dfp * (1.0 - zalureModifier);

    let nMin = ((effectiveMinAtp - effectiveDfp) / 5);
    if (nMin < 0) {
        nMin = 0;
    }
    let nMax = ((effectiveMaxAtp - effectiveDfp) / 5);
    if (nMax < 0) {
        nMax = 0;
    }
    return {
        nMin: nMin,
        nMax: nMax
    };
}

function calculateAccuracy(baseAta, attackType, special, comboModifier, totalEvp, range = 0) {
    if (attackType === 'NONE') {
        return 100;
    }
    if (attackType === 'SPECIAL' && special === 'TJS') {
        return 100;
    }
    let effectiveAta = baseAta * accuracyModifierForAttackType(attackType, special) * comboModifier;
    let accuracy = effectiveAta - (totalEvp * 0.2);
    let rangePenalty = 0.33 * range;
    accuracy = accuracy - rangePenalty;
    if (accuracy > 100) {
        accuracy = 100;
    }
    if (accuracy < 0) {
        accuracy = 0;
    }
    return accuracy;
}

function pushSort(colName) {
    if (colName === sortColumn) {
        if (sortAscending === null) {
            sortAscending = true;
        } else if (sortAscending === true) {
            sortAscending = false;
        } else {
            sortAscending = null;
        }
    } else {
        sortColumn = colName;
        sortAscending = true;
    }
    updateDamageTable();
}

function updateDamageTable() {
    let frozen = $('#frozenCheckbox').is(":checked");
    let paralyzed = $('#paralyzedCheckbox').is(":checked");
    let snGlitch = $('#ataGlitch').is(":checked");
    let base_ata = Number($('#ataInput').val());
    let evpModifier = getEvpModifier(frozen, paralyzed);
    let autoCombo = $('#autoCombo').is(":checked")
    // let opm = $('#opm_checkbox').is(":checked");

    let atpInput = {
        playerClass: $('#class-select').val(),
        classMinAtp: Number($('#classMinAtpInput').val()),
        classMaxAtp: Number($('#classMaxAtpInput').val()),
        minAtp: Number($('#minAtpInput').val()),
        maxAtp: Number($('#maxAtpInput').val()),
        areaPercent: Number($('#sphereInput').val()),
        useMaxDamageRoll: $('#maxDamageCheckbox').is(":checked"),
        shifta: Number($('#shiftaInput').val()),
        zalure: Number($('#zalureInput').val()),
    };
    let range = getRange();
    let comboInput = {
        a1Type: $('#attack1').val(),
        a1Hits: Number($('#hits1').val()),
        a2Type: $('#attack2').val(),
        a2Hits: Number($('#hits2').val()),
        a3Type: $('#attack3').val(),
        a3Hits: Number($('#hits3').val()),
    }
    const special = $('#special-select').val();

    let tbody = $('#combo-calc-table tbody');
    tbody.empty();
    let rows = [];
    for (let index in selectedEnemies) {
        let enemy = enemies[selectedEnemies[index]];
        let row = createMonsterRow(
            special, autoCombo, selectedWeapon, enemy,
            evpModifier, base_ata, snGlitch, atpInput, comboInput, range
        );
        rows.push(row);
    }

    rows.sort((a, b) => {
        if (sortAscending === null) {
            return 0;
        }
        if (sortColumn === "name") {
            if (sortAscending === true) {
                return a.name.localeCompare(b.name);
            } else {
                return b.name.localeCompare(a.name);
            }
        } else if (sortColumn === "damage") {
            if (sortAscending === true) {
                return a.percentDamage - b.percentDamage;
            } else {
                return b.percentDamage - a.percentDamage;
            }
        } else if (sortColumn === "accuracy") {
            if (sortAscending === true) {
                return a.overallAccuracy - b.overallAccuracy;
            } else {
                return b.overallAccuracy - a.overallAccuracy;
            }
        }
    })

    for (let index in rows) {
        tbody.append(appendMonsterRow(rows[index], range > 0));
    }
}

function getSetEffectAtp(weapon, frameName, barrierName) {
    let atpBonus = 0;
    if (frameName === "Thirteen" && weapon.name === "Diska of Braveman") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.5
    }
    if (frameName === "Crimson Coat" && (
        weapon.name === "Red Slicer" || weapon.name === "Red Dagger" || weapon.name === "Red Saber"
    )) {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.5
    }
    if (frameName === "Samurai Armor" && weapon.name === "Orotiagito") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.3
    }
    if (frameName === "Sweetheart (1)") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.15
    }
    if (frameName === "Sweetheart (2)") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.2
    }
    if (frameName === "Sweetheart (3)") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.25
    }
    return atpBonus;
}

function getSetEffectAta(weapon, frameName, barrierName, unitName) {
    let ataBonus = 0;
    if (barrierName === "Safety Heart" && weapon.name === "Rambling May") {
        ataBonus += 30;
    }
    if (frameName === "Thirteen" && weapon.name === "Diska of Braveman") {
        ataBonus += 30
    }
    if (frameName === "Crimson Coat" && (
        weapon.name === "Red Slicer" || weapon.name === "Red Dagger" || weapon.name === "Red Saber"
    )) {
        ataBonus += 22
    }
    if (frameName === "Samurai Armor" && weapon.name === "Orotiagito") {
        ataBonus += 20
    }
    if (unitName === "POSS1" && possWeapons.includes(weapon.name)) {
        ataBonus += 30
    }
    if (unitName === "POSS2" && possWeapons.includes(weapon.name)) {
        ataBonus += 60
    }
    if (unitName === "POSS3" && possWeapons.includes(weapon.name)) {
        ataBonus += 90
    }
    if (unitName === "POSS4" && possWeapons.includes(weapon.name)) {
        ataBonus += 120
    }
    return ataBonus;
}

function getFrameDataForWeapon(weapon, className) {
    const animation = weapon.animation;
    const classAnimation = classStats[className].animation;
    let animationFrameData = null;
    let animationSource = ""
    if (!!classSpecificFrameData[className]) {
        animationFrameData = classSpecificFrameData[className][animation];
        animationSource = " (class specific animation)";
    }
    if (!animationFrameData && classAnimation === "female") {
        animationFrameData = femaleFrameData[animation];
        animationSource = " (female animation)";
    }
    if (!animationFrameData) {
        animationFrameData = frameData[animation];
        animationSource = " (base animation)";
    }
    return {animationFrameData, animationSource}
}

function formatAccuracy(num) {
    return Math.floor(num * 100) / 100;
}

function updateTotalFrames() {
    const attack1 = $('#attack1').val();
    const attack2 = $('#attack2').val();
    const attack3 = $('#attack3').val();
    let className = $('#class-select').val();
    let frameDataForWeapon = getFrameDataForWeapon(selectedWeapon, className)
    let totalFrames = getFramesForCombo(attack1, attack2, attack3, frameDataForWeapon.animationFrameData)

    $('#total-frames').text("Total Frames: " + totalFrames + frameDataForWeapon.animationSource)
}

function getFramesForCombo(attack1, attack2, attack3, animationFrameData) {
    let totalFrames = 0
    if (attack1 === "NORMAL") {
        if (attack2 === "NONE" && attack3 === "NONE") {
            totalFrames = animationFrameData.n1;
        } else {
            totalFrames = animationFrameData.n1c;
        }
    } else if ((attack1 === "HEAVY" || attack1 === "SPECIAL")) {
        if (attack2 === "NONE" && attack3 === "NONE") {
            totalFrames = animationFrameData.h1;
        } else {
            totalFrames = animationFrameData.h1c;
        }
    }
    if (attack2 === "NORMAL") {
        if (attack3 === "NONE") {
            totalFrames += animationFrameData.n2;
        } else {
            totalFrames += animationFrameData.n2c;
        }
    } else if ((attack2 === "HEAVY" || attack2 === "SPECIAL")) {
        if (attack3 === "NONE") {
            totalFrames += animationFrameData.h2;
        } else {
            totalFrames += animationFrameData.h2c;
        }
    }

    if (attack3 === "NORMAL") {
        totalFrames += animationFrameData.n3;
    } else if ((attack3 === "HEAVY" || attack3 === "SPECIAL")) {
        totalFrames += animationFrameData.h3;
    }
    return totalFrames;
}
