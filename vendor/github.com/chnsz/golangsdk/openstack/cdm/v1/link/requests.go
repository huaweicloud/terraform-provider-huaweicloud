package link

import "github.com/chnsz/golangsdk"

const (
	GenericJdbcConnector   = "generic-jdbc-connector"
	ObsConnector           = "obs-connector"
	HdfsConnector          = "hdfs-connector"
	HbaseConnector         = "hbase-connector"
	HiveConnector          = "hive-connector"
	SftpConnector          = "sftp-connector"
	FtpConnector           = "ftp-connector"
	MongodbConnector       = "mongodb-connector"
	RedisConnector         = "redis-connector"
	KafkaConnector         = "kafka-connector"
	DisConnector           = "dis-connector"
	ElasticsearchConnector = "elasticsearch-connector"
	DliConnector           = "dli-connector"
	OpentsdbConnector      = "opentsdb-connector"
	DmsKafkaConnector      = "dms-kafka-connector"

	// thirdparty-obs-connector is deprecated now
	ThirdpartyObsConnector = "thirdparty-obs-connector"
)

type LinkCreateOpts struct {
	Links []Link `json:"links" required:"true"`
}

type Link struct {
	Name string `json:"name" required:"true"`
	// Connector name. The mappings between the connectors and links are as follows:
	// generic-jdbc-connector: link to a relational database;
	// obs-connector: link to OBS;
	// hdfs-connector: link to HDFS;
	// hbase-connector: link to HBase and link to CloudTable;
	// hive-connector: link to Hive;
	// ftp-connector/sftp-connector: link to an FTP or SFTP server;
	// mongodb-connector: link to MongoDB;
	// redis-connector: link to Redis;
	// kafka-connector: link to Kafka;
	// dis-connector: link to DIS;
	// elasticsearch-connector: link to Elasticsearch/Cloud Search Service;
	// dli-connector: link to DLI;
	// opentsdb-connector: link to CloudTable OpenTSDB;
	// dms-kafka-connector: link to DMS Kafka
	ConnectorName    string      `json:"connector-name" required:"true"`
	LinkConfigValues LinkConfigs `json:"link-config-values" required:"true"`
	Enabled          *bool       `json:"enabled,omitempty"`
	CreationUser     string      `json:"creation-user,omitempty"`
	CreationDate     *int        `json:"creation-date,omitempty"`
	UpdateUser       string      `json:"update-user,omitempty"`
	UpdateDate       *int        `json:"update-date,omitempty"`
}

type LinkConfigs struct {
	Configs []Configs `json:"configs,omitempty"`
}

type Configs struct {
	Inputs []Input `json:"inputs" required:"true"`
	Name   string  `json:"name" required:"true"`
}

type Input struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Type  string      `json:"type,omitempty"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func Create(c *golangsdk.ServiceClient, clusterId string, opts LinkCreateOpts) (*LinkCreateResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst LinkCreateResponse
	_, err = c.Post(createLinkURL(c, clusterId), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Get(c *golangsdk.ServiceClient, clusterId string, linkName string) (*LinkDetail, error) {
	var rst golangsdk.Result
	_, err := c.Get(showLinkURL(c, clusterId, linkName), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r LinkDetail
		rst.ExtractInto(&r)
		return &r, nil
	}
	return nil, err
}

func Update(c *golangsdk.ServiceClient, clusterId string, linkName string, opts LinkCreateOpts) (*LinkUpdateResponse, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst LinkUpdateResponse
	_, err = c.Put(updateLinkURL(c, clusterId, linkName), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func Delete(c *golangsdk.ServiceClient, clusterId string, linkName string) (*LinkDeleteResponse, error) {
	var rst LinkDeleteResponse
	_, err := c.DeleteWithResponse(deleteLinkURL(c, clusterId, linkName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}
