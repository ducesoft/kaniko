/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"os"

	"github.com/ducesoft/kaniko/pkg/dockerfile"
	"github.com/ducesoft/kaniko/pkg/util"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/sirupsen/logrus"
)

type RunMarkerCommand struct {
	BaseCommand
	cmd   *instructions.RunCommand
	Files []string
}

func (r *RunMarkerCommand) ExecuteCommand(config *v1.Config, buildArgs *dockerfile.BuildArgs) error {
	// run command `touch filemarker`
	logrus.Debugf("using new RunMarker command")
	prevFilesMap, _ := util.GetFSInfoMap("/", map[string]os.FileInfo{})
	if err := runCommandInExec(config, buildArgs, r.cmd); err != nil {
		return err
	}
	_, r.Files = util.GetFSInfoMap("/", prevFilesMap)

	logrus.Debugf("files changed %s", r.Files)
	return nil
}

// String returns some information about the command for the image config
func (r *RunMarkerCommand) String() string {
	return r.cmd.String()
}

func (r *RunMarkerCommand) FilesToSnapshot() []string {
	return r.Files
}

func (r *RunMarkerCommand) ProvidesFilesToSnapshot() bool {
	return true
}

// CacheCommand returns true since this command should be cached
func (r *RunMarkerCommand) CacheCommand(img v1.Image) DockerCommand {

	return &CachingRunCommand{
		img:       img,
		cmd:       r.cmd,
		extractFn: util.ExtractFile,
	}
}

func (r *RunMarkerCommand) MetadataOnly() bool {
	return false
}

func (r *RunMarkerCommand) RequiresUnpackedFS() bool {
	return true
}

func (r *RunMarkerCommand) ShouldCacheOutput() bool {
	return true
}

func (r *RunMarkerCommand) ShouldDetectDeletedFiles() bool {
	return true
}
