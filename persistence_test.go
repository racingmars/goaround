/*
 * File:	persistence_test.go
 *
 * Implements tests for the persistence.go functionality.
 *
 *
 * Copyright (c) 2013, Matthew R. Wilson <mwilson@mattwilson.org>
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

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"
)

// TestEmptyRoundtrip tests with an empty new database.
func TestEmptyRoundtrip(t *testing.T) {
	db := New(8, 17)
	doRoundtrip(db, t)
}

// TestDataRoundtrip tests with a new database with all elements filled with
// (fake) data.
func TestDataRoundtrip(t *testing.T) {
	db := New(30, 5)
	baseTime := time.Now()
	db.head = 1
	db.tail = 3
	db.currentStart = baseTime
	db.currentStop = baseTime.Add(5 * time.Minute)
	db.lastEntry = baseTime.Add(10 * time.Minute)
	for i, _ := range db.entries {
		db.entries[i] = float32(i * 7.0)
	}
	doRoundtrip(db, t)
}

// doRoundtrip will encode db to a gob, then decode it and make sure the data
// is the same, reporting errors to t.
func doRoundtrip(db *Db, t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	err := enc.Encode(db)
	if err != nil {
		t.Errorf("Error encoding: %v", err)
	}

	var newdb *Db = new(Db)
	err = dec.Decode(newdb)
	if err != nil {
		t.Errorf("Error decoding: %v", err)
	}

	if !newdb.equals(db) {
		t.Errorf("Encoded and decoded db do not match")
	}
}

// equals will tell you if the Dbs a and b hold equal data.
func (a *Db) equals(b *Db) bool {
	simpleValues := a.res == b.res &&
		a.head == b.head &&
		a.tail == b.tail &&
		a.currentStart == b.currentStart &&
		a.currentStop == b.currentStop &&
		a.lastEntry == b.lastEntry

	var entriesEqual bool = true

	for i, v := range a.entries {
		if b.entries[i] != v {
			entriesEqual = false
			break
		}
	}

	return simpleValues && entriesEqual
}
