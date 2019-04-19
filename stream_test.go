package gostream

import (
	"context"
	"testing"
)

func addInts(a ...interface{}) interface{} {
	s := 0
	for i := 0; i < len(a); i++ {
		s += a[i].(int)
	}
	return s
}

func mulIntsBy5(a ...interface{}) interface{} {
	s := 1
	for i := 0; i < len(a); i++ {
		s *= a[i].(int)
	}
	return s * 5
}

func BenchmarkAppendStringsStream(b *testing.B) {
	ctx := context.Background()
	defer ctx.Done()

	// repeat: 1, 2, 3, 1, 2, 3
	repeat := GenerateRepeatStream(ctx, 1, 2, 3)

	// samenum: 1000, 1000, 1000, ...
	samenum := GenerateRepeatStream(ctx, 1000)

	// random: 3, 3843029809, 11, ... (for example)
	random := GenerateRandIntsStream(ctx)

	// calc: 1*5 + 1, 2*5 + 2, 3*5 + 1,...
	calc := Map(ctx, addInts,
		Map(ctx, mulIntsBy5, GenerateSerialIntsStream(ctx)), // 1*5, 2*5, 3*5, ...
		GenerateRepeatStream(ctx, 1, 2))                     // 1, 2, 1, 2, 1, 2, ...

	merged := Merge(ctx, repeat, samenum, random, calc)

	for range Take(ctx, merged, b.N) {
	}
}

func TestGenerateRepeatStream(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	s := GenerateRepeatStream(ctx, 0, 1, 2)
	for i := 0; i < 100; i++ {
		x := <-s
		if x != i%3 {
			t.Errorf("expected:%d actual:%d\n", i%3, x)
			return
		}
	}
}

func TestGenerateRandIntsStream(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	s := GenerateRandIntsStream(ctx)
	for i := 0; i < 100; i++ {
		x := <-s
		if _, ok := x.(int); !ok {
			t.Errorf("bad type")
			return
		}
	}
}

func TestGenerateSerialIntsStream(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	s := GenerateSerialIntsStream(ctx)
	for i := 0; i < 100; i++ {
		x := <-s
		if x != i {
			t.Errorf("expected:%d actual:%d\n", i, x)
			return
		}
	}
}

func TestMap(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	s1 := GenerateSerialIntsStream(ctx)
	s2 := GenerateSerialIntsStream(ctx)
	s := Map(ctx, addInts, s1, s2)
	for i := 0; i < 100; i++ {
		x := <-s
		if x != 2*i {
			t.Errorf("expected:%d actual:%d\n", 2*i, x)
			return
		}
	}
}

func TestTake(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	s := Take(ctx, GenerateSerialIntsStream(ctx), 10)
	a := make([]interface{}, 0)
	for x := range s {
		a = append(a, x)
	}
	if len(a) != 10 {
		t.Errorf("expected length:%d actual:%d\n", 10, len(a))
	}
}

func TestMerge(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()

	s1 := Take(ctx, GenerateSerialIntsStream(ctx), 10)
	s2 := Take(ctx, GenerateSerialIntsStream(ctx), 10)
	s := Merge(ctx, s1, s2)
	a := make([]interface{}, 0)
	for x := range s {
		a = append(a, x)
	}
	if len(a) != 20 {
		t.Errorf("expected length:%d actual:%d\n", 20, len(a))
	}
}
