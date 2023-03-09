package totp

import (
	"SrvCat/config"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"net/url"
	"time"
)

func GenerateURI(name, secret string) string {
	s := base32.StdEncoding.EncodeToString([]byte(secret))
	auth := "totp/" + name + ":"
	q := make(url.Values)
	q.Add("secret", s)
	q.Add("issuer", name)
	return "otpauth://" + auth + "SrvCat?" + q.Encode()
}

func ValidateCode(code int) (bool, error) {
	t0 := int(time.Now().Unix() / 30)
	minT := t0 - 2
	maxT := t0 + 2
	for t := minT; t <= maxT; t++ {
		hash := hmac.New(sha1.New, []byte(config.Config.Machine.Secret))
		err := binary.Write(hash, binary.BigEndian, int64(t))
		if err != nil {
			return false, err
		}
		h := hash.Sum(nil)
		offset := h[19] & 0x0f
		truncated := binary.BigEndian.Uint32(h[offset : offset+4])
		truncated &= 0x7fffffff
		if int(truncated%1000000) == code {
			return true, nil
		}
	}
	return false, nil
}
