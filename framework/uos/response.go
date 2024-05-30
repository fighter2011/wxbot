package uos

type RobotInfoResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson struct {
		HideInputBarFlag  int
		StarFriend        int
		Sex               int
		AppAccountFlag    int
		VerifyFlag        int
		ContactFlag       int
		WebWxPluginSwitch int
		HeadImgFlag       int
		SnsFlag           int
		IsOwner           int
		MemberCount       int
		ChatRoomId        int
		UniFriend         int
		OwnerUin          int
		Statues           int
		AttrStatus        int64
		Uin               int64
		Province          string
		City              string
		Alias             string
		DisplayName       string
		KeyWord           string
		EncryChatRoomId   string
		UserName          string
		NickName          string
		HeadImgUrl        string
		RemarkName        string
		PYInitial         string
		PYQuanPin         string
		RemarkPYInitial   string
		RemarkPYQuanPin   string
		Signature         string
	} `json:"data"`
}

// ObjectInfoResp 对象可以是好友、群、公众号
type ObjectInfoResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson struct {
		Data struct {
			Account     string `json:"account"`
			Avatar      string `json:"avatar"`
			City        string `json:"city"`
			Country     string `json:"country"`
			Nickname    string `json:"nickname"`
			Province    string `json:"province"`
			Remark      string `json:"remark"`
			Sex         int    `json:"sex"`
			Signature   string `json:"signature"`
			SmallAvatar string `json:"small_avatar"`
			SnsPic      string `json:"sns_pic"`
			SourceType  int    `json:"source_type"`
			Status      int    `json:"status"`
			V1          string `json:"v1"`
			V2          string `json:"v2"`
			Wxid        string `json:"wxid"`
		} `json:"data"`
		Type int `json:"type"`
	} `json:"ReturnJson"`
}

// FriendsListResp 获取好友列表响应
type FriendsListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson []struct {
		HideInputBarFlag  int
		StarFriend        int
		Sex               int
		AppAccountFlag    int
		VerifyFlag        int
		ContactFlag       int
		WebWxPluginSwitch int
		HeadImgFlag       int
		SnsFlag           int
		IsOwner           int
		MemberCount       int
		ChatRoomId        int
		UniFriend         int
		OwnerUin          int
		Statues           int
		AttrStatus        int64
		Uin               int64
		Province          string
		City              string
		Alias             string
		DisplayName       string
		KeyWord           string
		EncryChatRoomId   string
		UserName          string
		NickName          string
		HeadImgUrl        string
		RemarkName        string
		PYInitial         string
		PYQuanPin         string
		RemarkPYInitial   string
		RemarkPYQuanPin   string
		Signature         string
	} `json:"data"`
}

// GroupListResp 获取群组列表响应
type GroupListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson []struct {
		Avatar      string `json:"avatar"`
		IsManager   int    `json:"is_manager"`
		ManagerWxid string `json:"manager_wxid"`
		Nickname    string `json:"nickname"`
		TotalMember int    `json:"total_member"`
		Wxid        string `json:"wxid"`
	} `json:"ReturnJson"`
}

// GroupMemberListResp 获取群成员列表响应
type GroupMemberListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson []struct {
		HideInputBarFlag  int
		StarFriend        int
		Sex               int
		AppAccountFlag    int
		VerifyFlag        int
		ContactFlag       int
		WebWxPluginSwitch int
		HeadImgFlag       int
		SnsFlag           int
		IsOwner           int
		MemberCount       int
		ChatRoomId        int
		UniFriend         int
		OwnerUin          int
		Statues           int
		AttrStatus        int64
		Uin               int64
		Province          string
		City              string
		Alias             string
		DisplayName       string
		KeyWord           string
		EncryChatRoomId   string
		UserName          string
		NickName          string
		HeadImgUrl        string
		RemarkName        string
		PYInitial         string
		PYQuanPin         string
		RemarkPYInitial   string
		RemarkPYQuanPin   string
		Signature         string
	} `json:"data"`
}

// SubscriptionListResp 获取订阅号列表响应
type SubscriptionListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson []struct {
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
		Wxid     string `json:"wxid"`
	} `json:"ReturnJson"`
}
