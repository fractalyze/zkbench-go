// Copyright 2026 The zkbench-go Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// ==============================================================================

package zkbench

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

// GetGitCommitSHA returns the current git commit SHA (first 12 characters).
// Returns "unknown" if not in a git repository.
func GetGitCommitSHA() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "unknown"
	}
	sha := strings.TrimSpace(string(out))
	if len(sha) > 12 {
		sha = sha[:12]
	}
	return sha
}

// ComputeHash returns the SHA-256 hex digest of the given bytes.
func ComputeHash(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}
