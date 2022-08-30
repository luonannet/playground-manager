package handler

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"playground_backend/common"
	"playground_backend/models"
	"time"

	"github.com/Authing/authing-go-sdk/lib/authentication"
	"github.com/Authing/authing-go-sdk/lib/constant"
	"github.com/Authing/authing-go-sdk/lib/management"
	"github.com/Authing/authing-go-sdk/lib/model"
	"github.com/astaxie/beego"
	beegoCtx "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	OAUTH2 = "oauth2"
	GITHUB = "github"
	WECHAT = "wechat"
)

var OIDCConfig *oauth2.Config

type GiteeTokenInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
	AuthCode     string `json:"code"`
}

type GiteeUserInfo struct {
	UserId              int64
	SubUid              string
	Name                string
	PhoneNumber         string
	PhoneNumberVerified int8
	Nickname            string
	Picture             string
	Email               string
	EmailVerified       int8
	AccessToken         string
	ExpirationTime      string
	GivenName           string
	FamilyName          string
	MiddleName          string
	PreferredNsername   string
	Profile             string
	Website             string
	Gender              string
	Birthdate           string
	Zoneinfo            string
	Locale              string
	Formatted           string
	StreetAddress       string
	Locality            string
	Region              string
	PostalCode          string
	Country             string
	CreatedAt           string
	UpdatedAt           string
	CurStatus           model.EnumUserStatus
	UserName            string
	Unionid             string
	LastLogin           string
	LastIP              string
	SignedUp            string
	Blocked             bool
	IsDeleted           bool
	Device              string
	Identity            []Identities
}

type Identities struct {
	Openid      string
	IdentityId  string
	Provider    string
	ExtIdpId    string
	UserIdInIdp string
	Company     string
	City        string
	Email       string
	idpUserInfo IdpUserInfo
}
type AuthingLoginUser struct {
	Birthdate  string `json:"birthdate,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Name       string `json:"name,omitempty"`
	Nickname   string `json:"nickname,omitempty"`
	Picture    string `json:"picture,omitempty"`
	UpdatedAT  string `json:"updated_at,omitempty"`
	Website    string `json:"website,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Sub        string `json:"sub,omitempty"`
}

type IdpUserInfo struct {
	Phone    string
	Email    string
	Name     string
	UserName string
	Avatar   string
}

var JwtString = []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR")

func init() {
	//use env value
	if os.Getenv("TOKEN_KEY") != "" {
		JwtString = []byte(os.Getenv("TOKEN_KEY"))
	}
}

func GiteeConstructor(res map[string]interface{}, giteeToken *GiteeTokenInfo) {
	if _, ok := res["access_token"]; ok {
		giteeToken.AccessToken = res["access_token"].(string)
		giteeToken.TokenType = res["token_type"].(string)
		giteeToken.RefreshToken = ""
		giteeToken.IdToken = res["id_token"].(string)
		giteeToken.Scope = res["scope"].(string)
		giteeToken.ExpiresIn = int64(res["expires_in"].(float64))
		giteeToken.CreatedAt = common.PraseTimeInt(common.GetCurTime())
	}
}

func GiteeUserConstructor(res map[string]interface{}, gui *GiteeUserInfo) {
	sub, ok := res["sub"]
	if ok && sub != nil {
		gui.SubUid = sub.(string)
	}
	name, ok := res["name"]
	if ok && name != nil {
		gui.Name = name.(string)
	}
	givenName, ok := res["given_name"]
	if ok && givenName != nil {
		gui.GivenName = givenName.(string)
	}
	familyName, ok := res["family_name"]
	if ok && familyName != nil {
		gui.FamilyName = familyName.(string)
	}
	middleName, ok := res["middle_name"]
	if ok && middleName != nil {
		gui.MiddleName = middleName.(string)
	}
	nickname, ok := res["nickname"]
	if ok && nickname != nil {
		gui.Nickname = nickname.(string)
	}
	preferredUsername, ok := res["preferred_username"]
	if ok && preferredUsername != nil {
		gui.PreferredNsername = preferredUsername.(string)
	}
	profile, ok := res["profile"]
	if ok && profile != nil {
		gui.Profile = profile.(string)
	}
	picture, ok := res["picture"]
	if ok && picture != nil {
		gui.Picture = picture.(string)
	}
	website, ok := res["website"]
	if ok && website != nil {
		gui.Website = website.(string)
	}
	email, ok := res["email"]
	if ok && email != nil {
		gui.Email = email.(string)
	}
	emailVerified, ok := res["email_verified"]
	if ok && emailVerified != nil {
		if emailVerified.(bool) {
			gui.EmailVerified = 1
		} else {
			gui.EmailVerified = 0
		}
	}
	gender, ok := res["gender"]
	if ok && gender != nil {
		gui.Gender = gender.(string)
	}
	birthdate, ok := res["birthdate"]
	if ok && birthdate != nil {
		gui.Birthdate = birthdate.(string)
	}
	zoneinfo, ok := res["zoneinfo"]
	if ok && zoneinfo != nil {
		gui.Zoneinfo = zoneinfo.(string)
	}
	locale, ok := res["locale"]
	if ok && locale != nil {
		gui.Locale = locale.(string)
	}
	phoneNumber, ok := res["phone_number"]
	if ok && phoneNumber != nil {
		gui.PhoneNumber = phoneNumber.(string)
	}
	phoneNumberVerified, ok := res["phone_number_verified"]
	if ok && phoneNumberVerified != nil {
		if phoneNumberVerified.(bool) {
			gui.PhoneNumberVerified = 1
		} else {
			gui.PhoneNumberVerified = 0
		}
	}
	address, ok := res["address"]
	if ok && address != nil {
		addressObj := address.(map[string]interface{})
		formatted, ok := addressObj["formatted"]
		if ok && formatted != nil {
			gui.Formatted = formatted.(string)
		}
		streetAddress, ok := addressObj["street_address"]
		if ok && streetAddress != nil {
			gui.StreetAddress = streetAddress.(string)
		}
		locality, ok := addressObj["locality"]
		if ok && locality != nil {
			gui.Locality = locality.(string)
		}
		region, ok := addressObj["region"]
		if ok && region != nil {
			gui.Region = region.(string)
		}
		postalCode, ok := addressObj["postal_code"]
		if ok && postalCode != nil {
			gui.PostalCode = postalCode.(string)
		}
		country, ok := addressObj["country"]
		if ok && country != nil {
			gui.Country = country.(string)
		}
	}
	updatedAt, ok := res["updated_at"]
	if ok && updatedAt != nil {
		gui.UpdatedAt = updatedAt.(string)
	}
}

