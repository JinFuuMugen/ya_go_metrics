package pool

import "testing"

type testObj struct {
	value int
}

func (o *testObj) Reset() {
	o.value = 0
}

func TestPool_GetPut(t *testing.T) {
	p := New(func() *testObj {
		return &testObj{}
	})

	obj1 := p.Get()
	if obj1 == nil {
		t.Fatal("get returned nil")
	}

	if obj1.value != 0 {
		t.Errorf("got value = %d, want 0", obj1.value)
	}

	obj1.value = 42
	p.Put(obj1)

	if obj1.value != 0 {
		t.Errorf("got value = %d, want 0", obj1.value)
	}

	obj2 := p.Get()
	if obj2 == nil {
		t.Fatal("get returned nil")
	}
	if obj2.value != 0 {
		t.Errorf("got value = %d, want 0", obj2.value)
	}
	if obj2 != obj1 {
		t.Errorf("got a different object")
	}

	obj3 := p.Get()
	if obj3 == nil {
		t.Fatal("expected a new object")
	}
	if obj3 == obj2 {
		t.Errorf("expected a new object distinct from obj2 when pool is empty")
	}
	p.Put(obj2)
	p.Put(obj3)

	gotA := p.Get()
	gotB := p.Get()
	if !((gotA == obj2 && gotB == obj3) || (gotA == obj3 && gotB == obj2)) {
		t.Errorf("gotA = %p, gotB = %p, want obj2 = %p and obj3 = %p", gotA, gotB, obj2, obj3)
	}
}
