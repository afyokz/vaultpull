package vault

import (
	"strings"
	"testing"
)

func TestComputeChecksum_Empty(t *testing.T) {
	result := ComputeChecksum(map[string]string{})
	if result.Algorithm != "sha256" {
		t.Errorf("expected sha256, got %s", result.Algorithm)
	}
	if result.KeyCount != 0 {
		t.Errorf("expected 0 keys, got %d", result.KeyCount)
	}
	if result.Digest == "" {
		t.Error("expected non-empty digest for empty input")
	}
}

func TestComputeChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r1 := ComputeChecksum(secrets)
	r2 := ComputeChecksum(secrets)
	if r1.Digest != r2.Digest {
		t.Errorf("checksum not deterministic: %s vs %s", r1.Digest, r2.Digest)
	}
}

func TestComputeChecksum_OrderIndependent(t *testing.T) {
	a := map[string]string{"KEY1": "v1", "KEY2": "v2"}
	b := map[string]string{"KEY2": "v2", "KEY1": "v1"}
	rA := ComputeChecksum(a)
	rB := ComputeChecksum(b)
	if rA.Digest != rB.Digest {
		t.Errorf("checksum should be order-independent: %s vs %s", rA.Digest, rB.Digest)
	}
}

func TestComputeChecksum_ChangesOnValueUpdate(t *testing.T) {
	before := map[string]string{"TOKEN": "abc"}
	after := map[string]string{"TOKEN": "xyz"}
	if ComputeChecksum(before).Digest == ComputeChecksum(after).Digest {
		t.Error("digest should differ when value changes")
	}
}

func TestVerifyChecksum_Match(t *testing.T) {
	secrets := map[string]string{"A": "1", "B": "2"}
	result := ComputeChecksum(secrets)
	if !VerifyChecksum(secrets, result.Digest) {
		t.Error("expected verification to pass with bare digest")
	}
	if !VerifyChecksum(secrets, "sha256:"+result.Digest) {
		t.Error("expected verification to pass with prefixed digest")
	}
}

func TestVerifyChecksum_Mismatch(t *testing.T) {
	secrets := map[string]string{"A": "1"}
	if VerifyChecksum(secrets, "deadbeef") {
		t.Error("expected verification to fail with wrong digest")
	}
}

func TestFormatChecksumReport_ContainsFields(t *testing.T) {
	secrets := map[string]string{"X": "y"}
	report := FormatChecksumReport(ComputeChecksum(secrets))
	for _, want := range []string{"sha256", "Digest", "Keys", "1"} {
		if !strings.Contains(report, want) {
			t.Errorf("report missing %q:\n%s", want, report)
		}
	}
}
