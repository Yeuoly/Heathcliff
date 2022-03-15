package runner

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Yeuoly/Heathcliff/types"
	"github.com/Yeuoly/Heathcliff/util"
)

func RunCompile(task *types.Task) (string, error) {
	switch task.Type {
	case types.TASK_TYPE_ACM_C:
		c_tmp_file_name := strings.Join([]string{
			"./tmp/",
			"acm_c_",
			strconv.FormatInt(time.Now().Unix(), 36),
			".c",
		},
			"",
		)

		out_tmp_file_name := c_tmp_file_name + ".out"

		err := os.WriteFile(c_tmp_file_name, task.Buf.Buf, os.ModeAppend)
		if err != nil {
			return "", errors.New("failed to write file")
		}

		_, err = runCmd("gcc", c_tmp_file_name, "-o", out_tmp_file_name, "-static")
		if err != nil {
			return "", errors.New("failed to compile C code")
		}

		if _, err := os.Stat(out_tmp_file_name); err != nil && !os.IsExist(err) {
			return "", errors.New("compile failed")
		}

		return out_tmp_file_name, nil
	}
	return "", nil
}

func RunExec(path string) (string, error) {
	chroot_path := "./chroot/" + util.RandStr(16)
	err := os.Mkdir(chroot_path, os.ModePerm)
	if err != nil {
		return "", err
	}

	_, err = runCmd("cp", path, chroot_path+"/exec")
	if err != nil {
		return "", err
	}

	return runCmd("chroot", chroot_path, "./exec")
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var err error

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return "", err
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	mreader := io.MultiReader(stdout, stderr)

	result, err := ioutil.ReadAll(mreader)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
