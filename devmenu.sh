tmux new-session -d -s DevMenu
go run ./cmd/main.go
tmux attach -t DevMenu
