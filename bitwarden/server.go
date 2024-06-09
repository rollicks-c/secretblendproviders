package bitwarden

import (
	"bytes"
	"fmt"
	"github.com/rollicks-c/configcove"
	"github.com/rollicks-c/secretblendproviders/internal/network"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type apiServer struct {
	dataDir string
}

const (
	appName = "sb.bitwarden"
)

func (s apiServer) start() error {

	// re-use existing
	pidCurrent, ok := configcove.Store(appName).LoadNumber("bitwarden.pid")
	if ok {
		if s.checkForRunningInstance(pidCurrent) {
			return nil
		}
	}

	// create new
	pidNew, err := s.startInstance()
	if err != nil {
		return err
	}
	configcove.Store(appName).SaveNumber("bitwarden.pid", pidNew)

	return nil
}

func (s apiServer) startInstance() (int, error) {

	// prep command
	port := network.FindFreePort(8087)
	cmd := exec.Command("bw", "serve", "--port", fmt.Sprintf("%d", port))
	cmd.Env = append(cmd.Env, fmt.Sprintf("BITWARDENCLI_APPDATA_DIR=%s", s.dataDir))
	errorReader := &bytes.Buffer{}
	cmd.Stderr = errorReader

	// own process group to not tear down after exit parent
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Setpgid = true

	// await start
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	if err := s.awaitStart(cmd); err != nil {
		return 0, err
	}
	if err := s.awaitReadiness(errorReader); err != nil {
		return 0, err
	}

	return cmd.Process.Pid, nil
}

func (s apiServer) checkForRunningInstance(pID int) bool {
	proc, err := os.FindProcess(pID)
	if err != nil {
		return false
	}
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return false
	}
	return true
}

func (s apiServer) awaitStart(cmd *exec.Cmd) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if cmd.Process != nil {
				return nil
			}
		case <-time.After(time.Second * 1):
			return fmt.Errorf("timeout waiting for server to start")
		}
	}
}

func (s apiServer) awaitReadiness(errOut *bytes.Buffer) error {
	<-time.After(time.Second * 1)
	if errOut.Len() > 0 {
		return fmt.Errorf("error starting server: %s", errOut.String())
	}
	return nil
}

func (s apiServer) getURL() string {
	return "http://localhost:8087"
}
