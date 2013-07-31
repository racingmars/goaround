/* File:	timebox_test.go
 *
 * Implements tests for timebox.go functionality
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
	"testing"
	"time"
)

var timeBoxTests = []struct {
	in     string
	chunk  int
	start  string
	stop   string
	format string
}{
	{"2013-01-02T10:04:10Z", 30, "2013-01-02T10:04:00Z", "2013-01-02T10:04:30Z", time.RFC3339},
	{"2013-01-02T10:04:29Z", 30, "2013-01-02T10:04:00Z", "2013-01-02T10:04:30Z", time.RFC3339},
	{"2013-01-02T10:04:30Z", 30, "2013-01-02T10:04:30Z", "2013-01-02T10:05:00Z", time.RFC3339},
	{"2013-01-02T10:04:10Z", 60, "2013-01-02T10:04:00Z", "2013-01-02T10:05:00Z", time.RFC3339},
	{"2013-01-02T10:04:10Z", 3600 /* one hour */, "2013-01-02T10:00:00Z", "2013-01-02T11:00:00Z", time.RFC3339},
	{"2013-01-02T10:04:10Z", 86400 /* one day */, "2013-01-02T00:00:00Z", "2013-01-03T00:00:00Z", time.RFC3339},
}

func TestTimeboxing(t *testing.T) {
	for i, tt := range timeBoxTests {
		in, _ := time.Parse(tt.format, tt.in)
		goodStart, _ := time.Parse(tt.format, tt.start)
		goodStop, _ := time.Parse(tt.format, tt.stop)
		start, stop := BoxTime(in, tt.chunk)
		if !start.Equal(goodStart) {
			t.Errorf("Test %d: got start %v, expected %v", i, start, goodStart)
		}
		if !stop.Equal(goodStop) {
			t.Errorf("Test %d: got stop %v, expected %v", i, stop, goodStop)
		}
	}
}
