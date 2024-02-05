package archeryapi

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"strings"
)

const (
	// SubmitUrl is the url to submit a workflow
	SubmitUrl = "/autoreview/"
	// SubmitOtherInstance is the url to submit a workflow to another instance
	SubmitOtherInstance = "/submitotherinstance/"
)

var validate *validator.Validate

type WorkflowService interface {
	Submit(request *WorkflowSubmitRequest) error
}

type WorkflowClient struct {
	apiClient *Client
}

// Submit submits a workflow
func (c WorkflowClient) Submit(request *WorkflowSubmitRequest) error {
	validate = validator.New()
	err := validate.Struct(request)
	if err != nil {
		return err
	}
	token, err := c.getMiddlewareToken()
	if err != nil {
		return err
	}
	request.CsrfMiddlewareToken = token
	marshal, err := json.Marshal(request)
	if err != nil {
		return err
	}
	params := map[string]string{}
	err = json.Unmarshal(marshal, &params)
	if err != nil {
		return err
	}
	_, err = c.apiClient.httpClient.R().
		SetFormData(params).
		Post(SubmitUrl)
	if err != nil {
		return err
	}
	return nil
}

// getMiddlewareToken gets the middleware token from html form
func (c WorkflowClient) getMiddlewareToken() (string, error) {
	resp, err := c.apiClient.httpClient.R().
		Get(SubmitOtherInstance)
	if err != nil {
		return "", err
	}
	csrfMiddlewareToken := resp.String()[strings.Index(resp.String(), "csrfmiddlewaretoken")+len("csrfmiddlewaretoken")+len(` value="`):]
	csrfMiddlewareToken = csrfMiddlewareToken[1:strings.Index(csrfMiddlewareToken, `">`)]
	return csrfMiddlewareToken, nil

}

type WorkflowSubmitRequest struct {
	CsrfMiddlewareToken string `json:"csrfmiddlewaretoken"`
	WorkflowId          string `json:"workflow_id"`
	SqlContent          string `json:"sql_content" validate:"required"`
	SqlUpload           string `json:"sql-upload"`
	WorkflowName        string `json:"workflow_name" validate:"required"`
	DemandUrl           string `json:"demand_url"`
	GroupName           string `json:"group_name" validate:"required"`
	InstanceName        string `json:"instance_name" validate:"required"`
	DbName              string `json:"db_name" validate:"required"`
	IsBackup            string `json:"is_backup" validate:"required"`
	RunDateStart        string `json:"run_date_start"`
	RunDateEnd          string `json:"run_date_end"`
	WorkflowAuditors    string `json:"workflow_auditors" validate:"required"`
}
