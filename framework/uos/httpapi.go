package uos

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/yqchilde/wxbot/framework/dean"
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
	UrlListGroup        string = "/api/v1/robot/group/list"
	UrlListGroupMembers string = "/api/v1/robot/group/member/list"
	UrlListMp           string = "/api/v1/robot/mps/list"
	UrlTextSend         string = "/api/v1/robot/text/send"
	UrlFileSend         string = "/api/v1/robot/file/send"
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
	var emoji dean.EmojiXml
	if err := xml.Unmarshal([]byte(msg.Content), &emoji); err != nil {
		return ""
	}
	return emoji.Emoji.Cdnurl
}

func (f *Framework) SendText(toWxId, text string) error {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlTextSend)
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
	//拼接回复用户名称
	content := fmt.Sprintf("回复:[%s], %s", toWxName, text)
	return f.SendText(toGroupWxId, content)
}

func (f *Framework) SendImage(toWxId, path string) error {
	return f.sendFile(toWxId, path, IMAGE)
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	panic("Not Support Yet!")
}

func (f *Framework) SendFile(toWxId, path string) error {
	return f.sendFile(toWxId, path, FILE)
}

func (f *Framework) SendVideo(toWxId, path string) error {
	return f.sendFile(toWxId, path, VIDEO)
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	return f.SendText(toWxId, path)
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return errors.New("Not Support Yet")
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	panic("Not Support Yet!")
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	return nil
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	return errors.New("Not Support Yet")
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	return errors.New("Not Support Yet")
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {

	return nil
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	//todo
	return nil
}

func (f *Framework) AgreeFriendVerify(message *robot.FriendVerifyMessage) error {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlAgreeUserVerify)
	payload := map[string]interface{}{
		"verifyContents": []string{"同意"},
		"recommendInfo":  message.RecommendInfo,
	}
	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[UOS] GetRobotInfo error: %v", err)
		return err
	}
	return nil
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	//todo
	return nil
}

func (f *Framework) sendFile(toWxId string, url string, fileType FileType) error {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, UrlFileSend)
	payload := map[string]interface{}{
		"wxId":     toWxId,
		"url":      url,
		"fileType": fileType,
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[UOS] SendFile error: %v", err.Error())
		return err
	}
	return nil
}
