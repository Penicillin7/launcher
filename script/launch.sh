#!/bin/bash

CS2ExecPath="$1"
Port="$2"

# Check if the CS2 executable exists
if [ ! -f "$CS2ExecPath" ]; then
  echo "CS2 executable not found at $CS2ExecPath"
  exit 1
fi

# Check if the CS2 executable is executable
if [ ! -x "$CS2ExecPath" ]; then
  echo "CS2 executable is not executable"
  exit 1
fi

# Check port is available
if lsof -Pi :$Port -sTCP:LISTEN -t >/dev/null; then
  echo "Port $Port is already in use"
  exit 1
fi

# Run the CS2 executable
$CS2ExecPath -dedicated -maxplayers 10 -console +game_type 0 +game_mode 1 +mapgroup mg_active +map de_dust2 -high -port $Port -ip 0.0.0.0
