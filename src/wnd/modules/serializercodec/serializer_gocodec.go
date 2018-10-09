package serializercodec

import (
	"wnd/modules"
	"wnd/utils/log"
	
	"encoding/binary"
	"errors"
	"reflect"
	"fmt"
	
	"github.com/ugorji/go/codec"
)

var (
	handle = new(codec.BincHandle)
)

type gocodecSerializer struct {
}

func New() modules.Serializer {
	
	return &gocodecSerializer{ }
}

func (this *gocodecSerializer) Priority() int8 {
	return -128
}

func (this *gocodecSerializer) ID() string {
	return "gocodecSerializer"
}

func (this *gocodecSerializer) Init() error {
	return nil
}

func (this *gocodecSerializer) Destroy() {}



func (this *gocodecSerializer) Serialize(i interface{}) (total []byte, e error) {
	log.Tracef("%v", reflect.TypeOf(i).String())
	
	name,idto := toDto(i)
	
	bname := []byte(name)
	lenname := len(bname)
	
	var b []byte
	e = codec.NewEncoderBytes(&b, handle).Encode(idto)
	
	total = make([]byte, lenname + len(b) + 8)
	
	binary.LittleEndian.PutUint32(total[:4], uint32(lenname))
	copy(total[4:], bname)
	binary.LittleEndian.PutUint32(total[4+lenname:8+lenname], uint32(len(b)))
	copy(total[8+lenname:], b)
	
	log.Debugf("serialized to %v bytes", len(total))
	
	return
}

func (this *gocodecSerializer) Deserialize(b []byte) (i interface{}, e error) {
	log.Tracef("%v bytes", len(b))
	
	lenname := binary.LittleEndian.Uint32(b[:4])
	name := string(b[4:4+lenname])
	
	i = fromName(name)
	
	lenb := binary.LittleEndian.Uint32(b[4+lenname:8+lenname])
	total := 8+lenname+lenb
	
	if len(b) < int(total) {
		return nil, errors.New(fmt.Sprintf("broken data (want %d, got %d)", total, len(b)))
	}
	
	e = codec.NewDecoderBytes(b[8+lenname:total], handle).Decode(i)
	
	if e == nil {
		log.Debugf("decoded to %v", reflect.TypeOf(i).String())
	}
	
	return
}