#!/bin/bash

tmux new-session -d -s cigareSession
tmux bind-key -n C-x kill-session

tmux split-window -h
tmux select-pane -t 1
tmux split-window -v

tmux send-keys -t cigareSession:0.1 'export $(cat dev.env | xargs) && air' C-m

tmux send-keys -t cigareSession:0.2 'make browser-refresh' C-m

tmux select-pane -t 2
tmux split-window -h

tmux send-keys -t cigareSession:0.3 'npx tailwindcss -i ./web/view/styles.css -o ./web/static/styles.css --watch' C-m

tmux resize-pane -t 2 -y 20%

tmux select-pane -t 0

tmux attach-session -t cigareSession
