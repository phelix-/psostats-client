'use strict';

let sortColumn = "";
let sortAscending = null;

const weapons = {
    "Unarmed": {
        name: "Unarmed",
        minAtp: 0,
        maxAtp: 0,
        ata: 0,
        grind: 0,
        maxHit: 0,
        maxAttr: 0,
        animation: "Fist",
        special: "None"
    },

    "Saber": {name: "Saber", animation: "Saber", minAtp: 40, maxAtp: 55, ata: 30, grind: 35},
    "Brand": {name: "Brand", animation: "Saber", minAtp: 80, maxAtp: 100, ata: 33, grind: 32},
    "Buster": {name: "Buster", animation: "Saber", minAtp: 120, maxAtp: 160, ata: 35, grind: 30},
    "Pallasch": {name: "Pallasch", animation: "Saber", minAtp: 170, maxAtp: 220, ata: 38, grind: 26},
    "Gladius": {name: "Gladius", animation: "Saber", minAtp: 240, maxAtp: 280, ata: 40, grind: 18},
    "Battledore": {name: "Battledore", animation: "Saber", minAtp: 1, maxAtp: 1, ata: 1, grind: 0},
    "Red Saber": {name: "Red Saber", animation: "Saber", minAtp: 450, maxAtp: 489, ata: 51, grind: 78},
    "Lavis Cannon": {
        name: "Lavis Cannon",
        animation: "Saber",
        minAtp: 730,
        maxAtp: 750,
        ata: 54,
        grind: 0,
        special: "Lavis"
    },
    "Excalibur": {
        name: "Excalibur",
        animation: "Saber",
        minAtp: 900,
        maxAtp: 950,
        ata: 60,
        grind: 0,
        special: "Berserk"
    },
    "Galatine": {name: "Galatine", animation: "Saber", minAtp: 990, maxAtp: 1260, ata: 77, grind: 9, special: "Spirit"},
    "ES Saber": {
        name: "ES Saber", animation: "Saber", minAtp: 150, maxAtp: 150, ata: 50, grind: 250, maxHit: 0, maxAttr: 0,
    },
    "ES Axe": {
        name: "ES Axe",
        animation: "Saber",
        minAtp: 200,
        maxAtp: 200,
        ata: 50,
        grind: 250,
        maxHit: 0,
        maxAttr: 0
    },

    "Sword": {name: "Sword", animation: "Sword", minAtp: 25, maxAtp: 60, ata: 15, grind: 46},
    "Gigush": {name: "Gigush", animation: "Sword", minAtp: 55, maxAtp: 100, ata: 18, grind: 32},
    "Breaker": {name: "Breaker", animation: "Sword", minAtp: 100, maxAtp: 150, ata: 20, grind: 18},
    "Claymore": {name: "Claymore", animation: "Sword", minAtp: 150, maxAtp: 200, ata: 23, grind: 16},
    "Calibur": {name: "Calibur", animation: "Sword", minAtp: 210, maxAtp: 255, ata: 25, grind: 10},
    "Flowen's Sword (3084)": {
        name: "Flowen's Sword (3084)",
        animation: "Sword",
        minAtp: 300,
        maxAtp: 320,
        ata: 34,
        grind: 85, special: "Spirit"
    },
    "Red Sword": {
        name: "Red Sword",
        animation: "Sword",
        minAtp: 400,
        maxAtp: 611,
        ata: 37,
        grind: 52,
        special: "Arrest"
    },
    "Chain Sawd": {
        name: "Chain Sawd",
        animation: "Sword",
        minAtp: 500,
        maxAtp: 525,
        ata: 36,
        grind: 15,
        special: "Gush"
    },
    "Zanba": {name: "Zanba", animation: "Sword", minAtp: 310, maxAtp: 438, ata: 38, grind: 38, special: "Berserk"},
    "Sealed J-Sword": {
        name: "Sealed J-Sword",
        animation: "Sword",
        minAtp: 420,
        maxAtp: 525,
        ata: 35,
        grind: 0,
        special: "Hell"
    },
    "Laconium Axe": {name: "Laconium Axe", animation: "Sword", minAtp: 700, maxAtp: 750, ata: 40, grind: 25},
    "Dark Flow": {
        name: "Dark Flow",
        animation: "Sword",
        minAtp: 756,
        maxAtp: 900,
        ata: 50,
        grind: 0,
        special: "Dark Flow",
        combo: {"attack1": "SPECIAL", "attack1Hits": 5, "attack2": "NONE", "attack3": "NONE"}
    },
    "Tsumikiri J-Sword": {
        name: "Tsumikiri J-Sword",
        animation: "Sword",
        minAtp: 700,
        maxAtp: 756,
        ata: 40,
        grind: 50,
        special: "TJS"
    },
    "TypeSW/J-Sword": {
        name: "TypeSW/J-Sword",
        animation: "Sword",
        minAtp: 100,
        maxAtp: 150,
        ata: 40,
        grind: 125,
        special: "Spirit"
    },
    "ES Sword": {
        name: "ES Sword",
        animation: "Sword",
        minAtp: 200,
        maxAtp: 200,
        ata: 35,
        grind: 250,
        maxHit: 0,
        maxAttr: 0
    },

    "Dagger": {
        name: "Dagger",
        animation: "Dagger",
        minAtp: 25,
        maxAtp: 40,
        ata: 20,
        grind: 65,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Knife": {
        name: "Knife",
        animation: "Dagger",
        minAtp: 50,
        maxAtp: 70,
        ata: 22,
        grind: 50,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Blade": {
        name: "Blade",
        animation: "Dagger",
        minAtp: 80,
        maxAtp: 100,
        ata: 24,
        grind: 35,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Edge": {
        name: "Edge",
        animation: "Dagger",
        minAtp: 105,
        maxAtp: 130,
        ata: 26,
        grind: 25,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Ripper": {
        name: "Ripper",
        animation: "Dagger",
        minAtp: 125,
        maxAtp: 160,
        ata: 28,
        grind: 15,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "S-Beat's Blade": {
        name: "S-Beat's Blade",
        animation: "Dagger",
        minAtp: 210,
        maxAtp: 220,
        ata: 35,
        grind: 15,
        maxHit: 50,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}, special: "None"
    },
    "P-Arms' Blade": {
        name: "P-Arms' Blade",
        animation: "Dagger",
        minAtp: 250,
        maxAtp: 270,
        ata: 34,
        grind: 25,
        maxHit: 50,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Red Dagger": {
        name: "Red Dagger",
        animation: "Dagger",
        minAtp: 245,
        maxAtp: 280,
        ata: 35,
        grind: 65,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "S-Red's Blade": {
        name: "S-Red's Blade",
        animation: "Dagger",
        minAtp: 340,
        maxAtp: 350,
        ata: 39,
        grind: 15,
        maxHit: 50,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Two Kamui": {
        name: "Two Kamui",
        animation: "Dagger",
        minAtp: 600,
        maxAtp: 650,
        ata: 50,
        grind: 0,
        maxHit: 0,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}, special: "None"
    },
    "Lavis Blade": {
        name: "Lavis Blade",
        animation: "Dagger",
        minAtp: 380,
        maxAtp: 450,
        ata: 40,
        grind: 0,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}, special: "Lavis"
    },
    "Daylight Scar": {
        name: "Daylight Scar",
        animation: "Dagger",
        minAtp: 500,
        maxAtp: 550,
        ata: 48,
        grind: 25,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}, special: "Berserk"
    },
    "ES Blade": {
        name: "ES Blade",
        animation: "Dagger",
        minAtp: 10,
        maxAtp: 10,
        ata: 35,
        grind: 200,
        maxHit: 0, maxAttr: 0,
        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}
    },

    "Gungnir": {name: "Gungnir", animation: "Partisan", minAtp: 150, maxAtp: 180, ata: 32, grind: 10},
    "Vjaya": {name: "Vjaya", animation: "Partisan", minAtp: 160, maxAtp: 220, ata: 36, grind: 15, special: "Vjaya"},
    "Tyrell's Parasol": {name: "Tyrell's Parasol", animation: "Partisan", minAtp: 250, maxAtp: 300, ata: 40, grind: 0},
    "Madam's Umbrella": {
        name: "Madam's Umbrella",
        animation: "Partisan",
        minAtp: 210,
        maxAtp: 280,
        ata: 40,
        grind: 0,
        special: "Berserk"
    },
    "Plantain Huge Fan": {
        name: "Plantain Huge Fan",
        animation: "Partisan",
        minAtp: 265,
        maxAtp: 300,
        ata: 38,
        grind: 9
    },
    "Asteron Belt": {name: "Asteron Belt", animation: "Partisan", minAtp: 380, maxAtp: 400, ata: 55, grind: 9},
    "Yunchang": {
        name: "Yunchang",
        animation: "Partisan",
        minAtp: 300,
        maxAtp: 350,
        ata: 49,
        grind: 25,
        special: "Berserk"
    },
    "ES Partisan": {
        name: "ES Partisan", animation: "Partisan", minAtp: 10, maxAtp: 10, ata: 40, grind: 200, maxHit: 0, maxAttr: 0,
    },
    "ES Scythe": {
        name: "ES Scythe", animation: "Partisan", minAtp: 10, maxAtp: 10, ata: 40, grind: 180, maxHit: 0, maxAttr: 0,
    },

    "Diska": {name: "Diska", animation: "Slicer", minAtp: 85, maxAtp: 105, ata: 25, grind: 10},
    "Diska of Braveman": {
        name: "Diska of Braveman",
        animation: "Slicer",
        minAtp: 150,
        maxAtp: 167,
        ata: 31,
        grind: 9,
        special: "Berserk"
    },
    "Slicer of Fanatic": {
        name: "Slicer of Fanatic",
        animation: "Slicer",
        minAtp: 340,
        maxAtp: 360,
        ata: 40,
        grind: 30,
        special: "Demon's"
    },
    "Red Slicer": {name: "Red Slicer", animation: "Slicer", minAtp: 190, maxAtp: 200, ata: 38, grind: 45},
    "Rainbow Baton": {name: "Rainbow Baton", animation: "Slicer", minAtp: 300, maxAtp: 320, ata: 40, grind: 24},
    "ES Slicer": {
        name: "ES Slicer", animation: "Slicer", minAtp: 10, maxAtp: 10, ata: 35, grind: 140, maxHit: 0, maxAttr: 0,
    },
    "ES J-Cutter": {
        name: "ES J-Cutter", animation: "Slicer", minAtp: 25, maxAtp: 25, ata: 35, grind: 150, maxHit: 0, maxAttr: 0,
    },

    "Demolition Comet": {
        name: "Demolition Comet",
        animation: "Double Saber",
        minAtp: 530,
        maxAtp: 530,
        ata: 38,
        grind: 25,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}
    },
    "Girasole": {
        name: "Girasole",
        animation: "Double Saber",
        minAtp: 500,
        maxAtp: 550,
        ata: 50,
        grind: 0,
        maxHit: 0,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}, special: "Lavis"
    },
    "Twin Blaze": {
        name: "Twin Blaze",
        animation: "Double Saber",
        minAtp: 300,
        maxAtp: 520,
        ata: 40,
        grind: 9,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}, special: "None"
    },
    "Meteor Cudgel": {
        name: "Meteor Cudgel",
        animation: "Double Saber",
        minAtp: 300,
        maxAtp: 560,
        ata: 42,
        grind: 15,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}
    },
    "Vivienne": {
        name: "Vivienne",
        animation: "Double Saber",
        minAtp: 575,
        maxAtp: 590,
        ata: 49,
        grind: 50,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}
    },
    "Black King Bar": {
        name: "Black King Bar",
        animation: "Double Saber",
        minAtp: 590,
        maxAtp: 600,
        ata: 43,
        grind: 80,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}
    },
    "Double Cannon": {
        name: "Double Cannon",
        animation: "Double Saber",
        minAtp: 620,
        maxAtp: 650,
        ata: 45,
        grind: 0,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3},
        special: "Lavis"
    },
    "ES Twin": {
        name: "ES Twin",
        animation: "Double Saber",
        minAtp: 50,
        maxAtp: 50,
        ata: 40,
        grind: 250,
        maxHit: 0, maxAttr: 0,
        combo: {"attack1Hits": 2, "attack2Hits": 1, "attack3Hits": 3}
    },

    "Toy Hammer": {name: "Toy Hammer", animation: "Katana", minAtp: 1, maxAtp: 400, ata: 53, grind: 0},
    "Raikiri": {name: "Raikiri", animation: "Katana", minAtp: 550, maxAtp: 560, ata: 30, grind: 0},
    "Orotiagito": {
        name: "Orotiagito",
        animation: "Katana",
        minAtp: 750,
        maxAtp: 800,
        ata: 55,
        grind: 0,
        maxHit: 0,
        special: "Lavis"
    },

    "Musashi": {
        name: "Musashi",
        animation: "Twin Sword",
        minAtp: 330,
        maxAtp: 350,
        ata: 35,
        grind: 40,
        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}, special: "Berserk"
    },
    "Yamato": {
        name: "Yamato",
        animation: "Twin Sword",
        minAtp: 380,
        maxAtp: 390,
        ata: 40,
        grind: 60,
        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}, special: "Blizzard"
    },
    "G-Assassin's Sabers": {
        name: "G-Assassin's Sabers",
        animation: "Twin Sword",
        minAtp: 350,
        maxAtp: 360,
        ata: 35,
        grind: 25,
        maxHit: 50,
        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Asuka": {
        name: "Asuka",
        animation: "Twin Sword",
        minAtp: 560,
        maxAtp: 570,
        ata: 50,
        grind: 30,
        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Sange & Yasha": {
        name: "Sange & Yasha",
        animation: "Twin Sword",
        minAtp: 640,
        maxAtp: 650,
        ata: 50,
        grind: 30,
        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}
    },
    "Jizai": {
        name: "Jizai",
        animation: "Twin Sword",
        minAtp: 800,
        maxAtp: 810,
        ata: 55,
        grind: 40,
        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}, special: "Hell"
    },
    "TypeSS/Swords": {
        name: "TypeSS/Swords",
        animation: "Twin Sword",
        minAtp: 150,
        maxAtp: 150,
        ata: 45,
        grind: 125,
        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}
    },
    "ES Swords": {
        name: "ES Swords",
        animation: "Twin Sword",
        minAtp: 180,
        maxAtp: 180,
        ata: 45,
        grind: 250,
        maxHit: 0, maxAttr: 0,

        combo: {"attack1Hits": 1, "attack2Hits": 2, "attack3Hits": 2}
    },

    "Raygun": {name: "Raygun", animation: "Handgun", minAtp: 150, maxAtp: 180, ata: 35, grind: 15},
    "Master Raven": {
        name: "Master Raven",
        animation: "Master Raven",
        minAtp: 350,
        maxAtp: 380,
        ata: 52,
        grind: 9,
        maxHit: 0,
        combo: {"attack1Hits": 3, "attack2": "NONE", "attack3": "NONE"}
    },
    "Last Swan": {
        name: "Last Swan",
        animation: "Last Swan",
        minAtp: 80,
        maxAtp: 90,
        ata: 32,
        grind: 9,
        maxHit: 0,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Heaven Striker": {
        name: "Heaven Striker",
        animation: "Handgun",
        minAtp: 550,
        maxAtp: 600,
        ata: 55,
        grind: 20,
        special: "Berserk"
    },

    "Laser": {name: "Laser", animation: "Rifle", minAtp: 200, maxAtp: 210, ata: 50, grind: 25},
    "Spread Needle": {name: "Spread Needle", animation: "Rifle", minAtp: 1, maxAtp: 110, ata: 40, grind: 40, special: "Seize"},
    "Bringer's Rifle": {
        name: "Bringer's Rifle",
        animation: "Rifle",
        minAtp: 330,
        maxAtp: 370,
        ata: 63,
        grind: 9,
        special: "Demon's",
        maxHit: 50
    },
    "Frozen Shooter": {name: "Frozen Shooter", animation: "Rifle", minAtp: 240, maxAtp: 250, ata: 60, grind: 9, special: "Lavis"},
    "Snow Queen": {
        name: "Snow Queen",
        animation: "Rifle",
        minAtp: 330,
        maxAtp: 350,
        ata: 60,
        grind: 18,
        combo: {"attack2": "NONE", "attack3": "NONE"}, special: "Lavis"
    },
    "Holy Ray": {name: "Holy Ray", animation: "Rifle", minAtp: 290, maxAtp: 300, ata: 70, grind: 40, special: "Arrest"},
    "ES Rifle": {
        name: "ES Rifle", animation: "Rifle", minAtp: 10, maxAtp: 10, ata: 60, grind: 220, maxHit: 0, maxAttr: 0,
    },
    "ES Needle": {
        name: "ES Needle", animation: "Rifle", minAtp: 10, maxAtp: 10, ata: 40, grind: 70, maxHit: 0, maxAttr: 0,
    },

    "Mechgun": {
        name: "Mechgun",
        animation: "Mechgun",
        minAtp: 2,
        maxAtp: 4,
        ata: 0,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Assault": {
        name: "Assault",
        animation: "Mechgun",
        minAtp: 5,
        maxAtp: 8,
        ata: 3,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Repeater": {
        name: "Repeater",
        animation: "Mechgun",
        minAtp: 5,
        maxAtp: 12,
        ata: 6,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Gatling": {
        name: "Gatling",
        animation: "Mechgun",
        minAtp: 5,
        maxAtp: 16,
        ata: 9,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Vulcan": {
        name: "Vulcan",
        animation: "Mechgun",
        minAtp: 5,
        maxAtp: 20,
        ata: 12,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Rocket Punch": {
        name: "Rocket Punch",
        animation: "Mechgun",
        minAtp: 50,
        maxAtp: 300,
        ata: 10,
        grind: 50,
        combo: {"attack1Hits": 3, "attack2": "NONE", "attack3": "NONE"}
    },
    "M&A60 Vise": {
        name: "M&A60 Vise",
        animation: "Mechgun",
        minAtp: 15,
        maxAtp: 25,
        ata: 15,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}, special: "Berserk"
    },
    "Red Mechgun": {
        name: "Red Mechgun",
        animation: "Mechgun",
        minAtp: 50,
        maxAtp: 50,
        ata: 25,
        grind: 30,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Yasminkov 9000M": {
        name: "Yasminkov 9000M",
        animation: "Mechgun",
        minAtp: 40,
        maxAtp: 80,
        ata: 27,
        grind: 10,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Rage de Feu": {
        name: "Rage de Feu",
        animation: "Mechgun",
        minAtp: 175,
        maxAtp: 185,
        ata: 40,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Guld Milla": {
        name: "Guld Milla",
        animation: "Mechgun",
        minAtp: 180,
        maxAtp: 200,
        ata: 30,
        grind: 9,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Mille Marteaux": {
        name: "Mille Marteaux",
        animation: "Mechgun",
        minAtp: 200,
        maxAtp: 220,
        ata: 45,
        grind: 12,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "Dual Bird": {
        name: "Dual Bird",
        animation: "Mechgun",
        minAtp: 200,
        maxAtp: 210,
        ata: 22,
        grind: 21,
        maxHit: 0,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "TypeME/Mechgun": {
        name: "TypeME/Mechgun",
        animation: "Mechgun",
        minAtp: 10,
        maxAtp: 10,
        ata: 20,
        grind: 30,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "ES Mechgun": {
        name: "ES Mechgun",
        animation: "Mechgun",
        minAtp: 10,
        maxAtp: 10,
        ata: 20,
        grind: 50,
        maxHit: 0, maxAttr: 0,
        combo: {"attack1Hits": 3, "attack2Hits": 3, "attack3Hits": 3}
    },
    "ES Psychogun": {
        name: "ES Psychogun", animation: "Mechgun", minAtp: 10, maxAtp: 10, ata: 20, grind: 50, maxHit: 0, maxAttr: 0,
    },
    "ES Punch": {
        name: "ES Punch",
        animation: "Mechgun",
        minAtp: 10,
        maxAtp: 10,
        ata: 40,
        grind: 250,
        maxHit: 0, maxAttr: 0,
        combo: {"attack1Hits": 3, "attack2": "NONE", "attack3": "NONE"}
    },

    "Shot": {name: "Shot", animation: "Shot", minAtp: 20, maxAtp: 25, ata: 27, grind: 20},
    "Spread": {name: "Spread", animation: "Shot", minAtp: 30, maxAtp: 50, ata: 28, grind: 20},
    "Cannon": {name: "Cannon", animation: "Shot", minAtp: 40, maxAtp: 80, ata: 30, grind: 15},
    "Launcher": {name: "Launcher", animation: "Shot", minAtp: 50, maxAtp: 110, ata: 31, grind: 15},
    "Arms": {name: "Arms", animation: "Shot", minAtp: 60, maxAtp: 140, ata: 33, grind: 10},
    "L&K38 Combat": {
        name: "L&K38 Combat",
        animation: "L&K38 Combat",
        minAtp: 150,
        maxAtp: 250,
        ata: 40,
        grind: 25,
        combo: {"attack1Hits": 5, "attack2": "NONE", "attack3": "NONE"}, special: "Burning"
    },
    "Rambling May": {
        name: "Rambling May",
        animation: "Shot",
        minAtp: 360,
        maxAtp: 450,
        ata: 45,
        grind: 0, maxHit: 0,

        combo: {"attack1Hits": 2, "attack2Hits": 2, "attack3Hits": 2}, special: "Chaos"
    },
    "Baranz Launcher": {
        name: "Baranz Launcher",
        animation: "Shot",
        maxHit: 50,
        minAtp: 230,
        maxAtp: 240,
        ata: 40,
        grind: 30
    },
    "Dark Meteor": {
        name: "Dark Meteor",
        animation: "Shot",
        minAtp: 150,
        maxAtp: 280,
        ata: 45,
        grind: 25,
        combo: {"attack2": "NONE", "attack3": "NONE"}, special: "Dark Flow"
    },
    "TypeSH/Shot": {name: "TypeSH/Shot", animation: "Shot", minAtp: 10, maxAtp: 10, ata: 40, grind: 60},
    "ES Shot": {
        name: "ES Shot", animation: "Shot", minAtp: 10, maxAtp: 10, ata: 40, grind: 125, maxHit: 0, maxAttr: 0,
    },
    "ES Bazooka": {
        name: "ES Bazooka", animation: "Shot", minAtp: 10, maxAtp: 10, ata: 40, grind: 250, maxHit: 0, maxAttr: 0,
    },

    "ES Launcher": {
        name: "ES Launcher",
        animation: "Launcher",
        minAtp: 10,
        maxAtp: 10,
        ata: 40,
        grind: 180,
        maxHit: 0,
        maxAttr: 0,
        special: "Berserk"
    },

    "Cannon Rouge": {
        name: "Cannon Rouge",
        animation: "Launcher",
        minAtp: 600,
        maxAtp: 750,
        ata: 45,
        grind: 30,
        combo: {"attack1Hits": 1, "attack2": "NONE", "attack3": "NONE"}
    },

    "Gal Wind": {
        name: "Gal Wind",
        animation: "Card",
        minAtp: 270,
        maxAtp: 310,
        ata: 40,
        grind: 15,
        maxHit: 50,
        combo: {"attack1Hits": 1, "attack2Hits": 1, "attack3Hits": 3}
    },
    "Guardianna": {
        name: "Guardianna",
        animation: "Card",
        minAtp: 200,
        maxAtp: 280,
        ata: 40,
        grind: 9,
        combo: {"attack1Hits": 1, "attack2Hits": 1, "attack3Hits": 3}
    },
    "ES Cards": {
        name: "ES Cards",
        animation: "Card",
        minAtp: 10,
        maxAtp: 10,
        ata: 45,
        grind: 150,
        maxHit: 0, maxAttr: 0,
        combo: {"attack1Hits": 1, "attack2Hits": 1, "attack3Hits": 3}
    }
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
    HUmar: {animation: "male", atp: 1397, ata: 200},
    HUnewearl: {animation: "female", atp: 1237, ata: 199},
    HUcast: {animation: "male", atp: 1639, ata: 191},
    HUcaseal: {animation: "female", atp: 1301, ata: 218},
    RAmar: {animation: "male", atp: 1260, ata: 249},
    RAmarl: {animation: "female", atp: 1145, ata: 241},
    RAcast: {animation: "male", atp: 1350, ata: 224},
    RAcaseal: {animation: "female", atp: 1175, ata: 231},
    FOmar: {animation: "male", atp: 1002, ata: 163},
    FOmarl: {animation: "female", atp: 872, ata: 170},
    FOnewm: {animation: "male", atp: 814, ata: 180},
    FOnewearl: {animation: "female", atp: 583, ata: 186}
};
const frames = {
    NONE: {atp: 0, ata: 0},
    THIRTEEN: {atp: 0, ata: 0},
    D_PARTS101: {atp: 35, ata: 0},
    SAMURAI: {atp: 0, ata: 0},
    CRIMSON_COAT: {atp: 0, ata: 0},
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
    if (special === 'Vjaya' || special === "Dark Flow" || special === "Lavis") {
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
    } else if (special === 'Dark Flow' || special === 'Lavis' || special === 'TJS') {
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

function createMonsterRow(
    special, autoCombo, weapon, enemy,
    evpModifier, base_ata, snGlitch, atpInput, comboInput
) {
    let modified_evp = enemy.evp * evpModifier;

    let baseDamage = calculateBaseDamage(atpInput, enemy);
    let damageToUse = atpInput.useMaxDamageRoll ? baseDamage.nMax : baseDamage.nMin;
    if (autoCombo) {
        return generateAutoCombo(special, weapon, enemy, modified_evp, base_ata, snGlitch, damageToUse, atpInput, comboInput)
    }

    let accuracyResult = getAccuracyForCombo(base_ata, comboInput, special, modified_evp, snGlitch);
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
        a1Damage: comboDamage.a1Damage,
        a1Type: comboInput.a1Type,
        a1Accuracy: accuracyResult.a1Accuracy,
        a2Type: comboInput.a2Type,
        a2Damage: comboDamage.a2Damage,
        a2Accuracy: accuracyResult.a2Accuracy,
        a3Type: comboInput.a3Type,
        a3Damage: comboDamage.a3Damage,
        a3Accuracy: accuracyResult.a3Accuracy,
    }
}

function getAccuracyForCombo(base_ata, comboInput, special, modified_evp, snGlitch) {
    let a1Accuracy = calculateAccuracy(base_ata, comboInput.a1Type, special, 1.0, modified_evp);
    let a2Accuracy = calculateAccuracy(base_ata, comboInput.a2Type, special, 1.3, modified_evp);
    let a3Accuracy = calculateAccuracy(base_ata, comboInput.a3Type, special, 1.69, modified_evp);

    // Account for SN glitch - I'm assuming optimistic case where they're able to glitch
    // if the accuracy is better but not if it's worse
    let glitchedA1Accuracy = a1Accuracy;
    if (snGlitch && a2Accuracy > a1Accuracy && comboInput.a2Type !== 'NONE') {
        glitchedA1Accuracy = a2Accuracy;
    }

    let glitchedA2Accuracy = a2Accuracy;
    if (snGlitch && a3Accuracy > a2Accuracy && comboInput.a3Type !== 'NONE') {
        glitchedA2Accuracy = a3Accuracy;
    }
    let overallAccuracy = Math.pow((glitchedA1Accuracy * 0.01), comboInput.a1Hits)
        * Math.pow((glitchedA2Accuracy * 0.01), comboInput.a2Hits)
        * Math.pow((a3Accuracy * 0.01), comboInput.a3Hits);
    overallAccuracy *= 100;
    return {overallAccuracy, a1Accuracy, a2Accuracy, a3Accuracy};
}

function getDamageForCombo(enemyHp, atpInput, comboInput, special, baseDamage) {
    let total = 0;
    let a1Damage = 0;
    let a2Damage = 0;
    let a3Damage = 0;
    let demonsMultiplier = getDemonsModifier(atpInput.playerClass);
    if (comboInput.a1Type === "SPECIAL" && special === "Demon's") {
        for (let i = 0; i < comboInput.a1Hits; i++) {
            a1Damage = (enemyHp - total) * demonsMultiplier;
            total += a1Damage;
        }
    } else {
        a1Damage = getDamageModifierForAttackType(comboInput.a1Type, special) * baseDamage;
        total += a1Damage * comboInput.a1Hits;
    }
    if (comboInput.a2Type === "SPECIAL" && special === "Demon's") {
        for (let i = 0; i < comboInput.a2Hits; i++) {
            a2Damage = (enemyHp - total) * demonsMultiplier;
            total += a2Damage;
        }
    } else {
        a2Damage = getDamageModifierForAttackType(comboInput.a2Type, special) * baseDamage;
        total += a2Damage * comboInput.a2Hits;
    }
    if (comboInput.a3Type === "SPECIAL" && special === "Demon's") {
        for (let i = 0; i < comboInput.a3Hits; i++) {
            a3Damage = (enemyHp - total) * demonsMultiplier;
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

function generateAutoCombo(
    special, weapon, enemy, modified_evp, base_ata, snGlitch, baseDamage, atpInput, comboInput
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
                let accuracyResult = getAccuracyForCombo(base_ata, comboInput, special, modified_evp, snGlitch);
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

function appendMosterRow(rowEntry) {
    let comboKill = rowEntry.comboDamage > rowEntry.hp;
    let damageBgColor = comboKill ? 'rgb(61,73,61)' : 'rgb(73,73,61)';

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
            'style': 'background: ' + damageBgColor + '; padding: 0.78571429em 0.78571429em; width: ' + rowEntry.percentDamage + '%',
            'text': rowEntry.comboDamage.toFixed(0),
            'title': rowEntry.comboDamage.toFixed(0) + '/' + rowEntry.hp
        }))))
        .append($('<td/>', {
            'data-label': 'accuracy',
            'text': rowEntry.overallAccuracy.toFixed(2) + '%',
            'style': rowEntry.overallAccuracy >= 100.0 ? 'background: rgba(150,255,150,0.1)' : 'background: rgba(255,150,150,0.1)'
        }))
        .append($('<td/>', {
            'data-label': 'a1-accuracy',
            'text': rowEntry.a1Damage.toFixed(0) + ' (' + rowEntry.a1Accuracy.toFixed(0) + '%)',
            'style': rowEntry.a1Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
        }))
        .append($('<td/>', {
            'data-label': 'a2-accuracy',
            'text': rowEntry.a2Damage.toFixed(0) + ' (' + rowEntry.a2Accuracy.toFixed(0) + '%)',
            'style': rowEntry.a2Type === 'NONE' ? 'color: rgba(255,255,255,0.3)' : 'color: rgba(255,255,255,0.9)'
        }))
        .append($('<td/>', {
            'data-label': 'a3-accuracy',
            'text': rowEntry.a3Damage.toFixed(0) + ' (' + rowEntry.a3Accuracy.toFixed(0) + '%)',
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
        baseAtp: Number($('#atpInput').val()),
        minAtp: Number($('#minAtpInput').val()),
        maxAtp: Number($('#maxAtpInput').val()),
        areaPercent: Number($('#sphereInput').val()),
        useMaxDamageRoll: $('#maxDamageCheckbox').is(":checked"),
        shifta: Number($('#shiftaInput').val()),
        zalure: Number($('#zalureInput').val()),
    };

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
            evpModifier, base_ata, snGlitch, atpInput, comboInput
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
        tbody.append(appendMosterRow(rows[index]))
    }
}

function getSetEffectAtp(weapon, frameName, barrierName) {
    let atpBonus = 0;
    if (frameName === "THIRTEEN" && weapon.name === "Diska of Braveman") {
        atpBonus += (weapon.minAtp + (2 * weapon.grind)) * 0.5
    }
    if (frameName === "CRIMSON_COAT" && (
        weapon.name === "Red Slicer" || weapon.name === "Red Dagger" || weapon.name === "Red Saber"
    )) {
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
    if (barrierName === "Safety Heart" && weapon.name === "Rambling May") {
        ataBonus += 30;
    }
    if (frameName === "THIRTEEN" && weapon.name === "Diska of Braveman") {
        ataBonus += 30
    }
    if (frameName === "CRIMSON_COAT" && (
        weapon.name === "Red Slicer" || weapon.name === "Red Dagger" || weapon.name === "Red Saber"
    )) {
        ataBonus += 22
    }
    if (frameName === "SAMURAI" && weapon.name === "Orotiagito") {
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