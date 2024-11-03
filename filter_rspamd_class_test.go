package main

import (
	"github.com/rstms/rspamd-classes/classes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVersion(t *testing.T) {
	SpamClasses, err := classes.New("")
	require.Nil(t, err)
	require.NotEmpty(t, classes.Version)
	require.NotNil(t, SpamClasses)
}
