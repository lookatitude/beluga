package main

import "testing"

func TestSetup(t *testing.T) {
	t.Run("Basic Test", func(t *testing.T) {
		if 1+1 != 2 {
			t.Error("Basic math test failed")
		}
	})
}