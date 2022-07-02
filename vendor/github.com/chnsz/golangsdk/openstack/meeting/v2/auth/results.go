package auth

type AuthResp struct {
	// Token string.
	AccessToken string `json:"accessToken"`
	// Login account type.
	// 72: (fixed) API call type.
	ClientType int `json:"clientType"`
	// The creation timestamp of the token, in milliseconds.
	CreateTime int `json:"createTime"`
	// The number of days the password is valid.
	DaysPwdAvailable int `json:"daysPwdAvailable"`
	// Token's expiration timestamp, in seconds.
	ExpireTime int `json:"expireTime"`
	// Whether to log in for the first time.
	// The first login means that the password has not been changed. When logging in for the first time, the system will remind the user that the password needs to be changed.
	// Default: false.
	FirstLogin bool `json:"firstLogin"`
	// Proxy authentication information.
	ProxyToken ProxyToken `json:"proxyToken"`
	// Whether the password has expired.
	// Default: false.
	PwdExpired bool `json:"pwdExpired"`
	// The creation timestamp of the Refresh Token, in milliseconds.
	RefreshCreateTime int `json:"refreshCreateTime"`
	// The expiration timestamp of the Refresh Token, in seconds.
	RefreshExpireTime int `json:"refreshExpireTime"`
	// Refresh Token string.
	RefreshToken string `json:"refreshToken"`
	// The validity period of the Refresh Token, in seconds.
	RefreshValidPeriod int `json:"refreshValidPeriod"`
	// User IP.
	TokenIp string `json:"tokenIp"`
	// Token type.
	// 0: User Access Token
	// 1: Conference control Token
	// 2: One-time Token
	TokenType int `json:"tokenType"`
	// User authentication information.
	User UserInfo `json:"user"`
	// The valid duration of the token, in seconds.
	ValidPeriod int `json:"validPeriod"`
	// Preempt login ID.
	// 0: non-preemptive
	// 1: preemption is not enabled
	ForceLoginInd int `json:"forceLoginInd"`
	// Whether to delay the deletion of the state.
	DelayDelete bool `json:"delayDelete"`
}

type ProxyToken struct {
	// The short token string of the proxy authentication server.
	AccessToken string `json:"accessToken"`
	// Whether to enable secondary routing.
	EnableRerouting bool `json:"enableRerouting"`
	// The long token string of the proxy authentication server.
	LongAccessToken string `json:"longAccessToken"`
	// Middle-End intranet address.
	MiddleEndInnerUrl string `json:"middleEndInnerUrl"`
	// Middle-End address。
	MiddleEndUrl string `json:"middleEndUrl"`
	// Token valid duration, unit: seconds.
	ValidPeriod int `json:"validPeriod"`
}

type UserInfo struct {
	// Administrator type.
	// 0: default administrator
	// 1: Ordinary administrator
	// 2: Non-administrators, that is, ordinary enterprise members, valid when "userType" is "2".
	AdminType int `json:"adminType"`
	// Application ID.
	AppId string `json:"appId"`
	// HUAWEI CLOUD account ID.
	CloudUserId string `json:"cloudUserId"`
	// business domain name.
	CompanyDomain string `json:"companyDomain"`
	// The enterprise ID to which the user belongs.
	CompanyId string `json:"companyId"`
	// Enterprise plan type.
	// 0: Enterprise Edition;
	// 5: free version;
	// 6: Professional Edition.
	CorpType int `json:"corpType"`
	// Identifies whether it is a free trial user.
	FreeUser bool `json:"freeUser"`
	// Identifies whether it is a grayscale user.
	GrayUser bool `json:"grayUser"`
	// Avatar link.
	HeadPictureUrl string `json:"headPictureUrl"`
	// Indicates whether to bind the mobile phone.
	IsBindPhone bool `json:"isBindPhone"`
	// User name.
	Name string `json:"name"`
	// User name in English.
	NameEn string `json:"nameEn"`
	// The number corresponds to HA1.
	NumberHA1 string `json:"numberHA1"`
	// User alias.
	Alias string `json:"alias1"`
	// Paid user machine account, used for smart screen login.
	PaidAccount string `json:"paidAccount"`
	// Paid user machine account password, used for smart screen login.
	PaidPassword string `json:"paidPassword"`
	// Machine password, used for smart screen login.
	Password string `json:"password"`
	// Local authentication.
	Realm string `json:"realm"`
	// The SIP number associated with the user.
	ServiceAccount string `json:"serviceAccount"`
	// The SP ID of the enterprise where the user belongs.
	SpId string `json:"spId"`
	// user status.
	// 0: normal;
	// 1: Disable.
	Status int `json:"status"`
	// Third-party user accounts.
	ThirdAccount string `json:"thirdAccount"`
	// tr069 account number.
	Tr069Account string `json:"tr069Account"`
	// HUAWEI CLOUD conference account.
	UcloginAccount string `json:"ucloginAccount"`
	// User UUID.
	UserId string `json:"userId"`
	// user type.
	// 1: SP management user
	// 2: Enterprise users
	// 3: Free registered user
	// 10: Enterprise device users
	// 11: anonymous user
	// 12: Smart screen users
	//13: IdeaHub user
	// 14: Electronic whiteboard user
	UserType int `json:"userType"`
	// Smart screen device ID.
	VisionAccount string `json:"visionAccount"`
	// Identifies whether it is a WeLink user.
	WeLinkUser bool `json:"weLinkUser"`
}
