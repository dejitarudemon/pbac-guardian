/*
Package tests contains tests for public API of base package.

Tests check Effect value validation and correctness of working
with ALLOW and DENY effects through public IsValid method.
The tests verify that only Effect_ALLOW and Effect_DENY are
considered valid effect values.
*/
package tests

import (
	"testing"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
TestEffectIsValid tests the public IsValid method for checking Effect value validity.

The test checks:
  - Validity of Effect_ALLOW ("allow")
  - Validity of Effect_DENY ("deny")
  - Invalidity of arbitrary strings
  - Invalidity of empty string
*/
func TestEffectIsValid(t *testing.T) {
	tests := []struct {
		name   string
		effect base.Effect
		want   bool
	}{
		{"valid allow", base.Effect_ALLOW, true},          // Valid allow effect
		{"valid deny", base.Effect_DENY, true},            // Valid deny effect
		{"invalid effect", base.Effect("invalid"), false}, // Invalid effect
		{"empty effect", base.Effect(""), false},          // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.effect.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
