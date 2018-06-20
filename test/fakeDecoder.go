package test

import (
	"github.com/nmarsollier/authgo/tools/db"
)

// FakeDecoder permite mockear mongo.DocumentResponse
func FakeDecoder(decoder func(v interface{}) error) db.Decoder {
	return &fakeDecoder{
		decoder: decoder,
	}
}

type fakeDecoder struct {
	err     error
	decoder func(v interface{}) error
}

func (fd *fakeDecoder) Decode(v interface{}) error {
	if fd.decoder != nil {
		return fd.decoder(v)
	}
	return fd.err
}
