package callbacks

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(src string) string {
	h := md5.New()
	_, _ = io.WriteString(h, src)
	return hex.EncodeToString(h.Sum(nil))
}
