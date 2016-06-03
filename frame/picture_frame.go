package frame

import (
	"bytes"
	"errors"
	"github.com/bogem/id3v2/util"
	"io"
	"os"
)

type PictureFramer interface {
	Framer

	Description() string
	SetDescription(string)

	MimeType() string
	SetMimeType(string)

	Picture() io.Reader
	SetPicture(io.Reader)

	PictureType() byte
	SetPictureType(byte)
}

type PictureFrame struct {
	description string
	mimeType    string
	picture     io.Reader
	pictureType byte
}

func (pf PictureFrame) Bytes() ([]byte, error) {
	b := bytesBufPool.Get().(*bytes.Buffer)
	b.Reset()
	defer bytesBufPool.Put(b)

	b.WriteByte(util.NativeEncoding)
	b.WriteString(pf.mimeType)
	b.WriteByte(0)
	if pf.pictureType < 0 || pf.pictureType > 20 {
		return nil, errors.New("Incorrect picture type in picture frame with description " + pf.description)
	}
	b.WriteByte(pf.pictureType)
	b.WriteString(pf.description)
	b.WriteByte(0)

	b.ReadFrom(pf.picture)
	if v, ok := pf.picture.(*os.File); ok {
		v.Close()
	}

	return b.Bytes(), nil
}

func (pf PictureFrame) Description() string {
	return pf.description
}

func (pf *PictureFrame) SetDescription(desc string) {
	pf.description = desc
}

func (pf PictureFrame) MimeType() string {
	return pf.mimeType
}

func (pf *PictureFrame) SetMimeType(mt string) {
	pf.mimeType = mt
}

func (pf PictureFrame) Picture() io.Reader {
	return pf.picture
}

func (pf *PictureFrame) SetPicture(rd io.Reader) {
	pf.picture = rd
}

func (pf *PictureFrame) SetPictureFromFile(name string) error {
	if pf.picture != nil {
		if v, ok := pf.picture.(*os.File); ok {
			v.Close()
		}
	}
	pictureFile, err := os.Open(name)
	pf.picture = pictureFile
	return err
}

func (pf PictureFrame) PictureType() byte {
	return pf.pictureType
}

func (pf *PictureFrame) SetPictureType(pt byte) {
	pf.pictureType = pt
}
