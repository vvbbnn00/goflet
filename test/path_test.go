package test

import (
	"github.com/stretchr/testify/assert"
	"goflet/util"
	"testing"
)

func TestMatchPath1(t *testing.T) {
	patterns := []string{
		"/a/*",
		"/b/",
		"/c",
	}
	assert.True(t, util.MatchPath("/a/b", patterns), "Path should match")
	assert.True(t, util.MatchPath("/b/", patterns), "Path should match")
	assert.True(t, util.MatchPath("/c", patterns), "Path should match")
	assert.True(t, util.MatchPath("/a/b/c", patterns), "Path should match")
	assert.False(t, util.MatchPath("/b", patterns), "Path should not match")
	assert.False(t, util.MatchPath("/d", patterns), "Path should not match")
}

func TestMatchPath2(t *testing.T) {
	patterns := []string{
		"/*",
	}
	assert.True(t, util.MatchPath("/a/b", patterns), "Path should match")
	assert.True(t, util.MatchPath("/b/", patterns), "Path should match")
	assert.True(t, util.MatchPath("/c", patterns), "Path should match")
	assert.True(t, util.MatchPath("/a/b/c", patterns), "Path should match")
	assert.False(t, util.MatchPath("b", patterns), "Path should not match")
}

func TestMatchPath3(t *testing.T) {
	patterns := []string{
		"*",
	}
	assert.True(t, util.MatchPath("/a/b", patterns), "Path should match")
	assert.True(t, util.MatchPath("/b/", patterns), "Path should match")
	assert.True(t, util.MatchPath("/c", patterns), "Path should match")
	assert.True(t, util.MatchPath("/a/b/c", patterns), "Path should match")
	assert.True(t, util.MatchPath("b", patterns), "Path should match")
}

func TestMatchPath4(t *testing.T) {
	patterns := []string{
		"/",
	}
	assert.False(t, util.MatchPath("/a/b", patterns), "Path should not match")
	assert.False(t, util.MatchPath("/b/", patterns), "Path should not match")
	assert.False(t, util.MatchPath("/c", patterns), "Path should not match")
	assert.False(t, util.MatchPath("/a/b/c", patterns), "Path should not match")
	assert.False(t, util.MatchPath("/b", patterns), "Path should not match")
	assert.False(t, util.MatchPath("/d", patterns), "Path should not match")
}

func TestMatchPath5(t *testing.T) {
	patterns := []string{
		"/a**",
	}
	assert.True(t, util.MatchPath("/a", patterns), "Path should match")
	assert.True(t, util.MatchPath("/a/b", patterns), "Path should match")
	assert.False(t, util.MatchPath("/b", patterns), "Path should not match")
}
