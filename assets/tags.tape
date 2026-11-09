# Run from the repo root:
#   vhs assets/tags.tape && gif2webp -lossy -q 70 -m 6 assets/tags.gif -o assets/tags.webp && rm assets/tags.gif
Output assets/tags.gif

Require go

Set Shell "bash"
Set FontFamily "JetBrains Mono"
Set FontSize 16
Set Width 760
Set Height 450
Set Padding 24
Set WindowBar Colorful
Set Theme {"name":"go-masker","black":"#484f58","red":"#ff7b72","green":"#7ee787","yellow":"#e3b341","blue":"#79c0ff","magenta":"#d2a8ff","cyan":"#4fd1e8","white":"#b1bac4","brightBlack":"#6e7681","brightRed":"#ffa198","brightGreen":"#56d364","brightYellow":"#e3b341","brightBlue":"#a5d6ff","brightMagenta":"#e2c5ff","brightCyan":"#b3f0ff","brightWhite":"#ffffff","background":"#0d1117","foreground":"#e6edf3","selection":"#0f3a47","cursor":"#4fd1e8"}
Set TypingSpeed 30ms

Hide
Type "export BASH_SILENCE_DEPRECATION_WARNING=1 && go build ./examples/... && export PS1='\[\e[36m\]$\[\e[0m\] ' && clear"
Enter
Sleep 2s
Show

Type `sed -n '/^type Account/,/^}/p' examples/tour/main.go`
Sleep 250ms
Enter
Sleep 2s
Type `go run ./examples/tour tags`
Sleep 250ms
Enter
Sleep 3.5s
