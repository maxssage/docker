package daemon

import (
	"fmt"
	"syscall"
)

// ContainerKill is the API entrypoint for Kill. Note Windows does not
// really use signals as the concept doesn't largely exist in Windows.
// However, HCS does support a graceful shutdown and a more forceful terminate
// of a container which can roughly be equated to SIGTERM and SIGKILL
// We default to graceful container shutdown.
func (daemon *Daemon) ContainerKill(name string, sig uint64) error {
	container, err := daemon.Get(name)
	if err != nil {
		return err
	}

	switch sig {
	case 0:
		sig = uint64(syscall.SIGTERM)
	case uint64(syscall.SIGTERM), uint64(syscall.SIGKILL):
	default:
		return fmt.Errorf("Windows only supports %d (graceful container shutdown) and %d (forced container kill) 'signals'", syscall.SIGTERM, syscall.SIGKILL)
	}

	if err := daemon.kill(container, int(sig), 10); err != nil {
		return err
	}

	container.logEvent("kill")
	return nil
}
