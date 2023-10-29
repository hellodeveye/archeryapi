package archeryapi

import (
	"encoding/json"
	"errors"

	"github.com/fatih/structs"
)

type DatabaseService interface {
	Query(sql string, opt ...QueryOption) (DatabaseQueryResponse, error)
}

type DatabaseClient struct {
	apiClient *Client
}

type QueryOption func(*QueryRequest)

func (c *DatabaseClient) Query(sql string, opt ...QueryOption) (DatabaseQueryResponse, error) {
	request := &QueryRequest{}
	for _, o := range opt {
		o(request)
	}
	request.SQLContent = sql
	params := map[string]string{}
	for k, v := range structs.Map(request) {
		params[k] = string(v.(string))
	}

	r, err := c.apiClient.httpClient.R().
		SetFormData(params).
		Post("/query/")
	if err != nil {
		return DatabaseQueryResponse{}, err
	}
	var result Result
	result.Data = &DatabaseQueryResponse{}

	if err := json.Unmarshal(r.Body(), &result); err != nil {
		return DatabaseQueryResponse{}, err
	}
	if result.Status != 0 {
		return DatabaseQueryResponse{}, errors.New(result.Msg)
	}
	return *(result.Data.(*DatabaseQueryResponse)), nil
}

type QueryRequest struct {
	InstanceName string `structs:"instance_name"`
	DBName       string `structs:"db_name"`
	SchemaName   string `structs:"schema_name"`
	TBName       string `structs:"tb_name"`
	SQLContent   string `structs:"sql_content"`
	LimitNum     string `structs:"limit_num"`
}

func WithInstanceName(i string) QueryOption {
	return func(c *QueryRequest) {
		c.InstanceName = i
	}
}

func WithDBName(d string) QueryOption {
	return func(c *QueryRequest) {
		c.DBName = d
	}
}

func WithSchemaName(d string) QueryOption {
	return func(c *QueryRequest) {
		c.SchemaName = d
	}
}

func WithTBName(d string) QueryOption {
	return func(c *QueryRequest) {
		c.TBName = d
	}
}

func WithSQLContent(d string) QueryOption {
	return func(c *QueryRequest) {
		c.SQLContent = d
	}
}

func WithLimitNum(d string) QueryOption {
	return func(c *QueryRequest) {
		c.LimitNum = d
	}
}

type DatabaseQueryResponse struct {
	FullSql             string          `json:"full_sql"`
	IsExecute           bool            `json:"is_execute"`
	Checked             interface{}     `json:"checked"`
	IsMasked            bool            `json:"is_masked"`
	QueryTime           float64         `json:"query_time"`
	MaskRuleHit         bool            `json:"mask_rule_hit"`
	MaskTime            string          `json:"mask_time"`
	Warning             interface{}     `json:"warning"`
	Error               interface{}     `json:"error"`
	IsCritical          bool            `json:"is_critical"`
	Rows                [][]interface{} `json:"rows"`
	ColumnList          []string        `json:"column_list"`
	Status              interface{}     `json:"status"`
	AffectedRows        int             `json:"affected_rows"`
	SecondsBehindMaster interface{}     `json:"seconds_behind_master"`
}
