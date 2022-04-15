package mmr_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/alleniversons/mmr"
	"github.com/alleniversons/mmr/db"
	"github.com/minio/blake2b-simd"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
	"testing"
)

type testElement struct {
	Id    uint64
	Text  string
	State byte
}

func TestMerkleObjects(t *testing.T) {
	m := mmr.Merkle(blake2b.New256)

	i := 0
	for i < 16 {
		te := &testElement{
			Id: uint64(i), Text: fmt.Sprintf("test %d", i),
			State: 0,
		}
		_, err := m.Add(te)
		assert.NoError(t, err)
		i++
	}

	te3 := &testElement{}
	fmt.Println(hex.EncodeToString(m.Root().Hash()))
	err := m.Get(3, te3)
	fmt.Println(te3)
	te4 := &testElement{
		Id: uint64(4), Text: fmt.Sprintf("test %d", 4),
		State: 0,
	}
	_,err = m.Add(te4)
	fmt.Println("增加后")
	fmt.Println(hex.EncodeToString(m.Root().Hash()))
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, "test 3", te3.Text)
	assert.Equal(t, uint64(3), te3.Id)

	d := db.Memory()
	if err := m.Save(d); err != nil {
		assert.NoError(t, err)
	}

	fmt.Printf("Root: %x\n", m.Root().Hash())
	m2, err := mmr.MerkleFromSource(blake2b.New256, m.Root().Hash(), d)
	if err != nil {
		t.Error(err)
		return
	}

	te5 := &testElement{}
	err = m2.Get(4, te5)
	assert.NoError(t, err)
	assert.Equal(t, "test 4", te4.Text)
	assert.Equal(t, uint64(4), te4.Id)
}

func BenchmarkMerkleSha3(b *testing.B) {
	m := mmr.Merkle(sha3.New256)
	for n := 0; n < b.N; n++ {
		te := &testElement{
			Id: uint64(n), Text: fmt.Sprintf("test %d", n),
			State: 0,
		}
		m.Add(te)
	}
}

func BenchmarkMerkleSha256(b *testing.B) {
	m := mmr.Merkle(sha256.New)
	for n := 0; n < b.N; n++ {
		te := &testElement{
			Id: uint64(n), Text: fmt.Sprintf("test %d", n),
			State: 0,
		}
		m.Add(te)
	}
}

func BenchmarkMerkleBlake(b *testing.B) {
	m := mmr.Merkle(blake2b.New256)

	for n := 0; n < b.N; n++ {
		te := &testElement{
			Id: uint64(n), Text: fmt.Sprintf("test %d", n),
			State: 0,
		}
		m.Add(te)
	}
}
