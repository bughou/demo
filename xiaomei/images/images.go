package images

import (
	"github.com/fatih/color"
	"github.com/lovego/xiaomei/utils"
	"github.com/lovego/xiaomei/utils/cmd"
	"github.com/lovego/xiaomei/xiaomei/images/access"
	"github.com/lovego/xiaomei/xiaomei/images/app"
	"github.com/lovego/xiaomei/xiaomei/images/web"
	"github.com/lovego/xiaomei/xiaomei/release"
)

var imagesMap = map[string]Image{
	`app`:    Image{`app`, app.Image{}},
	`web`:    Image{`web`, web.Image{}},
	`access`: Image{`access`, access.Image{}},
}

func Has(svcName string) bool {
	_, ok := imagesMap[svcName]
	return ok
}

func Run(svcName string, ports []string) error {
	image, ok := imagesMap[svcName]
	if !ok {
		return nil
	}
	return image.Run(ports)
}

func Build(svcName string, pull bool) error {
	if svcName == `` {
		return release.EachServiceDo(func(svcName string) error {
			return Build(svcName, pull)
		})
	}
	image, ok := imagesMap[svcName]
	if !ok {
		return nil
	}
	return image.Build(pull)
}

func Push(svcName string) error {
	if svcName == `` {
		return release.EachServiceDo(Push)
	}
	if _, ok := imagesMap[svcName]; !ok {
		return nil
	}
	utils.Log(color.GreenString(`pushing ` + svcName + ` image.`))
	_, err := cmd.Run(cmd.O{}, `docker`, `push`, release.ImageNameOf(svcName))
	return err
}
