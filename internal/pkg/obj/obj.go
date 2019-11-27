package obj

import "encoding/json"

func Copy(src, dstBase interface{}) (dst interface{}, err error) {
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bs, err := json.Marshal(src)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, dstBase)
	if err != nil {
		return
	}
	dst = dstBase
	return
}
