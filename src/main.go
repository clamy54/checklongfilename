package main

/*

   checklongfilenames  -- check if filenames exceed 260 characters

   Copyright (C) 2022  Cyril LAMY

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.


*/

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const prgversion = "1.0"

func DisplayHelp(prgname string) {
	fmt.Printf("Checklongfilenames v%s \n", prgversion)
	fmt.Printf("This program detect if some filenames (with their absolute path) exceed 260 characters (Windows MAX_PATH standard size)\n")
	fmt.Println()
	fmt.Println("Usage :")
	fmt.Printf("%s  <directory_name> \n", filepath.Base(prgname))
	fmt.Println("  - <directory_name> : start the scan from this directory.")
	fmt.Println()
}

func Scandir(path string) (err error) {
	var cptfiles, cptdirs, warnings int64 = 0, 0, 0
	var filename string

	dirtree := []string{}
	// Search for all files and sub-directories
	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		dirtree = append(dirtree, path)
		info, _ := os.Stat(path)
		if info.IsDir() {
			cptdirs++
		} else {
			cptfiles++
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Check for filenames > 260 chars
	for _, filename = range dirtree {
		if len(filename) > 260 {
			fmt.Printf("%s (%d chars)\n", filename, len(filename))
			warnings++
		}
	}

	// Scan summary
	fmt.Printf("%v files and %v directories scanned - %v files have an absolute filename length > 260 characters \n", cptfiles, cptdirs, warnings)
	return nil

}

func main() {

	if runtime.GOOS == "windows" {
		var lencmdline int
		var cmdline string = ""

		lencmdline = len(os.Args)

		if len(os.Args[1:]) > 0 {
			// add all command line arguments in one variable to deal with spaces in filenames
			// if command line arguments are unquoted (ex: c:\program files\ )
			for cpt, arg := range os.Args[1:] {
				cmdline = cmdline + arg
				if cpt != len(os.Args[1:]) {
					cmdline = cmdline + " "
				}
			}
		}

		if (lencmdline < 2) || (strings.TrimSpace(cmdline) == "/h") {
			DisplayHelp(os.Args[0])
		} else {
			// Get absolute path
			path, err := filepath.Abs(filepath.Clean(cmdline))
			if err != nil {
				fmt.Printf("Error : %s\n", err.Error())
			}

			info, err := os.Stat(path)

			if !os.IsNotExist(err) && info.IsDir() {
				// The specified directory exists

				fmt.Printf("Scanning directory %s for long path filenames ...\n", path)
				err = Scandir(path)
				if err != nil {
					fmt.Printf("Error : %s \n", err.Error())
				}
			} else {
				// The directory did not exists

				fmt.Printf("Error : %s is not a directory or did not exists\n", path)
			}
		}
	} else {
		fmt.Println("Error : this utility is made for running on Microsoft Windows Platform ! ")
	}
}
