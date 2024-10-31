package main

import (
	"fmt"
	"math"
	"testing"
)

func TestSpamLevelsNoFile(t *testing.T) {
	err := initClassLevels("")
	if err != nil {
		t.Errorf("initClassLevels failed: %v", err)
	}

	fmt.Printf("SpamClassLevels: %v\n", SpamClassLevels)
	if len(SpamClassLevels) != 1 {
		t.Errorf("lengh = %d; expected 1", len(SpamClassLevels))
	}
	expected := []SpamClass{{"ham", 0.0}, {"possible", 3.0}, {"probable", 10.0}, {"spam", math.MaxFloat32}}
	for key, value := range SpamClassLevels {
		if key != "default" {
			t.Errorf("key = %s; expected %s", key, "default")
		}
		if len(value) != len(expected) {
			t.Errorf("value = %v; expected %v", value, expected)
		}
		for i, level := range value {
			if level != expected[i] {
				t.Errorf("level[%d]= %v; expected %v", i, level, expected[i])
			}
		}
	}
}

func TestSpamLevelsConfigFile(t *testing.T) {
	err := initClassLevels("testdata/rspamd_classes.json")
	if err != nil {
		t.Errorf("initClassLevels failed: %v", err)
	}
	fmt.Printf("SpamClassLevels: %v\n", SpamClassLevels)

	if len(SpamClassLevels) != 2 {
		t.Errorf("length = %d; expected 2", len(SpamClassLevels))
	}

	defaultLevels, ok := SpamClassLevels["default"]
	if !ok {
		t.Errorf("missing default")
	}
	userLevels, ok := SpamClassLevels["username@example.org"]
	if !ok {
		t.Errorf("missing default")
	}
	if len(defaultLevels) != len(userLevels) {
		t.Errorf("length mismatch: default=%d user=%d", len(defaultLevels), len(userLevels))
	}
	for i, level := range defaultLevels {
		if level != userLevels[i] {
			t.Errorf("level[%d]= %v; expected %v", i, level, userLevels[i])
		}
	}
}
