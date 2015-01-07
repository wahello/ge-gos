package gslib

import (
	"api"
	"fmt"
	"gslib/gen_server"
	"gslib/routes"
	. "gslib/utils"
	"gslib/utils/packet"
	"net"
	"runtime"
)

type Player struct {
	playerId  string
	processed int
	Conn      net.Conn
}

type WrapHandler func() interface{}
type AsyncWrapHandler func()

/*
   GenServer Callbacks
*/
func (self *Player) Init(args []interface{}) (err error) {
	name := args[0].(string)
	fmt.Println("server ", name, " started!")
	self.playerId = name
	return nil
}

func (self *Player) HandleCast(args []interface{}) {
	method_name := args[0].(string)
	if method_name == "HandleRequest" {
		self.HandleRequest(args[1].([]byte), args[2].(net.Conn))
	} else if method_name == "HandleWrap" {
		self.HandleWrap(args[1].(WrapHandler))
	} else if method_name == "removeConn" {
		self.Conn = nil
	}
}

func (self *Player) HandleCall(args []interface{}) interface{} {
	method_name := args[0].(string)
	if method_name == "HandleWrap" {
		return self.HandleWrap(args[1].(WrapHandler))
	}
	return nil
}

func (self *Player) Terminate(reason string) (err error) {
	fmt.Println("callback Termiante!")
	return nil
}

/*
   IPC Methods
*/

func (self *Player) SystemInfo() int {
	return runtime.NumCPU()
}

func (self *Player) SendData(encode_method string, msg interface{}) {
	if self.Conn != nil {
		writer := api.Encode(encode_method, msg)
		writer.Send(self.Conn)
	}
}

func (self *Player) HandleRequest(data []byte, conn net.Conn) {
	self.Conn = conn
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("caught panic in player HandleRequest(): ", x)
		}
	}()
	reader := packet.Reader(data)
	protocol := reader.ReadUint16()
	decode_method := api.IdToName[protocol]
	handler, err := routes.Route(decode_method)
	if err == nil {
		params := api.Decode(decode_method, reader)
		encode_method, response := handler(self, params)
		writer := api.Encode(encode_method, response)

		self.processed++
		// INFO("Processed: ", self.processed, " Response Data: ", response_data)
		if self.Conn != nil {
			writer.Send(self.Conn)
		}
	} else {
		ERR(err)
	}
}

func (self *Player) HandleWrap(fun WrapHandler) interface{} {
	return fun()
}

func (self *Player) HandleAsyncWrap(fun AsyncWrapHandler) {
	fun()
}

func (self *Player) Wrap(targetPlayerId string, fun WrapHandler) (interface{}, error) {
	if self.playerId == targetPlayerId {
		return self.HandleWrap(fun), nil
	} else {
		return gen_server.Call(targetPlayerId, "HandleWrap", fun)
	}
}

func (self *Player) AsyncWrap(targetPlayerId string, fun AsyncWrapHandler) {
	if self.playerId == targetPlayerId {
		self.HandleAsyncWrap(fun)
	} else {
		gen_server.Cast(targetPlayerId, "HandleAsyncWrap", fun)
	}
}