func ProcUserDetail(gui *models.AuthUserInfo, aud *models.AuthUserDetail,
	giteeUser *GiteeUserInfo, authToken AuthToken, userDetailList []string) {
	if len(authToken.IdentityId) < 1 {
		idStr := ""
		if len(gui.PhoneNumber) > 0 && len(gui.Email) > 0 {
			idStr = gui.PhoneNumber + gui.Email
			aud.Provider = "PhoneNumber,Email"
		} else if len(gui.PhoneNumber) > 0 {
			idStr = gui.PhoneNumber
			aud.Provider = "PhoneNumber"
		} else {
			idStr = gui.Email
			aud.Provider = "Email"
		}
		if len(idStr) > 0 {
			data := []byte(idStr)
			has := md5.Sum(data)
			aud.IdentityId = fmt.Sprintf("%x", has)
			aud.UserName = giteeUser.UserName
			userDetailList = append(userDetailList, "UserName")
			aud.NickName = giteeUser.Nickname
			userDetailList = append(userDetailList, "Nickname")
			aud.Photo = giteeUser.Picture
			userDetailList = append(userDetailList, "Photo")
			aud.Email = giteeUser.Email
			userDetailList = append(userDetailList, "Email")
			audMod := models.AuthUserDetail{IdentityId: aud.IdentityId}
			queryErr := models.QueryAuthUserDetail(&audMod, "IdentityId")
			if audMod.UserDetailId > 0 {
				GetFieldName(aud, audMod)
				aud.UpdateTime = common.GetCurTime()
				userDetailList = append(userDetailList, "UpdateTime")
				upDetailErr := models.UpdateAuthUserDetail(aud, userDetailList...)
				if upDetailErr != nil {
					logs.Error("ProcUserDetail, upDetailErr: ", upDetailErr)
				}
			} else {
				aud.CreateTime = common.GetCurTime()
				id, inDetailErr := models.InsertAuthUserDetail(aud)
				if inDetailErr != nil {
					logs.Error("ProcUserDetail, inDetailErr: ", inDetailErr, ",id: ", id, queryErr)
				}
			}
		}
	} else {
		if len(giteeUser.Identity) > 0 {
			for _, idy := range giteeUser.Identity {
				if len(idy.IdentityId) < 1 {
					continue
				}
				if authToken.IdentityId == idy.IdentityId {
					gui.UserName = idy.idpUserInfo.UserName
					gui.Picture = idy.idpUserInfo.Avatar
					gui.NickName = idy.idpUserInfo.Name
					gui.Email = idy.idpUserInfo.Email
					gui.PhoneNumber = idy.idpUserInfo.Phone
					giteeUser.UserName = idy.idpUserInfo.UserName
					giteeUser.Picture = idy.idpUserInfo.Avatar
					giteeUser.Nickname = idy.idpUserInfo.Name
					giteeUser.PhoneNumber = idy.idpUserInfo.Phone
					giteeUser.Email = idy.idpUserInfo.Email
				}
				aud.IdentityId = idy.IdentityId
				aud.Openid = idy.Openid
				userDetailList = append(userDetailList, "Openid")
				aud.Provider = idy.Provider
				userDetailList = append(userDetailList, "Provider")
				aud.ExtIdpId = idy.ExtIdpId
				userDetailList = append(userDetailList, "ExtIdpId")
				aud.UserIdInIdp = idy.UserIdInIdp
				userDetailList = append(userDetailList, "UserIdInIdp")
				aud.UserName = idy.idpUserInfo.UserName
				userDetailList = append(userDetailList, "UserName")
				aud.NickName = idy.idpUserInfo.Name
				userDetailList = append(userDetailList, "NickName")
				aud.Photo = idy.idpUserInfo.Avatar
				userDetailList = append(userDetailList, "Photo")
				aud.Company = idy.Company
				userDetailList = append(userDetailList, "Company")
				aud.City = idy.City
				userDetailList = append(userDetailList, "City")
				aud.Email = idy.Email
				userDetailList = append(userDetailList, "Email")
				audMod := models.AuthUserDetail{IdentityId: idy.IdentityId}
				queryErr := models.QueryAuthUserDetail(&audMod, "IdentityId")
				if audMod.UserDetailId > 0 {
					GetFieldName(aud, audMod)
					aud.UpdateTime = common.GetCurTime()
					userDetailList = append(userDetailList, "UpdateTime")
					upDetailErr := models.UpdateAuthUserDetail(aud, userDetailList...)
					if upDetailErr != nil {
						logs.Error("ProcUserDetail, upDetailErr: ", upDetailErr)
					}
				} else {
					aud.CreateTime = common.GetCurTime()
					id, inDetailErr := models.InsertAuthUserDetail(aud)
					if inDetailErr != nil {
						logs.Error("ProcUserDetail, inDetailErr: ", inDetailErr, ",id: ", id, queryErr)
					}
				}
			}
		}
	}
}

