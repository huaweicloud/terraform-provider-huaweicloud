package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DetailDataInfo struct {
	// 平均响应时间

	AverageRespTime *float64 `json:"averageRespTime,omitempty"`
	// 平均接收字节数

	AvgRecBytes *float64 `json:"avgRecBytes,omitempty"`
	// 平均发送字节数

	AvgSentBytes *float64 `json:"avgSentBytes,omitempty"`
	// 事务平均响应时间

	AvgTranRespTime *float64 `json:"avgTranRespTime,omitempty"`
	// 用例Uri

	CaseUri *string `json:"caseUri,omitempty"`
	// 创建时间

	CreateTime *string `json:"createTime,omitempty"`
	// 最大并发数

	CurrentThreadNum *float64 `json:"currentThreadNum,omitempty"`
	// 详情id

	DetailId *string `json:"detailId,omitempty"`
	// 结束时间

	EndTime *string `json:"endTime,omitempty"`
	// 失败请求数

	ErrorCount *float64 `json:"errorCount,omitempty"`
	// ERROR级别的事件个数

	ErrorEventsCount *float64 `json:"errorEventsCount,omitempty"`
	// 断言失败

	FailedAssert *float64 `json:"failedAssert,omitempty"`
	// 其他失败

	FailedOthers *float64 `json:"failedOthers,omitempty"`
	// 解析失败

	FailedParsed *float64 `json:"failedParsed,omitempty"`
	// 连接被拒

	FailedRefused *float64 `json:"failedRefused,omitempty"`
	// 超时失败

	FailedTimeout *float64 `json:"failedTimeout,omitempty"`
	// id

	Id *string `json:"id,omitempty"`
	// 是否aw

	IsAW *bool `json:"isAW,omitempty"`
	// 最大响应时间

	Max *float64 `json:"max,omitempty"`
	// 最大接收字节数

	MaxRecBytes *float64 `json:"maxRecBytes,omitempty"`
	// 探底最大响应时间

	MaxRespTime *float64 `json:"maxRespTime,omitempty"`
	// 最大发送字节数

	MaxSentBytes *float64 `json:"maxSentBytes,omitempty"`
	// 事务最大响应时间

	MaxTranRespTime *float64 `json:"maxTranRespTime,omitempty"`
	// 最小响应时间

	Min *float64 `json:"min,omitempty"`
	// 最小带宽

	MinNetworkTraffic *float64 `json:"minNetworkTraffic,omitempty"`
	// 名字

	Name *string `json:"name,omitempty"`
	// 请求数

	Requests *float64 `json:"requests,omitempty"`
	// aw执行结果

	Result *float64 `json:"result,omitempty"`
	// 开始时间

	StartTime *string `json:"startTime,omitempty"`
	// 用例状态

	Status *float64 `json:"status,omitempty"`
	// 成功请求数

	SuccessCount *float64 `json:"successCount,omitempty"`
	// 成功率

	SuccessRate *float64 `json:"successRate,omitempty"`
	// 1xx请求数

	Sum1xx *float64 `json:"sum1xx,omitempty"`
	// 2xx请求数

	Sum2xx *float64 `json:"sum2xx,omitempty"`
	// 3xx请求数

	Sum3xx *float64 `json:"sum3xx,omitempty"`
	// 4xx请求数

	Sum4xx *float64 `json:"sum4xx,omitempty"`
	// 5xx请求数

	Sum5xx *float64 `json:"sum5xx,omitempty"`
	// 任务id_轮次

	TaskId *string `json:"taskId,omitempty"`
	// 任务id

	TaskProjectId *string `json:"taskProjectId,omitempty"`
	// 任务状态

	TaskStatus *float64 `json:"taskStatus,omitempty"`
	// 用例uri

	TestCaseUri *string `json:"testCaseUri,omitempty"`
	// tp50

	Tp50 *float64 `json:"tp50,omitempty"`
	// tp75

	Tp75 *float64 `json:"tp75,omitempty"`
	// tp90

	Tp90 *float64 `json:"tp90,omitempty"`
	// tp95

	Tp95 *float64 `json:"tp95,omitempty"`
	// tp99

	Tp99 *float64 `json:"tp99,omitempty"`
	// tps

	Tps *float64 `json:"tps,omitempty"`
	// 事务tps

	TranTPS *float64 `json:"tranTPS,omitempty"`
	// 事务id

	TransactionId *string `json:"transactionId,omitempty"`
	// 事务成功率

	TransactionSuccess *float64 `json:"transactionSuccess,omitempty"`
	// 事务成功率

	TransactionalSuccessRate *float64 `json:"transactionalSuccessRate,omitempty"`
	// 自定义事务tps

	TransactionalTps *float64 `json:"transactionalTps,omitempty"`
	// 自定义事务成功率

	TransactionalTpsSuccess *float64 `json:"transactionalTpsSuccess,omitempty"`
	// 事务数

	Transactions *float64 `json:"transactions,omitempty"`
	// 更新时间

	UpdateTime *string `json:"updateTime,omitempty"`
	// 分钟*并发数

	Vum *float64 `json:"vum,omitempty"`
	// 平均带宽

	AvgNetworkTraffic *float64 `json:"avgNetworkTraffic,omitempty"`
	// 最大带宽

	MaxNetworkTraffic *float64 `json:"maxNetworkTraffic,omitempty"`
}

func (o DetailDataInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DetailDataInfo struct{}"
	}

	return strings.Join([]string{"DetailDataInfo", string(data)}, " ")
}
