package classes

import (
	"fmt"
	"math"
	"testing"
)

func TestSpamLevelsNoFile(t *testing.T) {
	SpamClasses, err := New("")
	if err != nil {
		t.Errorf("initClassLevels failed: %v", err)
	}

	fmt.Printf("SpamClasses: %v\n", SpamClasses)
	if len(SpamClasses.Classes) != 1 {
		t.Errorf("lengh = %d; expected 1", len(SpamClasses.Classes))
	}
	expected := []SpamClass{{"ham", 0.0}, {"possible", 3.0}, {"probable", 10.0}, {"spam", math.MaxFloat32}}
	for key, value := range SpamClasses.Classes {
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
	SpamClasses, err := New("testdata/rspamd_classes.json")
	if err != nil {
		t.Errorf("classes.New failed: %v", err)
	}
	fmt.Printf("SpamClasses: %v\n", SpamClasses)

	if len(SpamClasses.Classes) != 2 {
		t.Errorf("length = %d; expected 2", len(SpamClasses.Classes))
	}

	defaultLevels, ok := SpamClasses.Classes["default"]
	if !ok {
		t.Errorf("missing default")
	}
	userLevels, ok := SpamClasses.Classes["username@example.org"]
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

func TestWrite(t *testing.T) {
	SpamClasses, err := New("testdata/rspamd_classes.json")
	if err != nil {
		t.Errorf("New failed: %v", err)
	}
	err = SpamClasses.Write("testdata/output.json")
	if err != nil {
		t.Errorf("Write failed: %v", err)
	}
	readback, err := New("testdata/output.json")
	if err != nil {
		t.Errorf("New failed: %v", err)
	}
	for key, classes := range SpamClasses.Classes {
		rclasses, ok := readback.Classes[key]
		if !ok {
			t.Errorf("readback key not found: %s\n", key)
		}
		fmt.Printf("key=%s classes: %v\n", key, classes)
		fmt.Printf("key=%s rclasses: %v\n", key, rclasses)
		for i, class := range classes {
			if class != rclasses[i] {
				t.Errorf("key=%s classes[%d] (%v) mismatches rclasses[%d] (%v)\n", key, i, class, i, rclasses[i])
			}
		}
	}
}
