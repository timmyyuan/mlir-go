package gen

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ResolveIncludeRoot resolves the MLIR include root, defaulting to llvm-config.
func ResolveIncludeRoot(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	llvmConfig := os.Getenv("LLVM_CONFIG")
	if llvmConfig == "" {
		llvmConfig = "llvm-config"
	}
	out, err := exec.Command(llvmConfig, "--includedir").Output()
	if err != nil {
		return "", fmt.Errorf("resolve include root via %s --includedir: %w", llvmConfig, err)
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("%s --includedir returned an empty path", llvmConfig)
	}
	return root, nil
}
