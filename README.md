# gocs-core

Go library containing the core utilities for analyzing Counter-Strike: Global Offensive demo files (*.dem files).
Currently supports the following features:
 - Generate shot accuracy stats for a player
 - Generate damage summary for a player (categorized by weapon and head-to-head against opponent players)
 - Generate K/D (categorized by weapon and head-to-head against opponent players)
 - Generate heatmap for K/D and bomb plant events for a player

More stats on the way. Feature suggestions are welcome.

<br/>
<br/>

## Web implementation
 - [gocs-grpc-server](https://github.com/AbhilashJN/gocs-grpc-server): A gRPC server which exposes the apis from this library as gRPC services.
 - [gocs-ui-web](https://github.com/AbhilashJN/gocs-ui-web):
 A ReactJS UI for this service, using [grpc-web](https://github.com/grpc/grpc-web).

<br/>
<br/>

## Desktop app implementation
[gocs-w](https://github.com/AbhilashJN/gocs-w): Cross platform desktop app built using Go and ReactJS using the [Wails framework](https://github.com/wailsapp/wails).
