// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"
	"github.com/kr/beanstalk"
)

// a key-value store backed by raft
type kvstore struct {
	proposeC chan<- string // channel for proposing updates
	mu       sync.RWMutex
	kvStore  map[string]string // current committed key-value pairs
}


// a key-value store backed by raft
type bstalk struct {
	proposeC chan<- string // channel for proposing updates
	mu       sync.RWMutex
	conn	*beanstalk.Conn
}

type kv struct {
	Action string
	Val string
}

func newKVStore(proposeC chan<- string, commitC <-chan *string, errorC <-chan error, id *int) *bstalk {
	var c *beanstalk.Conn
	//var err error

//	if (*id == 1) {
//		log.Printf("connecting to remote beanstalkd")
//		c, _ = beanstalk.Dial("tcp", "172.21.140.160:11300")
//	} else {
//		log.Printf("connecting to local")
		c, _ = beanstalk.Dial("tcp", "127.0.0.1:11300")
//	}
	s := &bstalk{proposeC: proposeC, conn: c}
	// replay log into key-value map
	s.readCommits(commitC, errorC)
	// read commits from raft into kvStore map until error
	go s.readCommits(commitC, errorC)
	return s
}

func (s *bstalk) Lookup(key string) (string, bool) {
	s.mu.RLock()
	log.Printf("Beanstalk GET command")

	id, body, _ := s.conn.Reserve(0)
	s.conn.Delete(id)
	log.Printf("Beanstalk GET command completed for %s", string(body))
	s.mu.RUnlock()
	ok := true
	return string(body[:]), ok
}

func (s *bstalk) Propose(k string, v string) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv{k, v}); err != nil {
		log.Fatal(err)
	}
	s.proposeC <- string(buf.Bytes())
}

func (s *bstalk) readCommits(commitC <-chan *string, errorC <-chan error) {
	for data := range commitC {
		if data == nil {
			// done replaying log; new data incoming
			return
		}

		var data_kv kv
		dec := gob.NewDecoder(bytes.NewBufferString(*data))
		if err := dec.Decode(&data_kv); err != nil {
			log.Fatalf("raftexample: could not decode message (%v)", err)
		}
		s.mu.Lock()
		log.Printf("Beanstalk put command with data %s", data_kv.Val)
		s.conn.Put([]byte(data_kv.Val), 1, 0, 0)
		s.mu.Unlock()
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}
