package checksum

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/TrueCloudLab/frostfs-api-go/v2/refs"
)

func ExampleCalculate() {
	payload := []byte{0, 1, 2, 3, 4, 5, 6}
	var cs Checksum

	Calculate(&cs, SHA256, payload)
	Calculate(&cs, TZ, payload)
}

func ExampleChecksum_WriteToV2() {
	var (
		csRaw [sha256.Size]byte
		csV2  refs.Checksum
		cs    Checksum
	)

	rand.Read(csRaw[:])
	cs.SetSHA256(csRaw)

	cs.WriteToV2(&csV2)

	fmt.Println(bytes.Equal(cs.Value(), csV2.GetSum()))
	// Output: true
}
