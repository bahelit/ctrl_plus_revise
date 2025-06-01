// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin || windows
// +build darwin windows

package clipboard_test

import (
	"testing"

	"github.com/go-vgo/robotgo/clipboard"
)

func TestCopyAndPaste(t *testing.T) {
	expected := "Hello Clipboard!"

	err := clipboard.WriteAll(expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := clipboard.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("want %s, got %s", expected, actual)
	}
}

func TestMultiCopyAndPaste(t *testing.T) {
	expected1 := "French: éèêëàùœç - Japanese: 日本語"
	expected2 := "Weird UTF-8: 💩☃"

	err := clipboard.WriteAll(expected1)
	if err != nil {
		t.Fatal(err)
	}

	actual1, err := clipboard.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if actual1 != expected1 {
		t.Errorf("want %s, got %s", expected1, actual1)
	}

	err = clipboard.WriteAll(expected2)
	if err != nil {
		t.Fatal(err)
	}

	actual2, err := clipboard.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if actual2 != expected2 {
		t.Errorf("want %s, got %s", expected2, actual2)
	}
}

func BenchmarkReadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		clipboard.ReadAll()
	}
}

func BenchmarkWriteAll(b *testing.B) {
	text := "Hello Clipboard!"
	for i := 0; i < b.N; i++ {
		clipboard.WriteAll(text)
	}
}
