package main

import (
	"sync"
)

// сюда писать код
var ok = true
var recieved uint32

type JobAndChanel struct {
	Job       job
	OutChanel chan interface{}
	InChanel  chan interface{}
}

func NewJobAndChanel(j job) JobAndChanel {
	c := make(chan interface{})
	jac := JobAndChanel{
		Job:       j,
		OutChanel: c,
	}
	return jac
}

func CreateJobs(jobs []job) []JobAndChanel {
	createdJobs := make([]JobAndChanel, 0, len(jobs))
	for _, j := range jobs {
		jac := NewJobAndChanel(j)
		createdJobs = append(createdJobs, jac)
	}
	return createdJobs
}



func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	createdJobs := CreateJobs(jobs)

	for idx, j := range createdJobs {
		wg.Add(1)
		go func(idx int, j JobAndChanel) {
			defer wg.Done()
			var in chan interface{}
			var out chan interface{}

			switch idx {
			case 0:
				in = nil
				out = j.OutChanel
			case len(createdJobs) - 1:
				in = createdJobs[idx-1].OutChanel
				out = j.OutChanel
			default:
				in = createdJobs[idx-1].OutChanel
				out = j.OutChanel
			}
			j.Job(in, out)
			close(out)
		}(idx, j)
	}

	wg.Wait()
}

