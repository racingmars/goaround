/*
 * File:	db.go
 *
 * Implements a round-robin time-series database.
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

import "fmt"
import "time"

type Db struct {
	res          int       // resolution - how many seconds elapse between successive entries
	entries      []float32 // the individual database entries
	head         int       // index of the beginning of the list. -1 means no data.
	tail         int       // index of the end of the list. -1 means no data.
	currentStart time.Time // beginning time of current bucket
	currentStop  time.Time // end time of current bucket
	lastEntry    time.Time // last update time
}

// New creates and returns a new Db with the specified resolution (in seconds)
// and capacity.
func New(resolution int, capacity int) *Db {
	db := new(Db)
	db.res = resolution
	db.entries = make([]float32, capacity)
	db.head = -1
	db.tail = -1
	return db
}

// Res returns the resolution of the database.
func (db *Db) Res() int {
	return db.res
}

// Capacity returns the capacity of the database.
func (db *Db) Capacity() int {
	return len(db.entries)
}

// Add will add value v to the database at the current time.
func (db *Db) Add(v float32) {
	db.AddAt(v, time.Now())
}

// AddAt will add a value, v, to the database at the specific time, t. Data will
// be consolidated (averaged) correctly to apply data with any timestamp into
// the defined timeboxes of the database.
func (db *Db) AddAt(v float32, t time.Time) {
	// Normalize everything to UTC
	t = t.UTC()

	// t0 and t1 is the time box that this data point should live in
	t0, t1 := BoxTime(t, db.res)
	t0, t1 = t0.UTC(), t1.UTC()

	// Check if database is empty. If so, create the first entry.
	if db.tail == -1 {
		db.tail = 0
		db.head = 0
		db.entries[db.tail] = v
		db.currentStart = t0
		db.currentStop = t1
		db.lastEntry = t
		return
	}

	// Are we trying to rewrite history?
	if t.Before(db.lastEntry) {
		// TODO: real error handling
		fmt.Println("Can't rewrite history.")
		return
	}

	// Are we still in tail's timebox?
	if t.Before(db.currentStop) {
		prevFill := float32(db.lastEntry.Sub(db.currentStart).Seconds())
		curDuration := float32(t.Sub(db.lastEntry).Seconds())
		oldval := db.entries[db.tail]
		newval := (oldval*prevFill + v*curDuration) / (prevFill + curDuration)
		db.entries[db.tail] = newval
		db.lastEntry = t
		return
	}

	// Have we moved exactly one timebox forward?
	if temp := db.currentStop.Add(time.Duration(db.res) * time.Second); t.Before(temp) {
		// First we need to apply whatever was left in the previous timebox
		prevFill := float32(db.lastEntry.Sub(db.currentStart).Seconds())
		curDuration := float32(db.currentStop.Sub(db.lastEntry).Seconds())
		oldval := db.entries[db.tail]
		newval := (oldval*prevFill + v*curDuration) / (prevFill + curDuration)
		db.entries[db.tail] = newval

		// Move the tail (which also updates the start and stop times)
		db.moveForward()

		// Apply reading to current (new) timebox/tail
		db.entries[db.tail] = v
		db.lastEntry = t

		return
	} else {
		// We've gone more than one timebox forward
		// Catch up to where we should be, filling in zeros in the missing slots
		for db.currentStop.Before(t) {
			db.moveForward()
			db.entries[db.tail] = 0
		}

		// Apply reading to current timebox/tail
		db.entries[db.tail] = v
		db.lastEntry = t
	}
}

// moveForward will increment the tail (and head if necessary) by one position
// and update the currentStart and currentStop time for the new timebox
func (db *Db) moveForward() {
	db.tail++
	if db.tail >= len(db.entries) {
		db.tail = 0
	}

	if db.tail == db.head {
		db.head++
		if db.head >= len(db.entries) {
			db.head = 0
		}
	}

	// Update times. Running through BoxTime() instead of directly adding
	// seconds so that floating point errors don't accumulate over time
	newTime := db.currentStop.Add(time.Duration(1) * time.Second)
	db.currentStart, db.currentStop = BoxTime(newTime, db.res)
}

// Len returns the length of actual data in the database [e.g. a database with
// a large capacity but not yet filled could have Len() < Capacity(), but Len()
// will never be greater than Capacity()].
func (db *Db) Len() int {
	if db.tail == -1 {
		return 0
	}

	if db.head <= db.tail {
		return db.tail - db.head + 1
	}

	if db.head > db.tail {
		return len(db.entries) - db.head + db.tail + 1
	}

	panic("It shouldn't be possible to get here.")
}

// Get returns the value at the indicated index. Index must not be outside the
// bounds of the current populated data [i.e. index must be less than Len(),
// even if Capacity() > Len()].
func (db *Db) Get(i int) float32 {
	if i >= db.Len() {
		panic("Index out of bounds.")
	}

	if db.head <= db.tail {
		return db.entries[db.head+i]
	}

	if db.head > db.tail {
		j := db.head + i
		if j < len(db.entries) {
			return db.entries[j]
		} else {
			return db.entries[j-len(db.entries)]
		}
	}

	panic("It shouldn't be possible to get here.")
}

func (db *Db) printDebug() {
	fmt.Println("---- DB Dump ------------------------------")
	fmt.Printf("res: %v, head: %v, tail: %v ", db.res, db.head, db.tail)
	fmt.Printf("cap: %v, len: %v\n", len(db.entries), db.Len())
	fmt.Printf("start: %v, stop: %v\n", db.currentStart.UTC(), db.currentStop.UTC())
	fmt.Printf("last: %v\n", db.lastEntry.UTC())
	fmt.Printf("data: %v\n", db.entries)
}
