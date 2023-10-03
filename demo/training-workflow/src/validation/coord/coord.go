package coord

import (
	"sync"
)

var lock = &sync.Mutex{}

type single struct {
	singleMutex *sync.Mutex
}

var singleInstance *single

func GetInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &single{
				singleMutex: &sync.Mutex{},
			}
		}
	}

	return singleInstance
}

func (s *single) GetMutex() *sync.Mutex {
	return s.singleMutex
}
