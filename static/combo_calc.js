'use strict';

const weapons = {
    "Unarmed": {                 minAtp: 0,   maxAtp: 0,   ata: 0,  grind: 0, maxHit: 0, maxAttr: 0, animation: "Fist", special: "None"},

    "Saber": {                   animation: "Saber", minAtp: 40,  maxAtp: 55,  ata: 30, grind: 35},
    "Brand": {                   animation: "Saber", minAtp: 80,  maxAtp: 100, ata: 33, grind: 32},
    "Buster": {                  animation: "Saber", minAtp: 120, maxAtp: 160, ata: 35, grind: 30},
    "Pallasch": {                animation: "Saber", minAtp: 170, maxAtp: 220, ata: 38, grind: 26},
    "Gladius": {                 animation: "Saber", minAtp: 240, maxAtp: 280, ata: 40, grind: 18},
    "Battledore": {              animation: "Saber", minAtp: 1,   maxAtp: 1,   ata: 1,  grind: 0},
    "Red Saber": {               animation: "Saber", minAtp: 450, maxAtp: 489, ata: 51, grind: 78},
    "Lame d'Argent": {           animation: "Saber", minAtp: 430, maxAtp: 465, ata: 40, grind: 35},
    "Lavis Cannon": {            animation: "Saber", minAtp: 730, maxAtp: 750, ata: 54, grind: 0},
    "Excalibur": {               animation: "Saber", minAtp: 900, maxAtp: 950, ata: 60, grind: 0},
    "Galatine": {                animation: "Saber", minAtp: 990, maxAtp: 1260,ata: 77, grind: 9},
    "ES Saber": {                animation: "Saber", minAtp: 150, maxAtp: 150, ata: 50, grind: 250},
    "ES Axe": {                  animation: "Saber", minAtp: 200, maxAtp: 200, ata: 50, grind: 250, maxHit: 0, maxAttr: 0},

    "Sword": {                   animation: "Sword", minAtp: 25,  maxAtp: 60,  ata: 15, grind: 46},
    "Gigush": {                  animation: "Sword", minAtp: 55,  maxAtp: 100, ata: 18, grind: 32},
    "Breaker": {                 animation: "Sword", minAtp: 100, maxAtp: 150, ata: 20, grind: 18},
    "Claymore": {                animation: "Sword", minAtp: 150, maxAtp: 200, ata: 23, grind: 16},
    "Calibur": {                 animation: "Sword", minAtp: 210, maxAtp: 255, ata: 25, grind: 10},
    "Flowen's Sword (3084)": {   animation: "Sword", minAtp: 300, maxAtp: 320, ata: 34, grind: 85},
    "Red Sword": {               animation: "Sword", minAtp: 400, maxAtp: 611, ata: 37, grind: 52},
    "Chain Sawd": {              animation: "Sword", minAtp: 500, maxAtp: 525, ata: 36, grind: 15},
    "Zanba": {                   animation: "Sword", minAtp: 310, maxAtp: 438, ata: 38, grind: 38},
    "Sealed J-Sword": {          animation: "Sword", minAtp: 420, maxAtp: 525, ata: 35, grind: 0},
    "Laconium Axe": {            animation: "Sword", minAtp: 700, maxAtp: 750, ata: 40, grind: 25},
    "Dark Flow": {               animation: "Sword", minAtp: 756, maxAtp: 900, ata: 50, grind: 0, special: "Dark Flow", combo: {   "attack1": "H",   "attack1Hits": 5,   "attack2": "NONE",   "attack3": "NONE" } },
    "Tsumikiri J-Sword": {       animation: "Sword", minAtp: 700, maxAtp: 756, ata: 40, grind: 50, special: "TJS"},
    "TypeSW/J-Sword": {          animation: "Sword", minAtp: 100, maxAtp: 150, ata: 40, grind: 125},
    "ES Sword": {                animation: "Sword", minAtp: 200, maxAtp: 200, ata: 35, grind: 250, maxHit: 0, maxAttr: 0},

    "Dagger": {                  animation: "Dagger", minAtp: 25,  maxAtp: 40,  ata: 20, grind: 65, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Knife": {                   animation: "Dagger", minAtp: 50,  maxAtp: 70,  ata: 22, grind: 50, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Blade": {                   animation: "Dagger", minAtp: 80,  maxAtp: 100, ata: 24, grind: 35, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Edge": {                    animation: "Dagger", minAtp: 105, maxAtp: 130, ata: 26, grind: 25, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Ripper": {                  animation: "Dagger", minAtp: 125, maxAtp: 160, ata: 28, grind: 15, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "S-Beat's Blade": {          animation: "Dagger", minAtp: 210, maxAtp: 220, ata: 35, grind: 15, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "P-Arms' Blade": {           animation: "Dagger", minAtp: 250, maxAtp: 270, ata: 34, grind: 25, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Red Dagger": {              animation: "Dagger", minAtp: 245, maxAtp: 280, ata: 35, grind: 65, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "S-Red's Blade": {           animation: "Dagger", minAtp: 340, maxAtp: 350, ata: 39, grind: 15, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Two Kamui": {               animation: "Dagger", minAtp: 600, maxAtp: 650, ata: 50, grind: 0, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Lavis Blade": {             animation: "Dagger", minAtp: 380, maxAtp: 450, ata: 40, grind: 0, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Daylight Scar": {           animation: "Dagger", minAtp: 500, maxAtp: 550, ata: 48, grind: 25, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "ES Blade": {                animation: "Dagger", minAtp: 10,  maxAtp: 10,  ata: 35, grind: 200, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},

    "Gungnir": {          animation: "Partisan", minAtp: 150, maxAtp: 180,   ata: 32, grind: 10},
    "Vjaya": {            animation: "Partisan", minAtp: 160, maxAtp: 220,   ata: 36, grind: 15, combo: {"attack1": "VJAYA", "attack2": "VJAYA", "attack3": "VJAYA"}},
    "Tyrell's Parasol": { animation: "Partisan", minAtp: 250, maxAtp: 300,   ata: 40, grind: 0},
    "Madam's Umbrella": { animation: "Partisan", minAtp: 210, maxAtp: 280,   ata: 40, grind: 0},
    "Plantain Huge Fan": {animation: "Partisan", minAtp: 265, maxAtp: 300,   ata: 38, grind: 9},
    "Asteron Belt": {     animation: "Partisan", minAtp: 380, maxAtp: 400,   ata: 55, grind: 9},
    "Yunchang": {         animation: "Partisan", minAtp: 300, maxAtp: 350,   ata: 49, grind: 25},
    "ES Partisan": {      animation: "Partisan", minAtp: 10,  maxAtp: 10,    ata: 40, grind: 200},
    "ES Scythe": {        animation: "Partisan", minAtp: 10,  maxAtp: 10,    ata: 40, grind: 180},

    "Diska": {                        animation: "Slicer", minAtp: 85,  maxAtp: 105, ata: 25, grind: 10},
    "Diska of Braveman": {            animation: "Slicer", minAtp: 150, maxAtp: 167, ata: 31, grind: 9},
    "Slicer of Fanatic": {            animation: "Slicer", minAtp: 340, maxAtp: 360, ata: 40, grind: 30},
    "Red Slicer": {                   animation: "Slicer", minAtp: 190, maxAtp: 200, ata: 38, grind: 45},
    "Rainbow Baton": {                animation: "Slicer", minAtp: 300, maxAtp: 320, ata: 40, grind: 24},
    "ES Slicer": {                    animation: "Slicer", minAtp: 10,  maxAtp: 10,  ata: 35, grind: 140},
    "ES J-Cutter": {                  animation: "Slicer", minAtp: 25,  maxAtp: 25,  ata: 35, grind: 150},

    "Demolition Comet": {               animation: "Double Saber", minAtp: 530,  maxAtp: 530,  ata: 38, grind: 25, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},
    "Girasole": {                       animation: "Double Saber", minAtp: 500,  maxAtp: 550,  ata: 50, grind: 0, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},
    "Twin Blaze": {                     animation: "Double Saber", minAtp: 300,  maxAtp: 520,  ata: 40, grind: 9, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},
    "Meteor Cudgel": {                  animation: "Double Saber", minAtp: 300,  maxAtp: 560,  ata: 42, grind: 15, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},
    "Vivienne": {                       animation: "Double Saber", minAtp: 575,  maxAtp: 590,  ata: 49, grind: 50, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},
    "Black King Bar": {                 animation: "Double Saber", minAtp: 590,  maxAtp: 600,  ata: 43, grind: 80, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},
    "Double Cannon": {                  animation: "Double Saber", minAtp: 620,  maxAtp: 650,  ata: 45, grind: 0, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},
    "ES Twin": {                        animation: "Double Saber", minAtp: 50,   maxAtp: 50,   ata: 40, grind: 250, combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}},

    "Toy Hammer": {                  animation: "Katana", minAtp: 1,   maxAtp: 400,  ata: 53, grind: 0},
    "Raikiri": {                     animation: "Katana", minAtp: 550, maxAtp: 560,  ata: 30, grind: 0},
    "Orotiagito": {                  animation: "Katana", minAtp: 750, maxAtp: 800,  ata: 55, grind: 0},

    "Musashi": {             animation: "Twin Sword", minAtp: 330, maxAtp: 350, ata: 35, grind: 40, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},
    "Yamato": {              animation: "Twin Sword", minAtp: 380, maxAtp: 390, ata: 40, grind: 60, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},
    "G-Assassin's Sabers": { animation: "Twin Sword", minAtp: 350, maxAtp: 360, ata: 35, grind: 25, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},
    "Asuka": {               animation: "Twin Sword", minAtp: 560, maxAtp: 570, ata: 50, grind: 30, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},
    "Sange & Yasha": {       animation: "Twin Sword", minAtp: 640, maxAtp: 650, ata: 50, grind: 30, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},
    "Jizai": {               animation: "Twin Sword", minAtp: 800, maxAtp: 810, ata: 55, grind: 40, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},
    "TypeSS/Swords": {       animation: "Twin Sword", minAtp: 150, maxAtp: 150, ata: 45, grind: 125, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},
    "ES Swords": {           animation: "Twin Sword", minAtp: 180, maxAtp: 180, ata: 45, grind: 250, combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}},

    "Raygun": {           animation: "Handgun", minAtp: 150, maxAtp: 180, ata: 35, grind: 15},
    "Master Raven": {     animation: "Master Raven", minAtp: 350, maxAtp: 380, ata: 52, grind: 9, combo: {"attack1Hits": 3, "attack2": "NONE", "attack3": "NONE"}},
    "Last Swan": {        animation: "Last Swan", minAtp: 80,  maxAtp: 90,  ata: 32, grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Heaven Striker": {   animation: "Handgun", minAtp: 550, maxAtp: 600, ata: 55, grind: 20},

    "Laser": {           minAtp: 200, maxAtp: 210, ata: 50, grind: 25},
    "Spread Needle": {   minAtp: 1,   maxAtp: 110, ata: 40, grind: 40},
    "Bringer's Rifle": { minAtp: 330, maxAtp: 370, ata: 63, grind: 9},
    "Frozen Shooter": {  minAtp: 240, maxAtp: 250, ata: 60, grind: 9},
    "Snow Queen": {      minAtp: 330, maxAtp: 350, ata: 60, grind: 18, combo: {"attack2": "NONE", "attack3": "NONE"}},
    "Holy Ray": {        minAtp: 290, maxAtp: 300, ata: 70, grind: 40},
    "ES Rifle": {        minAtp: 10,  maxAtp: 10,  ata: 60, grind: 220},
    "ES Needle": {       minAtp: 10,  maxAtp: 10,  ata: 40, grind: 70},

    "Mechgun": {         minAtp: 2,   maxAtp: 4,   ata: 0,  grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Assault": {         minAtp: 5,   maxAtp: 8,   ata: 3,  grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Repeater": {        minAtp: 5,   maxAtp: 12,  ata: 6,  grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Gatling": {         minAtp: 5,   maxAtp: 16,  ata: 9,  grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Vulcan": {          minAtp: 5,   maxAtp: 20,  ata: 12, grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Samba Maracas": {   minAtp: 5,   maxAtp: 10,  ata: 10, grind: 0, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Rocket Punch": {    minAtp: 50,  maxAtp: 300, ata: 10, grind: 50, combo: {"attack1Hits": 3, "attack2": "NONE", "attack3": "NONE"}},
    "M&A60 Vise": {      minAtp: 15,  maxAtp: 25,  ata: 15, grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "H&S25 Justice": {   minAtp: 15,  maxAtp: 30,  ata: 18, grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "L&K14 Combat": {    minAtp: 15,  maxAtp: 30,  ata: 18, grind: 20, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Twin Psychogun": {  minAtp: 35,  maxAtp: 40,  ata: 23, grind: 0},
    "Red Mechgun": {     minAtp: 50,  maxAtp: 50,  ata: 25, grind: 30, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Yasminkov 9000M": { minAtp: 40,  maxAtp: 80,  ata: 27, grind: 10, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Rage de Feu": {     minAtp: 175, maxAtp: 185, ata: 40, grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Guld Milla": {      minAtp: 180, maxAtp: 200, ata: 30, grind: 9, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Mille Marteaux": {  minAtp: 200, maxAtp: 220, ata: 45, grind: 12, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "Dual Bird": {       minAtp: 200, maxAtp: 210, ata: 22, grind: 21, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "TypeME/Mechgun": {  minAtp: 10,  maxAtp: 10,  ata: 20, grind: 30, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "ES Mechgun": {      minAtp: 10,  maxAtp: 10,  ata: 20, grind: 50, combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}},
    "ES Psychogun": {    minAtp: 10,  maxAtp: 10,  ata: 20, grind: 50},
    "ES Punch": {        minAtp: 10,  maxAtp: 10,  ata: 40, grind: 250, combo: {"attack1Hits": 3, "attack2": "NONE", "attack3": "NONE"}},

    "Shot": {            animation: "Shot", minAtp: 20,  maxAtp: 25, 	  ata: 27, grind: 20},
    "Spread": {          animation: "Shot", minAtp: 30,  maxAtp: 50, 	  ata: 28, grind: 20},
    "Cannon": {          animation: "Shot", minAtp: 40,  maxAtp: 80, 	  ata: 30, grind: 15},
    "Launcher": {        animation: "Shot", minAtp: 50,  maxAtp: 110, 	ata: 31, grind: 15},
    "Arms": {            animation: "Shot", minAtp: 60,  maxAtp: 140, 	ata: 33, grind: 10},
    "L&K38 Combat": {    animation: "L&K38 Combat", minAtp: 150, maxAtp: 250, 	ata: 40, grind: 25, combo: {"attack1Hits": 5, "attack2": "NONE", "attack3": "NONE" }},
    "Rambling May": {    animation: "Shot", minAtp: 360, maxAtp: 450, 	ata: 45, grind: 0, combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}},
    "Baranz Launcher": { animation: "Shot", maxHit: 50, minAtp: 230, maxAtp: 240, 	ata: 40, grind: 30},
    "Dark Meteor": {     animation: "Shot", minAtp: 150, maxAtp: 280, 	ata: 45, grind: 25, combo: {"attack2": "NONE", "attack3": "NONE" }},
    "TypeSH/Shot": {     animation: "Shot", minAtp: 10,  maxAtp: 10, 	  ata: 40, grind: 60},
    "ES Shot": {         animation: "Shot", minAtp: 10,  maxAtp: 10, 	  ata: 40, grind: 125},
    "ES Bazooka": {      animation: "Shot", minAtp: 10,  maxAtp: 10, 	  ata: 40, grind: 250},

    "Cannon Rouge": {     minAtp: 600, maxAtp: 750, ata: 45, grind: 30, combo: {"attack1Hits": 1, "attack2": "NONE", "attack3": "NONE"}},

    "Gal Wind": {         minAtp: 270,  maxAtp: 310, ata: 40, grind: 15, combo: {"attack1Hits": 1, "attack2Hits": 1, "attack3Hits": 3}},
    "Guardianna": {       minAtp: 200,  maxAtp: 280, ata: 40, grind: 9, combo: {"attack1Hits": 1, "attack2Hits": 1, "attack3Hits": 3}},
    "ES Cards": {         minAtp: 10,   maxAtp: 10,  ata: 45, grind: 150, combo: {"attack1Hits": 1, "attack2Hits": 1, "attack3Hits": 3}}
}

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
const classStats = {
    HUmar: {atp: 1397, ata: 200},
    HUnewearl: {atp: 1237, ata: 199},
    HUcast: {atp: 1639, ata: 191},
    HUcaseal: {atp: 1301, ata: 218},
    RAmar: {atp: 1260, ata: 249},
    RAmarl: {atp: 1145, ata: 241},
    RAcast: {atp: 1350, ata: 224},
    RAcaseal: {atp: 1175, ata: 231},
    FOmar: {atp: 1002, ata: 163},
    FOmarl: {atp: 872, ata: 170},
    FOnewm: {atp: 814, ata: 180},
    FOnewearl: {atp: 583, ata: 186}
};
const frameStats = {
    NONE: {atp: 0, ata: 0},
    THIRTEEN: {atp: 0, ata: 0},
    D_PARTS101: {atp: 35, ata: 0},
    SAMURAI: {atp: 0, ata: 0},
    CRIMSON: {atp: 0, ata: 0},
    SWEETHEART1: {atp: 0, ata: 0},
    SWEETHEART2: {atp: 0, ata: 0},
    SWEETHEART3: {atp: 0, ata: 0},
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
    if (special === 'Vjaya' || special === "Dark Flow") {
        return 0.7;
    } else {
        return 0.5;
    }
}

function getDamageModifierForAttackType(attackType, special) {
    if (attackType === 'NORMAL') {
        return 1.0;
    } else if (attackType === 'HEAVY') {
        return 1.89;
    } else if (attackType === 'SPECIAL') {
        return getSpecialDamageModifier(special);
    } else if (attackType === 'NONE') {
        return 0;
    }
}

function getSpecialDamageModifier(special) {
    if (special === 'Charge' || special === 'Spirit' || special === 'Berserk') {
        return 3.32;
    } else if (special === 'Vjaya') {
        return 5.56;
    } else if (special === 'Dark Flow') {
        return 1.89;
    } else {
        return 0;
    }
}

function getEvpModifier(frozen, paralyzed) {
    let modifier = 1.0;
    if (frozen) {
        modifier *= 0.7;
    }
    if (paralyzed) {
        modifier *= 0.85;
    }
    return modifier;
}

function createMonsterRow(enemy, evpModifier, base_ata, snGlitch, atpInput, comboInput) {
    let modified_evp = enemy.evp * evpModifier;
    const a1Type = $('#attack1').val();
    const a2Type = $('#attack2').val();
    const a3Type = $('#attack3').val();
    const special = $('#special-select').val();

    let a1Accuracy = calculateAccuracy(base_ata, a1Type, special, 1.0, modified_evp);
    let a2Accuracy = calculateAccuracy(base_ata, a2Type, special, 1.3, modified_evp);
    let a3Accuracy = calculateAccuracy(base_ata, a3Type, special, 1.69, modified_evp);

    let baseDamage = calculateBaseDamage(atpInput, enemy);
    let damageToUse = atpInput.useMaxDamageRoll ? baseDamage.nMax : baseDamage.nMin;

    // Account for SN glitch - I'm assuming optimistic case where they're able to glitch
    // if the accuracy is better but not if it's worse
    let glitchedA1Accuracy = a1Accuracy;
    if (snGlitch && a2Accuracy > a1Accuracy && a2Type !== 'NONE') {
        glitchedA1Accuracy = a2Accuracy;
    }

    let glitchedA2Accuracy = a2Accuracy;
    if (snGlitch && a3Accuracy > a2Accuracy && a3Type !== 'NONE') {
        glitchedA2Accuracy = a3Accuracy;
    }
    let overallAccuracy = Math.pow((glitchedA1Accuracy * 0.01), comboInput.a1Hits)
        * Math.pow((glitchedA2Accuracy * 0.01), comboInput.a2Hits)
        * Math.pow((a3Accuracy * 0.01), comboInput.a3Hits);
    overallAccuracy *= 100;

    let a1Damage = getDamageModifierForAttackType(a1Type, special) * damageToUse;
    let a2Damage = getDamageModifierForAttackType(a2Type, special) * damageToUse;
    let a3Damage = getDamageModifierForAttackType(a3Type, special) * damageToUse;

    let comboDamage = (a1Damage * comboInput.a1Hits) + (a2Damage * comboInput.a2Hits) + (a3Damage * comboInput.a3Hits);
    let comboKill = comboDamage > enemy.hp;
    let percentDamage = 100 * (comboDamage / enemy.hp);
    if (percentDamage > 100) {
        percentDamage = 100;
    }

    let damageBgColor = comboKill ? 'rgb(61,73,61)' : 'rgb(73,73,61)';

    return $('<tr/>')
        .append($('<th/>', {
            'colspan': 2,
            'data-label': 'monster',
            'text': enemy.name
        }))
        .append($('<td/>', {
            'style': 'padding: 0'
        }).append($('<div>', {
            'style': 'background: rgba(255,150,150,0.1)'
        }).append($('<div>', {
            'style': 'background: ' + damageBgColor + '; padding: 0.78571429em 0.78571429em; width: ' + percentDamage + '%',
            'text': comboDamage.toFixed(0),
            'title': comboDamage.toFixed(0) + '/' + enemy.hp
        }))))
        .append($('<td/>', {
            'data-label': 'accuracy',
            'text': overallAccuracy.toFixed(2) + '%',
            'style': overallAccuracy >= 100.0 ? 'background: rgba(150,255,150,0.1)' : 'background: rgba(255,150,150,0.1)'
        }))
        .append($('<td/>', {
            'data-label': 'a1-accuracy',
            'text': a1Damage.toFixed(0) + ' (' + a1Accuracy.toFixed(0) + '%)',
            'style': a1Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
        }))
        .append($('<td/>', {
            'data-label': 'a2-accuracy',
            'text': a2Damage.toFixed(0) + ' (' + a2Accuracy.toFixed(0) + '%)',
            'style': a2Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
        }))
        .append($('<td/>', {
            'data-label': 'a3-accuracy',
            'text': a3Damage.toFixed(0) + ' (' + a3Accuracy.toFixed(0) + '%)',
            'style': a3Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
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
    let minShiftaAtp = shiftaModifier * atpInput.baseAtp;
    let maxShiftaAtp = (shiftaModifier * atpInput.baseAtp) + ((atpInput.maxAtp - atpInput.minAtp) * shiftaModifier);

    let effectiveMinAtp = atpInput.baseAtp + minWeaponAtp + minShiftaAtp;
    let effectiveMaxAtp = atpInput.baseAtp + maxWeaponAtp + maxShiftaAtp;

    let effectiveDfp = enemy.dfp * (1.0 - zalureModifier);

    let nMin = ((effectiveMinAtp - effectiveDfp) / 5) * 0.9;
    if (nMin < 0) {
        nMin = 0;
    }
    let nMax = ((effectiveMaxAtp - effectiveDfp) / 5) * 0.9;
    if (nMax < 0) {
        nMax = 0;
    }
    return {
        nMin: nMin,
        nMax: nMax
    };
}

function calculateAccuracy(baseAta, attackType, special, comboModifier, totalEvp) {
    if (attackType === 'NONE') {
        return 100;
    }
    if (attackType === 'SPECIAL' && special === 'TJS') {
        return 100;
    }
    let effectiveAta = baseAta * accuracyModifierForAttackType(attackType, special) * comboModifier;
    let accuracy = effectiveAta - (totalEvp * 0.2);
    if (accuracy > 100) {
        accuracy = 100;
    }
    if (accuracy < 0) {
        accuracy = 0;
    }
    return accuracy;
}

function updateDamageTable() {
    let frozen = $('#frozenCheckbox').is(":checked");
    let paralyzed = $('#paralyzedCheckbox').is(":checked");
    let snGlitch = $('#ataGlitch').is(":checked");
    let base_ata = $('#ataInput').val();
    let evpModifier = getEvpModifier(frozen, paralyzed);
    // let opm = $('#opm_checkbox').is(":checked");

    let atpInput = {
        baseAtp: Number($('#atpInput').val()),
        minAtp: Number($('#minAtpInput').val()),
        maxAtp: Number($('#maxAtpInput').val()),
        areaPercent: Number($('#sphereInput').val()),
        useMaxDamageRoll: $('#maxDamageCheckbox').is(":checked"),
        shifta: Number($('#shiftaInput').val()),
        zalure: Number($('#zalureInput').val()),
    };

    let comboInput = {
        a1Hits: Number($('#hits1').val()),
        a2Hits: Number($('#hits2').val()),
        a3Hits: Number($('#hits3').val()),
    }

    var tbody = $('#combo-calc-table tbody');
    tbody.empty();
    const selectedEnemies = $('#enemy-select').val();
    for (let index in selectedEnemies) {
        var enemy = enemies[selectedEnemies[index]];
        var row = createMonsterRow(enemy, evpModifier, base_ata, snGlitch, atpInput, comboInput)
        // row.append($('<th scope="row" colspan="2">' + enemy.name + '</th>'))
        // row.append($('<td style="padding: 0"><div style="height: 100%; background-color: rgba(255,150,150,0.1)"><div style="padding: .75rem; width: 25%; height: 100%; background-color: rgb(73,73,61);">100</div></div></td>'))
        // row.append($('<td style="background-color: rgb(73,73,61)">85%</td>'))
        // row.append($('<td>1000 (85%)</td>'))
        // row.append($('<td>1000 (100%)</td>'))
        // row.append($('<td>1000 (100%)</td>'))
        tbody.append(row)
    }

    // $('#combo-calc-table').append(tbody)
    //
    // let enemyValues = $('#enemy').dropdown('get values');
    // if (!!enemyValues) {
    //     enemyValues.forEach(function(enemyName) {
    //         let enemy = enemiesByName[enemyName];
    //         $('#accuracy_table_body').append(createMonsterRow(enemy, evpModifier, base_ata, snGlitch, atpInput, comboInput));
    //     })
    // }
}

function applyPreset() {
    let playerClass = classStats[$('#playerClass').dropdown('get value')];
    let frameName = $('#frame').dropdown('get value');
    let frame = frameStats[frameName];
    let barrierName = $('#barrier').dropdown('get value');
    let barrier = barrierStats[barrierName];
    let hit = $('#hit_input').val();
    hit = !!hit ? Number(hit) : 0
    let weaponName = $('#weapon').dropdown('get value');
    let weapon = !!weaponName ? weaponsByName[weaponName] : weaponsByName['None'];
    let unitName = $('#unit').dropdown('get value');

    let setEffectAta = getSetEffectAta(weapon, frameName, barrierName, unitName);
    let setEffectAtp = getSetEffectAtp(weapon, frameName, barrierName);

    $('#ata_input').val(playerClass.ata + weapon.ata + hit + frame.ata + barrier.ata + setEffectAta);
    $('#min_atp_input').val(weapon.minAtp + (2 * weapon.grind) + setEffectAtp + barrier.atp + frame.atp);
    $('#max_atp_input').val(weapon.maxAtp + (2 * weapon.grind) + setEffectAtp + barrier.atp + frame.atp);
    $('#base_atp_input').val(playerClass.atp);
}

function applyWeaponStats() {
    let weaponName = $('#weapon').dropdown('get value');
    let weapon = !!weaponName ? weaponsByName[weaponName] : weaponsByName['None'];
    let hits = !!weapon.comboPreset && !!weapon.comboPreset.attack1Hits ? weapon.comboPreset.attack1Hits : 1;
    $('#hits1_input').val(hits).change();
    hits = !!weapon.comboPreset && !!weapon.comboPreset.attack2Hits ? weapon.comboPreset.attack2Hits : 1
    $('#hits2_input').val(hits).change();
    hits = !!weapon.comboPreset && !!weapon.comboPreset.attack3Hits ? weapon.comboPreset.attack3Hits : 1
    $('#hits3_input').val(hits).change();

    if (!!weapon.comboPreset && !!weapon.comboPreset.attack1) {
        $('#attack1_input').val(weapon.comboPreset.attack1).change();
    } else if ($('#attack1').dropdown('get value') === 'NONE') {
        $('#attack1_input').val('N').change();
    }
    if (!!weapon.comboPreset && !!weapon.comboPreset.attack2) {
        $('#attack2_input').val(weapon.comboPreset.attack2).change();
    } else if ($('#attack2').dropdown('get value') === 'NONE') {
        $('#attack2_input').val('N').change();
    }
    if (!!weapon.comboPreset && !!weapon.comboPreset.attack3) {
        $('#attack3_input').val(weapon.comboPreset.attack3).change();
    } else if ($('#attack3').dropdown('get value') === 'NONE') {
        $('#attack3_input').val('N').change();
    }
    applyPreset();
}

function getSetEffectAtp(weapon, frameName, barrierName) {
    let atpBonus = 0;
    if (frameName === "THIRTEEN" && weapon.name === "Diska of Braveman") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.5
    }
    if (frameName === "CRIMSON" && weapon.name === "Red Slicer") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.5
    }
    if (frameName === "SAMURAI" && weapon.name === "Orotiagito") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.3
    }
    if (frameName === "SWEETHEART1") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.15
    }
    if (frameName === "SWEETHEART2") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.2
    }
    if (frameName === "SWEETHEART3") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.25
    }
    return atpBonus;
}

function getSetEffectAta(weapon, frameName, barrierName, unitName) {
    let ataBonus = 0;
    if (barrierName === "SAFETY_HEART" && weapon.name === "Rambling May") {
        ataBonus += 30;
    }
    if (frameName === "THIRTEEN" && weapon.name === "Diska of Braveman") {
        ataBonus += 30
    }
    if (frameName === "CRIMSON" && weapon.name === "Red Slicer") {
        ataBonus += 22
    }
    if (frameName === "SAMURAI" && weapon.name === "Orotiagito") {
        ataBonus += 20
    }
    if (unitName === "POSS" && possWeapons.includes(weapon.name)) {
        ataBonus += 30
    }
    if (unitName === "2POSS" && possWeapons.includes(weapon.name)) {
        ataBonus += 60
    }
    if (unitName === "3POSS" && possWeapons.includes(weapon.name)) {
        ataBonus += 90
    }
    if (unitName === "4POSS" && possWeapons.includes(weapon.name)) {
        ataBonus += 120
    }
    return ataBonus;
}