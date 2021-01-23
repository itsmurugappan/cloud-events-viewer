package codec

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	nethttp "net/http"
	"strconv"
	"sync"

	"github.com/linkedin/goavro/v2"
)

type AvroDecoder struct {
	schemaURL string
	codecMap  map[int]*goavro.Codec
	mutex     sync.RWMutex
}

func NewAvroDecoder(schemaURL string) *AvroDecoder {
	if schemaURL != "" {
		return &AvroDecoder{
			schemaURL: schemaURL,
			codecMap:  make(map[int]*goavro.Codec),
		}
	}
	return nil
}

func (decoder *AvroDecoder) Decode(msg []byte) ([]byte, error) {
	if len(msg) < 5 {
		log.Println("Invalid message to decode")
		return msg, nil
	}
	schemaID := binary.BigEndian.Uint32(msg[1:5])
	if schemaID == 0 {
		log.Printf("Error getting the schema id %d\n", int(schemaID))
		return msg, nil
	}
	codec, err := decoder.GetCodec(int(schemaID))
	if err != nil {
		log.Printf("Error getting the schema for schemaid %d, %v\n", int(schemaID), err)
		return msg, nil
	}
	native, _, err := codec.NativeFromBinary(msg[5:])
	if err != nil {
		log.Printf("Error decoding message %v\n", err)
		return nil, err
	}
	value, err := codec.TextualFromNative(nil, native)
	if err != nil {
		log.Printf("Error decoding message %v\n", err)
		return nil, err
	}
	return value, nil
}

func (decoder *AvroDecoder) GetCodec(schemaID int) (*goavro.Codec, error) {
	decoder.mutex.RLock()
	if decoder.codecMap[schemaID] != nil {
		return decoder.codecMap[schemaID], nil
	}
	decoder.mutex.RUnlock()

	resp, err := nethttp.Get(fmt.Sprintf("%s/schemas/ids/%d", decoder.schemaURL, schemaID))
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var schemaResp map[string]interface{}
	if err := json.Unmarshal(body, &schemaResp); err != nil {
		return nil, err
	}

	schema, err := json.Marshal(schemaResp["schema"])
	if err != nil {
		return nil, err
	}

	schemaString, err := strconv.Unquote(string(schema))
	if err != nil {
		return nil, err
	}

	codec, err := goavro.NewCodec(schemaString)
	if err != nil {
		return nil, err
	}

	decoder.mutex.Lock()
	decoder.codecMap[schemaID] = codec
	decoder.mutex.Unlock()
	return codec, nil
}
