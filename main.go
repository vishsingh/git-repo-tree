package main

import "fmt"
import "io/ioutil"
import "os"
import "os/exec"

type DirectoryClassification int

const (
	NotGitDirectory DirectoryClassification = iota
	GitAutoCommitDirectory
	GitDirtyDirectory
	GitCleanDirectory
)

func IsRegular(m os.FileMode) bool {
	return m & os.ModeType == 0
}

func ClassifyDirectory(path string) DirectoryClassification {
	// it is git if:
	// .git is a directory
	// .git/HEAD is a file
	// git show-ref succeeds

	files, err := ioutil.ReadDir(path + "/.git")
	if err != nil {
		return NotGitDirectory
	}

	foundHEAD := false
	for i := 0; i < len(files); i++ {
		file := files[i]
		if IsRegular(file.Mode()) && file.Name() == "HEAD" {
			foundHEAD = true
			break
		}
	}
	if !foundHEAD {
		return NotGitDirectory
	}

	showrefcmd := exec.Command("/usr/bin/git", "show-ref")
	showrefcmd.Dir = path
	runerr := showrefcmd.Run()
	if runerr != nil {
		return NotGitDirectory
	}

	// this is where you do more analysis
	return GitCleanDirectory
}

func ProcessDirectory(full_path string, depth int) {
	// the directory must be classified as either
	// - not git at all
	// - git, in which case the recursion must stop
	//   - git-auto-commit
	//   - dirty
	//   - clean

	if ClassifyDirectory(full_path) == GitCleanDirectory {
		fmt.Printf("found git directory: %s\n", full_path)
	}
}

func RecurseInto(full_path string, depth int) {
	files, err := ioutil.ReadDir(full_path)
	if err != nil {
		return
	}

	ProcessDirectory(full_path, depth)

	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			RecurseInto(full_path + "/" + files[i].Name(), depth + 1)
		}
	}
}

func main() {
	// this program starts at the current directory
	// and recurses down, in sorted order, to find all the git repos

	RecurseInto(".", 0)
}
