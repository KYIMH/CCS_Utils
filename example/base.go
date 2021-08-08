/**
 * @Author KYIMH
 * @Description
 * @Date 2021/8/7 21:21
 **/

package main

import (
	"fmt"
	"github.com/KYIMH/CCS_Utils/mongo"
)

func main() {
	mongoCli := mongo.NewMongoClient()

	fmt.Printf("Succeed new %+v\n", mongoCli.Context)
}
