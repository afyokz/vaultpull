package vault

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// ChecksumResult holds the checksum for a set of secrets.
type ChecksumResult struct {
	Algorithm string
	Digest    string
	KeyCount  int
}

func (r ChecksumResult) String() string {
	return fmt.Sprintf("%s:%s (keys: %d)", r.Algorithm, r.Digest, r.KeyCount)
}

// ComputeChecksum produces a deterministic SHA-256 digest over the provided
// secrets map. Keys are sorted before hashing so the result is stable
// regardless of map iteration order.
func ComputeChecksum(secrets map[string]string) ChecksumResult {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		// Feed key=value\n into the hasher.
		fmt.Fprintf(h, "%s=%s\n", k, secrets[k])
	}

	digest := hex.EncodeToString(h.Sum(nil))
	return ChecksumResult{
		Algorithm: "sha256",
		Digest:    digest,
		KeyCount:  len(keys),
	}
}

// VerifyChecksum recomputes the checksum for secrets and compares it against
// the expected digest string ("sha256:<hex>" or bare hex). Returns true when
// they match.
func VerifyChecksum(secrets map[string]string, expected string) bool {
	result := ComputeChecksum(secrets)
	expected = strings.TrimPrefix(expected, "sha256:")
	return strings.EqualFold(result.Digest, expected)
}

// FormatChecksumReport returns a human-readable checksum report.
func FormatChecksumReport(result ChecksumResult) string {
	var sb strings.Builder
	sb.WriteString("Checksum Report\n")
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	fmt.Fprintf(&sb, "Algorithm : %s\n", result.Algorithm)
	fmt.Fprintf(&sb, "Digest    : %s\n", result.Digest)
	fmt.Fprintf(&sb, "Keys      : %d\n", result.KeyCount)
	return sb.String()
}
