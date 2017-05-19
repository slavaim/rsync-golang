/*
* Copyright(c) 2017 Slava Imameev
* BSD license
 */

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
$ sudo rsync --daemon
$ rsync -crv rsync://nouser@127.0.0.1/kb/userkb/  /work/temp/test/kb
*/

var debug = false
var debugoutput = true

func rsyncClient(remoteIPorName string, remotePathSrc string, localPathDst string) error {

	// the terminating / at the remote path is important, it suppress creating a new directory in destination
	// if it already exists
	serverConnectString := fmt.Sprintf("rsync://nouser@%s/kb/%s/", remoteIPorName, remotePathSrc)

	cmd := exec.Command("/usr/bin/rsync", "-crv", serverConnectString, localPathDst)

	// show rsync's output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if debug || debugoutput {
		fmt.Printf("Debug: \n%v\n", cmd.Args)
		if debug {
			return nil
		}
	}

	fmt.Printf("Running %v.\n", cmd.Args)
	return cmd.Run()
}

// the server must run with superuser privileges
func rsyncServer(localPath string) error {

	rsyncdFormat := "lock file = /var/run/rsync.lock\n" +
		"log file = /var/log/rsyncd.log\n" +
		"pid file = /var/run/rsyncd.pid\n" +
		"\n" +
		"[kb]\n" +
		"    path = %s\n" +
		"    read only = yes\n" +
		"    list = yes"

	rsyncdConf := fmt.Sprintf(rsyncdFormat, localPath)

	fo, err := os.Create("/etc/rsyncd.conf")
	if err != nil {
		return err
	}
	defer fo.Close()

	cmd := exec.Command("/usr/bin/rsync", "--daemon")

	// show rsync's output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if debug || debugoutput {
		fmt.Printf("Debug:\n%v\n", cmd.Args)
		fmt.Printf("Debug:\n%s\n", rsyncdConf)
		if debug {
			return nil
		}
	}

	_, err = io.Copy(fo, strings.NewReader(rsyncdConf))
	if err != nil {
		return err
	}

	fmt.Printf("Running %v.\n", cmd.Args)
	return cmd.Run()
}

func main() {

	err := rsyncClient("127.0.0.1", "user1KB", "/work/temp/test/kb")
	//err := rsyncServer("/work/test/kb")
	if err != nil {
		log.Fatal(err)
	}
}
