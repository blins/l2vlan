package main

type DebugLogger interface {
	Println(...interface{})
}

type Debugger struct {
	logger DebugLogger
}

func (self *Debugger) Println(v ...interface{}) {
	self.logger.Println(v...)
}

func (self *Debugger) SetLogger(l DebugLogger) {
	self.logger = l
}
