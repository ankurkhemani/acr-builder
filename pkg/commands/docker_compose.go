package commands

import (
	"fmt"

	build "github.com/Azure/acr-builder/pkg"
	"github.com/Azure/acr-builder/pkg/constants"
	"github.com/Azure/acr-builder/pkg/grok"
)

// dockerComposeSupportedFilenames needs to always be in sync with docker compose's SUPPORTED_FILENAMES in config.py
var dockerComposeSupportedFilenames = []string{
	"docker-compose.yml",
	"docker-compose.yaml",
}

// errNoDefaultDockerfile means that no default docker file is found when running FindDefaultDockerComposeFile
var errNoDefaultDockerfile = fmt.Errorf("No default docker-compose file found")

// findDefaultDockerComposeFile try and locate the default docker-compose file in the current working directory
func findDefaultDockerComposeFile(runner build.Runner) (string, error) {
	fs := runner.GetFileSystem()
	for _, defaultFile := range dockerComposeSupportedFilenames {
		exists, err := fs.DoesFileExist(defaultFile)
		if err != nil {
			return "", fmt.Errorf("Unexpected error while checking for default docker compose file: %s", err)
		}
		if exists {
			return defaultFile, nil
		}
	}
	return "", errNoDefaultDockerfile
}

// NewDockerComposeBuild creates a build target with defined docker-compose file
func NewDockerComposeBuild(path, projectDir string, buildArgs []string) build.Target {
	return &dockerComposeBuildTask{
		path:             path,
		projectDirectory: projectDir,
		buildArgs:        buildArgs,
	}
}

type dockerComposeBuildTask struct {
	path             string
	projectDirectory string
	buildArgs        []string
}

func (t *dockerComposeBuildTask) ScanForDependencies(runner build.Runner) ([]build.ImageDependencies, error) {
	env := runner.GetContext()
	var targetPath string
	if t.path != "" {
		targetPath = env.Expand(t.path)
	} else {
		var err error
		targetPath, err = findDefaultDockerComposeFile(runner)
		if err != nil {
			return []build.ImageDependencies{}, err
		}
	}
	return grok.ResolveDockerComposeDependencies(env, env.Expand(t.projectDirectory), targetPath)
}

func (t *dockerComposeBuildTask) Build(runner build.Runner) error {
	args := []string{}
	if t.path != "" {
		args = append(args, "-f", t.path)
	}
	args = append(args, "build")

	if t.projectDirectory != "" {
		args = append(args, "--project-directory", t.projectDirectory)
	}

	for _, buildArg := range t.buildArgs {
		args = append(args, "--build-arg", buildArg)
	}

	return runner.ExecuteCmd("docker-compose", args)
}

func (t *dockerComposeBuildTask) Export() []build.EnvVar {
	return []build.EnvVar{
		{
			Name:  constants.ExportsDockerComposeFile,
			Value: t.path,
		},
	}
}

func (t *dockerComposeBuildTask) Push(runner build.Runner) error {
	args := []string{}
	if t.path != "" {
		args = append(args, "-f", t.path)
	}
	args = append(args, "push")
	return runner.ExecuteCmd("docker-compose", args)
}
