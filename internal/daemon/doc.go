// Package daemon provides the core watch loop for the portwatch CLI.
//
// A Daemon is constructed with a [config.Config] and an [alert.Dispatcher],
// then started via Run. On each tick the daemon:
//
//  1. Scans the current open TCP/UDP ports via [ports.Scanner].
//  2. Filters out ports listed in Config.IgnorePorts.
//  3. Compares the result against the persisted [snapshot.Snapshot].
//  4. Dispatches an alert for every port that appeared or vanished.
//  5. Persists the new snapshot for the next tick.
//
// The loop runs until the supplied context is cancelled.
package daemon
