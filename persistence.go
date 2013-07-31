/*
 * File:	persistence.go
 *
 * Implements the ability to persist (save to disk and load from disk) the
 * database implemented in db.go.
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
	"errors"
	"time"
)


/*****************************************************************************/
// What follows is support to store Db structures as gobs. This is necessary
// because all of the fields in the Db struct are not exported. So we'll copy
// the data to a struct type that mirrors the exported Db struct, but is not
// itself exported and contains all exported fields, then encode that as a gob.
//
// This seems pretty hacky, but it's the best idea I have at the moment that
// doesn't end up involve implementing even more code and not just relying on
// the library functions to encode and decode stuff.
/*****************************************************************************/

type gobDb struct {
	Res          int
	Entries      []float32
	Head         int
	Tail         int
	CurrentStart time.Time
	CurrentStop  time.Time
	LastEntry    time.Time
}

const gobDbGobVersion byte = 1

// GobEncode implements the gob.GobEncoder interface.
func (db *Db) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	d := gobDb{db.res, db.entries, db.head, db.tail, db.currentStart,
		db.currentStop, db.lastEntry}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(gobDbGobVersion)
	if err != nil {
		return nil, err
	}

	err = enc.Encode(d)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GobDecode implements the gob.GobDecoder interface.
func (db *Db) GobDecode(b []byte) error {
	if len(b) == 0 {
		return errors.New("rrdb.GobDecode: no data")
	}

	var err error
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	var version byte
	err = dec.Decode(&version)
	if err != nil {
		return err
	}
	if version != gobDbGobVersion {
		return errors.New("rrdb.GobDecode: unknown version")
	}

	var d gobDb
	err = dec.Decode(&d)
	if err != nil {
		return err
	}

	db.res = d.Res
	db.entries = d.Entries
	db.head = d.Head
	db.tail = d.Tail
	db.currentStart = d.CurrentStart
	db.currentStop = d.CurrentStop
	db.lastEntry = d.LastEntry

	return nil
}
