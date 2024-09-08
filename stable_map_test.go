package stable_map

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestStableMap_Reset(t *testing.T) {
	m := NewStableMap[int, string]()
	assert.False(t, m.Has(1))

	m.Set(1, "one")
	assert.True(t, m.Has(1))
	m.Set(3, "three")

	assert.Equal(t, "one three", strings.Join(SequenceToSlice(m.Values()), " "))

	m.Reset()
	assert.False(t, m.Has(1))
	assert.False(t, m.Has(3))
	m.Set(3, "three")
	m.Set(1, "one")
	assert.Equal(t, "three one", strings.Join(SequenceToSlice(m.Values()), " "))
}

func TestStableMap_Delete(t *testing.T) {
	m := NewStableMap[int, string]()

	assert.Equal(t, 0, m.Len())

	m.Set(1, "one")
	assert.Equal(t, 1, m.Len())
	m.Set(3, "three")
	assert.Equal(t, 2, m.Len())

	assert.Equal(t, "one three", strings.Join(SequenceToSlice(m.Values()), " "))

	v, ok := m.Get(1)
	assert.True(t, ok)
	assert.Equal(t, "one", v)

	m.Delete(1)
	assert.False(t, m.Has(1))
	assert.Equal(t, 1, m.Len())

	v, ok = m.Get(1)
	assert.False(t, ok)
	assert.Equal(t, "", v)

	assert.Equal(t, "three", strings.Join(SequenceToSlice(m.Values()), " "))

	m.Set(1, "ONE")
	assert.Equal(t, 2, m.Len())
	assert.Equal(t, "three ONE", strings.Join(SequenceToSlice(m.Values()), " "))
}

func TestShouldInsertQuickly(t *testing.T) {

	m := NewStableMap[int, string]()

	s := time.Now()
	for i := 0; i < 1_000_000; i++ {
		m.Set(i, "")
	}
	assert.Less(t, time.Since(s), time.Millisecond*300)

	s = time.Now()
	for i := 0; i < 1_000_000; i++ {
		m.Delete(i)
	}
	assert.Less(t, time.Since(s), time.Millisecond*250)
}
