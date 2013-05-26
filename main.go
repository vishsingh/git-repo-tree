package main

import "fmt"
import "io/ioutil"
import "os"
import "os/exec"
import "bytes"
import "strings"

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

	return ClassifyGitDirectory(path)
}

func ClassifyGitDirectory(path string) DirectoryClassification {
	// auto-commit if git symbolic-ref -q HEAD returns the same as the contents of .git-auto-commit

	gitautocommit, err := ioutil.ReadFile(path + "/.git-auto-commit")
	if err == nil {
		symrefcmd := exec.Command("/usr/bin/git", "symbolic-ref", "-q", "HEAD")
		symrefcmd.Dir = path
		
		symrefout, err := symrefcmd.Output()
		if err == nil {
			if bytes.Equal(gitautocommit, symrefout) {
				return GitAutoCommitDirectory
			}
		}
	}

	// clean if git status --porcelain is empty

	statuscmd := exec.Command("/usr/bin/git", "status", "--porcelain")
	statuscmd.Dir = path

	statuscmdout, err := statuscmd.Output()
	if err == nil {
		if bytes.Equal(statuscmdout, nil) {
			return GitCleanDirectory
		}
	}

	return GitDirtyDirectory
}

// returns whether we should continue iterating
func ProcessDirectory(full_path string, depth int) bool {
	// the directory must be classified as either
	// - not git at all
	// - git, in which case the recursion must stop
	//   - git-auto-commit
	//   - dirty
	//   - clean

	dirclass := ClassifyDirectory(full_path)

	if dirclass == GitCleanDirectory {
		fmt.Printf("found git directory: %s\n", full_path)
	}

	if dirclass == GitAutoCommitDirectory {
		fmt.Printf("found git-auto-commit directory: %s\n", full_path)
	}

	if dirclass == GitDirtyDirectory {
		fmt.Printf("found DIRTY git directory: %s\n", full_path)
	}

	return dirclass == NotGitDirectory
}

func RecurseInto(full_path string, depth int, callback func(string, int) bool) {
	files, err := ioutil.ReadDir(full_path)
	if err != nil {
		return
	}

	shouldContinue := callback(full_path, depth)
	if !shouldContinue {
		return
	}

	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			RecurseInto(full_path + "/" + files[i].Name(), depth + 1, callback)
		}
	}
}

func AnyGitDirUnder(path string, gitDirs map[string]DirectoryClassification) bool {
	for gitDir, _ := range gitDirs {
		if gitDir == path {
			return true
		}

		if strings.HasPrefix(gitDir, path + "/") {
			return true
		}
	}

	return false
}

func main() {
	// this program starts at the current directory
	// and recurses down, in sorted order, to find all the git repos

	gitDirs := make(map[string]DirectoryClassification)

	collector := func (path string, depth int) bool {
		dirclass := ClassifyDirectory(path)

		if dirclass != NotGitDirectory {
			gitDirs[path] = dirclass
		}

		return dirclass == NotGitDirectory
	}

	RecurseInto(".", 0, collector)

	RecurseInto(".", 0, ProcessDirectory)
}
