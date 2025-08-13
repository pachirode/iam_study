package pumps

import "github.com/pachirode/iam_study/internal/pump/analytics"

type MongoType int

type CommonPumpConfig struct {
	filters               analytics.AnalyticsFilters
	timeout               int
	OmitDetailedRecording bool
}

type CSVConf struct {
	CSVDir string `mapstructure:"csv_dir"`
}

type BaseMongoConf struct {
	MongoURL                      string    `json:"mongo_url"                         mapstructure:"mongo_url"`
	MongoUseSSL                   bool      `json:"mongo_use_ssl"                     mapstructure:"mongo_use_ssl"`
	MongoSSLInsecureSkipVerify    bool      `json:"mongo_ssl_insecure_skip_verify"    mapstructure:"mongo_ssl_insecure_skip_verify"`
	MongoSSLAllowInvalidHostnames bool      `json:"mongo_ssl_allow_invalid_hostnames" mapstructure:"mongo_ssl_allow_invalid_hostnames"`
	MongoSSLCAFile                string    `json:"mongo_ssl_ca_file"                 mapstructure:"mongo_ssl_ca_file"`
	MongoSSLPEMKeyfile            string    `json:"mongo_ssl_pem_keyfile"             mapstructure:"mongo_ssl_pem_keyfile"`
	MongoDBType                   MongoType `json:"mongo_db_type"                     mapstructure:"mongo_db_type"`
}

type MongoConf struct {
	BaseMongoConf

	CollectionName            string `json:"collection_name"               mapstructure:"collection_name"`
	MaxInsertBatchSizeBytes   int    `json:"max_insert_batch_size_bytes"   mapstructure:"max_insert_batch_size_bytes"`
	MaxDocumentSizeBytes      int    `json:"max_document_size_bytes"       mapstructure:"max_document_size_bytes"`
	CollectionCapMaxSizeBytes int    `json:"collection_cap_max_size_bytes" mapstructure:"collection_cap_max_size_bytes"`
	CollectionCapEnable       bool   `json:"collection_cap_enable"         mapstructure:"collection_cap_enable"`
}
