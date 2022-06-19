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
		{Episode: 1, Name: "Mop-up Operation #1", Number: 101, Start: register(0), End: register(254)},
		{Episode: 1, Name: "Mop-up Operation #2", Number: 102, Start: register(0), End: register(254)},
		{Episode: 1, Name: "Mop-up Operation #3", Number: 103, Start: register(0), End: register(254)},
		{Episode: 1, Name: "Mop-up Operation #4", Number: 104, Start: register(0), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #1", Number: 1761, Start: register(210), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #2", Number: 1762, Start: register(210), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #3", Number: 1763, Start: register(210), End: register(254)},
		{Episode: 1, Name: "Sweep-up Operation #4", Number: 1764, Start: register(210), End: register(254)},
		{Episode: 1, Name: "Endless Nightmare #1", Number: 108, Start: warpIn(), End: register(30)},
		{Episode: 1, Name: "Endless Nightmare #2", Number: 109, Start: warpIn(), End: register(30)},
		{Episode: 1, Name: "Endless Nightmare #3", Number: 110, Start: warpIn(), End: register(30)},
		{Episode: 1, Name: "Endless Nightmare #4", Number: 111, Start: warpIn(), End: register(30)},
		{Episode: 1, Name: "Anomalous Ordeal #1", Number: 1810, Start: register(82), End: register(254)},
		{Episode: 1, Name: "Anomalous Ordeal #2", Number: 1811, Start: register(82), End: register(254)},
		// Today's Rate

		// ---- Retrieval ----
		// The Missing Maracas
		{Episode: 1, Name: "Lost HEAT SWORD", Number: 105, Start: warpIn(), End: register(15)},
		{Episode: 1, Name: "Lost ICE SPINNER", Number: 106, Start: warpIn(), End: register(15)},
		{Episode: 1, Name: "Lost SOUL BLADE", Number: 107, Start: warpIn(), End: register(18)},
		{Episode: 1, Name: "Lost HELL PALLASCH", Number: 120, Start: warpIn(), End: register(110)},
		// FoaM
		// Rappy's Holiday
		// RFR
		// Soul of Steel
		{Episode: 1, Name: "Forsaken Friends", Number: 907, Start: warpIn(), End: register(99)},
		// DR2.0
		// TMBTP

		// ---- Event ----
		{Episode: 1, Name: "Christmas Fiasco", Remap: remap("Christmas Fiasco Episode 1")},
		{Episode: 1, Name: "Christmas Fiasco Episode 1", Number: 900, Start: floorSwitch(4, 100), End: floorSwitch(10, 3)},
		{Episode: 1, Name: "Maximum Attack E: Forest", Number: 930, Start: floorSwitch(2, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(2, 1)},
			{Name: "Room 2", Trigger: floorSwitch(2, 2)},
			{Name: "Room 3", Trigger: floorSwitch(2, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 1, Name: "Maximum Attack E: Forest -MM-", Start: floorSwitch(2, 0), End: register(50),
			Splits: []Split{
				{Name: "Room 1", Trigger: floorSwitch(2, 99)},
				{Name: "Room 2", Trigger: floorSwitch(2, 3)},
				{Name: "Room 3", Trigger: floorSwitch(2, 2)},
				{Name: "Room 4", Trigger: register(50)},
			}},
		{Episode: 1, Name: "Maximum Attack E: Cave", Remap: remap("Maximum Attack E: Caves")},
		{Episode: 1, Name: "Maximum Attack E: Caves", Number: 931, Start: floorSwitch(4, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(4, 1)},
			{Name: "Room 2", Trigger: floorSwitch(4, 2)},
			{Name: "Room 3", Trigger: floorSwitch(4, 3)},
			{Name: "Room 4", Trigger: floorSwitch(4, 4)},
			{Name: "Room 5", Trigger: floorSwitch(4, 5)},
			{Name: "Room 6", Trigger: floorSwitch(4, 6)},
			{Name: "Room 7", Trigger: register(50)},
		}},
		{Episode: 1, Name: "Maximum Attack E: Caves -MM-", Start: floorSwitch(4, 0), End: register(50),
			Splits: []Split{
				{Name: "Room 1", Trigger: floorSwitch(4, 1)},
				{Name: "Room 2", Trigger: floorSwitch(4, 2)},
				{Name: "Room 3", Trigger: floorSwitch(4, 3)},
				{Name: "Room 4", Trigger: floorSwitch(4, 4)},
				{Name: "Room 5", Trigger: floorSwitch(4, 5)},
				{Name: "Room 6", Trigger: floorSwitch(4, 6)},
				{Name: "Room 7", Trigger: register(50)},
			}},
		{Episode: 1, Name: "Maximum Attack E: Mine", Remap: remap("Maximum Attack E: Mines")},
		{Episode: 1, Name: "Maximum Attack E: Mines", Number: 932, Start: floorSwitch(6, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(6, 1)},
			{Name: "Room 2", Trigger: floorSwitch(6, 2)},
			{Name: "Room 3", Trigger: floorSwitch(6, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 1, Name: "Maximum Attack E: Mines -MM-", Start: floorSwitch(6, 0), End: register(50),
			Splits: []Split{
				{Name: "Room 1", Trigger: floorSwitch(6, 3)},
				{Name: "Room 2", Trigger: floorSwitch(6, 4)},
				{Name: "Room 3", Trigger: floorSwitch(6, 2)},
				{Name: "Room 4", Trigger: register(50)},
			}},
		{Episode: 1, Name: "Maximum Attack E: Ruins", Number: 933, Start: floorSwitch(10, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(10, 1)},
			{Name: "Room 2", Trigger: floorSwitch(10, 2)},
			{Name: "Room 3", Trigger: floorSwitch(10, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 1, Name: "Maximum Attack E: Ruins -MM-", Start: floorSwitch(10, 0), End: register(50),
			Splits: []Split{
				{Name: "Room 1", Trigger: floorSwitch(10, 1)},
				{Name: "Room 2", Trigger: floorSwitch(10, 2)},
				{Name: "Room 3", Trigger: floorSwitch(10, 3)},
				{Name: "Room 4", Trigger: register(50)},
			}},
		// ---- Maximum Attack ----
		// "MAXIMUM ATTACK 1 Ver2"
		{Episode: 1, Name: "Maximum Attack 4 -1A-", Number: 144, Start: floorSwitch(4, 99), End: floorSwitch(10, 31)},
		{Episode: 1, Name: "Maximum Attack 4th Stage -1A-", Remap: remap("Maximum Attack 4 -1A-")},
		{Episode: 1, Name: "Maximum Attack 4 -1B-", Number: 145, Start: floorSwitch(4, 99), End: floorSwitch(10, 32)},
		{Episode: 1, Name: "Maximum Attack 4th Stage -1B-", Remap: remap("Maximum Attack 4 -1B-")},
		{Episode: 1, Name: "Maximum Attack 4 -1C-", Number: 146, Start: floorSwitch(4, 99), End: register(254)},
		{Episode: 1, Name: "Maximum Attack 4th Stage -1C-", Remap: remap("Maximum Attack 4 -1C-")},
		{Episode: 1, Name: "C1AM", Start: floorSwitch(10, 99), End: register(50),
			Splits: []Split{
				{Name: "Room 1", Trigger: floorSwitch(10, 4)},
				{Name: "Room 2", Trigger: floorSwitch(10, 5)},
				{Name: "Room 3", Trigger: floorSwitch(6, 53)},
				{Name: "Room 4", Trigger: floorSwitch(6, 220)},
				{Name: "Room 5", Trigger: floorSwitch(6, 52)},
				{Name: "Room 6", Trigger: floorSwitch(6, 21)},
				{Name: "Room 7", Trigger: floorSwitch(4, 35)},
				{Name: "Room 8", Trigger: floorSwitch(4, 30)},
				{Name: "Room 9", Trigger: floorSwitch(4, 45)},
				{Name: "Room 10", Trigger: floorSwitch(4, 60)},
				{Name: "Room 11", Trigger: floorSwitch(4, 121)},
				{Name: "Room 12", Trigger: floorSwitch(4, 22)},
				{Name: "Room 12", Trigger: floorSwitch(4, 9)},
				{Name: "Room 13", Trigger: register(50)},
			}},
		// MA1R
		{Episode: 1, Name: "Maximum Attack E: Episode 1", Number: 942, Start: floorSwitch(2, 0), End: register(254)},
		{Episode: 1, Name: "Maximum Attack S", Start: floorSwitch(5, 10), End: register(41)},
		{Episode: 1, Name: "Random Attack Xrd Stage", Number: 1303, Start: warpIn(), End: register(254)},
		{Episode: 1, Name: "Random Attack Xrd REV 1", Number: 1801, Start: warpIn(), End: register(254)},
		// ---- VR ----
		{Episode: 1, Name: "Towards the Future", Number: 118, Start: register(12), End: register(254)},
		{Episode: 1, Name: "Tyrell's Ego", Number: 161, Start: register(4), End: register(101)},
		{Episode: 1, Name: "総督の贈り物", Start: floorSwitch(6, 165), End: floorSwitch(6, 33)},
		// Labyrinthine Trial
		// Sugoroku
		// Sim 2.0
		// Mine Offensive
		{Episode: 1, Name: "Endless: Episode 1", Number: 1850, Start: register(50), End: register(248), ForceTerminal: true},

		// ---- Halloween ----
		{Episode: 1, Name: "Hollow Battlefield: Forest", Number: 1666, Start: warpIn(), End: register(0)},
		{Episode: 1, Name: "Hollow Battlefield: Cave", Number: 1667, Start: warpIn(), End: register(0)},
		{Episode: 1, Name: "Hollow Battlefield: Mine", Number: 1668, Start: warpIn(), End: register(0)},
		{Episode: 1, Name: "Hollow Battlefield: Ruins", Number: 1669, Start: warpIn(), End: register(0)},
		// ---- Challenge Mode ----
		{Episode: 1, Name: "Stage1", Remap: remap("1c1")},
		{Episode: 1, Name: "ステージ1", Remap: remap("1c1")},
		{Episode: 1, Name: "1c1", Start: register(111), End: register(240), CmodeStage: 1},
		{Episode: 1, Name: "Stage2", Remap: remap("1c2")},
		{Episode: 1, Name: "ステージ2", Remap: remap("1c2")},
		{Episode: 1, Name: "1c2", Start: register(111), End: register(240), CmodeStage: 2},
		{Episode: 1, Name: "Stage3", Remap: remap("1c3")},
		{Episode: 1, Name: "ステージ3", Remap: remap("1c3")},
		{Episode: 1, Name: "1c3", Start: register(111), End: register(240), CmodeStage: 3},
		{Episode: 1, Name: "Stage4", Remap: remap("1c4")},
		{Episode: 1, Name: "ステージ4", Remap: remap("1c4")},
		{Episode: 1, Name: "1c4", Start: register(111), End: register(240), CmodeStage: 4},
		{Episode: 1, Name: "Stage5", Remap: remap("1c5")},
		{Episode: 1, Name: "ステージ5", Remap: remap("1c5")},
		{Episode: 1, Name: "Stage5 ", Remap: remap("1c5")},
		{Episode: 1, Name: "1c5", Start: register(111), End: register(240), CmodeStage: 5},
		{Episode: 1, Name: "Stage6", Remap: remap("1c6")},
		{Episode: 1, Name: "ステージ6", Remap: remap("1c6")},
		{Episode: 1, Name: "1c6", Start: register(111), End: register(240), CmodeStage: 6},
		{Episode: 1, Name: "Stage7", Remap: remap("1c7")},
		{Episode: 1, Name: "ステージ7", Remap: remap("1c7")},
		{Episode: 1, Name: "1c7", Start: register(111), End: register(240), CmodeStage: 7},
		{Episode: 1, Name: "ステージ8", Remap: remap("1c8")},
		{Episode: 1, Name: "Stage8", Remap: remap("1c8")},
		{Episode: 1, Name: "1c8", Start: register(111), End: register(240), CmodeStage: 8},
		{Episode: 1, Name: "Stage9", Remap: remap("1c9")},
		{Episode: 1, Name: "ステージ9", Remap: remap("1c9")},
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
		{Episode: 2, Name: "Christmas Fiasco Episode 2", Number: 901, Start: floorSwitch(4, 100), End: floorSwitch(11, 3)},
		{Episode: 2, Name: "Maximum Attack E: Temple", Number: 934, Start: floorSwitch(1, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(1, 1)},
			{Name: "Room 2", Trigger: floorSwitch(1, 2)},
			{Name: "Room 3", Trigger: floorSwitch(1, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: Spaceship", Remap: remap("Maximum Attack E: Space")},
		{Episode: 2, Name: "Maximum Attack E: Space", Number: 935, Start: floorSwitch(3, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(3, 1)},
			{Name: "Room 2", Trigger: floorSwitch(3, 2)},
			{Name: "Room 3", Trigger: floorSwitch(3, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: CCA", Number: 936, Start: floorSwitch(5, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(5, 1)},
			{Name: "Room 2", Trigger: floorSwitch(5, 2)},
			{Name: "Room 3", Trigger: floorSwitch(5, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: Seabed", Number: 937, Start: floorSwitch(10, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(10, 1)},
			{Name: "Room 2", Trigger: floorSwitch(10, 2)},
			{Name: "Room 3", Trigger: floorSwitch(10, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 2, Name: "Maximum Attack E: Tower", Number: 938, Start: floorSwitch(17, 0), End: register(50), Splits: []Split{
			{Name: "1f", Trigger: floorSwitch(17, 1)},
			{Name: "2f", Trigger: floorSwitch(17, 2)},
			{Name: "3f", Trigger: floorSwitch(17, 3)},
			{Name: "4f", Trigger: floorSwitch(17, 4)},
			{Name: "5f", Trigger: register(50)},
		}},
		// ---- Extermination ----
		{Episode: 2, Name: "Malicious Uprising #1", Start: register(205), End: floorSwitch(2, 16)},
		{Episode: 2, Name: "Malicious Uprising #2", Start: register(205), End: floorSwitch(3, 12)},
		{Episode: 2, Name: "Malicious Uprising #3", Start: register(205), End: floorSwitch(8, 11)},
		{Episode: 2, Name: "Malicious Uprising #4", Start: register(205), End: floorSwitch(10, 14)},
		{Episode: 2, Name: "Malicious Uprising #5", Start: register(205), End: register(157)},
		{Episode: 2, Name: "Penumbral Surge #1", Number: 1821, Start: register(50), End: register(254), ForceTerminal: true},
		{Episode: 2, Name: "Penumbral Surge #2", Number: 1822, Start: register(15), End: register(254), ForceTerminal: true},
		{Episode: 2, Name: "Penumbral Surge #3", Number: 1823, Start: register(90), End: register(254), ForceTerminal: true},
		{Episode: 2, Name: "Penumbral Surge #4", Number: 1824, Start: register(51), End: register(254), ForceTerminal: true},
		{Episode: 2, Name: "Penumbral Surge #5", Number: 1825, Start: register(15), End: register(254), ForceTerminal: true},
		{Episode: 2, Name: "Penumbral Surge #6", Number: 1826, Start: register(15), End: register(254), ForceTerminal: true},
		{Episode: 2, Name: "Phantasmal World #1", Number: 233, Start: warpIn(), End: register(254)},
		{Episode: 2, Name: "Phantasmal World #2", Number: 234, Start: warpIn(), End: register(111)},
		{Episode: 2, Name: "Phantasmal World #3", Number: 235, Start: warpIn(), End: floorSwitch(11, 180)},
		{Episode: 2, Name: "Phantasmal World #4", Number: 236, Start: warpIn(), End: floorSwitch(16, 120)},
		{Episode: 2, Name: "Sweep-up Operation #5", Number: 1765, Start: register(210), End: register(254)},
		{Episode: 2, Name: "Sweep-up Operation #6", Number: 1766, Start: register(210), End: register(254)},
		{Episode: 2, Name: "Sweep-up Operation #7", Number: 1767, Start: register(210), End: register(254)},
		{Episode: 2, Name: "Sweep-up Operation #8", Number: 1768, Start: register(210), End: register(254), Splits: []Split{
			{Name: "First Rooms", Trigger: floorSwitch(10, 27)},
			{Name: "Delbiters", Trigger: floorSwitch(10, 23)},
			{Name: "Morfos Room", Trigger: floorSwitch(10, 34)},
			{Name: "Delbiters+", Trigger: floorSwitch(10, 36)},
			{Name: "Reco Room", Trigger: floorSwitch(10, 38)},
			{Name: "Small Rooms", Trigger: floorSwitch(10, 6)},
			{Name: "Final Room", Trigger: floorSwitch(10, 13)},
		}},
		{Episode: 2, Name: "Sweep-up Operation #9", Number: 1769, Start: register(210), End: register(254), Splits: []Split{
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
		{Episode: 2, Name: "Gal Da Val's Darkness", Number: 1309, Start: floorSwitch(3, 20), End: register(89)},
		{Episode: 2, Name: "CAL's Clock Challenge", Number: 1700, Start: floorSwitch(1, 40), End: register(254)},
		// ---- Retrieval ----
		{Episode: 2, Name: "Lost SHOCK RIFLE", Number: 1780, Start: warpIn(), End: register(15)},
		{Episode: 2, Name: "Lost BIND ASSAULT", Number: 1781, Start: warpIn(), End: register(15)},
		{Episode: 2, Name: "Lost FILL CANNON", Number: 1782, Start: warpIn(), End: register(15)},
		{Episode: 2, Name: "Lost DEMON'S RAILGUN", Number: 1783, Start: warpIn(), End: register(15)},
		{Episode: 2, Name: "Lost CHARGE VULCAN", Number: 1784, Start: warpIn(), End: register(15)},
		// Dream Messenger
		// Revisiting Darkness
		// Dolmolm Research

		// ---- Maximum Attack ----
		{Episode: 2, Name: "MAXIMUM ATTACK 2 Ver2", Number: 238, Start: floorSwitch(2, 12), End: register(123)},
		{Episode: 2, Name: "Maximum Attack 4th Stage -2A-", Remap: remap("Maximum Attack 4 -2A-")},
		{Episode: 2, Name: "Maximum Attack 4 -2A-", Number: 241, Start: floorSwitch(5, 99), End: floorSwitch(11, 29)},
		{Episode: 2, Name: "Maximum Attack 4th Stage -2B-", Remap: remap("Maximum Attack 4 -2B-")},
		{Episode: 2, Name: "Maximum Attack 4 -2B-", Number: 242, Start: floorSwitch(5, 99), End: floorSwitch(11, 29)},
		{Episode: 2, Name: "Maximum Attack 4th Stage -2C-", Remap: remap("Maximum Attack 4 -2C-")},
		{Episode: 2, Name: "Maximum Attack 4 -2C-", Number: 243, Start: floorSwitch(5, 99), End: floorSwitch(11, 29)},
		// Maximum Attack 4th Stage -2R-
		{Episode: 2, Name: "Maximum Attack E: VR", Number: 943, Start: floorSwitch(1, 0), End: register(254)},
		{Episode: 2, Name: "Maximum Attack E: GDV", Number: 944, Start: floorSwitch(5, 0), End: register(254)},
		{Episode: 2, Name: "Maximum Attack E: Gal Da Val", Number: 944, Start: floorSwitch(5, 0), End: register(254)},
		// Maximum Attack S
		// Random Attack Xrd Stage
		{Episode: 2, Name: "Random Attack Xrd REV 2", Number: 1802, Start: warpIn(), End: register(254)},
		// ---- Tower ----
		// The East Tower
		// The West Tower
		{Episode: 2, Name: "Raid on Central Tower", Ignore: true}, // Quest crashes if you're too quick to kill Olga :)
		{Episode: 2, Name: "The Military Strikes Back", Number: 1319, Start: warpIn(), End: register(121)},
		{Episode: 2, Name: "Twilight Sanctuary", Number: 1820, Start: register(50), End: register(254)},

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
		{Episode: 2, Name: "Respective Tomorrow", Number: 231, Start: register(84), End: register(98)},
		// ---- Halloween ----
		{Episode: 2, Name: "Hollow Reality: Temple", Number: 1670, Start: warpIn(), End: register(0)},
		{Episode: 2, Name: "Hollow Reality: Spaceship", Number: 1671, Start: warpIn(), End: register(0)},
		{Episode: 2, Name: "Hollow Phantasm: Jungle", Number: 1672, Start: warpIn(), End: register(0)},
		{Episode: 2, Name: "Hollow Phantasm: Seabed", Number: 1673, Start: warpIn(), End: register(0)},
		{Episode: 2, Name: "Hollow Phantasm: Tower", Number: 1674, Start: warpIn(), End: register(0)},
		// ---- Challenge Mode ----
		{Episode: 2, Name: "Stage1", Remap: remap("2c1")},
		{Episode: 2, Name: "ステージ１", Remap: remap("2c1")},
		{Episode: 2, Name: "2c1", Start: register(111), End: register(240), CmodeStage: 1},
		{Episode: 2, Name: "Stage2", Remap: remap("2c2")},
		{Episode: 2, Name: "ステージ２", Remap: remap("2c2")},
		{Episode: 2, Name: "2c2", Start: register(111), End: register(240), CmodeStage: 2},
		{Episode: 2, Name: "Stage3", Remap: remap("2c3")},
		{Episode: 2, Name: "ステージ３", Remap: remap("2c3")},
		{Episode: 2, Name: "2c3", Start: register(111), End: register(240), CmodeStage: 3},
		{Episode: 2, Name: "Stage4", Remap: remap("2c4")},
		{Episode: 2, Name: "ステージ４", Remap: remap("2c4")},
		{Episode: 2, Name: "2c4", Start: register(111), End: register(240), CmodeStage: 4},
		{Episode: 2, Name: "Stage5", Remap: remap("2c5")},
		{Episode: 2, Name: "ステージ５", Remap: remap("2c5")},
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
		{Episode: 4, Name: "War of Limits 1", Number: 811, Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 2", Number: 812, Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 3", Number: 813, Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 4", Number: 814, Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "War of Limits 5", Number: 815, Start: warpIn(), End: register(254)},
		{Episode: 4, Name: "New Mop-Up Operation #1", Number: 816, Start: register(205), End: register(157)},
		{Episode: 4, Name: "New Mop-Up Operation #2", Number: 817, Start: register(86), End: register(43)},
		{Episode: 4, Name: "New Mop-Up Operation #3", Number: 818, Start: register(205), End: register(254)},
		{Episode: 4, Name: "New Mop-Up Operation #4", Number: 819, Start: register(110), End: register(43)},
		{Episode: 4, Name: "New Mop-Up Operation #5", Number: 820, Start: register(86), End: register(43)},
		{Episode: 4, Name: "Sweep-up Operation #10", Number: 1770, Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #11", Number: 1771, Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #12", Number: 1772, Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #13", Number: 1773, Start: register(210), End: register(254)},
		{Episode: 4, Name: "Sweep-up Operation #14", Number: 1774, Start: register(210), End: register(254)},
		// ---- Retrieval ----
		{Episode: 4, Name: "Lost BERSERK BATON", Number: 1790, Start: warpIn(), End: register(15)},
		{Episode: 4, Name: "Lost SPIRIT STRIKER", Number: 1791, Start: warpIn(), End: register(15)},
		// ---- Event ----
		{Episode: 4, Name: "Christmas Fiasco", Remap: remap("Christmas Fiasco Episode 4")},
		{Episode: 4, Name: "Christmas Fiasco Episode 4", Number: 902, Start: floorSwitch(1, 100), End: floorSwitch(8, 3)},
		{Episode: 4, Name: "Maximum Attack E: Crater", Number: 939, Start: floorSwitch(2, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(2, 1)},
			{Name: "Room 2", Trigger: floorSwitch(2, 2)},
			{Name: "Room 3", Trigger: floorSwitch(2, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		{Episode: 4, Name: "Maximum Attack E: Desert", Number: 940, Start: floorSwitch(8, 0), End: register(50), Splits: []Split{
			{Name: "Room 1", Trigger: floorSwitch(8, 1)},
			{Name: "Room 2", Trigger: floorSwitch(8, 2)},
			{Name: "Room 3", Trigger: floorSwitch(8, 3)},
			{Name: "Room 4", Trigger: register(50)},
		}},
		// ---- Maximum Attack ----
		// MAXIMUM ATTACK 3 Ver2
		{Episode: 4, Name: "Maximum Attack 4 -4A-", Number: 303, Start: floorSwitch(5, 66), End: floorSwitch(8, 50)},
		{Episode: 4, Name: "Maximum Attack 4th Stage -4A-", Remap: remap("Maximum Attack 4 -4A-")},
		{Episode: 4, Name: "Maximum Attack 4 -4B-", Number: 304, Start: floorSwitch(5, 66), End: floorSwitch(8, 20)},
		{Episode: 4, Name: "Maximum Attack 4th Stage -4B-", Remap: remap("Maximum Attack 4 -4B-")},
		{Episode: 4, Name: "Maximum Attack 4th Stage -B-", Remap: remap("Maximum Attack 4 -4B-")},
		{Episode: 4, Name: "Maximum Attack 4 -4C-", Number: 305, Start: floorSwitch(5, 66), End: floorSwitch(8, 192)},
		{Episode: 4, Name: "Maximum Attack 4th Stage -4C-", Remap: remap("Maximum Attack 4 -4C-")},
		// Maximum Attack 4th Stage -4R-
		{Episode: 4, Name: "Maximum Attack E: Episode 4", Number: 945, Start: floorSwitch(2, 0), End: register(254)},
		// Maximum Attack S
		{Episode: 4, Name: "Random Attack Xrd REV 4", Number: 1803, Start: warpIn(), End: register(254)},
		// ---- Shop ----
		// Claire's Deal 5
		// The Beak's Cafe Ver.2

		// ---- VR ----
		{Episode: 4, Name: "Beyond the Horizon", Number: 313, Start: floorSwitch(1, 20), End: floorSwitch(8, 80)},
		// LOGiN presents 勇場のマッチレース

		// ---- Halloween ----
		{Episode: 4, Name: "Hollow Wasteland: Wilderness", Number: 1675, Start: warpIn(), End: register(0)},
		{Episode: 4, Name: "Hollow Wasteland: Desert", Number: 1676, Start: warpIn(), End: register(0)},
	}
}
