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

import "testing"

func TestGetGitCommitSHA(t *testing.T) {
	sha := GetGitCommitSHA()
	// Either a 12-char hex string or "unknown".
	if sha != "unknown" && len(sha) != 12 {
		t.Errorf("SHA should be 12 chars or 'unknown', got %q (len=%d)", sha, len(sha))
	}
}

func TestComputeHash(t *testing.T) {
	// SHA-256 of empty string.
	h := ComputeHash([]byte{})
	want := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if h != want {
		t.Errorf("got %q, want %q", h, want)
	}
}

func TestComputeHashDeterministic(t *testing.T) {
	data := []byte("hello world")
	h1 := ComputeHash(data)
	h2 := ComputeHash(data)
	if h1 != h2 {
		t.Error("ComputeHash should be deterministic")
	}
	if len(h1) != 64 {
		t.Errorf("SHA-256 hex should be 64 chars, got %d", len(h1))
	}
}
