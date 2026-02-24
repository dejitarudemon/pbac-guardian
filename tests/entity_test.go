/*
Package tests contains tests for public API of base package.

Tests check Entity value validation and correctness of working
with source and target entities through public IsValid method.
The tests verify that only Entity_SOURCE and Entity_TARGET are
considered valid entity values.
*/
package tests

import (
	"testing"

	"github.com/dejitarudemon/pbac-guardian/internal/base"
)

/*
TestEntityIsValid tests the public IsValid method for checking Entity value validity.

The test checks:
  - Validity of Entity_SOURCE ("source")
  - Validity of Entity_TARGET ("target")
  - Invalidity of arbitrary strings
  - Invalidity of empty string
*/
func TestEntityIsValid(t *testing.T) {
	tests := []struct {
		name   string
		entity base.Entity
		want   bool
	}{
		{"valid source", base.Entity_SOURCE, true},        // Valid source entity
		{"valid target", base.Entity_TARGET, true},        // Valid target entity
		{"invalid entity", base.Entity("invalid"), false}, // Invalid entity
		{"empty entity", base.Entity(""), false},          // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.entity.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
