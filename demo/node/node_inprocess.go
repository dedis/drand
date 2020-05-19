package node

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/drand/drand/core"
	"github.com/drand/drand/key"
	"github.com/drand/drand/log"
	"github.com/drand/drand/net"
	"github.com/drand/drand/protobuf/drand"
	"github.com/drand/drand/test"
	"github.com/kabukky/httpscerts"
)

type LocalNode struct {
	base       string
	i          int
	period     string
	publicPath string
	certPath   string
	// certificate key
	keyPath string
	// where all public certs are stored
	certFolder string
	logPath    string
	privAddr   string
	pubAddr    string
	ctrlAddr   string
	tls        bool

	log log.Logger

	daemon *core.Drand
}

func NewLocalNode(i int, period string, base string, tls bool) Node {
	nbase := path.Join(base, fmt.Sprintf("node-%d", i))
	os.MkdirAll(nbase, 0740)
	logPath := path.Join(nbase, "log")

	// make certificates for the node.
	err := httpscerts.Generate(
		path.Join(base, fmt.Sprintf("server-%d.crt", i)),
		path.Join(base, fmt.Sprintf("server-%d.key", i)),
		test.LocalHost())
	if err != nil {
		return nil
	}
	return &LocalNode{
		base:     nbase,
		i:        i,
		period:   period,
		tls:      tls,
		logPath:  logPath,
		log:      log.DefaultLogger,
		pubAddr:  test.FreeBind(true),
		privAddr: test.FreeBind(true),
		ctrlAddr: test.FreeBind(false),
	}
}

func (l *LocalNode) Start(certFolder string) error {
	opts := []core.ConfigOption{
		core.WithConfigFolder(l.base),
		core.WithPublicListenAddress(l.pubAddr),
		core.WithPrivateListenAddress(l.privAddr),
		core.WithControlPort(l.ctrlAddr),
	}
	var priv *key.Pair
	if l.tls {
		priv = key.NewTLSKeyPair(l.privAddr)
		opts = append(opts, core.WithTLS(
			path.Join(l.base, fmt.Sprintf("server-%d.crt", l.i)),
			path.Join(l.base, fmt.Sprintf("server-%d.key", l.i))))
	} else {
		priv = key.NewKeyPair(l.privAddr)
		opts = append(opts, core.WithInsecure())
	}
	conf := core.NewConfig(opts...)
	fs := key.NewFileStore(conf.ConfigFolder())
	fs.SaveKeyPair(priv)
	key.Save(path.Join(l.base, "public.toml"), priv.Public, false)
	drand, err := core.NewDrand(fs, conf)
	if err != nil {
		return err
	}
	l.daemon = drand
	return nil
}

func (l *LocalNode) PrivateAddr() string {
	return l.privAddr
}

func (l *LocalNode) PublicAddr() string {
	return l.pubAddr
}

func (l *LocalNode) Index() int {
	return l.i
}

func (l *LocalNode) RunDKG(nodes, thr int, timeout string, leader bool, leaderAddr string, beaconOffset int) *key.Group {
	cl, err := net.NewControlClient(l.ctrlAddr)
	if err != nil {
		l.log.Error("drand", "can't instantiate control client", "err", err)
		return nil
	}

	p, _ := time.ParseDuration(l.period)
	t, _ := time.ParseDuration(timeout)
	var grp *drand.GroupPacket
	if leader {
		grp, err = cl.InitDKGLeader(nodes, thr, p, t, nil, secretDKG, beaconOffset)
	} else {
		leader := net.CreatePeer(leaderAddr, l.tls)
		grp, err = cl.InitDKG(leader, nodes, thr, t, nil, secretDKG)
	}
	if err != nil {
		l.log.Error("drand", "dkg run failed", "err", err)
		return nil
	}
	kg, _ := key.GroupFromProto(grp)
	return kg
}

