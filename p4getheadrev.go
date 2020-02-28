package main

//	F.Demurger 2019-04
//  	1 args:
//			- depot file path and name
//
//      Option -v version
//      Option -u <user>
//
//      Returns the revision number of the head, status > 0 if there was an error
//
//
//	cross compilation AMD64:  env GOOS=windows GOARCH=amd64 go build p4getheadrev.go

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	var versionFlg bool
	var user string
	var p4Cmd string // p4 command path
	var err error

	const usageVersion   = "Display Version"
	const usageUser      = "specify a username"

  // Have to create a specific set, the default one is poluted by some test stuff from another lib (?!)
  checkFlags := flag.NewFlagSet("check", flag.ExitOnError)

	checkFlags.BoolVar(&versionFlg, "version", false, usageVersion)
	checkFlags.BoolVar(&versionFlg, "v", false, usageVersion + " (shorthand)")
	checkFlags.StringVar(&user, "user", "", usageUser)
	checkFlags.StringVar(&user, "u", "", usageUser + " (shorthand)")
	checkFlags.Usage = func() {
        fmt.Printf("Usage: %s <opt> <depot filename>\nReturns revision of head or status > 0 if err\n",os.Args[0])
        checkFlags.PrintDefaults()
    }

    // Check parameters
	checkFlags.Parse(os.Args[1:])

	if versionFlg {
		fmt.Printf("Version %s\n", "2020-02  v1.0.3")
		os.Exit(0)
	}

	// Check presence of p4 cli
	if p4Cmd, err = exec.LookPath("p4"); err != nil {
		fmt.Printf("P4 command line is not installed - %s\n", err)
		os.Exit(1)
	}

  // Parse the command parameters
  index := len(os.Args)
	depotFile :=  os.Args[index - 1]

	var out []byte
	if len(user) > 0 {
		// fmt.Printf(p4Cmd + " -u " + user + " files " + " " + depotFile + "\n")
		out, err = exec.Command(p4Cmd, "-u", user, "files", depotFile).Output()
	} else {
		// fmt.Printf(p4Cmd + " files " + depotFile + "\n")
		out, err = exec.Command(p4Cmd, "files", depotFile).Output()
	}
	if err != nil {
		fmt.Printf("P4 command line error - %s\n", err)
		os.Exit(2)
	}

	// Read version
	// e.g. //Project/dev/localization/afile_bulgarian.txt#8 - edit change 4924099 (utf16)
	idxBeg := strings.LastIndex(string(out),"#") + len("#")
	idxEnd := strings.LastIndex(string(out)," - ")
	// Check response to prevent out of bound index
	if idxBeg == -1 || idxEnd == -1 || idxBeg >= idxEnd {
		fmt.Printf("Format error in P4 response: %s\n", string(out))
		os.Exit(3)
	}
	// sRev := string(out[strings.LastIndex(string(out),"#") + len("#"):strings.LastIndex(string(out)," - ")])
	sRev := string(out[idxBeg:idxEnd])

	rev, err := strconv.Atoi(sRev) // Check format
	if err != nil {
		fmt.Printf("sRev=%s\n", sRev)
		fmt.Printf("Format err=%s\n", err)
		os.Exit(4)
	}

	fmt.Printf("%d\n",rev)
}