func GetFieldName(aud *models.AuthUserDetail, audMod models.AuthUserDetail) {
	aud.UserDetailId = audMod.UserDetailId
	aud.UserId = audMod.UserId
	aud.IdentityId = audMod.IdentityId
	aud.CreateTime = audMod.CreateTime
	aud.UpdateTime = audMod.UpdateTime
	aud.DeleteTime = audMod.DeleteTime
}

func ProcOauthData(giteeToken GiteeTokenInfo, giteeUser *GiteeUserInfo, token string, authToken AuthToken) int64 {
	userId := int64(0)
	userList := []string{}
	userDetailList := []string{}
	gui := models.AuthUserInfo{SubUid: giteeUser.SubUid}
	aud := models.AuthUserDetail{}
	queryErr := models.QueryAuthUserInfo(&gui, "SubUid")
	if queryErr != nil {
		return 0
	}
	if gui.UserId > 0 {
		userId = gui.UserId
		userList, userDetailList = CreateGiteeUserInfo(&gui, &aud, giteeUser, 2, token)
		upErr := models.UpdateAuthUserInfo(&gui, userList...)
		if upErr != nil {
			logs.Error("ProcOauthData, upErr: ", upErr)
		}
		aud.UserId = userId
	} else {
		CreateGiteeUserInfo(&gui, &aud, giteeUser, 1, token)
		id, inErr := models.InsertAuthUserInfo(&gui)
		if id > 0 {
			userId = id
		} else {
			logs.Error("inErr: ", inErr)
			return 0
		}
		aud.UserId = userId
	}
	if userId > 0 {
		ProcUserDetail(&gui, &aud, giteeUser, authToken, userDetailList)
		if len(authToken.IdentityId) > 0 {
			guif := models.AuthUserInfo{SubUid: giteeUser.SubUid}
			queryErr := models.QueryAuthUserInfo(&guif, "SubUid")
			if queryErr == nil {
				if len(gui.PhoneNumber) > 0 {
					guif.PhoneNumber = gui.PhoneNumber
				}
				if len(gui.Email) > 0 {
					guif.Email = gui.Email
				}
				guif.Picture = gui.Picture
				guif.UserName = gui.UserName
				guif.NickName = gui.NickName
				upErr := models.UpdateAuthUserInfo(&guif, "PhoneNumber", "Email", "Picture", "UserName", "NickName")
				if upErr != nil {
					logs.Error("ProcOauthData, upErr: ", upErr)
				}
				aud.UserId = userId
			}
		}
		gti := models.AuthTokenInfo{UserId: userId}
		getErr := models.QueryAuthTokenInfo(&gti, "UserId")
		if gti.Id > 0 {
			CreateAuthTokenInfo(&gti, giteeToken, 2, giteeToken.AuthCode)
			if len(giteeToken.AuthCode) > 1 {
				models.UpdateAuthTokenInfo(&gti, "UpdateTime", "AccessToken",
					"ExpiresIn", "Scope", "CreatedAt", "RefreshToken",
					"TokenType", "authCode", "IdToken")
			} else {
				models.UpdateAuthTokenInfo(&gti, "UpdateTime", "AccessToken",
					"ExpiresIn", "Scope", "CreatedAt",
					"RefreshToken", "TokenType", "IdToken")
			}
		} else {
			logs.Info("getErr: ", getErr)
			gti.UserId = userId
			CreateAuthTokenInfo(&gti, giteeToken, 1, giteeToken.AuthCode)
			models.InsertAuthTokenInfo(&gti)
		}
	}
	return userId
}

func CreateAuthTokenInfo(gti *models.AuthTokenInfo, giteeToken GiteeTokenInfo, flag int, authCode string) {
	if flag == 1 {
		gti.CreateTime = common.GetCurTime()
	} else {
		gti.UpdateTime = common.GetCurTime()
	}
	gti.AccessToken = giteeToken.AccessToken
	gti.ExpiresIn = giteeToken.ExpiresIn
	gti.Scope = giteeToken.Scope
	gti.CreatedAt = giteeToken.CreatedAt
	gti.RefreshToken = giteeToken.RefreshToken
	gti.IdToken = giteeToken.IdToken
	gti.TokenType = giteeToken.TokenType
	if len(authCode) > 0 {
		gti.AuthCode = authCode
	}
}

