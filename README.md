# docs-gen

```go
package main

import (
	"github.com/bakito/docs-gen/docs"
	"github.com/bakito/docs-gen/pkg/cli"
	"github.com/bakito/docs-gen/pkg/yaml"
)

const (
	cliStartMarker  = "<!-- cli-doc-start -->"
	cliEndMarker    = "<!-- cli-doc-end -->"
	envStartMarker  = "<!-- env-doc-start -->"
	envEndMarker    = "<!-- env-doc-end -->"
	yamlStartMarker = "<!-- yaml-doc-start -->"
	yamlEndMarker   = "<!-- yaml-doc-end -->"
)

func main() {
	docs.UpdateDocumentation("README.md",
env.UpdateDocumentation[types.Config](envStartMarker, envEndMarker),
cli.UpdateDocumentation(cliStartMarker, cliEndMarker, ".", "go", "run", ".", "--help"),
		yaml.UpdateDocumentation[Config](yamlStartMarker, yamlEndMarker),
	)
}

type Config struct {
	Cron                string        `docs:"Cron expression for the sync interval"                                          env:"CRON"                json:"cron,omitempty"              yaml:"cron,omitempty"`
	RunOnStart          bool          `docs:"Run the sync on startup"                                                        env:"RUN_ON_START"        json:"runOnStart,omitempty"        yaml:"runOnStart,omitempty"`
	PrintConfigOnly     bool          `docs:"Print current config only and stop the application"                             env:"PRINT_CONFIG_ONLY"   json:"printConfigOnly,omitempty"   yaml:"printConfigOnly,omitempty"`
	ContinueOnError     bool          `docs:"Continue sync on errors"                                                        env:"CONTINUE_ON_ERROR"   json:"continueOnError,omitempty"   yaml:"continueOnError,omitempty"`
	ClientTimeoutString string        `docs:"Define a custom http client timeout ^([0-9]+(\\.[0-9]+)?(ns|us|µs|ms|s|m|h))+$" env:"HTTP_CLIENT_TIMEOUT" json:"httpClientTimeout,omitempty" yaml:"httpClientTimeout,omitempty" faker:"oneof: 30s, 5m"`
	ClientTimeout       time.Duration `                                                                                                                json:"-"                           yaml:"-"`
	// Origin adguardhome instance
	Origin *AdGuardInstance `docs:"Origin instance" json:"origin" yaml:"origin"`
	// One single replica adguardhome instance
	Replica *AdGuardInstance `docs:"Single or replica instance (don't use in combination with replicas')" json:"replica,omitempty" yaml:"replica,omitempty"`
	// Multiple replica instances
	Replicas []AdGuardInstance `docs:"List or replica instances (don't use in combination with replicas')" faker:"slice_len=2" json:"replicas,omitempty" yaml:"replicas,omitempty"`
	API      API               `                                                                                               json:"api,omitempty"      yaml:"api,omitempty"`
	Features Features          `                                                                                               json:"features,omitempty" yaml:"features,omitempty"`
}
```

## README.md

```markdown
## env

<!-- env-doc-start -->
| Name | Type | Description |
| :--- | ---- |:----------- |
<!-- env-doc-end -->

## yaml

<!-- yaml-doc-start -->
...
<!-- yaml-doc-end -->

## cli

<!-- cli-doc-start -->
...
<!-- cli-doc-end -->
```

