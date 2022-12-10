package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Write your code here ...

func Advertise(freeFlowJobs ...job) {
	ready := make(chan struct{}, 1)
	var in, out chan any
	for _, j := range freeFlowJobs {
		in, out = out, make(chan any)
		go (func(j job, i, o chan any) {
			ready <- struct{}{}
			j(i, o)
			close(o)
		})(j, in, out)
		<-ready
	}
	<-out
}

func GetProfile(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		id := strconv.Itoa(data.(int))
		fp := FastPredict(id)
		wg.Add(1)
		go func(id, fp string) {
			out <- SlowPredict(id) + "-" + SlowPredict(fp)
			wg.Done()
		}(id, fp)
	}
	wg.Wait()
}

func GetGroup(in, out chan interface{}) {
	var th [6]string
	for i := 0; i < 6; i++ {
		th[i] = strconv.Itoa(i)
	}
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(d string) {
			res := ""
			for _, t := range th {
				res += SlowPredict(t + d)
			}
			out <- res
			wg.Done()
		}(data.(string))
	}
	wg.Wait()
}

func ConcatProfiles(in, out chan interface{}) {
	t := make([]string, 0)
	for data := range in {
		t = append(t, data.(string))
	}
	sort.Strings(t)
	out <- strings.Join(t, "_")
}
