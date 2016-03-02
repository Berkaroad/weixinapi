package weixinapi

import (
	"math/rand"
	"strconv"
	"time"
)

// <![CDATA[Text]]>
type CDATAText struct {
	Text string `xml:",innerxml"`
}

func NewCDATAText(text string) CDATAText {
	return CDATAText{Text: "<![CDATA[" + text + "]]>"}
}

// 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	code_serial := "0123456789abcdefghijklmnpqrstUVwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"
	random_str := ""
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		random_str += string(code_serial[rand.Intn(len(code_serial))])
	}
	return random_str
}

// 生成时间戳
func GenerateTimestamp() int {
	return int(time.Now().Unix())
}

// 生成时间戳
func GenerateTimestampString() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func HostToNetworkOrderInt32(inval int32) int32 {
	var outval int32
	for i := 0; i < 4; i++ {
		outval = (outval << 8) + ((inval >> (uint(i) * 8)) & 255)
	}
	return outval
}

func NetworkToHostOrderInt32(inval int32) int32 {
	outval := (int32(NetworkToHostOrderInt16(int16(inval)))&65535)<<16 | (int32(NetworkToHostOrderInt16(int16(inval>>16))) & 65535)
	return outval

}

func NetworkToHostOrderInt16(inval int16) int16 {
	outval := int16(int32(inval)&255<<8 | (int32(inval) >> 8 & 255))
	return outval
}

func ByteArrayToInt32(value []byte, position int) int32 {
	var result int32
	result = 0
	for i, val := range value[position:] {
		result += int32(val) << uint(8*i)
	}
	return result
}

func Int32ToByteArray(value int32) []byte {
	array := make([]byte, 4)
	array[3] = byte(value >> 24)
	array[2] = byte((value - (int32(array[3]) << 24)) >> 16)
	array[1] = byte((value - (int32(array[3]) << 24) - (int32(array[2]) << 16)) >> 8)
	array[0] = byte(value - (int32(array[3]) << 24) - (int32(array[2]) << 16) - (int32(array[1]) << 8))
	return array
}

func KCS7Encoder(text_length int) []byte {
	block_size := 32
	// 计算需要填充的位数
	amount_to_pad := block_size - (text_length % block_size)
	if amount_to_pad == 0 {
		amount_to_pad = block_size
	}
	// 获得补位所用的字符
	pad_chr := string(rune(byte(amount_to_pad & 0xff)))

	tmp := ""
	for index := 0; index < amount_to_pad; index++ {
		tmp += pad_chr
	}
	return []byte(tmp)
}
