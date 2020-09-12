package pg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnectDB(t *testing.T) {
	store, err := NewStore("postgresql://postgres:@db:5432/postgres?sslmode=disable")
	require.NoError(t, err)

	saved, err := store.Save("http://hello.world")
	require.NoError(t, err)
	require.Equal(t, "http://hello.world", saved.Original)

	loaded, err := store.GetByKey(saved.Key)
	require.NoError(t, err)

	require.Equal(t, saved.Key, loaded.Key)
	require.Equal(t, saved.Original, loaded.Original)
}
