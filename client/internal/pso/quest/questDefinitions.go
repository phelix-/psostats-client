package quest

func getAllQuests() []Quest {
	return []Quest{
		// ---------------------------------------------------------------------------------
		//     Episode 1
		// ---------------------------------------------------------------------------------

		// ---- Government ----
		// 1-1
		// 1-2
		// 1-3
		// 2-1
		// 2-2
		// 2-3
		// 2-4
		// 3-1
		// 3-2
		// 3-3
		// 4-1
		// 4-2
		// 4-3
		// 4-4
		// 4-5

		// ---- Solo Only ----
		// Gallon's Plan
		// TTP
		// Good Luck!
		// AOL CUP -Sunset Base-
		// Knight of Coral

		// ---- Side Story ----
		// Battle Training
		// Magnitude of Metal
		// Claiming a Stake
		// The Value of Money
		// Journalistic Pursuit
		// The Fake in Yellow
		// Native Research
		// Forest of Sorrow
		// Gran Squall
		// Addicting Food
		// The Lost Bride
		// Waterfall Tears
		// Black Paper
		// Secret Delivery
		// Soul of a Blacksmith
		// Letter from Lionel
		// The Grave's Butler
		// Knowing One's Heart
		// The Retired Hunter
		// Dr. Osto's Research
		// Unsealed Door
		// Gallon's Treachery
		// Doc's Secret Plan
		// Seek my Master
		// From the Depths
		// Central Dome Fire Swirl

		// ---- Extermination ----
		{Episode: 1, Name: "Mop-up Operation #1", Start: register(0), End: register(254)},
		{Episode: 1, Name: "Mop-up Operation #2", Start: register(0), End: register(254)},
		{Episode: 1, Name: "Mop-up Operation #3", Start: register(0), End: register(254)},
		{Episode: 1, Name: "Mop-up Operation #4", Start: register(0), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #1", Start: register(210), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #2", Start: register(210), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #3", Start: register(210), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #4", Start: register(210), End: register(254)},
		{Episode: 1, Name: "Endless Nightmare #1", Start: warpIn(), End: register(30)},
		{Episode: 1, Name: "Endless Nightmare #2", Start: warpIn(), End: register(30)},
		{Episode: 1, Name: "Endless Nightmare #3", Start: warpIn(), End: register(30)},
		{Episode: 1, Name: "Endless Nightmare #4", Start: warpIn(), End: register(30)},
		// Today's Rate

		// ---- Retrieval ----
		// The Missing Maracas
		{Episode: 1, Name: "Lost HEAT SWORD", Start: warpIn(), End: register(15)},
		{Episode: 1, Name: "Lost ICE SPINNER", Start: warpIn(), End: register(15)},
		{Episode: 1, Name: "Lost SOUL BLADE", Start: warpIn(), End: register(18)},
		{Episode: 1, Name: "Lost HELL PALLASCH", Start: warpIn(), End: register(110)},
		// FoaM
		// Rappy's Holiday
		// RFR
		// Soul of Steel
		{Episode: 1, Name: "Forsaken Friends", Start: warpIn(), End: register(99)},
		// DR2.0
		// TMBTP

		// ---- Event ----
		{Episode: 1, Name: "Christmas Fiasco", Remap: remap("Christmas Fiasco Episode 1")},
		{Episode: 1, Name: "Christmas Fiasco Episode 1", Start: floorSwitch(4, 100), End: floorSwitch(10, 3)},
		{Episode: 1, Name: "Maximum Attack E: Forest", Start: floorSwitch(2, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(2, 1)},
			{Name: "Room 2", Trigger: floorSwitch(2, 2)},
			{Name: "Room 3", Trigger: floorSwitch(2, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 1, Name: "Maximum Attack E: Caves", Start: floorSwitch(4, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(4, 1)},
			{Name: "Room 2", Trigger: floorSwitch(4, 2)},
			{Name: "Room 3", Trigger: floorSwitch(4, 3)},
			{Name: "Room 4", Trigger: floorSwitch(4, 4)},
			{Name: "Room 5", Trigger: floorSwitch(4, 5)},
			{Name: "Room 6", Trigger: floorSwitch(4, 6)},
			{Name: "Room 7", Trigger: register(50)},
		}},
		{Episode: 1, Name: "Maximum Attack E: Mines", Start: floorSwitch(6, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(6, 1)},
			{Name: "Room 2", Trigger: floorSwitch(6, 2)},
			{Name: "Room 3", Trigger: floorSwitch(6, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 1, Name: "Maximum Attack E: Ruins", Start: floorSwitch(10, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(10, 1)},
			{Name: "Room 2", Trigger: floorSwitch(10, 2)},
			{Name: "Room 3", Trigger: floorSwitch(10, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		// ---- Maximum Attack ----
		// "MAXIMUM ATTACK 1 Ver2"
		{Episode: 1, Name: "Maximum Attack 4 -1A-", Start: floorSwitch(4, 99), End: floorSwitch(10, 31)},
		{Episode: 1, Name: "Maximum Attack 4th Stage -1A-", Remap: remap("Maximum Attack 4 -1A-")},
		{Episode: 1, Name: "Maximum Attack 4 -1B-", Start: floorSwitch(4, 99), End: floorSwitch(10, 32)},
		{Episode: 1, Name: "Maximum Attack 4th Stage -1B-", Remap: remap("Maximum Attack 4 -1B-")},
		{Episode: 1, Name: "Maximum Attack 4 -1C-", Start: floorSwitch(4, 99), End: register(254)},
		{Episode: 1, Name: "Maximum Attack 4th Stage -1C-", Remap: remap("Maximum Attack 4 -1C-")},
		// MA1R
		{Episode: 1, Name: "Maximum Attack E: EP 1", Start: floorSwitch(2, 0), End: register(254)},
		{Episode: 1, Name: "Maximum Attack E: Episode 1", Start: floorSwitch(2, 0), End: register(254)},
		{Episode: 1, Name: "Maximum Attack S", Start: floorSwitch(5, 10), End: register(41)},
		{Episode: 1, Name: "Random Attack Xrd Stage", Start: warpIn(), End: register(254)},
		{Episode: 1, Name: "Random Attack Xrd REV 1", Start: warpIn(), End: register(254)},
		// ---- VR ----
		{Episode: 1, Name: "Towards the Future", Start: register(12), End: register(254)},
		{Episode: 1, Name: "Tyrell's Ego", Start: register(4), End: register(101)},
		// Labyrinthine Trial
		// Sugoroku
		// Sim 2.0
		// Mine Offensive
		{Episode: 1, Name: "Endless: Episode 1", Start: register(50), End: register(248), ForceTerminal: true},

		// ---- Halloween ----
		{Episode: 1, Name: "Hollow Battlefield: Caves", Start: warpIn(), End: register(0)},
		{Episode: 1, Name: "Hollow Battlefield: Ruins", Start: warpIn(), End: register(0)},
		// ---- Challenge Mode ----
		{Episode: 1, Name: "Stage1", Remap: remap("1c1")},
		{Episode: 1, Name: "1c1", Start: register(111), End: register(240), CmodeStage: 1},
		{Episode: 1, Name: "Stage2", Remap: remap("1c2")},
		{Episode: 1, Name: "1c2", Start: register(111), End: register(240), CmodeStage: 2},
		{Episode: 1, Name: "Stage3", Remap: remap("1c3")},
		{Episode: 1, Name: "1c3", Start: register(111), End: register(240), CmodeStage: 3},
		{Episode: 1, Name: "Stage4", Remap: remap("1c4")},
		{Episode: 1, Name: "1c4", Start: register(111), End: register(240), CmodeStage: 4},
		{Episode: 1, Name: "Stage5", Remap: remap("1c5")},
		{Episode: 1, Name: "1c5", Start: register(111), End: register(240), CmodeStage: 5},
		{Episode: 1, Name: "Stage6", Remap: remap("1c6")},
		{Episode: 1, Name: "1c6", Start: register(111), End: register(240), CmodeStage: 6},
		{Episode: 1, Name: "Stage7", Remap: remap("1c7")},
		{Episode: 1, Name: "1c7", Start: register(111), End: register(240), CmodeStage: 7},
		{Episode: 1, Name: "Stage8", Remap: remap("1c8")},
		{Episode: 1, Name: "1c8", Start: register(111), End: register(240), CmodeStage: 8},
		{Episode: 1, Name: "Stage9", Remap: remap("1c9")},
		{Episode: 1, Name: "1c9", Start: register(111), End: register(240), CmodeStage: 9},
		// ---- Shop ----
		{Episode: 1, Name: "Christmas Event Shop", Ignore: true},
		{Episode: 1, Name: "Anniversary Badge Shop", Ignore: true},

		// ---------------------------------------------------------------------------------
		//     Episode 4
		// ---------------------------------------------------------------------------------
		// --- Government ----
		// 5-1:Test/VR Temple 1
		// 5-2:Test/VR Temple 2
		// 5-3:Test/VR Temple 3
		// 5-4:Test/VR Temple 4
		// 5-5:Test/VR Temple 5
		// 6-1:Test/Spaceship 1
		// 6-2:Test/Spaceship 2
		// 6-3:Test/Spaceship 3
		// 6-4:Test/Spaceship 4
		// 6-5:Test/Spaceship 5
		// 7-1:From the Past
		// 7-2:Seeking Clues
		// 7-3:Silent Beach
		// 7-4:Central Control
		// 7-5:Isle of Mutants
		// 8-1:Below the Waves
		// 8-2:Desire's End
		// 8-3:Purple Lamplight

		// ---- Side Story ----
		// Seat of the Heart
		// Blue Star Memories

		// ---- Solo Only ----
		// Knight of Coral Advent
		// A New Hope

		// ---- Event ----
		// Festivity On The Beach
		{Episode: 2, Name: "Christmas Fiasco", Remap: remap("Christmas Fiasco Episode 2")},
		{Episode: 2, Name: "Christmas Fiasco Episode 2", Start: floorSwitch(4, 100), End: floorSwitch(11, 3)},
		{Episode: 2, Name: "Maximum Attack E: Temple", Start: floorSwitch(1, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(1, 1)},
			{Name: "Room 2", Trigger: floorSwitch(1, 2)},
			{Name: "Room 3", Trigger: floorSwitch(1, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: Spaceship", Remap: remap("Maximum Attack E: Space")},
		{Episode: 2, Name: "Maximum Attack E: Space", Start: floorSwitch(3, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(3, 1)},
			{Name: "Room 2", Trigger: floorSwitch(3, 2)},
			{Name: "Room 3", Trigger: floorSwitch(3, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: CCA", Start: floorSwitch(5, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(5, 1)},
			{Name: "Room 2", Trigger: floorSwitch(5, 2)},
			{Name: "Room 3", Trigger: floorSwitch(5, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: Seabed", Start: floorSwitch(10, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(10, 1)},
			{Name: "Room 2", Trigger: floorSwitch(10, 2)},
			{Name: "Room 3", Trigger: floorSwitch(10, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: Tower", Start: floorSwitch(17, 0), End: register(50), Splits: []Split{
			{Name: "1f", Trigger: floorSwitch(17, 1)},
			{Name: "2f", Trigger: floorSwitch(17, 2)},
			{Name: "3f", Trigger: floorSwitch(17, 3)},
			{Name: "4f", Trigger: floorSwitch(17, 4)},
			{Name: "5f", Trigger: register(50)},
		}},
		// ---- Extermination ----
		{Episode: 2, Name: "Phantasmal World #1", Start: warpIn(), End: register(254)},
		{Episode: 2, Name: "Phantasmal World #2", Start: warpIn(), End: register(111)},
		{Episode: 2, Name: "Phantasmal World #3", Start: warpIn(), End: floorSwitch(11, 180)},
		{Episode: 2, Name: "Phantasmal World #4", Start: warpIn(), End: floorSwitch(16, 120)},
		{Episode: 2, Name: "Sweep-up Operation #5", Start: register(210), End: register(254)},
		{Episode: 2, Name: "Sweep-up Operation #6", Start: register(210), End: register(254)},
		{Episode: 2, Name: "Sweep-up Operation #7", Start: register(210), End: register(254)},
		{Episode: 2, Name: "Sweep-up Operation #8", Start: register(210), End: register(254)},
		{Episode: 2, Name: "Sweep-up Operation #9", Start: register(210), End: register(254), Splits: []Split{
			{Name: "10f & 9f", Trigger: floorSwitch(17, 9)},
			{Name: "8f", Trigger: floorSwitch(17, 8)},
			{Name: "7f", Trigger: floorSwitch(17, 7)},
			{Name: "6f", Trigger: floorSwitch(17, 6)},
			{Name: "5f", Trigger: floorSwitch(17, 5)},
			{Name: "4f", Trigger: floorSwitch(17, 4)},
			{Name: "3f", Trigger: floorSwitch(17, 3)},
			{Name: "2f", Trigger: floorSwitch(17, 2)},
			{Name: "1f", Trigger: floorSwitch(17, 1)},
		}},
		{Episode: 2, Name: "Gal Da Val's Darkness", Start: floorSwitch(3, 20), End: register(89)},
		{Episode: 2, Name: "CAL's Clock Challenge", Start: floorSwitch(1, 40), End: register(254)},
		// ---- Retrieval ----
		{Episode: 2, Name: "Lost SHOCK RIFLE", Start: warpIn(), End: register(15)},
		{Episode: 2, Name: "Lost BIND ASSAULT", Start: warpIn(), End: register(15)},
		{Episode: 2, Name: "Lost FILL CANNON", Start: warpIn(), End: register(15)},
		{Episode: 2, Name: "Lost CHARGE VULCAN", Start: warpIn(), End: register(15)},
		// Dream Messenger
		// Revisiting Darkness
		// Dolmolm Research

		// ---- Maximum Attack ----
		{Episode: 2, Name: "MAXIMUM ATTACK 2 Ver2", Start: floorSwitch(2, 12), End: register(123)},
		{Episode: 2, Name: "Maximum Attack 4th Stage -2A-", Remap: remap("Maximum Attack 4 -2A-")},
		{Episode: 2, Name: "Maximum Attack 4 -2A-", Start: floorSwitch(5, 99), End: floorSwitch(11, 29)},
		{Episode: 2, Name: "Maximum Attack 4th Stage -2B-", Remap: remap("Maximum Attack 4 -2B-")},
		{Episode: 2, Name: "Maximum Attack 4 -2B-", Start: floorSwitch(5, 99), End: floorSwitch(11, 29)},
		{Episode: 2, Name: "Maximum Attack 4th Stage -2C-", Remap: remap("Maximum Attack 4 -2C-")},
		{Episode: 2, Name: "Maximum Attack 4 -2C-", Start: floorSwitch(5, 99), End: floorSwitch(11, 29)},
		// Maximum Attack 4th Stage -2R-
		{Episode: 2, Name: "Maximum Attack E: VR", Start: floorSwitch(1, 0), End: register(254)},
		{Episode: 2, Name: "Maximum Attack E: GDV", Start: floorSwitch(5, 0), End: register(254)},
		{Episode: 2, Name: "Maximum Attack E: Gal Da Val", Start: floorSwitch(5, 0), End: register(254)},
		// Maximum Attack S
		// Random Attack Xrd Stage
		{Episode: 2, Name: "Random Attack Xrd REV 2", Start: warpIn(), End: register(254)},
		// ---- Tower ----
		// The East Tower
		// The West Tower
		{Episode: 2, Name: "Raid on Central Tower", Ignore: true}, // Quest crashes if you're too quick to kill Olga :)
		{Episode: 2, Name: "The Military Strikes Back", Start: warpIn(), End: register(121)},
		{Episode: 2, Name: "Twilight Sanctuary", Start: register(50), End: register(254)},

		// ---- Shop ----
		{Episode: 2, Name: "Gallon's Shop", Ignore: true},
		{Episode: 2, Name: "Item Present", Ignore: true},
		{Episode: 2, Name: "Singing by the Beach", Ignore: true},
		{Episode: 2, Name: "Beach Laughter", Ignore: true},
		{Episode: 2, Name: "To The Deepest Blue -MA4 Venue-", Ignore: true},
		{Episode: 2, Name: "TA Reward Shop", Ignore: true},
		{Episode: 2, Name: "The Egg Shop", Ignore: true},

		// ---- VR ----
		// Reach for the Dream
		{Episode: 2, Name: "Respective Tomorrow", Start: register(84), End: register(98)},
		// ---- Halloween ----
		{Episode: 2, Name: "Hollow Reality: Temple", Start: warpIn(), End: register(0)},
		{Episode: 2, Name: "Hollow Reality: Spaceship", Start: warpIn(), End: register(0)},
		{Episode: 2, Name: "Hollow Phantasm: Jungle", Start: warpIn(), End: register(0)},
		{Episode: 2, Name: "Hollow Phantasm: Seabed", Start: warpIn(), End: register(0)},
		// ---- Challenge Mode ----
		{Episode: 2, Name: "Stage1", Remap: remap("2c1")},
		{Episode: 2, Name: "2c1", Start: register(111), End: register(240), CmodeStage: 1},
		{Episode: 2, Name: "Stage2", Remap: remap("2c2")},
		{Episode: 2, Name: "2c2", Start: register(111), End: register(240), CmodeStage: 2},
		{Episode: 2, Name: "Stage3", Remap: remap("2c3")},
		{Episode: 2, Name: "2c3", Start: register(111), End: register(240), CmodeStage: 3},
		{Episode: 2, Name: "Stage4", Remap: remap("2c4")},
		{Episode: 2, Name: "2c4", Start: register(111), End: register(240), CmodeStage: 4},
		{Episode: 2, Name: "Stage5", Remap: remap("2c5")},
		{Episode: 2, Name: "2c5", Start: register(111), End: register(240), CmodeStage: 5},

		// ---------------------------------------------------------------------------------
		//     Episode 4
		// ---------------------------------------------------------------------------------
		// ---- Government ----
		// 9-1:Missing Research
		// 9-2:Data Retrieval
		// 9-3:Reality & Truth
		// 9-4:Pursuit
		// 9-5:The Chosen (1/2),
		// 9-6:The Chosen (2/2),
		// 9-7:Sacred Ground
		// 9-8:The Final Cycle

		// ---- Side Story ----
		// Warrior's Pride
		// The Restless Lion
		// Pioneer Spirit
		// 荒野の果てに
		// Black Paper's Dangerous Deal
		// Black Paper's Dangerous Deal 2

		// ---- Extermination ----
		{Episode: 4, Name: "Point of Disaster", Start: warpIn(), End: register(233)},
		// The Robots' Reckoning
		{Episode: 4, Name: "War of Limits 1", Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 2", Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 3", Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 4", Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 5", Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "New Mop-Up Operation #1", Start: register(205), End: register(157)},
		{Episode: 4, Name: "New Mop-Up Operation #2", Start: register(86), End: register(43)},
		{Episode: 4, Name: "New Mop-Up Operation #3", Start: register(205), End: register(254)},
		{Episode: 4, Name: "New Mop-Up Operation #4", Start: register(110), End: register(43)},
		{Episode: 4, Name: "New Mop-Up Operation #5", Start: register(86), End: register(43)},
		{Episode: 4, Name: "Sweep-up Operation #10", Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #11", Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #12", Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #13", Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #14", Start: register(210), End: register(254)},
		// ---- Retrieval ----
		{Episode: 4, Name: "Lost BERSERK BATON", Start: warpIn(), End: register(15)},
		{Episode: 4, Name: "Lost SPIRIT STRIKER", Start: warpIn(), End: register(15)},
		// ---- Event ----
		{Episode: 4, Name: "Christmas Fiasco", Remap: remap("Christmas Fiasco Episode 4")},
		{Episode: 4, Name: "Christmas Fiasco Episode 4", Start: floorSwitch(1, 100), End: floorSwitch(8, 3)},
		{Episode: 4, Name: "Maximum Attack E: Crater", Start: floorSwitch(2, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(2, 1)},
			{Name: "Room 2", Trigger: floorSwitch(2, 2)},
			{Name: "Room 3", Trigger: floorSwitch(2, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 4, Name: "Maximum Attack E: Desert", Start: floorSwitch(8, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(8, 1)},
			{Name: "Room 2", Trigger: floorSwitch(8, 2)},
			{Name: "Room 3", Trigger: floorSwitch(8, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		// ---- Maximum Attack ----
		// MAXIMUM ATTACK 3 Ver2
		{Episode: 4, Name: "Maximum Attack 4 -4A-", Start: floorSwitch(5, 66), End: floorSwitch(8, 50)},
		{Episode: 4, Name: "Maximum Attack 4th Stage -4A-", Remap: remap("Maximum Attack 4 -4A-")},
		{Episode: 4, Name: "Maximum Attack 4 -4B-", Start: floorSwitch(5, 66), End: floorSwitch(8, 20)},
		{Episode: 4, Name: "Maximum Attack 4th Stage -4B-", Remap: remap("Maximum Attack 4 -4B-")},
		{Episode: 4, Name: "Maximum Attack 4th Stage -B-", Remap: remap("Maximum Attack 4 -4B-")},
		{Episode: 4, Name: "Maximum Attack 4 -4C-", Start: floorSwitch(5, 66), End: floorSwitch(8, 192)},
		{Episode: 4, Name: "Maximum Attack 4th Stage -4C-", Remap: remap("Maximum Attack 4 -4C-")},
		// Maximum Attack 4th Stage -4R-
		{Episode: 4, Name: "Maximum Attack E: Episode 4", Start: floorSwitch(2, 0), End: register(254)},
		// Maximum Attack S
		{Episode: 4, Name: "Random Attack Xrd REV 4", Start: warpIn(), End: register(254)},
		// ---- Shop ----
		// Claire's Deal 5
		// The Beak's Cafe Ver.2

		// ---- VR ----
		{Episode: 4, Name: "Beyond the Horizon", Start: floorSwitch(1, 20), End: register(254)},
		// LOGiN presents 勇場のマッチレース

		// ---- Halloween ----
		{Episode: 4, Name: "Hollow Wasteland: Wilderness", Start: warpIn(), End: register(0)},
	}
}
