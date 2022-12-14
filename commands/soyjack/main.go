package soyjack

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"vkbot/core"

	"github.com/SevereCloud/vksdk/v2/events"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const soyboy_file_path = "commands/soyjack/soyboy.png"
const gurba_file_path = "commands/soyjack/gurba.png"
const arkady_file_path = "commands/soyjack/arkady.png"
const zakhar_file_path = "commands/soyjack/zakhar.png"
const yuri_file_path = "commands/soyjack/yuri.png"
const nikita_file_path = "commands/soyjack/nikita.png"

type mode_data struct {
	Name string
	Mask []float64

	Wand *imagick.MagickWand

	Width  uint
	Height uint

	PosX int
	PosY int
}
type soyjack_pool []mode_data

var soyjacks soyjack_pool
var mwsoy, mwgurba, mwarkady, mwzakhar, mwyuri, mwnikita *imagick.MagickWand

func (s *soyjack_pool) Probe(name string) (w mode_data, err error) {
	for _, v := range *s {
		if v.Name == name {
			w = v

			return
		}
	}

	err = errors.New("ошибка: шаблон не найден")
	return
}

func (s *soyjack_pool) Help(obj *events.MessageNewObject) {
	names := []string{}
	for _, x := range *s {
		names = append(names, x.Name)
	}

	core.ReplySimple(obj, "доступные шаблоны: "+strings.Join(names, ", "))
}

func Register() core.Command {
	mwsoy = imagick.NewMagickWand()
	mwsoy.ReadImage(soyboy_file_path)

	mwgurba = imagick.NewMagickWand()
	mwgurba.ReadImage(gurba_file_path)

	mwarkady = imagick.NewMagickWand()
	mwarkady.ReadImage(arkady_file_path)

	mwzakhar = imagick.NewMagickWand()
	mwzakhar.ReadImage(zakhar_file_path)

	mwyuri = imagick.NewMagickWand()
	mwyuri.ReadImage(yuri_file_path)

	mwnikita = imagick.NewMagickWand()
	mwnikita.ReadImage(nikita_file_path)

	soyjacks = soyjack_pool{
		{
			Name: "сойбой",
			Mask: []float64{
				0, 0, 1, 20,
				1, 443, 33, 443,
				275, 443, 273, 423,
				275, 1, 238, 2,
			},
			Wand:   mwsoy,
			Width:  275,
			Height: 443,
			PosX:   25,
			PosY:   66,
		},
		{
			Name: "нс",
			Mask: []float64{
				0, 0, 3, 1,
				0, 414, 1, 412,
				595, 414, 595, 362,
				595, 0, 566, 6,
			},
			Wand:   mwgurba,
			Width:  595,
			Height: 414,
			PosX:   8,
			PosY:   1249,
		},
		{
			Name: "ас",
			Mask: []float64{
				0, 0, 41, 2,
				0, 436, 1, 356,
				607, 436, 607, 436,
				607, 0, 600, 93,
			},
			Wand:   mwarkady,
			Width:  607,
			Height: 436,
			PosX:   344,
			PosY:   1248,
		},
		{
			Name: "зс",
			Mask: []float64{
				0, 0, 21, 0,
				0, 359, 1, 310,
				584, 359, 570, 359,
				584, 0, 583, 21,
			},
			Wand:   mwzakhar,
			Width:  584,
			Height: 359,
			PosX:   1440,
			PosY:   606,
		},
		{
			Name: "юс",
			Mask: []float64{
				0, 0, 33, 2,
				0, 372, 2, 371,
				612, 372, 611, 363,
				612, 0, 606, 38,
			},
			Wand:   mwyuri,
			Width:  612,
			Height: 372,
			PosX:   430,
			PosY:   859,
		},
		// {
		// 	Name: "ннс",
		// 	Mask: []float64{
		// 		0, 0, 2, 307,
		// 		0, 695, 143, 685,
		// 		792, 695, 791, 552,
		// 		792, 0, 534, 0,
		// 	},
		// 	Wand:   mwnikita,
		// 	Width:  792,
		// 	Height: 695,
		// 	PosX:   665,
		// 	PosY:   85,
		// },
	}

	return core.Command{
		Aliases:     []string{"сой", "сойджек"},
		Description: "обмазать картинку соей",
		Handler:     handle,
	}
}

func handle(obj *events.MessageNewObject) (err error) {
	// atts := core.ExtractAttachments(obj, "photo,doc")
	atts := core.ExtractAttachments(obj, "photo")

	if len(atts) == 0 {
		core.ReplySimple(obj, core.ERR_NO_PICTURE)

		return
	}

	attachment := atts[0]

	link := ""

	switch attachment.Type {
	case "photo":
		link = attachment.Photo.MaxSize().URL
	case "doc":
		link = attachment.Doc.URL

		if attachment.Doc.Size > 30*1024*1024 {
			core.ReplySimple(obj, core.ERR_LARGE_GIF)
		}
	}

	response, err := http.Get(link)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	bt, err := io.ReadAll(response.Body)
	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	mw1 := imagick.NewMagickWand()
	mw1.ReadImageBlob(bt)

	data := mode_data{}

	if mw1.GetImageWidth() < mw1.GetImageHeight() {
		data, _ = soyjacks.Probe("сойбой")
	} else {
		data, _ = soyjacks.Probe("нс")
	}

	name := strings.Join(core.ExtractArguments(obj), " ")
	if name != "" {
		data, err = soyjacks.Probe(name)
		if err != nil {
			soyjacks.Help(obj)
			err = nil

			return
		}
	}

	mw1.ResizeImage(data.Width, data.Height, imagick.FILTER_UNDEFINED, 1)
	mw1.SetImageVirtualPixelMethod(imagick.VIRTUAL_PIXEL_TRANSPARENT)
	mw1.DistortImage(imagick.DISTORTION_PERSPECTIVE, data.Mask, false)

	mw2 := data.Wand.Clone()
	mw2.CompositeLayers(mw1, imagick.COMPOSITE_OP_DST_OVER, data.PosX, data.PosY)

	vkPhoto, err := core.GetStorage().Vk.UploadMessagesPhoto(0, bytes.NewReader(mw2.GetImageBlob()))

	mw1.Destroy()
	mw2.Destroy()

	if err != nil {
		core.ReplySimple(obj, core.ERR_UNKNOWN)

		return
	}

	core.ReplySimple(obj, "ваша картинка:", vkPhoto)

	return
}
