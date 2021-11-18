package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
	"io"
	"log"
)

const (
	StandardScryptN = 1 << 18
	StandardScryptP = 1
	keyHeaderKDF    = "scrypt"
)

type cryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams cipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type cipherparamsJSON struct {
	IV string `json:"iv"`
}

func main() {
	wallet := solana.NewWallet()
	log.Printf("Private Key: %s", wallet.PrivateKey.String())
	log.Printf("PubKey: %s", wallet.PrivateKey.String())

	encrypted, err := encryptPrivateKey(wallet.PrivateKey, "pass")
	if err != nil {
		panic(err)
	}
	b, _ := json.MarshalIndent(&encrypted, "", "  ")
	log.Printf("Encrypted\n%s", string(b))

	b, _ = json.Marshal(&encrypted)
	decrypted, err := decrypt(b, "pass")
	if err != nil {
		panic(err)
	}
	log.Printf("Decripted private key: %s", base58.Encode(decrypted))
	// Output
	// 2021/11/30 20:38:34 Private Key: 4KwdbpBqcD6tim6KNoXivstHdpxoB5ik7S4dLhBHUKt4LzXSKQ5ecyy49qZRTAH2ubaNuMENjM6QuTi9CmkJt9aT
	//2021/11/30 20:38:34 PubKey: 4KwdbpBqcD6tim6KNoXivstHdpxoB5ik7S4dLhBHUKt4LzXSKQ5ecyy49qZRTAH2ubaNuMENjM6QuTi9CmkJt9aT
	//2021/11/30 20:38:34 Encrypted
	//{
	//  "cipher": "aes-128-ctr",
	//  "ciphertext": "4f4d1df8996828430dc0fe39c729ad6f987f78207f4730222deb2a961fb5fe93bff7718fcf72cb986bde309a1445d8da3875b01281d0f02fd17248a485ece230",
	//  "cipherparams": {
	//    "iv": "bc7458e5a91eefe207569ab15243c34b"
	//  },
	//  "kdf": "scrypt",
	//  "kdfparams": {
	//    "dklen": 32,
	//    "n": 262144,
	//    "p": 1,
	//    "r": 8,
	//    "salt": "2e4041bb3defab2e9452622fab3f5c61846e2e6475acb301978962df5497537e"
	//  },
	//  "mac": "94bfb90d7dc921ed6450bcde946d802e9654b3922908c19059aa30f5c2380bd1"
	//}
	//2021/11/30 20:38:35 Decripted private key: 4KwdbpBqcD6tim6KNoXivstHdpxoB5ik7S4dLhBHUKt4LzXSKQ5ecyy49qZRTAH2ubaNuMENjM6QuTi9CmkJt9aT
}

func encryptPrivateKey(privateKey []byte, passphrase string) (*cryptoJSON, error) {
	var (
		scryptN     = StandardScryptN
		scryptP     = StandardScryptP
		scryptR     = 8
		scryptDKLen = 32
	)

	authArray := []byte(passphrase)

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		log.Println("reading from crypto/rand failed: " + err.Error())
		return nil, err
	}
	derivedKey, err := scrypt.Key(authArray, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		log.Println("derived from scrypt:" + err.Error())
		return nil, err
	}
	encryptKey := derivedKey[:16]

	iv := make([]byte, aes.BlockSize) // 16
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println("reading from crypto/rand failed: " + err.Error())
		return nil, err
	}

	cipherText, err := aesCTRXOR(encryptKey, privateKey, iv)

	// ORIGIN: keccak256
	hash := sha3.New256()
	hash.Write(derivedKey[16:32])
	hash.Write(cipherText)
	mac := hash.Sum(nil)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)

	cipherParamsJSON := cipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	return &cryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}, nil
}

func decrypt(keyjsonBytes []byte, passphrase string) (keyBytes []byte, err error) {
	var cryptoJson cryptoJSON
	if err := json.Unmarshal(keyjsonBytes, &cryptoJson); err != nil {
		log.Println("unmarshal cryptoJSON from bytes:" + err.Error())
		return nil, err
	}

	if cryptoJson.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("cipher not supported: %v", cryptoJson.Cipher)
	}

	mac, err := hex.DecodeString(cryptoJson.MAC)
	if err != nil {
		panic(err)
	}

	iv, err := hex.DecodeString(cryptoJson.CipherParams.IV)
	if err != nil {
		panic(err)
	}

	cipherText, err := hex.DecodeString(cryptoJson.CipherText)
	if err != nil {
		panic(err)
	}

	derivedKey, err := getKDFKey(cryptoJson, passphrase)
	if err != nil {
		panic(err)
	}

	hash := sha3.New256()
	hash.Write(derivedKey[16:32])
	hash.Write(cipherText)
	calculatedMAC := hash.Sum(nil)

	if !bytes.Equal(calculatedMAC, mac) {
		panic(err)
	}

	plainText, err := aesCTRXOR(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func aesCTRXOR(key, inText, iv []byte) ([]byte, error) {
	// AES-128 is selected due to size of encryptKey.
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

func getKDFKey(cryptoJSON cryptoJSON, auth string) ([]byte, error) {
	authArray := []byte(auth)
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}
	dkLen := ensureInt(cryptoJSON.KDFParams["dklen"])

	if cryptoJSON.KDF == keyHeaderKDF {
		n := ensureInt(cryptoJSON.KDFParams["n"])
		r := ensureInt(cryptoJSON.KDFParams["r"])
		p := ensureInt(cryptoJSON.KDFParams["p"])
		return scrypt.Key(authArray, salt, n, r, p, dkLen)

	} else if cryptoJSON.KDF == "pbkdf2" {
		c := ensureInt(cryptoJSON.KDFParams["c"])
		prf := cryptoJSON.KDFParams["prf"].(string)
		if prf != "hmac-sha256" {
			return nil, fmt.Errorf("Unsupported PBKDF2 PRF: %s", prf)
		}
		key := pbkdf2.Key(authArray, salt, c, dkLen, sha256.New)
		return key, nil
	}

	return nil, fmt.Errorf("Unsupported KDF: %s", cryptoJSON.KDF)
}

// TODO: can we do without this when unmarshalling dynamic JSON?
// why do integers in KDF params end up as float64 and not int after
// unmarshal?
func ensureInt(x interface{}) int {
	res, ok := x.(int)
	if !ok {
		res = int(x.(float64))
	}
	return res
}
