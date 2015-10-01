// +build windows

package windows

import (
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/execdriver"
	"github.com/microsoft/hcsshim"
)

// Terminate implements the exec driver Driver interface.
// This is only called from Register() on the daemon.
func (d *Driver) Terminate(p *execdriver.Command) error {
	logrus.Debugf("WindowsExec: Terminate() id=%s", p.ID)
	return kill(p.ID, p.ContainerPid, int(syscall.SIGKILL))
}

// Kill implements the exec driver Driver interface.
func (d *Driver) Kill(p *execdriver.Command, sig int) error {
	logrus.Debugf("WindowsExec: Kill() id=%s sig=%d", p.ID, sig)
	return kill(p.ID, p.ContainerPid, sig)
}

func kill(id string, pid int, sig int) error {
	logrus.Debugln("kill() ", id, pid, sig)
	var err error

	// Terminate Process
	if err = hcsshim.TerminateProcessInComputeSystem(id, uint32(pid)); err != nil {
		// Don't log warning if the process doesn't exist (ie has already exited)
		if (!strings.Contains(err.Error(), "2147942403")) && // 0x80070003
			(!strings.Contains(err.Error(), "2147943568")) { // 0x80070490
			logrus.Warnf("Kill - failed to terminate pid %d in %s: %q.", pid, id, err)
		}
		// Ignore errors regardless.
		err = nil
	}

	if sig == int(syscall.SIGKILL) || forceKill {
		// Terminate the compute system
		if err = hcsshim.TerminateComputeSystem(id); err != nil {
			logrus.Errorf("Failed to terminate %s - %q", id, err)
		}

	} else {
		// Shutdown the compute system
		if err = hcsshim.ShutdownComputeSystem(id); err != nil {
			logrus.Errorf("Failed to shutdown %s - %q", id, err)
		}
	}
	return err
}
