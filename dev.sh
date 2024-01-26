#!/bin/bash

# Start a new tmux session
tmux new-session -d -s cigareSession

# Bind 'Ctrl + x' to kill the session
tmux bind-key -n C-x kill-session

# Split the window horizontally for 'air' and 'make browser-refresh'
tmux split-window -h

# Split the right pane vertically
tmux select-pane -t 1
tmux split-window -v

# Run 'air' in the top right pane
tmux send-keys -t cigareSession:0.1 'air' C-m

# Run 'make browser-refresh' in the bottom right pane
tmux send-keys -t cigareSession:0.2 'make browser-refresh' C-m

# Split the bottom right pane horizontally to create a third row
tmux select-pane -t 2
tmux split-window -h

# Run 'npx tailwindcss -i ./web/view/styles.css -o ./web/static/styles.css --watch' in the new pane
tmux send-keys -t cigareSession:0.3 'npx tailwindcss -i ./web/view/styles.css -o ./web/static/styles.css --watch' C-m

# Resize the bottom right pane to be smaller
tmux resize-pane -t 2 -y 20%

# Select the left pane for typing
tmux select-pane -t 0

# Attach to the session
tmux attach-session -t cigareSession
