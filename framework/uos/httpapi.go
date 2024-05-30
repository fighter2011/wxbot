package uos

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

const (
	UrlRobotInfo        string = "/api/v1/robot/info"
	UrlListFriend       string = "/api/v1/robot/friends/list"
	UrlListGroup        string = "/api/v1/robot/group/member/list"
	UrlListGroupMembers string = "/api/v1/robot/group/member/list"
	UrlListMp           string = "/api/v1/robot/mps/list"
	UrlTextSend         string = "/api/v1/robot/text/send"
	UrlImageSend        string = "/api/v1/robot/image/send"
	UrlFileSend         string = "/api/v1/robot/file/send"
	UrlMusicSend        string = "/api/v1/robot/music/send"
	UrlEmojiSend        string = "/api/v1/robot/emoji/send"
	UrlInviteUserGroup  string = "/api/v1/robot/group/invite"
	UrlAgreeUserVerify  string = "/api/v1/robot/friend/verify"
)

type FileType int

const (
	IMAGE FileType = iota
	FILE
	VIDEO
)

func (f *Framework) msgFormat(msg string) string {
	buff := bytes.NewBuffer(make([]byte, 0, len(msg)*2))
	for _, r := range msg {
		if unicode.Is(unicode.Han, r) || unicode.IsLetter(r) {
			buff.WriteString(string(r))
			continue
		}
		switch utf8.RuneLen(r) {
		case 2, 3:
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x", r) + `]`)
		case 4:
			r1, r2 := utf16.EncodeRune(r)
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x]", r1))
			buff.WriteString(`[emoji=`)
			buff.WriteString(fmt.Sprintf("%04x]", r2))
		default:
			buff.WriteString(string(r))
		}
	}
	return strings.ReplaceAll(strings.ReplaceAll(buff.String(), "\r\n", "\r"), "\n", "\r")
}

func (f *Framework) GetRobotInfo() (*robot.User, error) {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlRobotInfo)
	var resp RobotInfoResp
	if err := NewRequest().Get(apiUrl).SetSuccessResult(&resp).Do().Err; err != nil {
		log.Errorf("[UOS] GetRobotInfo error: %v", err)
		return nil, err
	}

	return &robot.User{
		WxId:         resp.ReturnJson.UserName,
		WxNum:        resp.ReturnJson.NickName,
		Nick:         resp.ReturnJson.NickName,
		Country:      "",
		Province:     resp.ReturnJson.Province,
		City:         resp.ReturnJson.City,
		AvatarMinUrl: "",
		AvatarMaxUrl: "",
	}, nil
}

func (f *Framework) GetObjectInfo(wxId string) (*robot.User, error) {
	//todo
	return nil, nil
}

func (f *Framework) GetFriends(isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlListFriend)
	var resp FriendsListResp
	if err := NewRequest().Get(apiUrl).SetSuccessResult(&resp).Do().Err; err != nil {
		log.Errorf("[UOS] GetFriends error: %v", err)
		return nil, err
	}
	var friendsInfoList []*robot.User
	for _, res := range resp.ReturnJson {
		friendsInfoList = append(friendsInfoList, &robot.User{
			WxId:         res.UserName,
			WxNum:        res.UserName,
			Nick:         res.NickName,
			Remark:       res.RemarkName,
			NickBrief:    res.RemarkPYInitial,
			NickWhole:    res.RemarkPYQuanPin,
			Sign:         res.Signature,
			Country:      res.Province,
			Province:     res.Province,
			City:         res.City,
			AvatarMinUrl: res.HeadImgUrl,
			AvatarMaxUrl: res.HeadImgUrl,
			Sex:          strconv.Itoa(res.Sex),
		})
	}

	// 过滤系统用户
	var SystemUserWxId = map[string]struct{}{"medianote": {}, "newsapp": {}, "fmessage": {}, "floatbottle": {}}
	var filteredFriendInfo []*robot.User
	for i := range friendsInfoList {
		if _, ok := SystemUserWxId[friendsInfoList[i].WxId]; !ok {
			filteredFriendInfo = append(filteredFriendInfo, friendsInfoList[i])
		}
	}
	return filteredFriendInfo, nil
}

