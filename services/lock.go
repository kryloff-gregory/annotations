package services

import "sync"

type LockService struct {
	locks *sync.Map
}

func NewLockService() *LockService {
	return &LockService{locks: &sync.Map{}}
}

func (s *LockService) LockItem(item string) {
	lock, _ := s.locks.LoadOrStore(item, &sync.RWMutex{})
	lock.(*sync.RWMutex).Lock()
}

func (s *LockService) UnlockItem(item string) {
	lock, _ := s.locks.Load(item)
	lock.(*sync.RWMutex).Unlock()
}

func (s *LockService) RLockItem(item string) {
	lock, _ := s.locks.LoadOrStore(item, &sync.RWMutex{})
	lock.(*sync.RWMutex).RLock()
}

func (s *LockService) RUnlockItem(item string) {
	lock, _ := s.locks.Load(item)
	lock.(*sync.RWMutex).RUnlock()
}