func CreateGiteeUserInfo(gui *models.AuthUserInfo, aud *models.AuthUserDetail,
	giteeUser *GiteeUserInfo, flag int, token string) ([]string, []string) {
	updateList := make([]string, 0)
	updateDetailList := make([]string, 0)
	expirTime, _ := beego.AppConfig.Int("gitee::token_expir_time")
	newTime := time.Now().AddDate(0, 0, expirTime).Format(common.DATE_FORMAT)
	gui.SubUid = giteeUser.SubUid
	updateList = append(updateList, "SubUid")
	gui.Name = giteeUser.Name
	updateList = append(updateList, "Name")
	gui.UserName = giteeUser.UserName
	updateList = append(updateList, "UserName")
	gui.PhoneNumber = giteeUser.PhoneNumber
	updateList = append(updateList, "PhoneNumber")
	gui.PhoneNumberVerified = giteeUser.PhoneNumberVerified
	updateList = append(updateList, "PhoneNumberVerified")
	gui.AccessToken = token
	updateList = append(updateList, "AccessToken")
	gui.EmailVerified = giteeUser.EmailVerified
	updateList = append(updateList, "EmailVerified")
	gui.Email = giteeUser.Email
	updateList = append(updateList, "Email")
	gui.Picture = giteeUser.Picture
	updateList = append(updateList, "Picture")
	gui.NickName = giteeUser.Nickname
	updateList = append(updateList, "Nickname")
	if giteeUser.IsDeleted {
		gui.Status = 2
	} else {
		gui.Status = 1
	}
	updateList = append(updateList, "Status")
	gui.ExpirationTime = newTime
	updateList = append(updateList, "ExpirationTime")
	if flag == 1 {
		gui.CreateTime = common.GetCurTime()
		updateList = append(updateList, "CreateTime")
		aud.CreateTime = gui.CreateTime
		updateDetailList = append(updateDetailList, "CreateTime")
	} else {
		gui.UpdateTime = common.GetCurTime()
		updateList = append(updateList, "UpdateTime")
		aud.UpdateTime = gui.UpdateTime
		updateDetailList = append(updateDetailList, "CreateTime")
	}
	aud.UpdatedAt = giteeUser.UpdatedAt
	updateDetailList = append(updateDetailList, "UpdatedAt")
	aud.Country = giteeUser.Country
	updateDetailList = append(updateDetailList, "Country")
	aud.PostalCode = giteeUser.PostalCode
	updateDetailList = append(updateDetailList, "PostalCode")
	aud.Region = giteeUser.Region
	updateDetailList = append(updateDetailList, "Region")
	aud.Locality = giteeUser.Locality
	updateDetailList = append(updateDetailList, "Locality")
	aud.StreetAddress = giteeUser.StreetAddress
	updateDetailList = append(updateDetailList, "StreetAddress")
	aud.Formatted = giteeUser.Formatted
	updateDetailList = append(updateDetailList, "Formatted")
	aud.Locale = giteeUser.Locale
	updateDetailList = append(updateDetailList, "Locale")
	aud.Zoneinfo = giteeUser.Zoneinfo
	updateDetailList = append(updateDetailList, "Zoneinfo")
	aud.Birthdate = giteeUser.Birthdate
	updateDetailList = append(updateDetailList, "Birthdate")
	aud.Gender = giteeUser.Gender
	updateDetailList = append(updateDetailList, "Gender")
	aud.Website = giteeUser.Website
	updateDetailList = append(updateDetailList, "Website")
	aud.Profile = giteeUser.Profile
	updateDetailList = append(updateDetailList, "Profile")
	aud.PreferredNsername = giteeUser.PreferredNsername
	updateDetailList = append(updateDetailList, "PreferredNsername")
	aud.MiddleName = giteeUser.MiddleName
	updateDetailList = append(updateDetailList, "MiddleName")
	aud.FamilyName = giteeUser.FamilyName
	updateDetailList = append(updateDetailList, "FamilyName")
	aud.GivenName = giteeUser.GivenName
	updateDetailList = append(updateDetailList, "GivenName")
	return updateList, updateDetailList
}

type RespUserInfo struct {
	UserId    int64  `json:"userId"`
	NickName  string `json:"nickName"`
	Nickname  string `json:"nickname"`
	Sub       string `json:"sub"`
	AvatarUrl string `json:"avatarUrl"`
	Picture   string `json:"picture"`
	UserToken string `json:"userToken"`
	Email     string `json:"email"`
}

type AuthCode struct {
	AuthCode string `json:"code"`
}

type AuthToken struct {
	Id         string `json:"id"`
	IdentityId string `json:"federationIdentityId"`
}

type ReqIdPrams struct {
	Id         string
	IdentityId string
}

var (
	authing_app_secret      string
	authing_app_id          string
	authing_userpool_secret string
	authing_userpool_id     string
	authing_authingurl      string
	authing_redicturl       string
)

