package cmds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"
	"github.com/whiteblock/genesis-cli/pkg/validate"

	"github.com/spf13/cobra"
	"github.com/whiteblock/definition/schema"
	yaml "gopkg.in/yaml.v2"
)

type compose struct {
	Version  string
	Networks map[string]network
	Volumes  map[string]volume
	Services map[string]_service
}

type network struct{}

type volume struct {
	Driver, External string
	DriverOpts       map[string]string `yaml:"driver_opts"`
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

type _service struct {
	ContainerName                     string `yaml:"container_name"`
	Image                             string
	Networks, Ports, Volumes, Command []string
	VolumesFrom                       []string `yaml:"volumes_from"`
	DependsOn                         []string `yaml:"depends_on"` //build tree
	CapAdd                            []string `yaml:"cap_add"`
	Build                             build
	Environment                       map[string]string
	Deploy                            deploy
}

func mkSystemComponent(name string, serv _service) schema.SystemComponent {
	sys := schema.SystemComponent{
		Type:         name,
		PortMappings: serv.Ports,
		Count:        serv.Deploy.Replicas,
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

func mkService(volumes map[string]volume, name string, serv _service) schema.Service {
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

var composeCmd = &cobra.Command{
	Use:     "compose <file>",
	Short:   "Convert docker compose into a spec",
	Long:    `Convert docker compose into a spec`,
	Aliases: []string{},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		file, err := os.Open(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			util.ErrorFatal(err)
		}

		var comp compose
		err = yaml.Unmarshal(data, &comp)
		if err != nil {
			util.ErrorFatal(err)
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
							util.ErrorFatalf("dependency %s does not exist", dep)
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
			json.Unmarshal([]byte(`"infinite"`), &test.Timeout)
			if len(phases) > 1 {
				test.Phases = phases[1:]
			}
			root.Tests = []schema.Test{test}
		}
		data, err = yaml.Marshal(root)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print(string(data))

		res, err := validate.Schema(data)
		if err != nil {
			util.ErrorFatal(err)
		}
		if !res.Valid() {
			for _, schemaErr := range res.Errors() {
				util.Error(schemaErr.String())
			}
			os.Exit(1)
		}
		_, _, err = service.ProcessDefinitionFromBytes(data)
		if err != nil {
			util.ErrorFatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
