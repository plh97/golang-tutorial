package casbinPkg

import (
	"log"
	"os" // New Import
	"path/filepath"
	// Removed: "runtime"
	// Removed: "path/filepath" (still used, but no need for runtime)

	"github.com/casbin/casbin/v2"
)

func NewEnforcer() *casbin.Enforcer {
	// --- 1. Determine the Base Path ---
	// Get the current working directory, which is typically the project root
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	// --- 2. Construct Full Paths Relative to Root/config ---
	// Assumes: your-project-root/config/model.conf
	modelPath := filepath.Join(cwd, "config", "model.conf")
	// Assumes: your-project-root/config/policy.csv
	policyPath := filepath.Join(cwd, "config", "policy.csv")

	// --- 3. Initialize Enforcer ---
	e, err := casbin.NewEnforcer(modelPath, policyPath) // Use NewEnforcer for file-based policies
	if err != nil {
		// Log which paths failed to load
		log.Fatalf("failed to create Casbin enforcer. Check paths: Model=%s, Policy=%s. Error: %v", modelPath, policyPath, err)
	}

	e.LoadPolicy()
	log.Println("Casbin Enforcer successfully initialized from config/ paths.")
	
	// Ensure the log indicates the path used
	log.Printf("Loaded Casbin model from: %s", modelPath)
	
	return e
}