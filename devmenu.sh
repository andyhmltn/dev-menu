tmux new-session -d -s DevMenu
go run ./tmux/main.go
tmux attach -t DevMenu
