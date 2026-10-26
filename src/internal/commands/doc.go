// Package commands holds one Cobra verb per file, each registering itself on
// the root in its own init(), plus the server-side actions the `workstation config` items
// call (`server_*.go`) and hand to the registry at startup.
package commands
