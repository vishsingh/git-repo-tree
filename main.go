package main

import "fmt"
import "io/ioutil"

func RecurseInto(full_path string, depth int) {
	files, err := ioutil.ReadDir(full_path)
	if err != nil {
		return
	}

	for i := 0; i < len(files); i++ {
		file_name := files[i].Name()

		indent := ""
		for j := 0; j < depth; j++ {
			indent += " "
		}

		fmt.Println(indent, "file", i, "is called", file_name)

		if files[i].IsDir() {
			RecurseInto(full_path + "/" + file_name, depth + 1)
		}
	}
}

func main() {
	// this program starts at the current directory
	// and recurses down, in sorted order, to find all the git repos

	RecurseInto(".", 0)
}
