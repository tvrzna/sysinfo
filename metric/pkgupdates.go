package metric

import (
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/tvrzna/pkgtray/checker"
)

type pkgUpdates struct {
	disabled  bool
	lastCheck int64
	updates   int
	manager   *pkgManager
	mutex     *sync.Mutex
}

type pkgManager struct {
	name    string
	command string
	path    string
	pkgType checker.EnPkgManager
}

var pUpdates *pkgUpdates = &pkgUpdates{mutex: &sync.Mutex{}}

func LoadPkgUpdates() int {
	if pUpdates != nil && pUpdates.disabled {
		return 0
	}
	if pUpdates.mutex.TryLock() {
		go pUpdates.update()
	}
	return pUpdates.updates
}

func initPkgUpdates() {
	xbpsPkgMngr := pkgManager{"xbps", "xbps-install", "", checker.Xbps}
	pacmanPkgMngr := pkgManager{"pacman", "checkupdates", "", checker.Pacman}
	apkPkgMngr := pkgManager{"apk", "apk", "", checker.Apk}
	aptGetPkgManager := pkgManager{"apt-get", "apt-get", "", checker.Apt_get}

	pkgManagers := []pkgManager{xbpsPkgMngr, pacmanPkgMngr, apkPkgMngr, aptGetPkgManager}

	for _, pkgManager := range pkgManagers {
		path, err := exec.LookPath(pkgManager.command)
		if err == nil {
			pkgManager.path = path
			pUpdates.manager = &pkgManager
			break
		}
	}
	if pUpdates.manager == nil {
		log.Print("pkgupdates: no suitable package manager found")
		pUpdates.disabled = true
	}
}

func (p *pkgUpdates) update() {
	if p == nil || p.lastCheck == 0 {
		initPkgUpdates()
	}
	if p != nil && !p.disabled && p.lastCheck+600 <= time.Now().Unix() {
		p.updates = checker.CheckPackages(p.manager.pkgType, p.manager.path)
		p.lastCheck = time.Now().Unix()
	}
	p.mutex.Unlock()
}
