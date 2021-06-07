# Running

Configure `config.yaml` if desired.

`w` - write a game log file

`q` - quit

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
