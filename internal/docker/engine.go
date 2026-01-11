package docker

import (
	"fmt"
	"os/exec"
)

// Runner handles docker operations
type Runner struct{}

// NewRunner creates a new docker runner
func NewRunner() *Runner {
	return &Runner{}
}

// ExecuteComposeUp runs docker compose up
func (r *Runner) ExecuteComposeUp(project, dir string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "-p", project, "up", "-d", "--pull", "always", "--force-recreate")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// ExecuteComposeDown runs docker compose down
func (r *Runner) ExecuteComposeDown(project, dir string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "-p", project, "down", "-v", "--remove-orphans")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// ExecuteComposePull runs docker compose pull
func (r *Runner) ExecuteComposePull(project, dir string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "-p", project, "pull")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// ExecuteComposeRestart runs docker compose restart
func (r *Runner) ExecuteComposeRestart(project, dir string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "-p", project, "restart")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// ExecuteComposePs returns the status of containers
func (r *Runner) ExecuteComposePs(project, dir string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "-p", project, "ps", "--format", "json")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// ExecuteComposeLogs returns the logs of containers
func (r *Runner) ExecuteComposeLogs(project, dir string, tail int) ([]byte, error) {
	tailStr := fmt.Sprintf("%d", tail)
	cmd := exec.Command("docker", "compose", "-p", project, "logs", "--tail", tailStr, "--no-color")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// ExecuteComposeImages returns the images used by the services
func (r *Runner) ExecuteComposeImages(project, dir string) ([]byte, error) {
	cmd := exec.Command("docker", "compose", "-p", project, "images", "--format", "json")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// IsInstalled checks if docker and compose are available
func (r *Runner) IsInstalled() bool {
	cmd := exec.Command("docker", "compose", "version")
	err := cmd.Run()
	return err == nil
}

// ExecuteMongoScript executes a script in a mongo container
func (r *Runner) ExecuteMongoScript(host, script string) ([]byte, error) {
	// args for docker run
	// --rm: remove container after run
	// --network host: use host network to reach the mongo instance
	// mongo:latest: image to use
	// mongosh ...: command to run
	args := []string{
		"run", "--rm", "--network", "host",
		"mongo:latest",
		"mongosh",
		host,
		"--eval", script,
	}

	cmd := exec.Command("docker", args...)
	return cmd.CombinedOutput()
}
