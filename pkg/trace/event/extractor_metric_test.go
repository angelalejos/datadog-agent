package event

import (
	"math/rand"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/trace/pb"
	"github.com/DataDog/datadog-agent/pkg/trace/sampler"
	"github.com/DataDog/datadog-agent/pkg/trace/stats"
)

func createTestSpansWithEventRate(eventRate float64) []*stats.WeightedSpan {
	spans := make([]*stats.WeightedSpan, 1000)
	for i := range spans {
		spans[i] = &stats.WeightedSpan{Span: &pb.Span{TraceID: rand.Uint64(), Service: "test", Name: "test", Metrics: map[string]float64{}}}
		if eventRate >= 0 {
			spans[i].Metrics[sampler.KeySamplingRateEventExtraction] = eventRate
		}
	}
	return spans
}

func TestMetricBasedExtractor(t *testing.T) {
	tests := []extractorTestCase{
		// Name: <priority>/<extraction rate>
		{"none/missing", createTestSpansWithEventRate(-1), 0, -1},
		{"none/0", createTestSpansWithEventRate(0), 0, 0},
		{"none/0.5", createTestSpansWithEventRate(0.5), 0, 0.5},
		{"none/1", createTestSpansWithEventRate(1), 0, 1},
		{"1/missing", createTestSpansWithEventRate(-1), 1, -1},
		{"1/0", createTestSpansWithEventRate(0), 1, 0},
		{"1/0.5", createTestSpansWithEventRate(0.5), 1, 0.5},
		{"1/1", createTestSpansWithEventRate(1), 1, 1},
		// Priority 2 should have extraction rate of 1 so long as any extraction rate is set and > 0
		{"2/missing", createTestSpansWithEventRate(-1), 2, -1},
		{"2/0", createTestSpansWithEventRate(0), 2, 0},
		{"2/0.5", createTestSpansWithEventRate(0.5), 2, 1},
		{"2/1", createTestSpansWithEventRate(1), 2, 1},
	}

	for _, test := range tests {
		testExtractor(t, NewMetricBasedExtractor(), test)
	}
}