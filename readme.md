# Running

Configure `config.yaml` if desired.

`w` - write a game log file

`q` - quit

# Features

### Run Categorization

Runs will be marked as PB Category if any preparation has been done before starting the quest. Examples include:
* Mylla & Youlla (Twins PB) has been used before quest start
* Players have over 5% PB charged
* In quests without a starting console (eg. Mop-up and Sweep-up series)
    - Any player is below 95% HP
    - Shifta/Deband has been cast


# To Do

* Detect solo mode
* new player registration
* merge multiplayer games
* detect non-vanilla weapons
* quest splits?
* aggregate stats
  - player
  - overall
* ep4 boss trainwreck
* multiple units being combined
* pb graph
* death gunners https://discord.com/channels/209828685052641281/810649151058214963/848266650251952138

# Package Structure

    .
    ├── client                  # The PSO Stats Client
    │   ├── cmd                 # The main function for the client 
    │   └── internal            # Private packages for the client only 
    │       ├── client          # Main client logic
    │       ├── consoleui       # Draws current game state to the terminal
    │       ├── numbers         # Reads blocks pso-internal memory and parses into go primitives
    │       └── pso             # Interaction with PSO exe
    ├── pkg                     # Public go packages used by the client and server
    │   └── model               # Golang models representing public client and server data
    ├── server                  # The PSO Stats Server
    │   ├── cmd                 # The main function for the server 
    │   └── internal            # Private packages for the server only 
    │       ├── db              # Game database interaction
    │       ├── server          # TODO: ???
    │       └── userdb          # Database layer for users and guildcard mapping
    └── winres                  # Windows exe config

# Building client

```shell
# Generate syso files
go-winres make
mv rsrc_windows*.syso client/cmd/
# Build exe
cd client/cmd
go build -o psostats.exe
```