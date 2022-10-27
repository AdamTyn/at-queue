package queue

type Executor func() error

func (e Executor) Do() {
	if err := e(); err != nil {
		println(err.Error())
	}
}
