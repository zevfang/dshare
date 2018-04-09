package main

import "fmt"

//具体任务
type PayLoad struct {
	Name string
}

func (p *PayLoad) GetData() error {
	fmt.Printf("Hello Html. %s", p.Name)
	return nil
}


//任务对象
type Job struct {
	PayLoad PayLoad
}

var JobQueue chan Job

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				if err := job.PayLoad.GetData(); err != nil {
					fmt.Println("get html is err.")
				}
			case <-w.quit:
				fmt.Println("stop")
				return

			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

//调度器
type Dispatcher struct {
	WorkerPool chan chan Job
	maxWorkers int
}

func Prn()  {
	fmt.Println("daye")
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Job, maxWorkers),
		maxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		w := NewWorker(d.WorkerPool)
		w.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				//试图获得一个工人工作通道可用。
				jobChannel := <-d.WorkerPool

				//分派工作的工人工作通道
				jobChannel <- job
			}(job)
		}
	}
}
