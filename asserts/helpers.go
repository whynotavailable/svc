package asserts

import (
	"fmt"
	"testing"
)

func Assert(t *testing.T, cond bool, args ...any) {
	if !cond {
		t.Error(args...)
		t.FailNow()
	}
}

func throw(t *testing.T, msg string, args ...any) {
	if len(args) > 0 {
		t.Error(msg, args)
	} else {
		t.Error(msg)
	}
	t.FailNow()
}

func Eq[T comparable](t *testing.T, target T, actual T, extras ...any) {
	if target != actual {
		msg := fmt.Sprintf("target %#v does not match %#v", target, actual)
		throw(t, msg, extras...)
	}
}

func NoKey[TKey comparable, TVal any](t *testing.T, m map[TKey]TVal, key TKey, extras ...any) {
	if val, ok := m[key]; ok {
		msg := fmt.Sprintf("key %#v found when should be missing has %#v", key, val)
		throw(t, msg, extras...)
	}
}

func True(t *testing.T, b bool, extras ...any) {
	if !b {
		msg := "got false when should be true"
		throw(t, msg, extras...)
	}
}

func False(t *testing.T, b bool, extras ...any) {
	if b {
		msg := "got true when should be false"
		throw(t, msg, extras...)
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func NotNil(t *testing.T, obj any, extras ...any) {
	if obj == nil {
		msg := "got nil"
		throw(t, msg, extras...)
	}
}
