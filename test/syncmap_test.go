package test
import (
	"sync"
	"testing"
)

type Map struct {
	m map[int]int
	sync.RWMutex
}

type SMap struct {
	sm sync.Map
}

func (m *Map) Insert(i int, s int, wg *sync.WaitGroup) {
	m.Lock()
	m.m[i] = s
	m.Unlock()
	wg.Done()
}

func (sm *SMap) Insert(i int, s int, wg *sync.WaitGroup) {
	sm.sm.Store(i, s)
	wg.Done()
}

func (m *Map) Get(i int, wg *sync.WaitGroup) (s int) {
	defer wg.Done()
	m.RLock()
	s, ok := m.m[i]
	if ok {
		m.RUnlock()
		return s
	}
	m.RUnlock()
	return 0
}

func (sm *SMap) Get(i int, wg *sync.WaitGroup) (s int) {
	defer wg.Done()
	v, ok := sm.sm.Load(i)
	if ok {
		return v.(int)
	}
	return 0
}

func (m *Map) Delete(i int, wg *sync.WaitGroup) {
	m.Lock()
	delete(m.m, i)
	m.Unlock()
	wg.Done()
}

func (sm *SMap) Delete(i int, wg *sync.WaitGroup) {
	sm.sm.Delete(i)
	wg.Done()
}

func operateMap(m *Map, work int) {
	wg := sync.WaitGroup{}
	wg.Add(work*2 + 3)
	go func() {
		defer wg.Done()
		for i := 0; i < work; i++ {
			go m.Insert(i, i, &wg)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < work; i++ {
			if i%4 == 0 {
				wg.Add(1)
				go m.Delete(i, &wg)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < work; i++ {
			go m.Get(i, &wg)
		}
	}()

	wg.Wait()

}

func operateSyncMap(sm *SMap, work int) {
	wg := sync.WaitGroup{}
	wg.Add(work*2 + 3)
	go func() {
		defer wg.Done()
		for i := 0; i < work; i++ {
			go sm.Insert(i, i, &wg)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < work; i++ {
			if i%4 == 0 {
				wg.Add(1)
				go sm.Delete(i, &wg)
			}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < work; i++ {
			go sm.Get(i, &wg)
		}
	}()
	wg.Wait()

}
func BenchmarkOperateSyncMap8Work(b *testing.B) {
	sm := SMap{
		sm: sync.Map{},
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			operateSyncMap(&sm, 8)
		}

	})
}
func BenchmarkOperateMap8Work(b *testing.B) {
	m := Map{
		m: make(map[int]int, 0),
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			operateMap(&m, 8)
		}

	})
}

func BenchmarkOperateSyncMap256Work(b *testing.B) {
	sm := SMap{
		sm: sync.Map{},
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			operateSyncMap(&sm, 256)
		}

	})
}
func BenchmarkOperateMap256Work(b *testing.B) {
	m := Map{
		m: make(map[int]int, 0),
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			operateMap(&m, 256)
		}

	})
}

