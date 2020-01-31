package config

import (
	"fmt"
)

type URI struct {
	StopTestURI        string `mapstructure:"-"`
	StopDefURI         string `mapstructure:"-"`
	RunTestURI         string `mapstructure:"-"`
	CreateTestURI      string `mapstructure:"-"`
	ForkURI            string `mapstructure:"-"`
	TestInfoURI        string `mapstructure:"-"`
	ExecURI            string `mapstructure:"-"`
	MultipathUploadURI string `mapstructure:"-"`
	GetOrgURI          string `mapstructure:"-"`
	LogsURI            string `mapstructure:"-"`
	StatusURI          string `mapstructure:"-"`
	TestsURI           string `mapstructure:"-"`
	PrepareExecURI     string `mapstructure:"-"`
	AttachExecURI      string `mapstructure:"-"`
	RunDetachURI       string `mapstructure:"-"`
	ListContainersURI  string `mapstructure:"-"`
}

var (
	APIBase          = "/api/v1"
	TestExecutionAPI = APIBase + "/testexecution"
	RegistrarAPI     = APIBase + "/registrar"
	FilesAPI         = APIBase + "/files"
	LogsAPI          = APIBase + "/logs"
	ContainerAPI     = APIBase + "/container"
)

var DefaultURI = URI{
	StopTestURI:   TestExecutionAPI + "/stop/test/%s",       //testid
	StopDefURI:    TestExecutionAPI + "/stop/definition/%s", //def id
	RunTestURI:    TestExecutionAPI + "/run/%s/%s",          //org def
	StatusURI:     TestExecutionAPI + "/status/%s",          //testid
	TestsURI:      TestExecutionAPI + "/organizations/%s/tests",
	ForkURI:       TestExecutionAPI + "/fork/%s/%s", //org def
	CreateTestURI: TestExecutionAPI + "/run/%s",     //org
	TestInfoURI:   TestExecutionAPI + "/info/tests/%s",

	GetOrgURI: RegistrarAPI + "/organization/%s",

	MultipathUploadURI: FilesAPI + "/organizations/%s/definitions",

	LogsURI: LogsAPI + "/data",

	PrepareExecURI:    ContainerAPI + "/%s/exec",        //POST {testid}
	AttachExecURI:     ContainerAPI + "/%s/exec/attach", //UPGRADE {testid}
	RunDetachURI:      ContainerAPI + "/%s/exec/run",    //POST {testid}
	ListContainersURI: ContainerAPI + "/%s/list",        //GET {testid}
}

func (uri URI) PrepareExecURL(tid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.PrepareExecURI, tid)
}
func (uri URI) AttachExecURL(tid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.AttachExecURI, tid)
}
func (uri URI) RunDetachURL(tid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.RunDetachURI, tid)
}
func (uri URI) ListContainersURL(tid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.ListContainersURI, tid)
}
