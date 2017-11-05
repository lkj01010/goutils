package math
import (
	"bytes"
//	"github.com/lkj01010/log"
)

//百度:哈希算法
//http://blog.csdn.net/djinglan/article/details/8812934
//以下2种基本是一样的,只是seed不同,100000大概重复1~2个

// SDBMHash
func StringToHash(name string) int32 {
	if len(name) > 0 {
		var hash int32

		reader := bytes.NewBufferString(name)

		for {
			b, err := reader.ReadByte()
			if err != nil {
				break
			}
			if b == 0 {
				break
			}
			// equivalent to: hash = 65599*hash + int32(b);
			hash = int32(b) + (hash << 6) + (hash << 16) - hash
//			log.Debugf("b %v h %v", b, hash)
		}
		return hash
	}
	return 0
}

// BKDRHash
func StringToHash2(name string) int32 {
	if len(name) > 0 {
		var seed, hash int32
		seed = 131
		hash = 0

		reader := bytes.NewBufferString(name)

		for {
			b, err := reader.ReadByte()
			if err != nil {
				break
			}
			if b == 0 {
				break
			}
			hash = hash * seed + int32(b)
//			log.Debugf("b %v h %v", b, hash)
		}
		//如果不需要int32可以不进行求余
		return (hash & 0x7FFFFFFF)
	}
	return 0
}

