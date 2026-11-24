package main

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	log.Println("Starting E2E test...")

	log.Println("Waiting for the hello deployment to become ready")
	out, err := exec.Command("kubectl", "rollout", "status", "deployment/hello", "--timeout=300s").CombinedOutput()
	if err != nil {
		t.Fatalf("deployment rollout did not succeed: %s", out)
	}
	log.Printf("Rollout status: %s", out)

	log.Println("Checking hello pod phase")
	out, err = exec.Command("kubectl", "get", "pods", "-l", "app=hello", "-o", "jsonpath={.items[*].status.phase}").CombinedOutput()
	if err != nil {
		t.Fatalf("kubectl get pods failed: %s", out)
	}
	log.Println(out)

	if !strings.Contains(string(out), "Running") {
		t.Errorf("hello pod is not running: %s", out)
	}
}
