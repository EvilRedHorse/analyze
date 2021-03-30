package lockcheck

import (
	"go/types"
)

// Below are several stubs to aid in the unit testing of the lockcheck package

type fooType struct{}

func (f *fooType) String() string         { return "fooType" }
func (f *fooType) Underlying() types.Type { return f }

type fooMutex struct{}

func (f *fooMutex) String() string         { return "fooMutex" }
func (f *fooMutex) Underlying() types.Type { return f }

type fooMutexNot struct{}

func (f *fooMutexNot) String() string         { return "fooMutexNot" }
func (f *fooMutexNot) Underlying() types.Type { return f }

type syncMutex struct{}

func (s *syncMutex) String() string         { return "sync.Mutex" }
func (s *syncMutex) Underlying() types.Type { return s }

type syncRWMutex struct{}

func (s *syncRWMutex) String() string         { return "sync.RWMutex" }
func (s *syncRWMutex) Underlying() types.Type { return s }

type syncWaitGroup struct{}

func (s *syncWaitGroup) String() string         { return "sync.WaitGroup" }
func (s *syncWaitGroup) Underlying() types.Type { return s }

type siaSyncRWMutex struct{}

func (s *siaSyncRWMutex) String() string         { return "gitlab.com/NebulousLabs/Sia/sync.RWMutex" }
func (s *siaSyncRWMutex) Underlying() types.Type { return s }

type siaSyncTryMutex struct{}

func (s *siaSyncTryMutex) String() string         { return "gitlab.com/NebulousLabs/Sia/sync.TryMutex" }
func (s *siaSyncTryMutex) Underlying() types.Type { return s }

type siaSyncTryRWMutex struct{}

func (s *siaSyncTryRWMutex) String() string         { return "gitlab.com/NebulousLabs/Sia/sync.TryRWMutex" }
func (s *siaSyncTryRWMutex) Underlying() types.Type { return s }

type threadgroup struct{}

func (t *threadgroup) String() string         { return "gitlab.com/NebulousLabs/threadgroup.ThreadGroup" }
func (t *threadgroup) Underlying() types.Type { return t }
