package uos

type Message struct {
	isAt    bool `json:"isAt"`
	AppInfo struct {
		Type  int    `json:"type"`
		AppID string `json:"appID"`
	} `json:"appInfo"`
	AppMsgType            int                    `json:"appMsgType"`
	HasProductId          int                    `json:"hasProductId"`
	ImgHeight             int                    `json:"imgHeight"`
	ImgStatus             int                    `json:"imgStatus"`
	ImgWidth              int                    `json:"imgWidth"`
	ForwardFlag           int                    `json:"forwardFlag"`
	MsgType               int                    `json:"msgType"`
	Status                int                    `json:"status"`
	StatusNotifyCode      int                    `json:"statusNotifyCode"`
	SubMsgType            int                    `json:"subMsgType"`
	VoiceLength           int                    `json:"voiceLength"`
	CreateTime            int64                  `json:"createTime"`
	NewMsgId              int64                  `json:"newMsgId"`
	PlayLength            int64                  `json:"playLength"`
	MediaId               string                 `json:"mediaId"`
	MsgId                 string                 `json:"msgId"`
	EncryFileName         string                 `json:"encryFileName"`
	FileName              string                 `json:"fileName"`
	FileSize              string                 `json:"fileSize"`
	Content               string                 `json:"content"`
	FromUserName          string                 `json:"fromUserName"`
	OriContent            string                 `json:"oriContent"`
	StatusNotifyUserName  string                 `json:"statusNotifyUserName"`
	Ticket                string                 `json:"ticket"`
	ToUserName            string                 `json:"toUserName"`
	Url                   string                 `json:"url"`
	senderUserNameInGroup string                 `json:"senderUserNameInGroup"`
	RecommendInfo         RecommendInfo          `json:"recommendInfo"`
	item                  map[string]interface{} `json:"item"`
}

// RecommendInfo 一些特殊类型的消息会携带该结构体信息
type RecommendInfo struct {
	OpCode     int
	Scene      int
	Sex        int
	VerifyFlag int
	AttrStatus int64
	QQNum      int64
	Alias      string
	City       string
	Content    string
	NickName   string
	Province   string
	Signature  string
	Ticket     string
	UserName   string
}
