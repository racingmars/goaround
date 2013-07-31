/*
 * File:	db_test.go
 *
 * Implements tests for the db.go functionality
 *
 *
 * Copyright (c) 2013, Matthew R. Wilson <mwilson@mattwilson.org>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met: 
 * 
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer. 
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution. 
 * 
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package goaround

import "testing"
import "time"

func TestCreation(t *testing.T) {
	res := 5
	capacity := 27
	db := New(res, capacity)

	if x := db.Res(); x != res {
		t.Errorf("db.Res() = %v, want %v", x, res)
	}

	if x := db.Capacity(); x != capacity {
		t.Errorf("db.Capacity() = %v, want %v", x, capacity)
	}

	if db.Len() != 0 {
		t.Errorf("New db is not empty")
	}
}

func TestSimplePopulation(t *testing.T) {
	res := 5
	capacity := 10
	db := New(res, capacity)

	db.Add(0)

	if x := db.Len(); x != 1 {
		t.Errorf("db.Len() = %v, want 1", x)
	}

	if x := db.Get(0); x != 0 {
		t.Errorf("db.Get(0) = %v, want 0", x)
	}
}

// This test exercises timebox averaging as well as rolling over the database
func TestComplexPopulation(t *testing.T) {
	var data = []struct {
		t string
		v float32
	}{
		{"2013-01-01T08:10:01Z", 5},
		{"2013-01-01T08:10:30Z", 5},
		{"2013-01-01T08:10:45Z", 5},
		{"2013-01-01T08:11:00Z", 5},
		{"2013-01-01T08:11:15Z", 10},
		{"2013-01-01T08:11:35Z", 15},
		{"2013-01-01T08:11:40Z", 8},
		{"2013-01-01T08:11:42Z", 305},
		{"2013-01-01T08:12:04Z", 10},
		{"2013-01-01T08:13:34Z", 20},
		{"2013-01-01T08:14:05Z", 30},
		{"2013-01-01T08:14:35Z", 30},
		{"2013-01-01T08:15:20Z", 20},
	}

	var expectedResults = []struct {
		i int
		v float32
	}{
		{0, 5},
		{1, 12.5},
		{2, 30.166666667},
		{3, 10},
		{4, 0},
		{5, 0},
		{6, 28.666666667},
		{7, 30},
		{8, 21.666666667},
		{9, 20},
	}

	res := 30
	capacity := 10
	db := New(res, capacity)

	for _, v := range data {
		tm, _ := time.Parse(time.RFC3339, v.t)
		db.AddAt(v.v, tm)
	}

	for _, v := range expectedResults {
		if result := db.Get(v.i); result != v.v {
			t.Errorf("db.Get(%d) returned %v, expected %v", v.i, result, v.v)
		}
	}
}
