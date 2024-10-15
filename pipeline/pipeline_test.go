// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package pipeline_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/certmichelin/azurehound/v2/pipeline"
)

func TestBatch(t *testing.T) {

	done := make(chan interface{})
	in := make(chan string)

	go func() {
		in <- "foo"
		in <- "bar"

		in <- "bazz"
		time.Sleep(5 * time.Millisecond)

		in <- "buzz"

		close(in)
	}()

	batches := map[int]int{}
	i := 0
	for batch := range pipeline.Batch(done, in, 2, 5*time.Millisecond) {
		batches[i] = len(batch)
		i++
		fmt.Println(batch)
	}

	if len(batches) != 3 {
		t.Errorf("got %v, want %v", len(batches), 3)
	}

	if length, ok := batches[0]; !ok || length != 2 {
		t.Errorf("got %v, want %v", length, 2)
	}

	if length, ok := batches[1]; !ok || length != 1 {
		t.Errorf("got %v, want %v", length, 1)
	}

	if length, ok := batches[2]; !ok || length != 1 {
		t.Errorf("got %v, want %v", length, 1)
	}
}

func TestDemux(t *testing.T) {

	var (
		done  = make(chan interface{})
		in    = make(chan string)
		wg    sync.WaitGroup
		count int
	)

	go func() {
		defer close(in)
		in <- "foo"
		in <- "bar"
		in <- "bazz"
		in <- "buzz"
	}()

	outs := pipeline.Demux(done, in, 2)
	wg.Add(len(outs))
	for i := range outs {
		out := outs[i]
		go func() {
			defer wg.Done()
			for s := range out {
				fmt.Println(s)
				count++
			}
		}()
	}

	wg.Wait()
	if count != 4 {
		t.Fail()
	}

}
