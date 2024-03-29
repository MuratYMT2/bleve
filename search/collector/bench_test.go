//  Copyright (c) 2016 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/MuratYMT2/bleve/v2/search"
	index "github.com/blevesearch/bleve_index_api"
)

type createCollector func() search.Collector

func benchHelper(numOfMatches int, cc createCollector, b *testing.B) {
	matches := make([]*search.DocumentMatch, 0, numOfMatches)
	for i := 0; i < numOfMatches; i++ {
		matches = append(
			matches, &search.DocumentMatch{
				IndexInternalID: index.IndexInternalID(strconv.Itoa(i)),
				Score:           rand.Float64(),
			},
		)
	}

	b.ResetTimer()

	for run := 0; run < b.N; run++ {
		searcher := &stubSearcher{
			matches: matches,
		}
		collector := cc()
		err := collector.Collect(context.Background(), searcher, &stubReader{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
