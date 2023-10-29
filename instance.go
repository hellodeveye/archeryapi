package archeryclient

import (
	"encoding/json"
	"strconv"

	"github.com/fatih/structs"
)

type ResourceType string

const (
	DatabaseResouceType ResourceType = "database"
	TableResourceType   ResourceType = "table"
	ColumnResourceType  ResourceType = "column"
)

type InstanceClient struct {
	apiClient *Client
}

type InstanceService interface {
	GetInstances() ([]Instance, error)
	GetInstanceResource(req *InstanceResourceQueryRequest) ([]string, error)
	Describetable(req *InstanceResourceQueryRequest) (DescribeTable, error)
}

func (c *InstanceClient) GetInstances() ([]Instance, error) {
	r, err := c.apiClient.httpClient.R().
		SetQueryParam("tag_codes[]", "can_read").
		Get("/group/user_all_instances/")
	if err != nil {
		return nil, err
	}
	var result Result
	instances := make([]Instance, 0)
	result.Data = &instances

	if err := json.Unmarshal(r.Body(), &result); err != nil {
		return nil, err
	}
	return instances, nil
}

func (c *InstanceClient) GetInstanceResource(req *InstanceResourceQueryRequest) ([]string, error) {
	params := map[string]string{}
	for k, v := range structs.Map(req) {
		params[k] = string(v.(string))
	}
	r, err := c.apiClient.httpClient.R().
		SetQueryParams(params).
		Get("/instance/instance_resource/")
	if err != nil {
		return nil, err
	}
	var result Result
	instances := make([]string, 0)
	result.Data = &instances

	if err := json.Unmarshal(r.Body(), &result); err != nil {
		return nil, err
	}
	return instances, nil
}

func (c *InstanceClient) Describetable(req *InstanceResourceQueryRequest) (DescribeTable, error) {
	r, err := c.apiClient.httpClient.R().
		SetFormData(map[string]string{
			"instance_name": req.InstanceName,
			"db_name":       req.DbName,
			"schema_name":   req.SchemaName,
			"tb_name":       req.TbName,
			"resource_type": req.ResourceType.String(),
		}).
		Post("/instance/describetable/")
	if err != nil {
		return DescribeTable{}, err
	}
	var result Result
	var describetable DescribeTable
	result.Data = &describetable
	if err := json.Unmarshal(r.Body(), &result); err != nil {
		return DescribeTable{}, err
	}
	return describetable, nil
}

type InstanceResourceQueryRequest struct {
	InstanceName string       `json:"instance_name"`
	DbName       string       `json:"db_name"`
	SchemaName   string       `json:"schema_name"`
	TbName       string       `json:"tb_name"`
	ResourceType ResourceType `json:"resource_type"`
}

type Instance struct {
	DbType       string `json:"db_type"`
	Id           int    `json:"id"`
	InstanceName string `json:"instance_name"`
	Type         string `json:"type"`
}

// Instance ToString
func (i Instance) String() string {
	return i.InstanceName + " " + i.DbType + " " + i.Type + " " + strconv.Itoa(i.Id)
}

func (r ResourceType) String() string {
	return string(r)
}

type DescribeTable struct {
	FullSql      string     `json:"full_sql"`
	IsExecute    bool       `json:"is_execute"`
	Checked      string     `json:"checked"`
	IsMasked     bool       `json:"is_masked"`
	QueryTime    string     `json:"query_time"`
	MaskRuleHit  bool       `json:"mask_rule_hit"`
	MaskTime     string     `json:"mask_time"`
	Warning      string     `json:"warning"`
	Error        string     `json:"error"`
	IsCritical   bool       `json:"is_critical"`
	Rows         [][]string `json:"rows"`
	ColumnList   []string   `json:"column_list"`
	Status       string     `json:"status"`
	AffectedRows int        `json:"affected_rows"`
}

func (d DescribeTable) String() string {
	return d.FullSql + " " + strconv.FormatBool(d.IsExecute) + " " + d.Checked + " " + strconv.FormatBool(d.IsMasked) + " " + d.QueryTime + " " + strconv.FormatBool(d.MaskRuleHit) + " " + d.MaskTime + " " + d.Warning + " " + d.Error + " " + strconv.FormatBool(d.IsCritical) + " " + d.Status + " " + strconv.Itoa(d.AffectedRows)
}
