package app

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"exchangeapp/internal/repository"
)

type fakeCounter struct {
	popDirtyIDs []uint
	popDirtyErr error

	getVals map[uint]int64
	getErrs map[uint]error

	popCalls  int
	getCalls  []uint
	markCalls []uint
}

func (f *fakeCounter) PopDirty(limit int) ([]uint, error) {
	f.popCalls++
	return f.popDirtyIDs, f.popDirtyErr
}

func (f *fakeCounter) GetLikeCount(threadID uint) (int64, error) {
	f.getCalls = append(f.getCalls, threadID)
	if f.getErrs != nil {
		if err, ok := f.getErrs[threadID]; ok {
			return 0, err
		}
	}
	if f.getVals != nil {
		if val, ok := f.getVals[threadID]; ok {
			return val, nil
		}
	}
	return 0, errors.New("missing value")
}

func (f *fakeCounter) MarkDirty(threadID uint) error {
	f.markCalls = append(f.markCalls, threadID)
	return nil
}

type setCall struct {
	id  uint
	val int64
}

type fakeWriter struct {
	calls   []setCall
	failIDs map[uint]error
}

func (f *fakeWriter) SetLikeCount(threadID uint, value int64) error {
	f.calls = append(f.calls, setCall{id: threadID, val: value})
	if f.failIDs != nil {
		if err, ok := f.failIDs[threadID]; ok {
			return err
		}
	}
	return nil
}

func TestLikeCountFlusherFlushOnceBatchZeroUsesDefault(t *testing.T) {
	counter := &fakeCounter{}
	flusher := NewLikeCountFlusher(counter, &fakeWriter{}, 0, time.Second)
	flusher.flushOnce()

	if counter.popCalls != 1 {
		t.Fatalf("expected PopDirty called once, got %d", counter.popCalls)
	}
}

func TestLikeCountFlusherFlushOncePopDirtyError(t *testing.T) {
	counter := &fakeCounter{popDirtyErr: errors.New("boom")}
	flusher := NewLikeCountFlusher(counter, &fakeWriter{}, 10, time.Second)
	flusher.flushOnce()

	if counter.popCalls != 1 {
		t.Fatalf("expected 1 PopDirty call, got %d", counter.popCalls)
	}
	if len(counter.getCalls) != 0 {
		t.Fatalf("expected no GetLikeCount calls, got %d", len(counter.getCalls))
	}
	if len(counter.markCalls) != 0 {
		t.Fatalf("expected no MarkDirty calls, got %d", len(counter.markCalls))
	}
}

func TestLikeCountFlusherFlushOnceMixedResults(t *testing.T) {
	counter := &fakeCounter{
		popDirtyIDs: []uint{1, 2, 3, 4},
		getVals: map[uint]int64{
			1: 10,
			4: 7,
		},
		getErrs: map[uint]error{
			2: repository.ErrLikeCountNotFound,
			3: errors.New("read error"),
		},
	}
	writer := &fakeWriter{failIDs: map[uint]error{4: errors.New("write error")}}
	flusher := NewLikeCountFlusher(counter, writer, 10, time.Second)
	flusher.flushOnce()

	wantCalls := []setCall{{id: 1, val: 10}, {id: 4, val: 7}}
	if !reflect.DeepEqual(writer.calls, wantCalls) {
		t.Fatalf("unexpected SetLikeCount calls: %+v", writer.calls)
	}

	wantDirty := []uint{3, 4}
	if !reflect.DeepEqual(counter.markCalls, wantDirty) {
		t.Fatalf("unexpected MarkDirty calls: %+v", counter.markCalls)
	}
}
