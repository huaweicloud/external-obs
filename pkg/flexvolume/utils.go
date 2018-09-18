/*
Copyright 2018 The Kubernetes Authors.

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

package flexvolume

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// formatResult is used to format
func formatResult(fields map[string]interface{}, err error) string {
	var data map[string]interface{}
	if err != nil {
		data = map[string]interface{}{
			"status":  "Failure",
			"message": err.Error(),
		}
	} else {
		data = map[string]interface{}{
			"status": "Success",
		}
		for k, v := range fields {
			data[k] = v
		}
	}
	s, err := json.Marshal(data)
	if err != nil {
		panic("error marshalling the data")
	}
	return string(s) + "\n"
}

// execIsMounted is used to judge
func execIsMounted(mountDir string) (string, error) {
	// CONTAINER = `docker ps --filter "label=mountpath=${mount_dir}" --format "{{.ID}}"`
	output, err := exec.Command("docker",
		"ps",
		"--filter",
		"label=mountpath="+mountDir,
		"--format",
		"{{.ID}}").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker ps failed: %v", err)
	}

	// return
	str := strings.Replace(string(output), "\n", "", -1)
	if str == "" {
		return "0", nil
	} else {
		return "1", nil
	}
}

// execMount defines
func execMount(mountDir, strOBSAccessKey, strOBSSecretKey, strOBSBucket, strOBSEndpoint string) error {
	// mkdir -p ${mount_dir}
	_, err := exec.Command("mkdir",
		"-p",
		mountDir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkdir failed: %v", err)
	}

	// DOCKER_OUT=`docker run -d --privileged -l mountpath=${mount_dir} -e OBSAccessKey=${OBSAccessKey} -e OBSSecretKey=${OBSSecretKey}
	// -v ${mount_dir}:/mnt/mountpoint:shared --cap-add SYS_ADMIN quay.io/huaweicloud/obs-mountvolume ${OBSBucket} /mnt/mountpoint
	// -o passwd_file=/etc/passwd-s3fs -o url=${OBSEndpoint} -d -d -f -o f2 -o curldbg`
	_, err = exec.Command("docker",
		"run",
		"-d",
		"--privileged",
		"-l",
		"mountpath="+mountDir,
		"-e",
		"OBSAccessKey="+strOBSAccessKey,
		"-e",
		"OBSSecretKey="+strOBSSecretKey,
		"-v",
		mountDir+":/mnt/mountpoint:shared",
		"--cap-add",
		"SYS_ADMIN",
		"quay.io/huaweicloud/obs-mountvolume",
		strOBSBucket,
		"/mnt/mountpoint",
		"-o",
		"passwd_file=/etc/passwd-s3fs",
		"-o",
		"url="+strOBSEndpoint,
		"-d",
		"-d",
		"-f",
		"-o",
		"f2",
		"-o",
		"curldbg").CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run failed: %v", err)
	}
	return nil
}

// execUnmount is used to judge
func execUnmount(mountDir string) error {
	// CONTAINER=`docker ps --filter "label=mountpath=${mount_dir}" --format "{{.ID}}"`
	output, err := exec.Command("docker",
		"ps",
		"--filter",
		"label=mountpath="+mountDir,
		"--format",
		"{{.ID}}").CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker ps failed: %v", err)
	}

	// docker rm ${CONTAINER} -f
	str := strings.Replace(string(output), "\n", "", -1)
	_, err = exec.Command("docker",
		"rm",
		str,
		"-f").CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker rm failed: %v", err)
	}

	// umount -l ${mount_dir}
	_, err = exec.Command("umount",
		"-l",
		mountDir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("umount failed: %v", err)
	}

	// rmdir ${mount_dir}
	_, err = exec.Command("rm",
		"-rf",
		mountDir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("rmdir failed: %v", err)
	}

	return nil
}
