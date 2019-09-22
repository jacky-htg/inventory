package databasetest

import (
	"bytes"
	"os/exec"
	"testing"
)

// StartContainer runs a mysql container to execute commands.
func StartContainer(t *testing.T) {
	t.Helper()

	cmd := exec.Command("docker", "run", "-d", "--name", "rebel_mysql", "--publish", "33060:3306", "--env", "MYSQL_ROOT_PASSWORD=1234", "--env", "MYSQL_DATABASE=rebel_db", "mysql:8")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start docker : %v", err)
	}

}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T) {
	t.Helper()

	if err := exec.Command("docker", "container", "rm", "-f", "rebel_mysql").Run(); err != nil {
		t.Fatalf("could not stop mysql container: %v", err)
	}
}
