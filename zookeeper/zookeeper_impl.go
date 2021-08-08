//Auther: scola
//Date: 2021/07/22 20:27
//Description:
//Σ(っ °Д °;)っ

package zookeeper

import (
	"fmt"
	"github.com/KYIMH/CCS_Utils/share/enum/staict_const"
	"github.com/pochard/zkutils"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type ZkClientImpl struct {
	//zk host
	Host []string
	//listen zk path
	WatchPath []string
	//zk conn
	Conn *zk.Conn
	//watcher
	KeepWatcher *zkutils.KeepWatcher
	//TODO: authentication, need to be added later
	Auth string
	// zk event
	ZkEvent <-chan zk.Event
	//config data channel
	ConfigChan chan []byte
}

//create new zk client
func NewZkClient(host, watchPath []string, auth string) *ZkClientImpl {
	ch := make(chan []byte, 1)

	return &ZkClientImpl{
		Host:       host,
		WatchPath:  watchPath,
		Auth:       auth,
		ConfigChan: ch,
	}
}

//connect to zk and new a watcher
func (z *ZkClientImpl) Connect() (err error) {
	z.Conn, z.ZkEvent, err = zk.Connect(z.Host, time.Second*5)
	if err != nil {
		fmt.Println("zk connect error!")
		return err
	}
	z.KeepWatcher = zkutils.NewKeepWatcher(z.Conn)
	return nil
}

//close zk connection
func (z *ZkClientImpl) Close() {
	z.Conn.Close()
}

//set a watcher
func (z *ZkClientImpl) SetWatch() {
	for _, path := range z.WatchPath {
		var tmpPath = path
		go z.KeepWatcher.WatchData(tmpPath, func(data []byte, err error) {
			if err != nil {
				fmt.Println("watch error:", err)
			}
			//TODO:监听数据变动后序操作
			fmt.Printf("path: %s, value: %s\n", tmpPath, string(data))

			//update global_server config
			if tmpPath == staict_const.ConfigPath {
				z.ConfigChan <- data
				//configuration.UpdateServerConfig(data)
			}
		})
	}
}

//handle zk data change
func (z *ZkClientImpl) watchHandler() {
	for {
		select {
		case event := <-z.ZkEvent:
			switch event.Type {
			// node data change
			case zk.EventNodeDataChanged:
				//todo something
				fmt.Print("Zk data changed!\n")
			case zk.EventNodeDeleted:
				// todo something
				fmt.Print("Zk node deleted!\n")
			case zk.EventSession:
				// todo something
				fmt.Print("Got session event!\n")
			}
		}
	}
}

//get node data
func (z *ZkClientImpl) GetNodeData(path string) ([]byte, error) {
	data, _, err := z.Conn.Get(path)
	return data, err
}

//listen zk connection chan
func (z *ZkClientImpl) Watch() {
	//z.watchHandler()
	//go z.watchHandler()
	z.SetWatch()
}
