package client

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/Alekssasho/GopherInvaders/core"
)

func StartClient() {
	// newEntity := core.GameEntity{BasicEntity: ecs.NewBasic()}
	// newEntity.SpaceComponent = common.SpaceComponent{
	// 	Position: engo.Point{10, 10},
	// 	Width:    64,
	// 	Height:   64,
	// }
	newEntity := core.GameEntity{ID: 0, X: 10.0, Y: 52.0}

	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	checkError(err)

	fmt.Println("client connection made")

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	for n := 0; n < 10; n++ {
		fmt.Println("client send data")
		err := encoder.Encode(newEntity)
		checkError(err)
		var entity core.GameEntity
		fmt.Println("client receive data")
		err = decoder.Decode(&entity)
		checkError(err)
		fmt.Println(entity.String())
	}

	os.Exit(0)

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
