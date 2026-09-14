module go.abhg.dev/goldmark/mermaid/demo

go 1.25

toolchain go1.25.5

replace go.abhg.dev/goldmark/mermaid => ../

require (
	github.com/yuin/goldmark/v2 v2.0.1
	go.abhg.dev/goldmark/mermaid v0.6.0
)

require github.com/yuin/goldmark v1.7.16 // indirect
