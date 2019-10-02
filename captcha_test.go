// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"bytes"
	"github.com/satori/go.uuid"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c == "" {
		t.Errorf("expected id, got empty string")
	}
}

func TestVerify(t *testing.T) {
	id := New()
	if Verify(id, []byte{0, 0}) {
		t.Errorf("verified wrong captcha")
	}
	id = New()
	d := globalStore.Get(id, false) // cheating
	if !Verify(id, d) {
		t.Errorf("proper captcha not verified")
	}
}

func TestReload(t *testing.T) {
	id := New()
	d1 := globalStore.Get(id, false) // cheating
	Reload(id, false)
	d2 := globalStore.Get(id, false) // cheating again
	if bytes.Equal(d1, d2) {
		t.Errorf("reload didn't work: %v = %v", d1, d2)
	}
}

func TestForceReloadWithNoDiditalsCached(t *testing.T) {
	id := uuid.NewV4().String()
	d1 := globalStore.Get(id, false) // cheating
	if len(d1) != 0 {
		t.Errorf("there should be no digital cached %v", d1)
	}
	Reload(id, true)
	d2 := globalStore.Get(id, false) // cheating again
	if len(d2) != DefaultLen {
		t.Errorf("didn't generate proper digitals %v", d2)
	}
	if bytes.Equal(d1, d2) {
		t.Errorf("reload didn't work: %v = %v", d1, d2)
	}
}

func TestRandomDigits(t *testing.T) {
	d1 := RandomDigits(20)
	for _, v := range d1 {
		if v >= 12 {
			t.Errorf("digits not in range 0-9: %v", d1)
		}
	}
	d2 := RandomDigits(20)
	if bytes.Equal(d1, d2) {
		t.Errorf("digits seem to be not random")
	}
}
