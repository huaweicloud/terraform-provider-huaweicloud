package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CaseReportDetail struct {

	// 最大成功率检查点
	MaxSuccessRateCheckPoint *float64 `json:"MaxSuccessRateCheckPoint,omitempty"`

	// 别名
	Alias *string `json:"alias,omitempty"`

	// 平均响应时间
	AverageRespTime *float64 `json:"averageRespTime,omitempty"`

	// 平均响应时间检查点
	AverageRespTimeCheckPoint *float64 `json:"averageRespTimeCheckPoint,omitempty"`

	// 平均响应时间检查结果
	AverageRespTimeCheckRes *bool `json:"averageRespTimeCheckRes,omitempty"`

	// 平均带宽
	AvgNetworkTraffic *float64 `json:"avgNetworkTraffic,omitempty"`

	// 平均下行带宽
	AvgRecBytes *int32 `json:"avgRecBytes,omitempty"`

	// 平均下行带宽检查点
	AvgRecBytesCheckPoint *float64 `json:"avgRecBytesCheckPoint,omitempty"`

	// 平均下行带宽检查结果
	AvgRecBytesCheckRes *bool `json:"avgRecBytesCheckRes,omitempty"`

	// 平均上行带宽
	AvgSentBytes *int32 `json:"avgSentBytes,omitempty"`

	// 平均上行带宽检查点
	AvgSentBytesCheckPoint *float64 `json:"avgSentBytesCheckPoint,omitempty"`

	// 平均上行带宽检查结果
	AvgSentBytesCheckRes *bool `json:"avgSentBytesCheckRes,omitempty"`

	// 事务平均响应时间
	AvgTranRespTime *float64 `json:"avgTranRespTime,omitempty"`

	// 事务平均响应时间检查点
	AvgTranRespTimeCheckPoint *float64 `json:"avgTranRespTimeCheckPoint,omitempty"`

	// 事务平均响应时间检查结果
	AvgTranRespTimeCheckRes *bool `json:"avgTranRespTimeCheckRes,omitempty"`

	// 请求id
	AwId *string `json:"awId,omitempty"`

	// 用例Uri
	CaseUri *string `json:"caseUri,omitempty"`

	// 所有检查点结果的汇总结果
	CheckpointResult *bool `json:"checkpointResult,omitempty"`

	// cpu最大使用率
	CpuUsage *float64 `json:"cpuUsage,omitempty"`

	// cpu平均使用率
	CpuUsageAvg *float64 `json:"cpuUsageAvg,omitempty"`

	// cpu平均使用率检查点
	CpuUsageAvgCheckPoint *float64 `json:"cpuUsageAvgCheckPoint,omitempty"`

	// cpu平均使用率检查结果
	CpuUsageAvgCheckRes *bool `json:"cpuUsageAvgCheckRes,omitempty"`

	// cpu最大使用率检查点
	CpuUsageCheckPoint *float32 `json:"cpuUsageCheckPoint,omitempty"`

	// cpu最大使用率检查结果
	CpuUsageCheckRes *bool `json:"cpuUsageCheckRes,omitempty"`

	// 创建时间
	CreateTime *string `json:"createTime,omitempty"`

	// 最大并发数
	CurrentThreadNum *int32 `json:"currentThreadNum,omitempty"`

	// 数据类型(case/aw/transaction)
	DatumType *int32 `json:"datumType,omitempty"`

	// dcs平均时延
	DcsLatencyAvg *float64 `json:"dcsLatencyAvg,omitempty"`

	// dcs平均时延检查点
	DcsLatencyAvgCheckPoint *float64 `json:"dcsLatencyAvgCheckPoint,omitempty"`

	// dcs平均时延检查结果
	DcsLatencyAvgCheckRes *bool `json:"dcsLatencyAvgCheckRes,omitempty"`

	// dcs最大时延
	DcsLatencyMax *float64 `json:"dcsLatencyMax,omitempty"`

	// dcs最大时延检查点·
	DcsLatencyMaxCheckPoint *float64 `json:"dcsLatencyMaxCheckPoint,omitempty"`

	// dcs最大时延检查结果
	DcsLatencyMaxCheckRes *bool `json:"dcsLatencyMaxCheckRes,omitempty"`

	// dcs最小时延
	DcsLatencyMin *float64 `json:"dcsLatencyMin,omitempty"`

	// dcs最小时延检查点
	DcsLatencyMinCheckPoint *float64 `json:"dcsLatencyMinCheckPoint,omitempty"`

	// dcs最小时延检查结果
	DcsLatencyMinCheckRes *bool `json:"dcsLatencyMinCheckRes,omitempty"`

	// 用例/aw/事务在数据库中dc_case_aw表的主键ID
	DetailId *string `json:"detailId,omitempty"`

	// 磁盘最大读取速度
	DiskRead *float64 `json:"diskRead,omitempty"`

	// 磁盘平均读取速度
	DiskReadAvg *float64 `json:"diskReadAvg,omitempty"`

	// 磁盘平均读取速度检查点
	DiskReadAvgCheckPoint *float64 `json:"diskReadAvgCheckPoint,omitempty"`

	// 磁盘平均读取速度检查结果
	DiskReadAvgCheckRes *bool `json:"diskReadAvgCheckRes,omitempty"`

	// 磁盘最大读取速度检查点
	DiskReadCheckPoint *float64 `json:"diskReadCheckPoint,omitempty"`

	// 磁盘最大读取速度检查结果
	DiskReadCheckRes *bool `json:"diskReadCheckRes,omitempty"`

	// 磁盘最大使用率
	DiskUsage *float64 `json:"diskUsage,omitempty"`

	// 磁盘平均使用率
	DiskUsageAvg *float64 `json:"diskUsageAvg,omitempty"`

	// 磁盘平均使用率检查点
	DiskUsageAvgCheckPoint *float64 `json:"diskUsageAvgCheckPoint,omitempty"`

	// 磁盘平均使用率检查结果
	DiskUsageAvgCheckRes *bool `json:"diskUsageAvgCheckRes,omitempty"`

	// 磁盘最大使用率检查点
	DiskUsageCheckPoint *float64 `json:"diskUsageCheckPoint,omitempty"`

	// 磁盘最大使用率检查结果
	DiskUsageCheckRes *bool `json:"diskUsageCheckRes,omitempty"`

	// 磁盘最大写入速度
	DiskWrite *float64 `json:"diskWrite,omitempty"`

	// 磁盘平均写入速度
	DiskWriteAvg *float64 `json:"diskWriteAvg,omitempty"`

	// 磁盘平均写入速度检查点
	DiskWriteAvgCheckPoint *float64 `json:"diskWriteAvgCheckPoint,omitempty"`

	// 磁盘平均写入速度检查结果
	DiskWriteAvgCheckRes *bool `json:"diskWriteAvgCheckRes,omitempty"`

	// 磁盘最大写入速度检查点
	DiskWriteCheckPoint *float64 `json:"diskWriteCheckPoint,omitempty"`

	// 磁盘最大写入速度检查结果
	DiskWriteCheckRes *bool `json:"diskWriteCheckRes,omitempty"`

	// 运行时长
	Duration *int32 `json:"duration,omitempty"`

	// 结束时间
	EndTime *string `json:"endTime,omitempty"`

	// 错误数
	ErrorCount *int32 `json:"errorCount,omitempty"`

	// 错误事件数
	ErrorEventsCount *int32 `json:"errorEventsCount,omitempty"`

	// 断言失败数
	FailedAssert *int32 `json:"failedAssert,omitempty"`

	// 其他失败数
	FailedOthers *int32 `json:"failedOthers,omitempty"`

	// 解析失败数
	FailedParsed *int32 `json:"failedParsed,omitempty"`

	// 失败原因
	FailedReason *string `json:"failedReason,omitempty"`

	// 连接拒绝失败数
	FailedRefused *int32 `json:"failedRefused,omitempty"`

	// 连接超时失败数
	FailedTimeout *int32 `json:"failedTimeout,omitempty"`

	// 用例在数据库中dc_testcase表的主键id
	Id *string `json:"id,omitempty"`

	// 是否是aw
	IsAW *bool `json:"isAW,omitempty"`

	// 迭代uri
	IterationUri *string `json:"iterationUri,omitempty"`

	// 来源于设计服务的监控数据
	KpiMonitor *string `json:"kpiMonitor,omitempty"`

	// 最大响应时间
	Max *int32 `json:"max,omitempty"`

	// 平均响应时间
	MaxAvgTime *float64 `json:"maxAvgTime,omitempty"`

	// 平均响应时间检查点
	MaxAvgTimeCheckPoint *float64 `json:"maxAvgTimeCheckPoint,omitempty"`

	// 平均响应时间检查结果
	MaxAvgTimeCheckRes *bool `json:"maxAvgTimeCheckRes,omitempty"`

	// 流量峰值
	MaxNetworkTraffic *int32 `json:"maxNetworkTraffic,omitempty"`

	// 最大下行带宽
	MaxRecBytes *int32 `json:"maxRecBytes,omitempty"`

	// 最大下行带宽检查点
	MaxRecBytesCheckPoint *float64 `json:"maxRecBytesCheckPoint,omitempty"`

	// 最大下行带宽检查结果
	MaxRecBytesCheckRes *bool `json:"maxRecBytesCheckRes,omitempty"`

	// 最大响应时间
	MaxRespTime *int32 `json:"maxRespTime,omitempty"`

	// 最大响应时间检查点
	MaxRespTimeCheckPoint *float64 `json:"maxRespTimeCheckPoint,omitempty"`

	// 最大响应时间检查结果
	MaxRespTimeCheckRes *bool `json:"maxRespTimeCheckRes,omitempty"`

	// 最大RPS
	MaxRps *int32 `json:"maxRps,omitempty"`

	// 最大上行带宽
	MaxSentBytes *int32 `json:"maxSentBytes,omitempty"`

	// 最大上行带宽检查点
	MaxSentBytesCheckPoint *float64 `json:"maxSentBytesCheckPoint,omitempty"`

	// 最大上行带宽检查结果
	MaxSentBytesCheckRes *bool `json:"maxSentBytesCheckRes,omitempty"`

	// 最大成功率
	MaxSuccessRate *float64 `json:"maxSuccessRate,omitempty"`

	// 最大成功率检查结果
	MaxSuccessRateCheckRes *bool `json:"maxSuccessRateCheckRes,omitempty"`

	// 最大线程数
	MaxThreadNum *float64 `json:"maxThreadNum,omitempty"`

	// 最大线程数检查点
	MaxThreadNumCheckPoint *float64 `json:"maxThreadNumCheckPoint,omitempty"`

	// 最大线程数检查结果
	MaxThreadNumCheckRes *bool `json:"maxThreadNumCheckRes,omitempty"`

	// 最大TPS
	MaxTps *float64 `json:"maxTps,omitempty"`

	// 最大TPS检查点
	MaxTpsCheckPoint *float64 `json:"maxTpsCheckPoint,omitempty"`

	// 最大TPS检查结果
	MaxTpsCheckRes *bool `json:"maxTpsCheckRes,omitempty"`

	// 最大事务响应时间
	MaxTranRespTime *float64 `json:"maxTranRespTime,omitempty"`

	// 最大事务响应时间检查点
	MaxTranRespTimeCheckPoint *float64 `json:"maxTranRespTimeCheckPoint,omitempty"`

	// 最大事务响应时间检查结果
	MaxTranRespTimeCheckRes *bool `json:"maxTranRespTimeCheckRes,omitempty"`

	// 最大内存使用率
	MemoryUsage *float64 `json:"memoryUsage,omitempty"`

	// 平均内存使用率
	MemoryUsageAvg *float64 `json:"memoryUsageAvg,omitempty"`

	// 平均内存使用率检查点
	MemoryUsageAvgCheckPoint *float64 `json:"memoryUsageAvgCheckPoint,omitempty"`

	// 平均内存使用率检查结果
	MemoryUsageAvgCheckRes *bool `json:"memoryUsageAvgCheckRes,omitempty"`

	// 最大内存使用率检查点
	MemoryUsageCheckPoint *float64 `json:"memoryUsageCheckPoint,omitempty"`

	// 最大内存使用率检查结果
	MemoryUsageCheckRes *bool `json:"memoryUsageCheckRes,omitempty"`

	// 最小响应时间
	Min *int32 `json:"min,omitempty"`

	// 流量谷值
	MinNetworkTraffic *int32 `json:"minNetworkTraffic,omitempty"`

	// 压力模式
	Mode *string `json:"mode,omitempty"`

	// 监控峰值时间
	MonitorPeakTime *float64 `json:"monitorPeakTime,omitempty"`

	// 监控峰值时间检查点
	MonitorPeakTimeCheckPoint *float64 `json:"monitorPeakTimeCheckPoint,omitempty"`

	// 监控峰值时间检查结果
	MonitorPeakTimeCheckRes *bool `json:"monitorPeakTimeCheckRes,omitempty"`

	// 监控结果
	MonitorResult *float64 `json:"monitorResult,omitempty"`

	// 监控结果检查点
	MonitorResultCheckPoint *float64 `json:"monitorResultCheckPoint,omitempty"`

	// 监控结果检查结果
	MonitorResultCheckRes *bool `json:"monitorResultCheckRes,omitempty"`

	// 用例/aw/事务名
	Name *string `json:"name,omitempty"`

	// 网络最大接收数据速度
	NetworkRead *float64 `json:"networkRead,omitempty"`

	// 网络平均接收数据速度
	NetworkReadAvg *float64 `json:"networkReadAvg,omitempty"`

	// 网络平均接收数据速度检查点
	NetworkReadAvgCheckPoint *float64 `json:"networkReadAvgCheckPoint,omitempty"`

	// 网络平均接收数据速度检查结果
	NetworkReadAvgCheckRes *bool `json:"networkReadAvgCheckRes,omitempty"`

	// 网络最大接收数据速度检查点
	NetworkReadCheckPoint *float64 `json:"networkReadCheckPoint,omitempty"`

	// 网络最大接收数据速度检查结果
	NetworkReadCheckRes *bool `json:"networkReadCheckRes,omitempty"`

	// 网络最大写入数据速度
	NetworkWrite *float64 `json:"networkWrite,omitempty"`

	// 网络平均写入数据速度
	NetworkWriteAvg *float64 `json:"networkWriteAvg,omitempty"`

	// 网络平均写入数据速度检查点
	NetworkWriteAvgCheckPoint *float64 `json:"networkWriteAvgCheckPoint,omitempty"`

	// 网络平均写入数据速度检查结果
	NetworkWriteAvgCheckRes *bool `json:"networkWriteAvgCheckRes,omitempty"`

	// 网络最大写入数据速度检查点
	NetworkWriteCheckPoint *float64 `json:"networkWriteCheckPoint,omitempty"`

	// 网络最大写入数据速度检查结果
	NetworkWriteCheckRes *bool `json:"networkWriteCheckRes,omitempty"`

	// 峰值负载状态
	PeakLoadStatus *float64 `json:"peakLoadStatus,omitempty"`

	// 峰值负载状态检查点
	PeakLoadStatusCheckPoint *float64 `json:"peakLoadStatusCheckPoint,omitempty"`

	// 峰值负载状态检查结果
	PeakLoadStatusCheckRes *bool `json:"peakLoadStatusCheckRes,omitempty"`

	PeakMetric *PeakMetric `json:"peakMetric,omitempty"`

	// 工程ID
	ProjectId *string `json:"projectId,omitempty"`

	// 协议
	Protocols *[]string `json:"protocols,omitempty"`

	// 请求数
	Requests *int32 `json:"requests,omitempty"`

	// 用例结果
	Result *int32 `json:"result,omitempty"`

	// 用例结果日志
	ResultLog *string `json:"resultLog,omitempty"`

	// 执行轮次
	Round *int32 `json:"round,omitempty"`

	// 是否存储全量数据到CSS
	SaveAllData *bool `json:"saveAllData,omitempty"`

	// 服务ID
	ServiceId *string `json:"serviceId,omitempty"`

	// 阶段
	Stage *int32 `json:"stage,omitempty"`

	// 开始时间
	StartTime *string `json:"startTime,omitempty"`

	// 任务状态
	Status *int32 `json:"status,omitempty"`

	StreamingMediaVo *StreamingMediaReport `json:"streamingMediaVo,omitempty"`

	// 成功数
	SuccessCount *int32 `json:"successCount,omitempty"`

	// 成功率
	SuccessRate *int32 `json:"successRate,omitempty"`

	// 成功率检查点
	SuccessRateCheckPoint *float64 `json:"successRateCheckPoint,omitempty"`

	// 成功率检查结果
	SuccessRateCheckRes *bool `json:"successRateCheckRes,omitempty"`

	// 1XX响应码数量
	Sum1xx *int32 `json:"sum1xx,omitempty"`

	// 2XX响应码数量
	Sum2xx *int32 `json:"sum2xx,omitempty"`

	// 3XX响应码数量
	Sum3xx *int32 `json:"sum3xx,omitempty"`

	// 4XX响应码数量
	Sum4xx *int32 `json:"sum4xx,omitempty"`

	// 5XX响应码数量
	Sum5xx *int32 `json:"sum5xx,omitempty"`

	// 任务ID
	TaskId *string `json:"taskId,omitempty"`

	// 任务名
	TaskName *string `json:"taskName,omitempty"`

	// 任务项目ID
	TaskProjectId *string `json:"taskProjectId,omitempty"`

	// 任务状态
	TaskStatus *int32 `json:"taskStatus,omitempty"`

	// 用例基线uri
	TestCaseUri *string `json:"testCaseUri,omitempty"`

	// TP50
	Tp50 *int32 `json:"tp50,omitempty"`

	// TP50检查点
	Tp50CheckPoint *float64 `json:"tp50CheckPoint,omitempty"`

	// TP50检查结果
	Tp50CheckRes *bool `json:"tp50CheckRes,omitempty"`

	// TP75
	Tp75 *int32 `json:"tp75,omitempty"`

	// TP75检查点
	Tp75CheckPoint *float64 `json:"tp75CheckPoint,omitempty"`

	// TP75检查结果
	Tp75CheckRes *bool `json:"tp75CheckRes,omitempty"`

	// TP85
	Tp85 *int32 `json:"tp85,omitempty"`

	// TP85检查点
	Tp85CheckPoint *float64 `json:"tp85CheckPoint,omitempty"`

	// TP85检查结果
	Tp85CheckRes *bool `json:"tp85CheckRes,omitempty"`

	// TP90
	Tp90 *int32 `json:"tp90,omitempty"`

	// TP90检查点
	Tp90CheckPoint *float64 `json:"tp90CheckPoint,omitempty"`

	// TP90检查结果
	Tp90CheckRes *bool `json:"tp90CheckRes,omitempty"`

	// TP95
	Tp95 *int32 `json:"tp95,omitempty"`

	// TP95检查点
	Tp95CheckPoint *float64 `json:"tp95CheckPoint,omitempty"`

	// TP95检查结果
	Tp95CheckRes *bool `json:"tp95CheckRes,omitempty"`

	// TP99
	Tp99 *int32 `json:"tp99,omitempty"`

	// TP99.9
	Tp999 *int32 `json:"tp999,omitempty"`

	// TP99.99
	Tp9999 *int32 `json:"tp9999,omitempty"`

	// TP99.99检查点
	Tp9999CheckPoint *float64 `json:"tp9999CheckPoint,omitempty"`

	// TP99.99检查结果
	Tp9999CheckRes *bool `json:"tp9999CheckRes,omitempty"`

	// TP99.9检查点
	Tp999CheckPoint *float64 `json:"tp999CheckPoint,omitempty"`

	// TP99.9检查结果
	Tp999CheckRes *bool `json:"tp999CheckRes,omitempty"`

	// TP99检查点
	Tp99CheckPoint *float64 `json:"tp99CheckPoint,omitempty"`

	// TP99检查结果
	Tp99CheckRes *bool `json:"tp99CheckRes,omitempty"`

	// TPS
	Tps *float64 `json:"tps,omitempty"`

	// TPS检查点
	TpsCheckPoint *float64 `json:"tpsCheckPoint,omitempty"`

	// TPS检查结果
	TpsCheckRes *bool `json:"tpsCheckRes,omitempty"`

	// 平均TPS
	TranTPS *float64 `json:"tranTPS,omitempty"`

	// 平均TPS检查点
	TranTPSCheckPoint *float64 `json:"tranTPSCheckPoint,omitempty"`

	// 平均TPS检查结果
	TranTPSCheckRes *bool `json:"tranTPSCheckRes,omitempty"`

	// 事务ID
	TransactionId *string `json:"transactionId,omitempty"`

	// 事务成功数
	TransactionSuccess *float64 `json:"transactionSuccess,omitempty"`

	// 事务数
	Transactions *float64 `json:"transactions,omitempty"`

	// 事务数检查点
	TransactionsCheckPoint *float64 `json:"transactionsCheckPoint,omitempty"`

	// 事务数检查结果
	TransactionsCheckRes *bool `json:"transactionsCheckRes,omitempty"`

	// 更新时间
	UpdateTime *string `json:"updateTime,omitempty"`

	// aw的http url
	Url *string `json:"url,omitempty"`

	// 反应实时vuser数据
	UserConcur *int32 `json:"userConcur,omitempty"`

	// 分支uri
	VersionUri *string `json:"versionUri,omitempty"`
}

func (o CaseReportDetail) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CaseReportDetail struct{}"
	}

	return strings.Join([]string{"CaseReportDetail", string(data)}, " ")
}
