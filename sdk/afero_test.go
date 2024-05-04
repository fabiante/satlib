package sdk

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAferoFullBaseFsPath(t *testing.T) {
	inner := afero.NewBasePathFs(afero.NewMemMapFs(), "/tmp")
	outer := afero.NewBasePathFs(inner, "/something")

	assert.Equal(t, "/tmp/something", afero.FullBaseFsPath(outer.(*afero.BasePathFs), ""))
	assert.Equal(t, "/tmp/something", afero.FullBaseFsPath(outer.(*afero.BasePathFs), "/"))
	assert.Equal(t, "/tmp/something/brt.pdf", afero.FullBaseFsPath(outer.(*afero.BasePathFs), "brt.pdf"))
}

func TestFullFsPath(t *testing.T) {
	t.Run("works on single base path fs", func(t *testing.T) {
		fs := NewNestedFs(afero.NewMemMapFs(), "/tmp")
		x := FullFsPath(fs, "abc.pdf")
		require.Equal(t, "/tmp/abc.pdf", x)
	})

	t.Run("works on nested base path fs", func(t *testing.T) {
		inner := NewNestedFs(afero.NewMemMapFs(), "/tmp")
		outer := NewNestedFs(inner, "/something")
		x := FullFsPath(outer, "abc.pdf")
		require.Equal(t, "/tmp/something/abc.pdf", x)
	})
}
