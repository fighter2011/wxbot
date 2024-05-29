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
	URL_ROBOT_INFO  string = "/api/v1/robot/info"
	URL_LIST_FRIEND string = "/api/v1/robot/friends/list"
	URL_LIST_GRUOP  string = "/api/v1/robot/group/member/list"
	URL_LIST_MP     string = "/api/v1/robot/mps/list"
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
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, URL_ROBOT_INFO)
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
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, URL_LIST_FRIEND)
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
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, URL_LIST_GRUOP)
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

func (f *Framework) GetGroupMembers(groupWxId string, isRefresh bool) ([]*robot.User, error) {
	//todo
	return nil, nil
}

func (f *Framework) GetMPs(isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s%s", f.ApiUrl, URL_LIST_MP)
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
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0001",
		"data": map[string]interface{}{
			"wxid": toGroupWxId,
			"msg":  fmt.Sprintf("[@,wxid=%s,nick=%s,isAuto=true] %s", toWxId, toWxName, f.msgFormat(text)),
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendTextAndAt error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0010",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendImage error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0012",
		"data": map[string]interface{}{
			"wxid":    toWxId,
			"title":   title,
			"content": desc,
			"jumpUrl": jumpUrl,
			"path":    imageUrl,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendShareLink error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendFile(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0011",
		"data": map[string]interface{}{
			"wxid": toWxId,
			"path": path,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendFile error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendVideo(toWxId, path string) error {
	log.Errorf("[Dean] SendVideo not support")
	return errors.New("SendVideo not support")
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	log.Errorf("[Dean] SendEmoji not support")
	return errors.New("SendEmoji not support")
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
	apiUrl := fmt.Sprintf("%s/DaenWxHook/client/", f.ApiUrl)
	payload := map[string]interface{}{
		"type": "Q0013",
		"data": map[string]interface{}{
			"wxid":     toWxId,
			"title":    title,
			"content":  content,
			"jumpPath": jumpPath,
			"gh":       ghId,
			"path":     imagePath,
		},
	}

	if err := NewRequest().Post(apiUrl).SetBody(payload).Do().Err; err != nil {
		log.Errorf("[Dean] SendMiniProgram error: %v", err.Error())
		return err
	}
	return nil
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
