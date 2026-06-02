package changer

import (
	"math/rand"
	"sync"
	"time"

	"github.com/robot3human0/gnome-background-changer/internal/config"
	"github.com/robot3human0/gnome-background-changer/internal/wallpaper"
)

type Changer struct {
	mtx      sync.Mutex
	cfg      *config.Config
	files    []string
	idx      int
	ticker   *time.Ticker
	stopCh   chan struct{}
	running  bool
	onError  func(error)
	onChange func(string) 		// called with path of new wallpaper
	setter   func(string) error	// injectable dependency for setting wallpaper
}

func NewChanger(cfg *config.Config, onError func(error), onChange func(string), setter func(string)error) *Changer {
	return &Changer {
		cfg:      cfg,
		onError:  onError,
		onChange: onChange,
		stopCh:   make(chan struct{}),
		setter:   setter,
	}
}

func (ch *Changer) IsRunning() bool {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()
	return ch.running
}

func (ch *Changer) Start() {
	ch.mtx.Lock()
	if ch.running {
		ch.mtx.Unlock()
		return
	}
	ch.running = true
	ch.stopCh = make(chan struct{})
	ch.mtx.Unlock()

	if err := ch.reload(); err != nil {
		ch.onError(err)
		return
	}

	// Show current wallpaper immediately on start
	ch.applyNext()

	go ch.loop()
}

func (ch *Changer) Stop() {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()
	if ch.running {
		close(ch.stopCh)
		ch.running = false
	}
}

func (ch *Changer) Next() {
	ch.mtx.Lock()
	if err := ch.reloadLocked(); err != nil {
		ch.mtx.Unlock()
		ch.onError(err)
		return
	}
	path := ch.pickNext()
	ch.mtx.Unlock()

	if err := ch.setter(path); err != nil {
		ch.onError(err)
		return
	}

	if ch.onChange != nil {
		ch.onChange(path)
	}

	// Reset ticker so the interval restarts from now
	ch.mtx.Lock()
	if ch.ticker != nil {
		ch.ticker.Reset(ch.cfg.Interval)
	}
	ch.mtx.Unlock()
}

func (ch *Changer) UpdateConfig(cfg config.Config) {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()
	*ch.cfg = cfg
}

func (ch *Changer) reload() error {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()

	files, err := wallpaper.ListWallpapers(ch.cfg.FolderPath)
	if err != nil {
		return err
	}
	ch.files = files
	if ch.idx >= len(ch.files) {
		ch.idx = 0
	}

	return nil
}

func (ch *Changer) loop() {
	ch.mtx.Lock()
	interval := ch.cfg.Interval
	ch.ticker = time.NewTicker(interval)
	ticker := ch.ticker
	stopCh := ch.stopCh
	ch.mtx.Unlock()

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ch.applyNext()

			// Update ticker if interval changed
			ch.mtx.Lock()
			if ch.cfg.Interval != interval {
				interval = ch.cfg.Interval
				ticker.Reset(interval)
			}
			ch.mtx.Unlock()

		case <-stopCh:
			return
		}
	}	
}

func (ch *Changer) applyNext() {
	ch.mtx.Lock()
	if err := ch.reloadLocked(); err != nil {
		ch.mtx.Unlock()
		ch.onError(err)
		return
	}
	path := ch.pickNext()
	ch.mtx.Unlock()

	
	if err := ch.setter(path); err != nil {
		ch.onError(err)
		return
	}

	if ch.onChange != nil {
		ch.onChange(path)
	}
}

// reloadLocked reloads file list. Caller must hold mtx.
func (ch *Changer) reloadLocked() error {
	files, err := wallpaper.ListWallpapers(ch.cfg.FolderPath)
	if err != nil {
		return err
	}
	ch.files = files
	if ch.idx >= len(files) {
		ch.idx = 0
	}
	return nil
}

// picks and advances the idx. Caller must hold mtx
func (ch *Changer) pickNext() string {
	if len(ch.files) == 0 {
		return ""
	}
	if ch.cfg.Random {
		ch.idx = rand.Intn(len(ch.files))
	} else {
		ch.idx = (ch.idx + 1) % len(ch.files)
	}
	ch.cfg.LastIndex = ch.idx
	return ch.files[ch.idx]
}
