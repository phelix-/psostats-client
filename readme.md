# Running

Configure `config.yaml` if desired. If localUi is enabled and given network access the ui will be at `http://localhost:8081/` and you can view raw game data at `http://localhost:8081/game/info`

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
* equipped mags
* separate games
* detect non-vanilla weapons
* boss kill times
* CMode QoL
  - Floor names are wrong I think
* weapon/armor timeline? show unit/shield switches?

# Package Structure

    .
    ├── client                  # The PSO Stats Client
    │   ├── cmd                 # The main function for the client 
    │   └── internal            # Private packages for the client only 
    │       ├── client          # TODO: ???
    │       ├── consoleui       # TODO: ???
    │       ├── numbers         # TODO: ???
    │       └── pso             # TODO: ???
    ├── pkg                     # Public go packages used by the client and server
    │   └── model               # Golang models representing public client and server data
    ├── server                  # The PSO Stats Server
    │   ├── cmd                 # The main function for the server 
    │   └── internal            # Private packages for the server only 
    │       └── server          # TODO: ???
    └── winres                  # TODO: fill
