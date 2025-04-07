Dev Menu
---
- TODO: Readme
- TMUX Go Bindings
- Dev Menu via configuration

TODO
---
- titles
``
printf '\033]2;%s\033\\' 'title goes herasd'
tmux set -g pane-border-format "#{?pane_active,#[reverse],}#{pane_index}#[default] \"#{pane_title}\""
``