func GetAuthCode() {
	clientSecret := beego.AppConfig.String("gitee::client_secret")
	clientId := beego.AppConfig.String("gitee::client_id")
	authenticationClient := authentication.NewClient(clientId, clientSecret)
	authenticationClient.Protocol = constant.OIDC
	authenticationClient.TokenEndPointAuthMethod = constant.None
	req := model.OidcParams{
		AppId:       clientId,
		RedirectUri: "https://test.playground.osinfra.cn/api/playground/oauth2/callback",
		Nonce:       "test",
	}
	resp, err := authenticationClient.BuildAuthorizeUrlByOidc(req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}

func GetAuthToken(clientId, clientSecret string, gtk *GiteeTokenInfo) error {
	resValue := make(map[string]interface{}, 0)
	authenticationClient := authentication.NewClient(clientId, clientSecret)
	resp, err := authenticationClient.GetAccessTokenByCode(gtk.AuthCode)
	if err != nil {
		logs.Error("GetAuthToken, err: ", err.Error())
		return err
	} else {
		logs.Info("GetAuthToken, resp: ", resp)
		err = json.Unmarshal([]byte(resp), &resValue)
		if err != nil {
			logs.Error("jsonErr: ", err.Error())
			return err
		}
	}
	GiteeConstructor(resValue, gtk)
	return nil
}

func UserConstructor(user *model.User, gui *GiteeUserInfo) {
	if user.Token != nil {
		// gui.AccessToken = *user.Token
	}
	if user.Name != nil {
		gui.Name = *user.Name
	}
	if user.Website != nil {
		gui.Website = *user.Website
	}
	if user.Email != nil {
		gui.Email = *user.Email
	}
	if user.Phone != nil {
		gui.PhoneNumber = *user.Phone
	}
	if user.Nickname != nil {
		gui.Nickname = *user.Nickname
	}
	gui.SubUid = user.Id
	if user.Birthdate != nil {
		gui.Birthdate = *user.Birthdate
	}
	if user.Locality != nil {
		gui.Locality = *user.Locality
	}
	if user.Region != nil {
		gui.Region = *user.Region
	}
	if user.Formatted != nil {
		gui.Formatted = *user.Formatted
	}
	if user.Gender != nil {
		gui.Gender = *user.Gender
	}
	if user.Photo != nil {
		gui.Picture = *user.Photo
	}
	if user.GivenName != nil {
		gui.GivenName = *user.GivenName
	}
	if user.FamilyName != nil {
		gui.FamilyName = *user.FamilyName
	}
	if user.MiddleName != nil {
		gui.MiddleName = *user.MiddleName
	}
	if user.PreferredUsername != nil {
		gui.PreferredNsername = *user.PreferredUsername
	}
	if user.Profile != nil {
		gui.Profile = *user.Profile
	}
	if user.Zoneinfo != nil {
		gui.Zoneinfo = *user.Zoneinfo
	}
	if user.Locale != nil {
		gui.Locale = *user.Locale
	}
	if user.StreetAddress != nil {
		gui.StreetAddress = *user.StreetAddress
	}
	if user.PostalCode != nil {
		gui.PostalCode = *user.PostalCode
	}
	if user.Country != nil {
		gui.Country = *user.Country
	}
	if user.UpdatedAt != nil {
		gui.UpdatedAt = *user.UpdatedAt
	}
	if user.CreatedAt != nil {
		gui.CreatedAt = *user.CreatedAt
	}
	if user.EmailVerified != nil {
		if *user.EmailVerified {
			gui.EmailVerified = 1
		} else {
			gui.EmailVerified = 0
		}
	} else {
		gui.EmailVerified = 0
	}
	if user.PhoneVerified != nil {
		if *user.PhoneVerified {
			gui.PhoneNumberVerified = 1
		} else {
			gui.PhoneNumberVerified = 0
		}
	} else {
		gui.PhoneNumberVerified = 0
	}
	if user.TokenExpiredAt != nil {
		gui.ExpirationTime = *user.TokenExpiredAt
	}
	if user.Username != nil {
		gui.UserName = *user.Username
	}
	if user.Blocked != nil {
		gui.Blocked = *user.Blocked
	}
	if user.Status != nil {
		gui.CurStatus = *user.Status
	}
	if user.Device != nil {
		gui.Device = *user.Device
	}
	if user.IsDeleted != nil {
		gui.IsDeleted = *user.IsDeleted
	}
	if user.LastIP != nil {
		gui.LastIP = *user.LastIP
	}
	if user.LastLogin != nil {
		gui.LastLogin = *user.LastLogin
	}
	if user.SignedUp != nil {
		gui.SignedUp = *user.SignedUp
	}
	if user.Unionid != nil {
		gui.Unionid = *user.Unionid
	}
	if user.Identities != nil && len(user.Identities) > 0 {
		idsList := make([]Identities, len(user.Identities))
		for _, idy := range user.Identities {
			ids := Identities{}
			if idy.Openid != nil {
				ids.Openid = *idy.Openid
			}
			if idy.UserIdInIdp != nil {
				ids.UserIdInIdp = *idy.UserIdInIdp
			}
			if idy.Id != nil {
				ids.IdentityId = *idy.Id
			}
			if idy.ExtIdpId != nil {
				ids.ExtIdpId = *idy.ExtIdpId
			}
			if idy.Provider != nil {
				ids.Provider = *idy.Provider
			}
			var idpUserInfo IdpUserInfo
			if user.Phone != nil {
				idpUserInfo.Phone = *user.Phone
			}
			if user.Email != nil {
				idpUserInfo.Email = *user.Email
			}
			if user.Username != nil {
				idpUserInfo.UserName = *user.Username
			}
			if user.Nickname != nil {
				idpUserInfo.Name = *user.Nickname
			} else if user.Name != nil {
				idpUserInfo.Name = *user.Name
			}
			if user.Photo != nil {
				idpUserInfo.Avatar = *user.Photo
			}
			userErr := getCorrectUserInfo(idy, &idpUserInfo)
			if userErr == nil {
				ids.idpUserInfo = idpUserInfo
			}
			idsList = append(idsList, ids)
		}
		gui.Identity = idsList
	} else {
		logs.Error("user.Identities: ", user.Identities)
	}
	logs.Info("The return result of authing is:", gui)
}

func getCorrectUserInfo(idy *model.Identity, idpUserInfo *IdpUserInfo) error {
	jsonByte, err := json.Marshal(idy.UserInfoInIdp)
	if err != nil {
		return err
	}
	userInfoInIdpJson, err := simplejson.NewJson(jsonByte)
	provider := *idy.Provider
	if err != nil {
		return err
	}
	if provider == OAUTH2 {
		name, _ := userInfoInIdpJson.Get("middleName").String()
		username, _ := userInfoInIdpJson.Get("familyName").String()
		avatar, _ := userInfoInIdpJson.Get("photo").String()
		idpUserInfo.Name = name
		idpUserInfo.UserName = username
		idpUserInfo.Avatar = avatar
	} else if provider == GITHUB {
		name, _ := userInfoInIdpJson.Get("nickname").String()
		username, _ := userInfoInIdpJson.Get("username").String()
		avatar, _ := userInfoInIdpJson.Get("photo").String()
		idpUserInfo.Name = name
		idpUserInfo.UserName = username
		idpUserInfo.Avatar = avatar
	} else if provider == WECHAT {
		name, _ := userInfoInIdpJson.Get("nickname").String()
		avatar, _ := userInfoInIdpJson.Get("photo").String()
		idpUserInfo.Name = name
		idpUserInfo.Avatar = avatar
	}
	return nil
}

func GetAuthUserBySub(userPoolId, userPoolSecret, subId string, gui *GiteeUserInfo) error {
	client := management.NewClient(userPoolId, userPoolSecret)
	resp, err := client.Detail(subId)
	if err != nil {
		logs.Error("GetAuthUserBySub, err: ", err.Error())
		return err
	} else {
		logs.Info("GetAuthUserBySub, resp: ", resp)
		UserConstructor(resp, gui)
	}
	return nil
}

func GetAuthUser(clientId, clientSecret, authToken string, gui *GiteeUserInfo) error {
	resValue := make(map[string]interface{}, 0)
	authenticationClient := authentication.NewClient(clientId, clientSecret)
	resp, err := authenticationClient.GetUserInfoByAccessToken(authToken)
	if err != nil {
		logs.Error("GetAuthUser, err: ", err.Error())
		return err
	} else {
		logs.Info("GetAuthUser, resp: ", resp)
		err = json.Unmarshal([]byte(resp), &resValue)
		if err != nil {
			logs.Error("jsonErr: ", err.Error())
			return err
		}
	}
	GiteeUserConstructor(resValue, gui)
	return nil
}

func CheckAuthingIdToken(authToken AuthToken, rip *ReqIdPrams) error {
	// 1. get environment variables
	clientSecret := beego.AppConfig.String("gitee::client_secret")
	clientId := beego.AppConfig.String("gitee::client_id")
	req := model.ValidateTokenRequest{
		AccessToken: "",
		//IdToken:     authToken.AuthToken,
	}
	resValue := make(map[string]interface{}, 0)
	authenticationClient := authentication.NewClient(clientId, clientSecret)
	resp, err := authenticationClient.ValidateToken(req)
	if err != nil {
		logs.Error("CheckAuthingIdToken, err: ", err.Error())
		return err
	} else {
		logs.Info("CheckAuthingIdToken, resp: ", resp)
		err = json.Unmarshal([]byte(resp), &resValue)
		if err != nil {
			logs.Error("jsonErr: ", err.Error())
			return err
		}
		if sub, ok := resValue["sub"]; ok {
			rip.Id = sub.(string)
		}
		if id, ok := resValue["id"]; ok {
			rip.Id = id.(string)
		}
		if identityId, ok := resValue["federationIdentityId"]; ok {
			rip.IdentityId = identityId.(string)
		}
	}
	return nil
}

func SaveAuthUserInfo(authCode AuthCode, rui *RespUserInfo, gui *GiteeUserInfo) error {
	// 1. get environment variables
	clientSecret := beego.AppConfig.String("gitee::client_secret")
	clientId := beego.AppConfig.String("gitee::client_id")
	// 2. define variable value
	var gtk GiteeTokenInfo
	var authToken AuthToken
	gtk.AuthCode = authCode.AuthCode
	// 3. Obtain token information based on authorization code
	tokenErr := GetAuthToken(clientId, clientSecret, &gtk)
	if tokenErr != nil {
		logs.Error("tokenErr: ", tokenErr)
	}

	if len(gtk.AccessToken) > 1 {
		// 4. Get user information
		userErr := GetAuthUser(clientId, clientSecret, gtk.AccessToken, gui)
		if userErr != nil {
			logs.Error("userErr: ", userErr)
			GetAuthUserFromDb(gtk, rui, gui)
			return userErr
		}
		if len(gui.SubUid) > 0 {
			// 5. Store user information
			saveErr := SaveAuthUser(rui, gtk, gui, authToken)
			if saveErr != nil {
				logs.Error("saveErr: ", saveErr)
				GetAuthUserFromDb(gtk, rui, gui)
				return saveErr
			}
		} else {
			GetAuthUserFromDb(gtk, rui, gui)
		}
	} else {
		GetAuthUserFromDb(gtk, rui, gui)
	}
	return nil
}

func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SaveAuthUserByToken(rip ReqIdPrams, rui *RespUserInfo, gui *GiteeUserInfo, authToken AuthToken) error {
	// 1. get environment variables
	userPoolSecret := beego.AppConfig.String("gitee::userpool_secret")
	userPoolId := beego.AppConfig.String("gitee::userpool_id")
	// 2. define variable value
	err := GetAuthUserBySub(userPoolId, userPoolSecret, rip.Id, gui)
	if err != nil {
		gui.SubUid = rip.Id
		GetAuthUserFromDbBySubId(rui, gui)
		return err
	}

	var gtk GiteeTokenInfo
	gtk.AccessToken = gui.AccessToken
	if len(gtk.AccessToken) > 1 {
		// 4. Get user information
		if len(gui.SubUid) > 0 {
			// 5. Store user information
			gtk.IdToken = gui.AccessToken
			// gtk.AccessToken = ""

			saveErr := SaveAuthUser(rui, gtk, gui, authToken)
			if saveErr != nil {
				logs.Error("saveErr: ", saveErr)
				gui.SubUid = rip.Id
				GetAuthUserFromDbBySubId(rui, gui)
				return saveErr
			}
		} else {
			gui.SubUid = rip.Id
			GetAuthUserFromDbBySubId(rui, gui)
		}
	} else {
		gui.SubUid = rip.Id
		GetAuthUserFromDbBySubId(rui, gui)
	}
	return nil
}

func GetAuthUserFromDbBySubId(rui *RespUserInfo, guu *GiteeUserInfo) {
	gui := models.AuthUserInfo{SubUid: guu.SubUid}
	queryErr := models.QueryAuthUserInfo(&gui, "SubUid")
	if queryErr == nil {
		rui.UserId = gui.UserId
		rui.UserToken = gui.AccessToken
		rui.NickName = gui.UserName
		rui.AvatarUrl = gui.Picture
		guu.Picture = gui.Picture
		guu.SubUid = gui.SubUid
		guu.Name = gui.Name
		guu.Nickname = gui.NickName
		guu.PhoneNumber = gui.PhoneNumber
		guu.Email = gui.Email
		rui.Email = gui.Email
		guu.UserId = gui.UserId
	}
	aud := models.AuthUserDetail{UserId: gui.UserId}
	queryErr = models.QueryAuthUserDetail(&aud, "UserId")
	if queryErr == nil {
		guu.Gender = aud.Gender
		guu.Formatted = aud.Formatted
		guu.Region = aud.Region
		guu.Locality = aud.Locality
		guu.Birthdate = aud.Birthdate
	}
}

func GetAuthUserFromDb(gtk GiteeTokenInfo, rui *RespUserInfo, guu *GiteeUserInfo) {
	gti := models.AuthTokenInfo{}
	getErr := errors.New("")
	if len(gtk.AuthCode) > 1 {
		gti.AuthCode = gtk.AuthCode
		getErr = models.QueryAuthTokenInfo(&gti, "AuthCode")
	} else {
		gti.AccessToken = gtk.AccessToken
		getErr = models.QueryAuthTokenInfo(&gti, "AccessToken")
	}
	if getErr == nil {
		gui := models.AuthUserInfo{UserId: gti.UserId}
		queryErr := models.QueryAuthUserInfo(&gui, "UserId")
		if queryErr == nil {
			rui.UserId = gui.UserId
			rui.UserToken = gui.AccessToken
			rui.NickName = gui.NickName
			rui.AvatarUrl = gui.Picture
			guu.Picture = gui.Picture
			guu.SubUid = gui.SubUid
			guu.Name = gui.Name
			guu.Nickname = gui.NickName
			guu.PhoneNumber = gui.PhoneNumber
			guu.Email = gui.Email
			rui.Email = gui.Email
			guu.UserId = gui.UserId
		}
		aud := models.AuthUserDetail{UserId: gti.UserId}
		queryErr = models.QueryAuthUserDetail(&aud, "UserId")
		if queryErr == nil {
			guu.Gender = aud.Gender
			guu.Formatted = aud.Formatted
			guu.Region = aud.Region
			guu.Locality = aud.Locality
			guu.Birthdate = aud.Birthdate
		}
	}
}

func SaveAuthUser(rui *RespUserInfo, gtk GiteeTokenInfo, gui *GiteeUserInfo, authToken AuthToken) error {
	// token, terr := common.GenToken(strconv.Itoa(int(rui.UserId)), gtk.AccessToken)
	token := gtk.AccessToken
	userId := ProcOauthData(gtk, gui, token, authToken)
	gui.UserId = userId
	CreateRespUserInfo(rui, gtk, gui)
	rui.UserId = userId
	rui.UserToken = token

	return nil
}

func CreateRespUserInfo(rui *RespUserInfo, giteeToken GiteeTokenInfo, giteeUser *GiteeUserInfo) {
	rui.Email = giteeUser.Email
	rui.NickName = giteeUser.UserName
	rui.AvatarUrl = giteeUser.Picture
}

//CheckToken Check whether the token is legal
func GetGiteeUserData(gui *models.AuthUserInfo, rui *RespUserInfo) bool {
	queryErr := models.QueryAuthUserInfo(gui, "AccessToken", "UserId", "ExpirationTime")
	if gui.UserId > 0 {
		now := time.Now().Format(common.DATE_FORMAT)
		logs.Info("token: now: ", now, ",expir: ", gui.ExpirationTime)
		// if now > gui.ExpirationTime {
		// 	return false
		// }
	} else {
		logs.Error("queryErr: ", queryErr)
		return false
	}
	GetUserInfoByUserId(gui, rui)
	return true
}

func CheckToken(gui *models.AuthUserInfo) bool {
	queryErr := models.QueryAuthUserInfo(gui, "AccessToken", "UserId", "ExpirationTime")
	if gui.UserId > 0 {
		now := time.Now().Format(common.DATE_FORMAT)
		logs.Info("token: now: ", now, ",expir: ", gui.ExpirationTime)
		// if now > gui.ExpirationTime {
		// 	return false
		// }
	} else {
		logs.Error("queryErr: ", queryErr)
		return false
	}
	return true
}

func GetUserInfoByReshToken(userId int64, token string, rui *RespUserInfo) {
	clientSecret := beego.AppConfig.String("gitee::client_secret")
	clientId := beego.AppConfig.String("gitee::client_id")
	var gtk GiteeTokenInfo
	var gui GiteeUserInfo
	var authToken AuthToken
	gti := models.AuthTokenInfo{UserId: userId}
	models.QueryAuthTokenInfo(&gti, "UserId")
	if gti.Id > 0 {
		gtk.AuthCode = gti.AuthCode
		if len(gti.AccessToken) > 1 {
			// 4. Get user information
			userErr := GetAuthUser(clientId, clientSecret, gti.AccessToken, &gui)
			if userErr != nil {
				logs.Error("userErr: ", userErr)
				GetAuthUserFromDb(gtk, rui, &gui)
			} else {
				// 5. Store user information
				if len(gui.SubUid) > 0 {
					userId := ProcOauthData(gtk, &gui, token, authToken)
					gui.UserId = userId
					CreateRespUserInfo(rui, gtk, &gui)
					rui.UserId = userId
					rui.UserToken = token
				} else {
					GetAuthUserFromDb(gtk, rui, &gui)
				}
			}
		} else {
			GetAuthUserFromDb(gtk, rui, &gui)
		}
		// Save key information to file
		if rui.UserId > 0 {
			userStr := ""
			userJson, jsonErr := json.Marshal(gui)
			if jsonErr == nil {
				userStr = string(userJson)
			}
			sd := StatisticsData{UserId: rui.UserId,
				OperationTime: common.GetCurTime(), EventType: "Query login information", State: "success",
				StateMessage: "success", Body: userStr}
			sdErr := StatisticsLog(sd)
			if sdErr != nil {
				logs.Error("SaveAuthUserInfo, post, sdErr: ", sdErr)
			}
		}
	}
}

func GetUserInfoByUserId(aui *models.AuthUserInfo, rui *RespUserInfo) {
	var gui GiteeUserInfo
	rui.NickName = aui.UserName
	rui.Email = aui.Email
	rui.UserId = aui.UserId
	rui.AvatarUrl = aui.Picture
	rui.UserToken = aui.AccessToken
	GetAuthUserFromDbBySubId(rui, &gui)
	// Save key information to file
	if rui.UserId > 0 {
		userStr := ""
		userJson, jsonErr := json.Marshal(gui)
		if jsonErr == nil {
			userStr = string(userJson)
		}
		sd := StatisticsData{UserId: rui.UserId,
			OperationTime: common.GetCurTime(), EventType: "Query login information", State: "success",
			StateMessage: "success", Body: userStr}
		sdErr := StatisticsLog(sd)
		if sdErr != nil {
			logs.Error("SaveAuthUserInfo, post, sdErr: ", sdErr)
		}
	}
}

type RequestParameter struct {
	ResourceId   string `json:"resourceId"`
	CourseId     string `json:"courseId"`
	ChapterId    string `json:"chapterId"`
	Backend      string `json:"backend"`
	TemplatePath string `json:"templatePath"`
	UserId       string `json:"userId"`
	ContactEmail string `json:"contactEmail"`
	Token        string `json:"token"`
	ForceDelete  int    `json:"forceDelete"`
}
type ResData struct {
	ResInfo interface{} `json:"instanceInfo"`
	Mesg    string      `json:"message"`
	Code    int         `json:"code"`
}

func Authorize(ctx *beegoCtx.Context) {
	authString := ctx.Input.Header("token")

	if len(authString) == 0 {
		authString = ctx.Input.Query("token")
	}
	if len(authString) == 0 {
		var rp RequestParameter
		json.Unmarshal(ctx.Input.RequestBody, &rp)
		if len(rp.Token) > 0 {
			authString = rp.Token
		}
	}

	token := new(jwt.Token)
	token.Valid = false
	var err error
	if len(authString) > 0 {
		token, err = jwt.Parse(authString, func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtString), nil
		})
		if err != nil {
			ctx.Output.JSON("Authority authentication failed 1", false, false)
			return
		}

	}
	if token.Valid {
		userinfo, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.Output.JSON("Authority authentication failed 2", false, false)
			return
		}

		idtokenEncrypt := fmt.Sprintf("%s", userinfo["data"])
		idtokenStr, _ := base64.StdEncoding.DecodeString(idtokenEncrypt)
		var idtokenBytes []byte
		idtokenBytes, err = common.AesDeCrypt([]byte(idtokenStr), JwtString)
		if err != nil {
			ctx.Output.JSON("AesDeCrypt  1", false, false)
			return
		}
		//---解析idtoken的payload
		token, _ = jwt.Parse(string(idtokenBytes), func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtString), nil
		})
		idtokenPayload, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.Output.JSON("Authority authentication failed 5", false, false)
			return
		}
		userId := idtokenPayload["sub"].(string)
		ctx.Input.SetData("me", userId)
		ctx.Input.SetData("idtoken", string(idtokenBytes))
	} else {
		// ctx.Abort(http.StatusForbidden, "非法用户")
		ctx.WriteString("非法用户")
		return
	}
}

