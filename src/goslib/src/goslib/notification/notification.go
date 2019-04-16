package notification

import (
	"fmt"
	"github.com/appleboy/go-fcm"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"goslib/gen_server"
	"goslib/logger"
	"log"
)

const IOS_SERVER = "__ios_notification_server__"
const FCM_SERVER = "__android_notification_server__"

const (
	CHANNEL_IOS = iota
	CHANNEL_ANDROID
)

/*
   GenServer Callbacks
*/
type Server struct {
	bundleId  string
	iosClient *apns2.Client
	fcmClient *fcm.Client
}

func StartIOS(bundleId string, iosCertP12Path string, iosP12Password string, isProduction bool) {
	gen_server.Start(IOS_SERVER, new(Server), IOS_SERVER, bundleId, iosCertP12Path, iosP12Password, isProduction)
}

func StartFCM(fcmAPIKey string) {
	gen_server.Start(FCM_SERVER, new(Server), FCM_SERVER, fcmAPIKey)
}

func Send(channel int, deviceToken string, content string) {
	switch channel {
	case CHANNEL_IOS:
		gen_server.Cast(IOS_SERVER, &SendIOSParams{deviceToken, content})
	case CHANNEL_ANDROID:
		gen_server.Cast(FCM_SERVER, &SendGPParams{deviceToken, content})
	}
}

func (self *Server) Init(args []interface{}) (err error) {
	category := args[0].(string)
	switch category {
	case IOS_SERVER:
		self.bundleId = args[1].(string)
		iosCertP12Path := args[2].(string)
		iosP12Password := args[3].(string)
		isProduction := args[4].(bool)

		cert, err := certificate.FromP12File(iosCertP12Path, iosP12Password)
		if err != nil {
			logger.ERR("Start ios push failed: ", err)
			return err
		}

		if isProduction {
			self.iosClient = apns2.NewClient(cert).Production()
		} else {
			self.iosClient = apns2.NewClient(cert).Development()
		}
		break
	case FCM_SERVER:
		fcmAPIKey := args[1].(string)
		self.fcmClient, err = fcm.NewClient(fcmAPIKey)
		if err != nil {
			logger.ERR("Start FCM failed: ", err)
			return err
		}
	}

	return nil
}

func (self *Server) HandleCast(msg interface{}) {
	switch params := msg.(type) {
	case *SendIOSParams:
		self.sendIOS(params)
		break
	case *SendGPParams:
		self.sendGP(params)
		break
	}
}

type SendIOSParams struct {
	deviceToken string
	content string
}
func (self *Server) sendIOS(params *SendIOSParams) {
	notification := &apns2.Notification{}
	notification.DeviceToken = params.deviceToken
	notification.Topic = self.bundleId
	notification.Payload = payload.NewPayload().Alert(params.content).Badge(1)

	res, err := self.iosClient.Push(notification)

	if err != nil {
		logger.ERR("send ios push failed: ", err)
	}
	if res.Sent() {
		log.Println("Sent:", res.ApnsID)
	} else {
		fmt.Printf("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
	}
}

type SendGPParams struct {
	deviceToken string
	content string
}
func (self *Server) sendGP(params *SendGPParams) {
	msg := &fcm.Message{
		To: params.deviceToken,
		Notification: &fcm.Notification{
			Body:  params.content,
			Badge: "1",
			Sound: "default",
		},
	}
	_, err := self.fcmClient.Send(msg)
	if err != nil {
		logger.ERR("send fcm push failed: ", err)
	}
}

func (self *Server) HandleCall(msg interface{}) (interface{}, error) {
	return nil, nil
}

func (self *Server) Terminate(reason string) (err error) {
	return nil
}
