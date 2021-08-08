/**
 * @Author KYIMH
 * @Description
 * @Date 2021/8/7 16:30
 **/

package zookeeper

//zookeeper client operators
type ZkClient interface {
	Connect() (err error)
	Close()
}

//zookeeper node and data operators
type ZkDal interface {
	SetWatch()
	watchHandler()
	GetNodeData(path string) ([]byte, error)
	Watch()
}
