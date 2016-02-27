package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zabawaba99/gomutate/mutants"
)

func TestParseOptions(t *testing.T) {
	os.Args = []string{
		"binary",
		"-d",
		"-m", "conditionals-boundary",
	}

	opts := parseOptions()
	require.NotNil(t, opts)
	assert.True(t, opts.Debug)

	require.Len(t, opts.Mutator, 1)
	assert.Equal(t, opts.Mutator[0], "conditionals-boundary")
}

func TestGetMutators(t *testing.T) {
	opts := Options{
		Mutator: []string{"conditionals-boundary", "lol-wat?"},
	}

	mutators := opts.getMutators()
	require.Len(t, mutators, 1)
	assert.EqualValues(t, mutators[0], &mutants.ConditionalsBoundary{})
}

func TestGetPkgs(t *testing.T) {
	require.NoError(t, os.Mkdir("sub", 0777))
	require.NoError(t, os.Mkdir("_sub", 0777))
	defer func() {
		os.RemoveAll("sub")
		os.RemoveAll("_sub")
	}()

	opts := Options{
		packages: []string{".", ".", "./...", "./foo"},
	}

	pkgs := opts.getPkgs()

	assert.Len(t, pkgs, 3)
	assert.Contains(t, pkgs, "")
	assert.Contains(t, pkgs, "foo")
	assert.Contains(t, pkgs, "sub")
}

func TestDedup(t *testing.T) {
	dup := []string{".", ".", "lolz", "bar", "bar"}
	result := dedup(dup)

	assert.Len(t, result, 3)
	assert.Contains(t, result, ".")
	assert.Contains(t, result, "lolz")
	assert.Contains(t, result, "bar")
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		Value   string
		Prefixs []string
		Result  bool
	}{
		{
			Value:   "fooy",
			Prefixs: []string{"a", "f"},
			Result:  true,
		},
		{
			Value:   "_gaga",
			Prefixs: []string{"_"},
			Result:  true,
		},
		{
			Value:   "asdasd",
			Prefixs: []string{},
			Result:  false,
		},
		{
			Value:   "yahoo",
			Prefixs: []string{"o", "h", "a"},
			Result:  false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Result, hasPrefix(test.Value, test.Prefixs...))
	}
}
