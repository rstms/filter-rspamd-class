package main

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/rstms/rspamd-classes/classes"
    )

func TestVersion(t *testing.T) {
    SpamClasses, err := classes.New("")
	require.Nil(t, err)
	require.NotEmpty(t, classes.Version)
	require.NotNil(t, SpamClasses)
}
