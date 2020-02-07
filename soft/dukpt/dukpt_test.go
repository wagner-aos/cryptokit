package dukpt

import (
	"crypto/cipher"
	"crypto/des"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBdk = []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10}
var testKsn = []byte{0xFF, 0xFF, 0x98, 0x76, 0x54, 0x32, 0x10, 0xE0, 0x00, 0x08}
var testIpek = []byte{0x6a, 0xc2, 0x92, 0xfa, 0xa1, 0x31, 0x5b, 0x4d, 0x85, 0x8a, 0xb3, 0xa3, 0xd7, 0xd5, 0x93, 0x3a}
var testPek = []byte{0x27, 0xf6, 0x6d, 0x52, 0x44, 0xff, 0x62, 0x1e, 0xaa, 0x6f, 0x61, 0x20, 0xed, 0xeb, 0x42, 0x80}
var testCiphertext = []byte{0xC2, 0x5C, 0x1D, 0x11, 0x97, 0xD3, 0x1C, 0xAA, 0x87, 0x28, 0x5D, 0x59, 0xA8, 0x92, 0x04, 0x74, 0x26, 0xD9, 0x18, 0x2E, 0xC1, 0x13, 0x53, 0xC0, 0x51, 0xAD, 0xD6, 0xD0, 0xF0, 0x72, 0xA6, 0xCB, 0x34, 0x36, 0x56, 0x0B, 0x30, 0x71, 0xFC, 0x1F, 0xD1, 0x1D, 0x9F, 0x7E, 0x74, 0x88, 0x67, 0x42, 0xD9, 0xBE, 0xE0, 0xCF, 0xD1, 0xEA, 0x10, 0x64, 0xC2, 0x13, 0xBB, 0x55, 0x27, 0x8B, 0x2F, 0x12}

func TestLogs(t *testing.T) {

	t.Logf("testBdk = %x", testBdk)
	t.Logf("testKsn = %x", testKsn)
	t.Logf("testIpek = %x", testIpek)
	t.Logf("testPek = %x", testPek)
	t.Logf("testCiphertext = %x", testCiphertext)
	assert.Equal(t, true, false)
}

func TestEncodeDecodeKsn(t *testing.T) {
	encodedKsn := make([]byte, 10)

	ksn := DecodeKsn(testKsn)

	assert.Equal(t, []byte{0xFF, 0xFF, 0x98, 0x76, 0x54}, ksn.Ksi)
	assert.Equal(t, []byte{0x32, 0x10, 0xe0}, ksn.Trsm)
	assert.Equal(t, 8, ksn.Counter)

	t.Log("KSN DECODED:")
	t.Logf("ksi = %x", ksn.Ksi)
	t.Logf("Trsm = %x", ksn.Trsm)
	t.Logf("Count = %x", ksn.Counter)

	EncodeKsn(encodedKsn, ksn)

	t.Log("KSN ENCODED:")
	t.Logf("ksi = %x", encodedKsn)

	assert.Equal(t, testKsn, "encodedKsn")
}

func TestKeyDerivation(t *testing.T) {
	ipek, err := DeriveIpekFromBdk(testBdk, testKsn)

	assert.Nil(t, err)
	assert.NotNil(t, ipek)
	assert.Equal(t, testIpek, ipek, "Derived IPEK should be correct")

	pek, err := DerivePekFromIpek(ipek, testKsn)

	t.Logf("PEK %x", pek)

	assert.Nil(t, err)
	assert.NotNil(t, pek)
	assert.Equal(t, testPek, pek, "Derived PEK should be correct")
}

func TestDecryption(t *testing.T) {
	pek, err := DerivePekFromBdk(testBdk, testKsn)

	assert.Nil(t, err)
	assert.NotNil(t, pek)
	assert.Equal(t, testPek, pek, "Derived PEK should be correct")

	tdes, _ := des.NewTripleDESCipher(buildTdesKey(pek))
	cbc := cipher.NewCBCDecrypter(tdes, make([]byte, 8))

	result := make([]byte, len(testCiphertext))

	cbc.CryptBlocks(result, testCiphertext)
	assert.Equal(t, "%B5452300551227189^HOGAN/PAUL      ^08043210000000725000000?\x00\x00\x00\x00", string(result))
}

func TestECBDecryption(t *testing.T) {

	plainText := []byte("WagnerAOS")
	hexText := []byte("5761676E6572414F5300000000000000")
	key := []byte("BBF47BAC9F09FF312BE2DB35DC73C58A")

	encryptedText := EncryptDecryptAes128Ecb(hexText, key, true)
	decryptedText := EncryptDecryptAes128Ecb(encryptedText, key, false)

	t.Logf("plainText: %s", plainText)
	t.Logf("hexText: %x", hexText)
	t.Logf("decriptedText: %s", decryptedText)
	t.Logf("encriptedText: %x", encryptedText)

	assert.Equal(t, decryptedText, hexText)
}
