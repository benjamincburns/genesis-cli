package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/whiteblock/genesis-cli/pkg/util"

	yamlC "github.com/ghodss/yaml"
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/schema"
	yaml "gopkg.in/yaml.v2"
)

type compose struct {
	Version  string
	Networks map[string]interface{}
	Volumes  map[string]interface{}
	Services map[string]_service
}

type build struct {
	Context    string
	Dockerfile string
}

type limits struct {
	CPUs   string
	Memory string
}
type resources struct {
	Limits limits
}

type deploy struct {
	Resources resources
	Replicas  int64
}

type environment map[string]string

func (env *environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var arr []string
	err := unmarshal(&arr)
	if err != nil {
		return unmarshal(env)
	}
	(*env) = environment{}
	for i := range arr {
		kv := strings.SplitN(arr[i], "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("Invalid environment entry %s", arr[i])
		}
		(*env)[kv[0]] = kv[1]
	}
	return nil
}

type _service struct {
	ContainerName                     string `yaml:"container_name"`
	Image                             string
	Networks, Ports, Volumes, Command []string
	VolumesFrom                       []string `yaml:"volumes_from"`
	DependsOn                         []string `yaml:"depends_on"` //build tree
	CapAdd                            []string `yaml:"cap_add"`
	Build                             build
	Environment                       environment
	Deploy                            deploy
}

func mkSystemComponent(name string, serv _service) schema.SystemComponent {
	sys := schema.SystemComponent{
		Name:         name,
		Type:         name,
		PortMappings: serv.Ports,
		Count:        serv.Deploy.Replicas,
	}
	if sys.Count == 0 {
		sys.Count = 1
	}
	for i := range sys.PortMappings {
		if !strings.Contains(sys.PortMappings[i], ":") {
			sys.PortMappings[i] = sys.PortMappings[i] + ":" + sys.PortMappings[i]
		}
	}

	for _, network := range serv.Networks {
		sys.Resources.Networks = append(sys.Resources.Networks, schema.Network{Name: network})
	}

	return sys
}

func mkService(volumes map[string]interface{}, name string, serv _service) schema.Service {
	out := schema.Service{
		Name:        name,
		Image:       serv.Image,
		Environment: serv.Environment,
	}
	if serv.Command != nil {
		out.Script.Inline = strings.Join(serv.Command, " ")
	}
	for _, v := range serv.Volumes {
		pieces := strings.Split(v, ":")
		if _, exists := volumes[pieces[0]]; exists {
			out.Volumes = append(out.Volumes, schema.Volume{
				Name:  pieces[0],
				Path:  pieces[1],
				Scope: schema.GlobalScope,
			})
		} else if stat, err := os.Lstat(pieces[0]); err == nil { //do I need stuff from there?
			if !stat.Mode().IsRegular() && !stat.IsDir() {
				util.Errorf("cannot upload %s, not a normal file", pieces[0])
				continue
			}
			out.InputFiles = append(out.InputFiles, schema.InputFile{
				SourcePath:      pieces[0],
				DestinationPath: pieces[1],
			})
		} else {
			out.Volumes = append(out.Volumes, schema.Volume{
				Name:  pieces[0],
				Path:  pieces[1],
				Scope: schema.LocalScope,
			})
		}

	}
	out.Resources.Memory = serv.Deploy.Resources.Limits.Memory
	return out
}

func SchemaFromCompose(data []byte) (schema.RootSchema, error) {
	var comp compose
	err := yaml.Unmarshal(data, &comp)
	if err != nil {
		return schema.RootSchema{}, err
	}
	root := schema.RootSchema{}
	//Register the services
	for name, serv := range comp.Services {
		root.Services = append(root.Services, mkService(comp.Volumes, name, serv))
	}
	//build the dependency tree
	phases := []schema.Phase{}
	have := map[string]int{}

	for len(have) < len(comp.Services) {
		phase := schema.Phase{}
		for name, serv := range comp.Services {
			if _, ok := have[name]; ok {
				continue
			}
			if len(serv.DependsOn) == 0 {
				have[name] = len(phases)
				phase.System = append(phase.System, mkSystemComponent(name, serv))
				continue
			}
			broken := false
			for _, dep := range serv.DependsOn {
				val, haveit := have[dep]
				if !haveit || val >= len(phases) {
					//check if it actually exists
					if _, exists := comp.Services[dep]; !exists {
						return schema.RootSchema{}, fmt.Errorf("dependency %s does not exist", dep)
					}
					broken = true
					break
				}
			}
			if broken {
				continue
			}
			have[name] = len(phases)
			phase.System = append(phase.System, mkSystemComponent(name, serv))
		}
		phases = append(phases, phase)
	}
	for i := range phases {
		phases[i].Name = fmt.Sprintf("phase%d", i)
	}
	if len(phases) > 0 {
		test := schema.Test{
			Name:        "compose",
			System:      phases[0].System,
			Description: "This was auto-generated from a docker compose file",
		}
		if len(phases) > 1 {
			test.Phases = phases[1:]
			test.Phases[len(test.Phases)-1].Duration = command.InfiniteDuration
		} else {
			test.WaitFor = command.InfiniteDuration
		}
		root.Tests = []schema.Test{test}
	}
	return root, nil
}

func SchemaYAMLFromCompose(data []byte) ([]byte, error) {
	root, err := SchemaFromCompose(data)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(root)
}

func MustSchemaYAMLFromCompose(data []byte) []byte {
	out, err := SchemaYAMLFromCompose(data)
	if err != nil {
		util.ErrorFatal(err)
	}
	return out
}

func MustSchemaJSONFromCompose(data []byte) []byte {

	out, err := yamlC.YAMLToJSON(MustSchemaYAMLFromCompose(data))
	if err != nil {
		util.ErrorFatal(err)
	}
	return out
}