var authingClient *authentication.Client

func InitAuthing() error {
	authing_app_secret = beego.AppConfig.String("gitee::client_secret")
	authing_app_id = beego.AppConfig.String("gitee::client_id")
	authing_userpool_secret = beego.AppConfig.String("gitee::userpool_secret")
	authing_userpool_id = beego.AppConfig.String("gitee::userpool_id")
	authing_authingurl = beego.AppConfig.String("gitee::callback_url")
	authing_redicturl = beego.AppConfig.String("gitee::redirect_url")
	authingClient = authentication.NewClient(authing_app_id, authing_app_secret)

	ctx := context.Background()
	oidcProvider, err := oidc.NewProvider(ctx, authing_redicturl+"/oidc")
	if err != nil {
		log.Println(err, "--------InitAuthing-------:", authing_authingurl, "-------------", authing_redicturl)
		return err
	}
	OIDCConfig = &oauth2.Config{
		ClientID:     authing_app_id,
		ClientSecret: authing_app_secret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  authing_authingurl,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "external_id", "phone"},
	}
	return nil
}

func GetUserInfoByToken(access_token string) (userinfo *AuthingLoginUser, err error) {
	resp, err := http.Get(authing_redicturl + "/oidc/me?access_token=" + access_token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	userinfo = new(AuthingLoginUser)
	err = json.Unmarshal(respDataBytes, userinfo)
	if err != nil {
		return nil, err
	}
	return userinfo, nil

}

func SetCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

//GetJwtString GetJwtString
func GetJwtString(expire int, data, key string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	now := time.Now()
	claims["exp"] = now.Add(time.Hour * time.Duration(expire)).Unix()
	claims["iat"] = now.Unix()
	// claims["id"] = id
	claims["data"] = data
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(key))
	return tokenString, err
}

type LoginUserData struct {
	IdToken string            `json:"idtoken"`
	Token   string            `json:"token"`
	User    *AuthingLoginUser `json:"user"`
}

func GetTokenFromAuthing(code string) (*LoginUserData, error) {

	respStr, err := authingClient.GetAccessTokenByCode(code)
	if err != nil {
		return nil, err
	}
	var accessToken GiteeTokenInfo

	err = json.Unmarshal([]byte(respStr), &accessToken)
	if err != nil {
		logs.Error("jsonErr: ", err.Error())
		return nil, err
	}
	if accessToken.AccessToken == "" {
		return nil, fmt.Errorf(respStr)
	}
	userInfo, err := GetUserInfoByToken(accessToken.AccessToken)
	if err != nil {

		return nil, err
	}

	encrypted, err := common.AesEcrypt([]byte(accessToken.IdToken), JwtString)
	if err != nil {
		return nil, err
	}
	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)
	token, err := GetJwtString(72, encryptedBase64, string(JwtString))
	if err != nil {
		return nil, err
	}
	result := new(LoginUserData)
	result.User = userInfo
	result.Token = token
	return result, nil
}
