/*
Copyright 2020 Kazım SARIKAYA

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

package term

import (
	"github.com/creack/pty"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
	"io"
	klog "k8s.io/klog/v2"
	"os"
	"os/exec"
	"os/signal"
)

func CreateTerminal() error {
	c := exec.Command("/bin/sh")

	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)

	defer func() {
		close(ch)
		_ = ptmx.Close()
	}()

	signal.Notify(ch, unix.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				klog.V(5).Error(err, "error resizing pty")
			}
		}
	}()
	ch <- unix.SIGWINCH

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}