func (f *Framework) GetGroups(isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlListGroup)
	var resp FriendsListResp
	if err := NewRequest().Get(apiUrl).SetSuccessResult(&resp).Do().Err; err != nil {
		log.Errorf("[UOS] GetGroups error: %v", err)
		return nil, err
	}
	var friendsInfoList []*robot.User
	for _, res := range resp.ReturnJson {
		friendsInfoList = append(friendsInfoList, &robot.User{
			WxId:         res.UserName,
			WxNum:        res.UserName,
			Nick:         res.NickName,
			Remark:       res.RemarkName,
			NickBrief:    res.RemarkPYInitial,
			NickWhole:    res.RemarkPYQuanPin,
			Sign:         res.Signature,
			Country:      res.Province,
			Province:     res.Province,
			City:         res.City,
			AvatarMinUrl: res.HeadImgUrl,
			AvatarMaxUrl: res.HeadImgUrl,
			Sex:          strconv.Itoa(res.Sex),
		})
	}
	return friendsInfoList, nil
}

func (f *Framework) GetGroupMembers(groupWxId string, isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlListGroupMembers)
	var resp GroupMemberListResp
	if err := NewRequest().Get(apiUrl).SetSuccessResult(&resp).Do().Err; err != nil {
		log.Errorf("[UOS] GetGroupMembers error: %v", err)
		return nil, err
	}
	if err := NewRequest().Get(apiUrl).SetSuccessResult(&resp).Do().Err; err != nil {
		log.Errorf("[UOS] GetGroupMembers error: %v", err.Error())
		return nil, err
	}
	var groupMemberInfoList []*robot.User
	for _, res := range resp.ReturnJson {
		groupMemberInfoList = append(groupMemberInfoList, &robot.User{
			WxId: res.UserName,
			Nick: res.NickName,
		})
	}
	return groupMemberInfoList, nil
}

func (f *Framework) GetMPs(isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlListMp)
	var resp FriendsListResp
	if err := NewRequest().Get(apiUrl).SetSuccessResult(&resp).Do().Err; err != nil {
		log.Errorf("[UOS] GetRobotInfo error: %v", err)
		return nil, err
	}
	var friendsInfoList []*robot.User
	for _, res := range resp.ReturnJson {
		friendsInfoList = append(friendsInfoList, &robot.User{
			WxId:         res.UserName,
			WxNum:        res.UserName,
			Nick:         res.NickName,
			Remark:       res.RemarkName,
			NickBrief:    res.RemarkPYInitial,
			NickWhole:    res.RemarkPYQuanPin,
			Sign:         res.Signature,
			Country:      res.Province,
			Province:     res.Province,
			City:         res.City,
			AvatarMinUrl: res.HeadImgUrl,
			AvatarMaxUrl: res.HeadImgUrl,
			Sex:          strconv.Itoa(res.Sex),
		})
	}
	return friendsInfoList, nil
}

func (f *Framework) GetMemePictures(msg *robot.Message) string {
	return ""
}

func (f *Framework) SendText(toWxId, text string) error {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlFileSend)
	payload := map[string]interface{}{
		"wxId":    toWxId,
		"content": text,
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[UOS] SendText error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	panic("Not Support Yet!")
}

func (f *Framework) SendImage(toWxId, path string) error {
	return f.sendFile(toWxId, []string{path}, VIDEO)
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	panic("Not Support Yet!")
}

func (f *Framework) SendFile(toWxId, path string) error {
	return f.sendFile(toWxId, []string{path}, FILE)
}

func (f *Framework) SendVideo(toWxId, path string) error {
	return f.sendFile(toWxId, []string{path}, VIDEO)
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	return f.SendText(toWxId, path)
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0014",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"name":     name,
			"author":   author,
			"app":      app,
			"jumpUrl":  jumpUrl,
			"musicUrl": musicUrl,
			"imageUrl": coverUrl,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendMusic error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	panic("Not Support Yet!")
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	return nil
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	log.Errorf("[Dean] SendMessageRecordXML not support")
	return errors.New("SendMessageRecordXML not support, please use SendMessageRecord")
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	log.Errorf("[Dean] SendFavorites not support")
	return errors.New("SendFavorites not support")
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {

	return nil
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	//todo
	return nil
}

func (f *Framework) AgreeFriendVerify(v3, v4, scene string) error {
	//todo
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	//todo
	return nil
}

func (f *Framework) sendFile(toWxId string, urls []string, fileType FileType) error {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlImageSend)
	payload := map[string]interface{}{
		"wxId":     toWxId,
		"urls":     urls,
		"fileType": fileType,
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[UOS] SendFile error: %v", err.Error())
		return err
	}
	return nil
}
