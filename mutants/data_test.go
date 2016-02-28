package mutants

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDataLoad(t *testing.T) {
	s := "test_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	require.NoError(t, os.MkdirAll(s, 0777))
	defer os.RemoveAll(s)

	js := `
	{
		"filename":"/foo.go",
		"line_number":2,
		"original":">",
		"mutation":"<=",
		"type":"conditionals-boundary",
		"killed":true
	}`
	f, err := os.Create(filepath.Join(s, DataFileName))
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(js)
	require.NoError(t, err)

	data := new(Data)
	err = data.Load(s)
	require.NoError(t, err)

	assert.Equal(t, "/foo.go", data.Filename)
	assert.EqualValues(t, 2, data.LineNumber)
	assert.Equal(t, ">", data.Original)
	assert.Equal(t, "<=", data.Mutation)
	assert.Equal(t, "conditionals-boundary", data.Type)
	assert.Equal(t, true, data.Killed)
}

func TestDataSave(t *testing.T) {
	s := "test_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	require.NoError(t, os.MkdirAll(s, 0777))
	defer os.RemoveAll(s)

	data := Data{
		Filename:   "/foo.go",
		LineNumber: 2,
		Original:   "<",
		Mutation:   ">=",
		Type:       "foobar",
		Killed:     false,
	}

	err := data.Save(s)
	require.NoError(t, err)

	f, err := os.Open(filepath.Join(s, DataFileName))
	require.NoError(t, err)
	defer f.Close()

	var loaded Data
	err = json.NewDecoder(f).Decode(&loaded)
	require.NoError(t, err)

	assert.Equal(t, data, loaded)
}
