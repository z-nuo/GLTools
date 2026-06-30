package glcrypto

import "testing"

func TestMD5Hex(t *testing.T) {
	if got := MD5Hex("abc"); got != "900150983cd24fb0d6963f7d28e17f72" {
		t.Fatalf("MD5Hex() = %s", got)
	}
}

func TestSHA256Hex(t *testing.T) {
	if got := SHA256Hex("abc"); got != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatalf("SHA256Hex() = %s", got)
	}
}

func TestHMACSHA256Hex(t *testing.T) {
	got := HMACSHA256Hex("message", "secret")
	want := "8b5f48702995c1598c573db1e21866a9b825d4a794d169d7060a03605796360b"
	if got != want {
		t.Fatalf("HMACSHA256Hex() = %s, want %s", got, want)
	}
}

func TestBase64EncodeAndDecode(t *testing.T) {
	encoded := Base64Encode("hello")
	if encoded != "aGVsbG8=" {
		t.Fatalf("Base64Encode() = %s", encoded)
	}

	decoded, err := Base64Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != "hello" {
		t.Fatalf("Base64Decode() = %s, want hello", decoded)
	}
}

func TestBase64DecodeRejectsInvalidInput(t *testing.T) {
	if _, err := Base64Decode("not base64!"); err == nil {
		t.Fatal("Base64Decode() error = nil, want error")
	}
}
