package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SQLiteFileNotExist(t *testing.T) {
	ds, err := New("testdata/chinook2.db")

	require.EqualError(t, err, "stat testdata/chinook2.db: no such file or directory")
	require.Nil(t, ds)
}

func Test_SQLiteDirectory(t *testing.T) {
	ds, err := New("testdata")

	require.EqualError(t, err, `"testdata" isn't a file`)
	require.Nil(t, ds)
}

func Test_SQLite(t *testing.T) {
	ds, err := New("testdata/chinook.db")

	require.NoError(t, err)
	require.NotNil(t, ds)
	require.NoError(t, ds.Ping())

	tests := []struct {
		name string
		do   func(*testing.T)
	}{
		{
			name: "list schemas",
			do: func(t *testing.T) {
				schemas, err := ds.ListSchemas()

				assert.NoError(t, err)
				assert.ElementsMatch(t, schemas, []string{"main"})
			},
		},
		{
			name: "list tables",
			do: func(t *testing.T) {
				tables, err := ds.ListTables("")

				assert.NoError(t, err)
				assert.ElementsMatch(t, tables, []string{"albums", "sqlite_sequence", "artists", "customers", "employees", "genres", "invoices", "invoice_items", "media_types", "playlists", "playlist_track", "tracks", "sqlite_stat1"})
			},
		},
		{
			name: "preview table",
			do: func(t *testing.T) {
				albums, err := ds.PreviewTable("", "albums")

				assert.NoError(t, err)
				assert.Len(t, albums, 11)
				assert.Len(t, albums[0], 3)
			},
		},
		{
			name: "describe table",
			do: func(t *testing.T) {
				table, err := ds.DescribeTable("", "albums")

				assert.NoError(t, err)
				assert.Len(t, table, 2)
				assert.Contains(t, *table[1][0], `CREATE TABLE "albums"`)
			},
		},
		{
			name: "query table",
			do: func(t *testing.T) {
				query, err := ds.Query("", "SELECT Name from playlists WHERE PlaylistId < 6")

				assert.NoError(t, err)
				assert.Len(t, query, 6)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.do(t)
		})
	}
}
