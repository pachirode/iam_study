package idutil

import (
	"crypto/rand"

	"github.com/sony/sonyflake"
	hashIds "github.com/speps/go-hashids"

	"github.com/pachirode/iam_study/pkg/utils/iputil"
	"github.com/pachirode/iam_study/pkg/utils/stringutil"
)

const (
	Alphabet62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.MachineID = func() (uint16, error) {
		ip := iputil.GetLocalIP()

		return uint16([]byte(ip)[2])<<8 + uint16([]byte(ip)[3]), nil
	}

	sf = sonyflake.NewSonyflake(st)
}

func GetIntID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	return id
}

func GetInstanceID(uid uint64, prefix string) string {
	hd := hashIds.NewData()
	hd.Alphabet = Alphabet36
	hd.MinLength = 6
	hd.Salt = "xml12x"

	h, err := hashIds.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	id, err := h.Encode([]int{int(uid)})
	if err != nil {
		panic(err)
	}

	return prefix + stringutil.Reverse(id)
}

func GetUUID36(prefix string) string {
	id := GetIntID()
	hd := hashIds.NewData()
	hd.Alphabet = Alphabet36

	h, err := hashIds.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	uuid, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + stringutil.Reverse(uuid)
}

func randString(letters string, n int) string {
	output := make([]byte, n)
	randomNess := make([]byte, n)

	_, err := rand.Read(randomNess)
	if err != nil {
		panic(err)
	}

	length := len(letters)

	for pos := range output {
		random := randomNess[pos]
		randomPos := random % uint8(length)
		output[pos] = letters[randomPos]
	}

	return string(output)
}

func NewSecretID() string {
	return randString(Alphabet62, 36)
}

func NewSecretKey() string {
	return randString(Alphabet62, 32)
}