func (l *LocalNode) GetGroup() *key.Group {
	cl, err := net.NewControlClient(l.ctrlAddr)
	if err != nil {
		l.log.Error("drand", "can't instantiate control client", "err", err)
		return nil
	}
	grp, err := cl.GroupFile()
	if err != nil {
		l.log.Error("drand", "can't  get group", "err", err)
		return nil
	}
	group, err := key.GroupFromProto(grp)
	if err != nil {
		l.log.Error("drand", "can't deserialize group", "err", err)
		return nil
	}
	return group
}

func (l *LocalNode) RunReshare(nodes, thr int, oldGroup string, timeout string, leader bool, leaderAddr string, beaconOffset int) *key.Group {
	cl, err := net.NewControlClient(l.ctrlAddr)
	if err != nil {
		l.log.Error("drand", "can't instantiate control client", "err", err)
		return nil
	}

	t, _ := time.ParseDuration(timeout)
	var grp *drand.GroupPacket
	if leader {
		grp, err = cl.InitReshareLeader(nodes, thr, t, secretReshare, oldGroup, beaconOffset)
	} else {
		leader := net.CreatePeer(leaderAddr, l.tls)
		grp, err = cl.InitReshare(leader, nodes, thr, t, secretReshare, oldGroup)
	}
	if err != nil {
		l.log.Error("drand", "reshare failed", "err", err)
		return nil
	}
	kg, _ := key.GroupFromProto(grp)
	return kg
}

func (l *LocalNode) GetCokey(group string) bool {
	cl, err := net.NewControlClient(l.ctrlAddr)
	if err != nil {
		l.log.Error("drand", "can't instantiate control client", "err", err)
		return false
	}
	key, err := cl.CollectiveKey()
	if err != nil {
		l.log.Error("drand", "can't get cokey", "err", err)
		return false
	}
	sdist := hex.EncodeToString(key.GetCoKey())
	fmt.Printf("\t- Node %s has cokey %s\n", l.PrivateAddr(), sdist[10:14])
	return true
}

func (l *LocalNode) Ping() bool {
	cl, err := net.NewControlClient(l.ctrlAddr)
	if err != nil {
		l.log.Error("drand", "can't instantiate control client", "err", err)
		return false
	}
	if err := cl.Ping(); err != nil {
		l.log.Error("drand", "can't ping", "err", err)
		return false
	}
	return true
}

func (l *LocalNode) GetBeacon(groupPath string, round uint64) (resp *drand.PublicRandResponse, cmd string) {
	c := core.NewGrpcClient()
	if l.tls {
		m := net.NewCertManager()
		m.Add(path.Join(l.base, fmt.Sprintf("server-%d.crt", l.i)))
		c = core.NewGrpcClientFromCert(m)
	}

	group := l.GetGroup()
	pk := group.PublicKey

	var err error
	cmd = "unused"
	if round == 0 {
		resp, err = c.LastPublic(l.privAddr, pk, l.tls)
	} else {
		resp, err = c.Public(l.privAddr, pk, l.tls, int(round))
	}
	if err != nil {
		l.log.Error("drand", "can't get becon", "err", err)
	}
	return
}

func (l *LocalNode) WriteCertificate(p string) {
	if l.tls {
		exec.Command("cp", path.Join(l.base, fmt.Sprintf("server-%d.crt", l.i)), p)
	}
}

func (l *LocalNode) WritePublic(p string) {
	exec.Command("cp", path.Join(l.base, "public.toml"), p)
}

func (l *LocalNode) Stop() {
	cl, err := net.NewControlClient(l.ctrlAddr)
	if err != nil {
		l.log.Error("drand", "can't instantiate control client", "err", err)
		return
	}
	_, err = cl.Shutdown()
	if err != nil {
		l.log.Error("drand", "failed to shutdown", "err", err)
	}
	<-l.daemon.WaitExit()
}

func (l *LocalNode) PrintLog() {
	fmt.Printf("[-] Printing logs of node %s:\n", l.PrivateAddr())
	buff, err := ioutil.ReadFile(l.logPath)
	if err != nil {
		fmt.Printf("[-] Can't read logs !\n\n")
		return
	}
	os.Stdout.Write([]byte(buff))
	fmt.Println()
}
