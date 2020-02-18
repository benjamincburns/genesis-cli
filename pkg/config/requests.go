package config

import (
	"fmt"
)

type URI struct {
	StopTestURI             string `mapstructure:"-"`
	StopDefURI              string `mapstructure:"-"`
	RunTestURI              string `mapstructure:"-"`
	CreateTestURI           string `mapstructure:"-"`
	ForkURI                 string `mapstructure:"-"`
	TestInfoURI             string `mapstructure:"-"`
	ExecURI                 string `mapstructure:"-"`
	MultipathUploadURI      string `mapstructure:"-"`
	GetOrgURI               string `mapstructure:"-"`
	LogsURI                 string `mapstructure:"-"`
	StatusURI               string `mapstructure:"-"`
	TestsURI                string `mapstructure:"-"`
	PrepareExecURI          string `mapstructure:"-"`
	AttachExecURI           string `mapstructure:"-"`
	RunDetachURI            string `mapstructure:"-"`
	ListContainersURI       string `mapstructure:"-"`
	GetSelfURI              string `mapstructure:"-"`
	GetOrgRoleURI           string `mapstructure:"-"`
	CheckAdminURI           string `mapstructure:"-"`
	CheckMemberURI          string `mapstructure:"-"`
	LimitsURI               string `mapstructure:"-"`
	CreateUserURI           string `mapstructure:"-"`
	BillingHealthURI        string `mapstructure:"-"`
	FeaturedOrgsURI         string `mapstructure:"-"`
	GetOrgProfileURI        string `mapstructure:"-"`
	UpdateOrgProfileURI     string `mapstructure:"-"`
	CreateOrgProfileURI     string `mapstructure:"-"`
	UpdateOrgFeaturedURI    string `mapstructure:"-"`
	UpdateUserSuperAdminURI string `mapstructure:"-"`
	ContainerLogsURI        string `mapstructure:"-"`
}

var (
	APIBase          = "/api/v1"
	TestExecutionAPI = APIBase + "/testexecution"
	RegistrarAPI     = APIBase + "/registrar"
	FilesAPI         = APIBase + "/files"
	LogsAPI          = APIBase + "/logs"
	ContainerAPI     = APIBase + "/container"
	BillingAPI       = APIBase + "/billing"
	AdminAPI         = APIBase + "/admin"
)

const Product = "prod_GgJOE3a7Adfopv"

var DefaultURI = URI{
	StopTestURI:   TestExecutionAPI + "/stop/test/%s",       //testid
	StopDefURI:    TestExecutionAPI + "/stop/definition/%s", //def id
	RunTestURI:    TestExecutionAPI + "/run/%s/%s",          //org def
	StatusURI:     TestExecutionAPI + "/status/%s",          //testid
	TestsURI:      TestExecutionAPI + "/organizations/%s/tests",
	ForkURI:       TestExecutionAPI + "/fork/%s/%s", //org def
	CreateTestURI: TestExecutionAPI + "/run/%s",     //org
	TestInfoURI:   TestExecutionAPI + "/info/tests/%s",

	GetOrgURI:           RegistrarAPI + "/organization/%s",
	GetSelfURI:          RegistrarAPI + "/user",                 //GET
	GetOrgRoleURI:       RegistrarAPI + "/organization/%s/user", //GET
	CheckAdminURI:       RegistrarAPI + "/check/iam/%s",         //GET
	CheckMemberURI:      RegistrarAPI + "/check/member/%s",      //GET
	CreateUserURI:       RegistrarAPI + "/user",                 //POST
	FeaturedOrgsURI:     RegistrarAPI + "/featured",
	GetOrgProfileURI:    RegistrarAPI + "/organization/%s/profile", //GET {org}
	UpdateOrgProfileURI: RegistrarAPI + "/organization/%s/profile", //PUT {org}
	CreateOrgProfileURI: RegistrarAPI + "/organization/%s/profile", //POST {org}

	UpdateOrgFeaturedURI:    AdminAPI + "/organization/%s/featured", // POST|DELETE {org}
	UpdateUserSuperAdminURI: AdminAPI + "/user/%s/admin",            // POST|DELETE {user}

	MultipathUploadURI: FilesAPI + "/organizations/%s/definitions",

	LogsURI: LogsAPI + "/data",

	PrepareExecURI:    ContainerAPI + "/%s/exec",        //POST {testid}
	AttachExecURI:     ContainerAPI + "/%s/exec/attach", //UPGRADE {testid}
	RunDetachURI:      ContainerAPI + "/%s/exec/run",    //POST {testid}
	ListContainersURI: ContainerAPI + "/%s/list",        //GET {testid}
	ContainerLogsURI:  ContainerAPI + "/%s/logs/%s/%s",  //GET {testid, container, lines}

	LimitsURI:        BillingAPI + "/limits/%s/%s", //GET {org, product}
	BillingHealthURI: BillingAPI + "/health",
}

func (uri URI) PrepareExecURL(tid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.PrepareExecURI, tid)
}

func (uri URI) ContainerLogsURL(tid, cntr, lines string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.ContainerLogsURI, tid, cntr, lines)
}

func (uri URI) AttachExecURL(tid string) string {
	return "http://exec-dev.whiteblock.io" + fmt.Sprintf(uri.AttachExecURI, tid)
}
func (uri URI) RunDetachURL(tid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.RunDetachURI, tid)
}
func (uri URI) ListContainersURL(tid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.ListContainersURI, tid)
}
func (uri URI) GetSelfURL() string       { return conf.APIEndpoint() + uri.GetSelfURI }
func (uri URI) CreateUserURL() string    { return conf.APIEndpoint() + uri.CreateUserURI }
func (uri URI) BillingHealthURL() string { return conf.APIEndpoint() + uri.BillingHealthURI }
func (uri URI) FeaturedOrgsURL() string  { return conf.APIEndpoint() + uri.FeaturedOrgsURI }
func (uri URI) GetOrgProfileURL(id string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.GetOrgProfileURI, id)
}
func (uri URI) UpdateOrgProfileURL(id string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.UpdateOrgProfileURI, id)
}
func (uri URI) CreateOrgProfileURL(id string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.CreateOrgProfileURI, id)
}
func (uri URI) UpdateOrgFeaturedURL(oid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.UpdateOrgFeaturedURI, oid)
}
func (uri URI) UpdateUserSuperAdminURL(uid string) string {
	return conf.APIEndpoint() + fmt.Sprintf(uri.UpdateUserSuperAdminURI, uid)
}
