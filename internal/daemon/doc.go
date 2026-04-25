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
//
// # Tick interval
//
// The interval between scans is controlled by Config.Interval. If the scan
// or alert dispatch takes longer than the interval, the next tick begins
// immediately without drift accumulation (i.e. the ticker is reset, not
// queued).
//
// # Error handling
//
// Transient scan errors are logged and skipped; the snapshot is not updated
// so that ports from the previous successful scan are retained for comparison
// on the next tick. Persistent alert-dispatch errors are logged but do not
// stop the loop.
package daemon
