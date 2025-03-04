// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"io/ioutil"
	"testing"
)

func BenchmarkNewAudio(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		NewAudio( d, "")
	}
}

func BenchmarkAudioWriteTo(b *testing.B) {
	b.StopTimer()
	d := RandomDigits(DefaultLen)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		a := NewAudio( d, "")
		n, _ := a.WriteTo(ioutil.Discard)
		b.SetBytes(n)
	}
}
