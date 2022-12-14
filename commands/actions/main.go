package actions

import (
	"fmt"
	"strconv"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"обнять", "поцеловать", "опустить"},
		Description: "различные действия в отношении другого пользователя (RP команды)",
		NoPrefix:    true,
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	s := core.GetStorage()

	if obj.Message.PeerID == obj.Message.FromID {
		return
	}

	enabled, _ := s.Db.Get(s.Ctx, fmt.Sprintf("rp.%d.enabled", obj.Message.PeerID)).Result()
	if enabled == "false" {
		return
	}

	id := core.GetMention(obj)
	if id <= 0 {
		return
	}

	b := params.NewUsersGetBuilder()

	b.UserIDs([]string{strconv.Itoa(obj.Message.FromID)})
	b.Fields([]string{"sex"})

	res, err := s.Vk.UsersGet(b.Params)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	postfix := ""
	if res[0].Sex == 1 {
		postfix = "а"
	}

	action := ""
	switch strings.ToLower(strings.Split(obj.Message.Text, " ")[0]) {
	case "обнять":
		action = "обнял" + postfix + " ❤"
	case "поцеловать":
		action = "поцеловал" + postfix + " 😘"
	case "опустить":
		action = "опустил" + postfix + " 😝"
	}

	b.Fields([]string{})

	core.SendSimple(obj, "* [id"+
		strconv.Itoa(obj.Message.FromID)+
		"|"+
		core.GetNicknameOrFullName(obj.Message.FromID)+
		"] "+
		action+
		" [id"+
		strconv.Itoa(id)+
		"|"+
		core.GetNicknameOrFullName(id)+
		"]")

	return
}
