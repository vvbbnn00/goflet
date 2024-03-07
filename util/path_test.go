package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchPath1(t *testing.T) {
	patterns := []string{
		"/a/*",
		"/b/",
		"/c",
	}
	assert.True(t, MatchPath("/a/b", patterns), "Path should match")
	assert.True(t, MatchPath("/b/", patterns), "Path should match")
	assert.True(t, MatchPath("/c", patterns), "Path should match")
	assert.True(t, MatchPath("/a/b/c", patterns), "Path should match")
	assert.False(t, MatchPath("/b", patterns), "Path should not match")
	assert.False(t, MatchPath("/d", patterns), "Path should not match")
}

func TestMatchPath2(t *testing.T) {
	patterns := []string{
		"/*",
	}
	assert.True(t, MatchPath("/a/b", patterns), "Path should match")
	assert.True(t, MatchPath("/b/", patterns), "Path should match")
	assert.True(t, MatchPath("/c", patterns), "Path should match")
	assert.True(t, MatchPath("/a/b/c", patterns), "Path should match")
	assert.False(t, MatchPath("b", patterns), "Path should not match")
}

func TestMatchPath3(t *testing.T) {
	patterns := []string{
		"*",
	}
	assert.True(t, MatchPath("/a/b", patterns), "Path should match")
	assert.True(t, MatchPath("/b/", patterns), "Path should match")
	assert.True(t, MatchPath("/c", patterns), "Path should match")
	assert.True(t, MatchPath("/a/b/c", patterns), "Path should match")
	assert.True(t, MatchPath("b", patterns), "Path should match")
}

func TestMatchPath4(t *testing.T) {
	patterns := []string{
		"/",
	}
	assert.False(t, MatchPath("/a/b", patterns), "Path should not match")
	assert.False(t, MatchPath("/b/", patterns), "Path should not match")
	assert.False(t, MatchPath("/c", patterns), "Path should not match")
	assert.False(t, MatchPath("/a/b/c", patterns), "Path should not match")
	assert.False(t, MatchPath("/b", patterns), "Path should not match")
	assert.False(t, MatchPath("/d", patterns), "Path should not match")
}

func TestMatchPath5(t *testing.T) {
	patterns := []string{
		"/a**",
	}
	assert.True(t, MatchPath("/a", patterns), "Path should match")
	assert.True(t, MatchPath("/a/b", patterns), "Path should match")
	assert.False(t, MatchPath("/b", patterns), "Path should not match")
}
