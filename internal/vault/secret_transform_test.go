package vault

import (
	"testing"
)

func TestNewTransformRule_MutuallyExclusive(t *testing.T) {
	_, err := NewTransformRule(WithUppercase(), WithLowercase())
	if err == nil {
		t.Fatal("expected error for uppercase+lowercase")
	}
}

func TestTransform_NoOp(t *testing.T) {
	r, _ := NewTransformRule()
	in := map[string]string{"KEY": "val"}
	out := r.Apply(in)
	if out["KEY"] != "val" {
		t.Errorf("expected val, got %s", out["KEY"])
	}
}

func TestTransform_Uppercase(t *testing.T) {
	r, _ := NewTransformRule(WithUppercase())
	out := r.Apply(map[string]string{"db_host": "localhost"})
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST key")
	}
}

func TestTransform_Lowercase(t *testing.T) {
	r, _ := NewTransformRule(WithLowercase())
	out := r.Apply(map[string]string{"DB_HOST": "localhost"})
	if _, ok := out["db_host"]; !ok {
		t.Error("expected db_host key")
	}
}

func TestTransform_PrefixSuffix(t *testing.T) {
	r, _ := NewTransformRule(WithPrefix("APP_"), WithSuffix("_V1"))
	out := r.Apply(map[string]string{"KEY": "val"})
	if _, ok := out["APP_KEY_V1"]; !ok {
		t.Errorf("expected APP_KEY_V1, got keys: %v", out)
	}
}

func TestTransform_CombinedUppercasePrefix(t *testing.T) {
	r, _ := NewTransformRule(WithUppercase(), WithPrefix("SVC_"))
	out := r.Apply(map[string]string{"secret": "x"})
	if _, ok := out["SVC_SECRET"]; !ok {
		t.Errorf("expected SVC_SECRET, got %v", out)
	}
}

func TestTransform_EmptyInput(t *testing.T) {
	r, _ := NewTransformRule(WithPrefix("X_"))
	out := r.Apply(map[string]string{})
	if len(out) != 0 {
		t.Error("expected empty output")
	}
}
