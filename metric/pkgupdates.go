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
}

type pkgManager struct {
	name    string
	command string
	path    string
	pkgType checker.EnPkgManager
}

var pUpdates *pkgUpdates

func LoadPkgUpdates(wg *sync.WaitGroup, bundle *Bundle) {
	if pUpdates == nil {
		initPkgUpdates()
	}
	if pUpdates.disabled {
		wg.Done()
	}

	if pUpdates.lastCheck+600 <= time.Now().Unix() {
		pUpdates.updates = checker.CheckPackages(pUpdates.manager.pkgType, pUpdates.manager.path)
		pUpdates.lastCheck = time.Now().Unix()
	}
	bundle.Updates = pUpdates.updates

	wg.Done()
}

func initPkgUpdates() {
	p := &pkgUpdates{}
	xbpsPkgMngr := pkgManager{"xbps", "xbps-install", "", checker.Xbps}
	pacmanPkgMngr := pkgManager{"pacman", "checkupdates", "", checker.Pacman}
	apkPkgMngr := pkgManager{"apk", "apk", "", checker.Apk}
	aptGetPkgManager := pkgManager{"apt-get", "apt-get", "", checker.Apt_get}

	pkgManagers := []pkgManager{xbpsPkgMngr, pacmanPkgMngr, apkPkgMngr, aptGetPkgManager}

	for _, pkgManager := range pkgManagers {
		path, err := exec.LookPath(pkgManager.command)
		if err == nil {
			pkgManager.path = path
			p.manager = &pkgManager
			break
		}
	}
	if p.manager == nil {
		log.Print("pkgupdates: no suitable package manager found")
		p.disabled = true
	}

	pUpdates = p
}
