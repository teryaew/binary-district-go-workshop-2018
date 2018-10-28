package pool

type Pool struct{
  MaxWorkers int

  work chan func()
  sem chan struct{}
}

func NewPool(parallelism int) *Pool {
   return &Pool{
     sem: make(chan struct{}, parallelism),
     work: make(chan func()),
   }
}
func (p *Pool) Exec(task func()) {
  select {
  case p.work <- task:
  default:
    select {
    case p.sem <- struct{}{}:
      go p.worker(task)
    case p.work <- task:
    }
  }
}

func (p *Pool) Close() {
  close(p.work)
  // Wait for all workers are done
  for i := 0; i < cap(p.sem); i++ {
    p.sem <- struct{}{}
  }
}

func (p *Pool) worker(task func()) {
  defer func() { <-p.sem }()
  var ok bool
  for {
    task()
    if task, ok = <-p.work; !ok {
      return
    }
  }
}
