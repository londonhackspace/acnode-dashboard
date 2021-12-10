package usagelogs

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type PictureTaker interface {
	TakePicture(camid int) (string, error)
}

type ZMPictureTaker struct {
	zmurl      string
	imagestore string
}

func CreateZMPictureTaker(zmurl string, imageStore string) PictureTaker {
	return &ZMPictureTaker{
		zmurl:      zmurl,
		imagestore: imageStore,
	}
}

func (zmpt *ZMPictureTaker) TakePicture(camid int) (string, error) {
	bakedUrl := strings.Replace(zmpt.zmurl, "{ID}", strconv.Itoa(camid), 1)

	resp, err := http.Get(bakedUrl)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Non-success code: " + strconv.Itoa(resp.StatusCode))
	}

	imgKey := strconv.Itoa(camid) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"

	f, err := os.Create(zmpt.imagestore + "/" + imgKey)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	return imgKey, nil
}
