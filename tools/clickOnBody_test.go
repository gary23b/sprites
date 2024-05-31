package tools

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClickOnBody(t *testing.T) {
	b := NewTouchCollisionBody()
	require.NotNil(t, b)

	b.AddCircleBody(-10, 0, 5)
	b.AddRectangleBody(10, 20, 0, 10)

	// Add an invalid rec
	b.AddRectangleBody(10, -10, 10, -10)

	require.False(t, b.IsMouseClickInBody(100, 0))
	require.False(t, b.IsMouseClickInBody(0, 0))
	require.False(t, b.IsMouseClickInBody(-4, 0))
	require.True(t, b.IsMouseClickInBody(-6, 0))

	require.False(t, b.IsMouseClickInBody(9.9, 1))
	require.True(t, b.IsMouseClickInBody(10.1, 1))

	// Now move and rotate
	b.Pos(10, 10)
	b.Angle(-math.Pi / 2)

	require.False(t, b.IsMouseClickInBody(100, 0))
	require.False(t, b.IsMouseClickInBody(0, 0))
	require.False(t, b.IsMouseClickInBody(10+0, 10+4.9))
	require.True(t, b.IsMouseClickInBody(10+0, 10+5.1))

	require.False(t, b.IsMouseClickInBody(10+1, 10-9.9))
	require.True(t, b.IsMouseClickInBody(10+1, 10-10.1))

	x, y := b.GetMousePosRelativeToOriginalSprite(20, 10)
	require.InDelta(t, 0, x, 1e-6)
	require.InDelta(t, 10, y, 1e-6)

	// test clone
	b2 := b.Clone()
	require.False(t, b2.IsMouseClickInBody(10+0, 10+4.9))
	require.True(t, b2.IsMouseClickInBody(10+0, 10+5.1))

	require.False(t, b2.IsMouseClickInBody(10+1, 10-9.9))
	require.True(t, b2.IsMouseClickInBody(10+1, 10-10.1))
}
