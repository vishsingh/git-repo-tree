package main

//import "fmt"
import "io/ioutil"

func ProcessDirectory(full_path string) {
	
}

func RecurseInto(full_path string, depth int) {
	files, err := ioutil.ReadDir(full_path)
	if err != nil {
		return
	}

	ProcessDirectory(full_path)

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
