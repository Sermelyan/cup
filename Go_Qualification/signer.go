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
	for _, jobToSchedule := range freeFlowJobs {
		in, out = out, make(chan any)
		go (func(j job, i, o chan any) {
			ready <- struct{}{}
			j(i, o)
			close(o)
		})(jobToSchedule, in, out)
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
		go func(i, f string) {
			out <- genProfile(i, f)
			wg.Done()
		}(id, fp)
	}
	wg.Wait()
}

func GetGroup(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(d string) {
			out <- genGroup(d)
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

const groupCount = 6

type group [groupCount]string

var th group

func init() {
	for i := 0; i < groupCount; i++ {
		th[i] = strconv.Itoa(i)
	}
}

func genGroup(data string) string {
	res := group{}
	wg := &sync.WaitGroup{}
	wg.Add(groupCount)
	for i := 0; i < groupCount; i++ {
		go func(pos int) {
			res[pos] = SlowPredict(th[pos] + data)
			wg.Done()
		}(i)
	}
	wg.Wait()
	return strings.Join(res[:], "")
}

func genProfile(id, fp string) string {
	idc, fpc := make(chan string, 1), make(chan string, 1)
	go func() {
		idc <- SlowPredict(id)
	}()
	go func() {
		fpc <- SlowPredict(fp)
	}()
	return <-idc + "-" + <-fpc
}
