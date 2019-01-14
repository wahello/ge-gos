package gosconf

import (
	"google.golang.org/grpc"
	"time"
)

const IS_DEBUG = false

// 支持的通信协议
const (
	AGENT_PROTOCOL_TCP = iota
	AGENT_PROTOCOL_WS
)

// 支持的编码方式
const (
	PROTOCOL_ENCODING_RAW = iota
	PROTOCOL_ENCODING_JSON
	PROTOCOL_ENCODING_PB
)

const (
	AGENT_PROTOCOL = AGENT_PROTOCOL_WS
	AGENT_ENCODING = PROTOCOL_ENCODING_JSON
)

var REDIS_CLUSTERS = []string{
	"redis-cluster-0.redis-cluster.default.svc.cluster.local",
	"redis-cluster-1.redis-cluster.default.svc.cluster.local",
	"redis-cluster-2.redis-cluster.default.svc.cluster.local",
	"redis-cluster-3.redis-cluster.default.svc.cluster.local",
	"redis-cluster-4.redis-cluster.default.svc.cluster.local",
	"redis-cluster-5.redis-cluster.default.svc.cluster.local",
}

const (
	WORLD_SERVER_IP = "gos-world-service.default.svc.cluster.local"

	GAME_DOMAIN = "gos-game-service.default.svc.cluster.local"
)

const (
	TCP_READ_TIMEOUT      = 60 * time.Second
	HEARTBEAT             = 5 * time.Second
	SERVICE_DEAD_DURATION = 16
	RPC_REQUEST_TIMEOUT   = 5 * time.Second
	AGENT_CCU_MAX         = 20000
	GAME_CCU_MAX          = 6000
)

const (
	TIMERTASK_CHECK_DURATION  = time.Second // seconds
	TIMERTASK_TASKS_PER_CHECK = 100         // max tasks fetched per check
	TIMERTASK_MAX_RETRY       = 3           // retry 3 times
)

// Keys for retrive redis data
const (
	RK_GAME_APP_IDS = "__GAME_APP_IDS__"
	RK_SCENE_IDS    = "__SCENE_IDS__"
)

/*
HTTP Servers
*/

const (
	AUTH_SERVICE_PORT = "3000"
)

/*
TCP Servers
*/

type TCP struct {
	Packet  uint8
	Network string
	Address string
}

var TCP_SERVER_CONNECT_APP = &TCP{
	Packet:  2,
	Network: "tcp",
	Address: ":4000",
}

/*
Rpc Servers
*/
type Rpc struct {
	ListenNet   string // must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket"
	ListenAddr  string
	DialAddress string
	DialOptions []grpc.DialOption
}

var RPC_FOR_CONNECT_APP_MGR = &Rpc{
	ListenNet:   "tcp4",
	ListenAddr:  ":50051",
	DialAddress: WORLD_SERVER_IP + ":50051",
	DialOptions: []grpc.DialOption{grpc.WithInsecure()},
}

var RPC_FOR_GAME_APP_MGR = &Rpc{
	ListenNet:   "tcp4",
	ListenAddr:  ":50052",
	DialAddress: WORLD_SERVER_IP + ":50052",
	DialOptions: []grpc.DialOption{grpc.WithInsecure()},
}

// Dial addresses are dynamic dispatched by GameAppMgr
var RPC_FOR_GAME_APP_STREAM = &Rpc{
	ListenNet:   "tcp4",
	ListenAddr:  ":50053",
	DialAddress: "",
	DialOptions: []grpc.DialOption{grpc.WithInsecure()},
}
