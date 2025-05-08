# glick
A go project which focus on productivity by easily letting you create CU tasks

## Installation (Linux)
Build the project
1. go build -o glick main.go
2. sudo mv glick /usr/local/bin

If you want to create a shortcut to start it in i3

3. Create your launcher:  nvim ~/scripts/launch_glick.sh

```sh
gnome-terminal --hide-menubar --geometry=80x1 --title="glick" -- bash -c "set +H; glick; exec bash"

# Give the terminal time to open
sleep 0.1
# Resize and move the window (floating only)
i3-msg '[title="glick"] floating enable, move position center'
```
4. chmod +x launch_my_hello.sh
5. cd ~/.config/i3/config
6. bindsym $mod+t exec --no-startup-id ~/scripts/launch_glick.sh
