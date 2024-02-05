package archeryapi

import (
	"encoding/json"
	"errors"

	"github.com/fatih/structs"
	"github.com/mcuadros/go-defaults"
)

type DatabaseService interface {
	Query(sql string, request *QueryRequest) (DatabaseQueryResponse, error)
}

type DatabaseClient struct {
	apiClient *Client
}

func (c *DatabaseClient) Query(sql string, request *QueryRequest) (DatabaseQueryResponse, error) {
	request.SQLContent = sql
	defaults.SetDefaults(request)
	params := map[string]string{}
	for k, v := range structs.Map(request) {
		params[k] = v.(string)
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
	LimitNum     string `structs:"limit_num" default:"100"`
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
