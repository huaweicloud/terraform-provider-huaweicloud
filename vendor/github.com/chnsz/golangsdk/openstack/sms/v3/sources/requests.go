package sources

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// ListOpts 查询源端服务器列表
type ListOpts struct {
	// 源端服务器名称
	Name string `q:"name"`
	// 源端服务器ID
	Id string `q:"id"`
	// 源端服务器IP地址
	Ip string `q:"ip"`
	// 源端服务器状态
	State string `q:"state"`
	// 是否查询失去连接的源端
	Connected bool `q:"connected"`
	// 根据迁移周期查询
	MigrationCycle string `q:"migration_cycle"`
	// 迁移项目id,填写该参数将查询迁移项目下的所有虚拟机
	MigateProjectId string `q:"migproject"`
	// 需要查询的企业项目id
	EnterpriseProjectId string `q:"enterprise_project_id"`

	// 每一页记录的源端服务器数量,0表示用默认值 200
	Limit int `q:"limit"`
	// 偏移量,默认值0
	Offset int `q:"offset"`
}

// List is a method to query the list of the source servers using given opts.
func List(c *golangsdk.ServiceClient, opts ListOpts) ([]SourceServer, error) {
	url := listURL(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := SourceServerPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractSourceServers(pages)
}

// Get 查询指定ID的源端服务器
func Get(c *golangsdk.ServiceClient, id string) (*SourceServer, error) {
	var rst golangsdk.Result
	_, rst.Err = c.Get(getURL(c, id), &rst.Body, nil)

	var r SourceServer
	err := rst.ExtractInto(&r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
