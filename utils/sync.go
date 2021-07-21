package utils

import "github.com/8treenet/freedom"

func Async(name string, worker freedom.Worker, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				strace := string(GetStackTrace())
				worker.Logger().Errorf("%s recover:%v \n%s", name, err, strace)
			}
		}()

		worker.DeferRecycle()
		f()
	}()
}
