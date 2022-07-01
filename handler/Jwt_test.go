package handler

import (
	"encoding/base64"
	"fmt"
	"playground_backend/common"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestJwtTest(t *testing.T) {
	payload := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMTUiLCJhdWQiOiI2MjFkZTg4YzQwYzgyOGMyMjk2Y2QxY2MiLCJpYXQiOjE2NTg2NDM3NjEsImV4cCI6MTY1OTg1MzM2MSwiaXNzIjoiaHR0cHM6Ly90cnltZS5hdXRoaW5nLmNuL29pZGMiLCJub25jZSI6IjcwMjY3NjMzMDQyMjg3NTEiLCJuYW1lIjoiIiwibmlja25hbWUiOiJsdW9uYW5jb20zIiwicGljdHVyZSI6Imh0dHBzOi8vb2JzLXhpaGUtYmVpamluZzQub2JzLmNuLW5vcnRoLTQubXlodWF3ZWljbG91ZC5jb20veGloZS1pbWcvZGVmYXVsdF9hdmF0YXIvbWFuLTQucG5nIiwid2Vic2l0ZSI6IiIsImJpcnRoZGF0ZSI6IiIsImdlbmRlciI6IiIsInVwZGF0ZWRfYXQiOiIyMDIyLTA3LTI0VDA2OjIyOjQxLjE0M1oiLCJlbWFpbCI6Imx1b25hbmNvbUBxcS5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGhvbmVfbnVtYmVyIjoiMTM1NTEzMjI0ODIiLCJwaG9uZV9udW1iZXJfdmVyaWZpZWQiOnRydWUsImFkZHJlc3MiOnsiY291bnRyeSI6IiIsInBvc3RhbF9jb2RlIjoiIn19.nhrCLI_3CxlNyQzwNC3yOlK4uO2LRy9N8ohsDRbxeKo"
	AesKey := []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR") // 对称秘钥长度必须是16的倍数

	t.Logf("加密前%v\n", payload)
	encrypted, err := common.AesEcrypt([]byte(payload), AesKey)
	if err != nil {
		t.Fatal(" AesEcrypt error ", err)
	}
	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)
	t.Logf("加密后%v \n", (encryptedBase64))
	jwtString, _ := GetJwtString(72, encryptedBase64, string(AesKey))

	token := new(jwt.Token)
	token.Valid = false

	if jwtString != "" {
		token, err = jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			return (AesKey), nil
		})
		if err != nil {
			t.Fatal(err)
			return
		}
	}
	if token.Valid {
		userinfo, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("token.Valid")
			return
		}
		userinfoData := fmt.Sprintf("%s", userinfo["data"])
		t.Logf("jwt 的data 内容:%v\n", (userinfoData))
		userinfoOrigin, _ := base64.StdEncoding.DecodeString(userinfoData)
		idtokenBytes, err := common.AesDeCrypt([]byte(userinfoOrigin), AesKey)
		if err != nil {
			t.Fatal(" AesDeCrypt error ", err)
		}
		t.Logf("解密后%v\n", string(idtokenBytes))

		token, _ = jwt.Parse(string(idtokenBytes), func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtString), nil
		})

		idtokenPayload, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("token.Valid")
			return
		}
		t.Log(idtokenPayload["email"])
		time.Sleep(time.Second)

	} else {
		t.Fatal("AesDeCrypt error3:")
	}

}

func TestJwtAes(t *testing.T) {
	inputStr := "abcdefg123"
	AesKey := []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR") // 对称秘钥长度必须是16的倍数

	jwtString, _ := GetJwtString(72, string(inputStr), string(AesKey))
	jwtList := strings.Split(jwtString, ".")

	payload := jwtList[1]
	t.Logf("加密前%v\n", payload)
	encrypted, err := common.AesEcrypt([]byte(payload), AesKey)
	if err != nil {
		t.Fatal(" AesEcrypt error ", err)
	}
	t.Logf("加密后%v\n", (encrypted))
	jwtString = fmt.Sprintf("%s.%s.%s", jwtList[0], (encrypted), jwtList[2])
	//---------------------------------------
	t.Logf("---%v\n", jwtString)
	jwtList = strings.Split(jwtString, ".")
	encrypted = []byte(jwtList[1])
	// t.Logf("2对比%v", []byte(temp))
	originPayLoad, err := common.AesDeCrypt(encrypted, AesKey)
	if err != nil {
		t.Fatal(" AesDeCrypt error ", err)
	}
	t.Logf("解密后%v\n", string(originPayLoad))
	jwtString = fmt.Sprintf("%s.%s.%s", jwtList[0], (originPayLoad), jwtList[2])

	token := new(jwt.Token)
	token.Valid = false

	if jwtString != "" {
		token, err = jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			return AesKey, nil
		})
		if err != nil {
			t.Fatal(err)
			return
		}
	}
	if token.Valid {
		userinfo, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("token.Valid")
			return
		}
		t.Log("---1--------\n")
		t.Log(userinfo["dt"])
		dataStr := userinfo["dt"].(string)
		t.Log("----2-------\n")
		t.Log((encrypted))
		t.Log("----3-----------\n")
		// t.Log((temp))
		// t.Log("----4----------\n")
		t.Log([]byte(dataStr))

		time.Sleep(time.Second)

	} else {
		t.Fatal("AesDeCrypt error3:")
	}
	return

}

func TestDemoAes(t *testing.T) {
	text := "123嘎达杀死对方"                                  // 你要加密的数据
	AesKey := []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR") // 对称秘钥长度必须是16的倍数
	fmt.Printf("明文: %s\n秘钥: %s\n", text, string(AesKey))
	encrypted, err := common.AesEcrypt([]byte(text), AesKey)
	if err != nil {
		panic(err)
	}
	origin, err := common.AesDeCrypt(encrypted, AesKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("解密后明文: %s\n", string(origin))
}
