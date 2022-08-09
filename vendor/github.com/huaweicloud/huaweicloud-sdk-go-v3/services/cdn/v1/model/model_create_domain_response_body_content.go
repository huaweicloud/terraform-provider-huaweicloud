package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// 创建域名返回信息
type CreateDomainResponseBodyContent struct {

	// 加速域名ID。
	Id *string `json:"id,omitempty"`

	// 加速域名。
	DomainName *string `json:"domain_name,omitempty"`

	// 域名业务类型:-web:网站加速；-download:文件下载加速；-video:点播加速；-wholeSite:全站加速。
	BusinessType *string `json:"business_type,omitempty"`

	// 域名服务范围，若为mainland_china，则表示服务范围为中国大陆；若为outside_mainland_china，则表示服务范围为中国大陆境外；若为global，则表示服务范围为全球。
	ServiceArea *string `json:"service_area,omitempty"`

	// 域名所属用户的domain_id。
	UserDomainId *string `json:"user_domain_id,omitempty"`

	// 加速域名状态。取值意义：online表示“已开启”、offline表示“已停用”、configuring表示“配置中”、configure_failed表示“配置失败”、checking表示“审核中”、check_failed表示“审核未通过”、deleting表示“删除中”。
	DomainStatus *string `json:"domain_status,omitempty"`

	// 加速域名对应的CNAME。
	Cname *string `json:"cname,omitempty"`

	// 源站信息
	Sources *[]Sources `json:"sources,omitempty"`

	DomainOriginHost *DomainOriginHost `json:"domain_origin_host,omitempty"`

	// 是否开启HTTPS加速。
	HttpsStatus *int32 `json:"https_status,omitempty"`

	// 域名创建时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	CreateTime *int64 `json:"create_time,omitempty"`

	// 域名修改时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	ModifyTime *int64 `json:"modify_time,omitempty"`

	// 封禁状态（0代表未禁用；1代表禁用）。
	Disabled *int32 `json:"disabled,omitempty"`

	// 锁定状态（0代表未锁定；1代表锁定）。
	Locked *int32 `json:"locked,omitempty"`

	// range状态（\"off\"/\"on\"）。
	RangeStatus *string `json:"range_status,omitempty"`

	// follow302状态（\"off\"/\"on\"）。
	FollowStatus *string `json:"follow_status,omitempty"`

	// 是否暂停源站回源。
	OriginStatus *string `json:"origin_status,omitempty"`

	// 自动刷新预热（0代表关闭；1代表打开）
	AutoRefreshPreheat *int32 `json:"auto_refresh_preheat,omitempty"`
}

func (o CreateDomainResponseBodyContent) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDomainResponseBodyContent struct{}"
	}

	return strings.Join([]string{"CreateDomainResponseBodyContent", string(data)}, " ")
}
