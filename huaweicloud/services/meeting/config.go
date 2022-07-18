package meeting

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/meeting/v1/auth"
	v2auth "github.com/chnsz/golangsdk/openstack/meeting/v2/auth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type AuthOpts struct {
	AccountName        string
	Password           string
	AppId              string
	AppKey             string
	CorpId             string
	UserId             string
	TokenEffectiveTime int
	CurrentToken       string
}

// ConferenceTokens is a map used to storage all tokens for the conference management.
var (
	ConferenceTokens = map[string]string{}
	tokenRWMutex     sync.RWMutex
)

const (
	BasicAlgorithm  = "HMAC-SHA256"
	DefaultEndpoint = "https://api.meeting.huaweicloud.com"
)

func buildTokenEnvKey(d *schema.ResourceData) string {
	var corpId, userId string
	if val, ok := d.GetOk("corp_id"); ok {
		corpId = val.(string)
	}
	if val, ok := d.GetOk("user_id"); ok {
		userId = val.(string)
	}
	elements := []byte(d.Get("account_name").(string) + ":" + d.Get("account_password").(string) + ":" +
		d.Get("app_id").(string) + ":" + d.Get("app_key").(string) + ":" + corpId + ":" + userId)
	return base64.StdEncoding.EncodeToString(elements)
}

func buildAuthOpts(d *schema.ResourceData) AuthOpts {
	result := AuthOpts{
		AccountName:        d.Get("account_name").(string),
		Password:           d.Get("account_password").(string),
		AppId:              d.Get("app_id").(string),
		AppKey:             d.Get("app_key").(string),
		TokenEffectiveTime: 12,
		CurrentToken:       GetTokenFromEnv(buildTokenEnvKey(d)),
	}
	if val, ok := d.GetOk("corp_id"); ok {
		result.CorpId = val.(string)
	}
	if val, ok := d.GetOk("user_id"); ok {
		result.UserId = val.(string)
	}
	return result
}

func buildAuthOptsByState(state *terraform.ResourceState) AuthOpts {
	result := AuthOpts{
		AccountName:        state.Primary.Attributes["account_name"],
		Password:           state.Primary.Attributes["account_password"],
		AppId:              state.Primary.Attributes["app_id"],
		AppKey:             state.Primary.Attributes["app_key"],
		TokenEffectiveTime: 12,
	}
	if val, ok := state.Primary.Attributes["corp_id"]; ok {
		result.CorpId = val
	}
	if val, ok := state.Primary.Attributes["user_id"]; ok {
		result.UserId = val
	}
	return result
}

// GetTokenFromEnv is a method to obtain token from global map using given key.
func GetTokenFromEnv(key string) string {
	tokenRWMutex.Lock()
	defer tokenRWMutex.Unlock()
	val, ok := ConferenceTokens[key]
	if ok {
		return val
	}
	return ""
}

func isExpireTimeValid(expireTime int) bool {
	return expireTime >= 12 && expireTime <= 24
}

func getMeetingEndpoint(endpoints map[string]string, obj string) string {
	if endpoint, ok := endpoints["meeting"]; ok {
		return endpoint
	}
	return DefaultEndpoint
}

// NewMeetingV1Client is a method to general a client with ver.1 meeting endpoint.
func NewMeetingV1Client(conf *config.Config) *golangsdk.ServiceClient {
	return common.NewCustomClient(false, getMeetingEndpoint(conf.Endpoints, "meeting"), "v1")
}

// NewMeetingV2Client is a method to general a client with ver.2 meeting endpoint.
func NewMeetingV2Client(conf *config.Config) *golangsdk.ServiceClient {
	return common.NewCustomClient(false, getMeetingEndpoint(conf.Endpoints, "meeting"), "v2")
}

func isValidToken(client *golangsdk.ServiceClient, token string) bool {
	opts := auth.ValidateOpts{
		Token:           token,
		NeedGenNewToken: false,
	}
	_, err := auth.ValidateToken(client, opts)
	if err == nil {
		return true
	}
	if errCode, ok := err.(golangsdk.ErrDefault401); ok && errCode.Body != nil {
		var apiErr auth.ErrResponse
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr == nil && apiErr.Code == "USG.000000004" {
			return false
		}
	}
	log.Printf("[ERROR] Illegal response: %v", err)
	return false
}

// NewMeetingToken is a method to general a new token for conference control.
func NewMeetingToken(conf *config.Config, state interface{}) (string, error) {
	if res, ok := state.(*schema.ResourceData); ok {
		return getMeetingToken(conf, buildAuthOpts(res))
	}
	return getMeetingToken(conf, buildAuthOptsByState(state.(*terraform.ResourceState)))
}

func getMeetingToken(conf *config.Config, opt AuthOpts) (string, error) {
	if opt.AccountName != "" && opt.Password != "" {
		authorization := buildAuthorization(opt.AccountName, opt.Password)
		client := NewMeetingV1Client(conf)
		authOpts := auth.AuthOpts{
			Account:    opt.AccountName,
			ClientType: 72,
		}
		resp, err := auth.GetToken(client, authOpts, authorization)
		if err != nil {
			return "", err
		}
		return resp.AccessToken, nil
	}
	if opt.AppId != "" && opt.AppKey != "" && isExpireTimeValid(opt.TokenEffectiveTime) {
		client := NewMeetingV2Client(conf)
		if isValidToken(client, opt.CurrentToken) {
			return opt.CurrentToken, nil
		}
		authorization, nonce, expireTime := buildAppAuthorization(opt.AppId, opt.AppKey, opt.CorpId, opt.UserId, opt.TokenEffectiveTime)
		authOpts := v2auth.AuthOpts{
			AppId:      opt.AppId,
			ClientType: 72,
			ExpireTime: utils.Int(expireTime),
			Nonce:      nonce,
			CorpId:     opt.CorpId,
			UserId:     opt.UserId,
		}
		resp, err := v2auth.GetToken(client, authOpts, authorization)
		if err != nil {
			return "", err
		}
		return resp.AccessToken, nil
	}
	return "", fmt.Errorf("invalid auth information")
}

func buildAuthorization(username, password string) string {
	authInfo := []byte(username + ":" + password)
	auth := "Basic " + base64.StdEncoding.EncodeToString(authInfo)
	return auth
}

// buildAppAccessToken returns a authorization string according to the SP/Administrator account informations.
// Note: make sure the number is not a negative integer or a big integer.
func buildAppAuthorization(appId, appKey, corpInfo, userId string, validPeriod int) (auth, nonce string, expireTime int) {
	nonce = utils.RandomString(64)
	duration := validPeriod * 60 * 60
	expireTime = int(time.Now().Unix()) + duration

	var accountInfo = userId
	if corpInfo != "" {
		accountInfo = corpInfo + ":" + accountInfo
	}
	hmac256, _ := hmacsha256([]byte(appKey), appId+":"+accountInfo+":"+strconv.Itoa(expireTime)+":"+nonce)
	auth = fmt.Sprintf("%s signature=%s", BasicAlgorithm, hex.EncodeToString(hmac256))
	return
}

func hmacsha256(key []byte, data string) ([]byte, error) {
	h := hmac.New(sha256.New, key)
	if _, err := h.Write([]byte(data)); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
