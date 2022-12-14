package tts

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
)

func Register() core.Command {
	return core.Command{
		Aliases:     []string{"ттс"},
		Description: "озвучить текст",
		Handler:     handle,
	}
}

func getgTTSReader(text, lang string) (b io.Reader, err error) {
	q := url.Values{}
	q.Set("ie", "UTF-8")
	q.Set("total", "1")
	q.Set("idx", "0")
	q.Set("client", "tw-ob")
	q.Set("tl", lang)
	q.Set("ttsspeed", "1")
	q.Set("q", text)
	q.Set("textlen", strconv.Itoa(len(text)))

	u := &url.URL{
		Scheme:   "https",
		Host:     "translate.google.com",
		Path:     "translate_tts",
		RawQuery: q.Encode(),
	}

	response, err := http.Get(u.String())
	if err != nil {
		return
	}

	b = response.Body

	return
}

func handle(obj *events.MessageNewObject) (err error) {
	txt := strings.Join(core.ExtractArguments(obj), " ")

	if txt == "" {
		core.ReplySimple(obj, "ошибка: не указан текст")

		return
	}

	if len(txt) > 200 {
		core.ReplySimple(obj, "ошибка: максимальное количество символов не должно превышать 200")

		return
	}

	data, err := getgTTSReader(txt, "ru")

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	m, err := core.GetStorage().Vk.UploadMessagesDoc(obj.Message.PeerID, "audio_message", "gs.wav", "", data)

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша озвучка:", m.AudioMessage)

	return
}
