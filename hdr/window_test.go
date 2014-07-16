package hdr

import "testing"

func TestWindowedHistogram(t *testing.T) {
	w := NewWindowedHistogram(2, 1, 1000, 3)

	for i := 0; i < 100; i++ {
		w.Current.RecordValue(int64(i))
	}
	w.Rotate()

	for i := 100; i < 200; i++ {
		w.Current.RecordValue(int64(i))
	}
	w.Rotate()

	for i := 200; i < 300; i++ {
		w.Current.RecordValue(int64(i))
	}

	if v := w.Merge().ValueAtQuantile(50); v != 199 {
		t.Errorf("Median was %v, but expected 199", v)
	}
}

func BenchmarkWindowedHistogramRecordAndRotate(b *testing.B) {
	w := NewWindowedHistogram(3, 1, 10000000, 3)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := w.Current.RecordValue(100); err != nil {
			b.Fatal(err)
		}

		if i%100000 == 1 {
			w.Rotate()
		}
	}
}

func BenchmarkWindowedHistogramMerge(b *testing.B) {
	w := NewWindowedHistogram(3, 1, 10000000, 3)
	for i := 0; i < 10000000; i++ {
		if err := w.Current.RecordValue(100); err != nil {
			b.Fatal(err)
		}

		if i%100000 == 1 {
			w.Rotate()
		}
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Merge()
	}
}