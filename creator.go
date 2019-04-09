package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "io"
    "bufio"
    "log"
    "regexp"
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
	var arg = os.Args[1]
	creating_folder(arg + "/examples")
    creating_folder(arg + "/test")


    // making test/automatically_generated_test.go
    duplication_files(arg + "/test/automatically_generated" + "_test.go", "templates/test")


	// making examples/main.tf
    duplication_files(arg + "/examples/main.tf", "templates/examples_main_header")
	text_addition(arg + "/examples/main.tf",   "\nmodule \"example_module\" {\n")
    vars_list := get_namelist("variable \"", "\" {", arg + "/variables.tf")
    var bigest_len = 0
    for i := 0; i < len(vars_list); i++ {
    	if bigest_len < len(vars_list[i]) {
    		bigest_len = len(vars_list[i])
    	}
    }
	for i := 0; i < len(vars_list); i++ {
		text_addition(arg + "/examples/main.tf", "  " + vars_list[i] + spase_generator(bigest_len - len(vars_list[i])) + " = \"${var." + vars_list[i] + "}\"\n")
	}
	text_addition(arg  + "/examples/main.tf", "  source" + spase_generator(bigest_len - len("source")) + " = \"../\"\n}")


	// making examples/outputs.tf
    duplication_files(arg + "/examples/outputs.tf", "templates/examples_outputs")
	outputs_list := get_namelist("output \"", "\" {", arg + "/outputs.tf")
	for i := 0; i < len(outputs_list); i++ {
		text_addition(arg + "/examples/outputs.tf", "\noutput \"" + outputs_list[i] + "\" {\n  value = \"${module.example_module." + outputs_list[i] + "}\"\n}\n")
	}


	// making examples/variables.tf
    duplication_files(arg + "/examples/variables.tf", arg + "/variables.tf")
    text_addition(arg  + "/examples/variables.tf", reading_file("templates/examples_variables"))


    // correcting .gitignore
    text_addition(arg  + "/.gitignore", reading_file("templates/gitignore"))


    // correcting .travis.yml
    text_addition(arg  + "/.travis.yml",    reading_file("templates/travis"))
}

func spase_generator(count int) string {
	var spases = ""
	for i := 0; i < count; i++ {
		spases = spases + " "
	}
	return spases
}

func creating_folder(path string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        fmt.Println("folder " + path + " created")
        os.MkdirAll(path, os.ModePerm)
    } else {
        fmt.Println("folder " + path + " already exists")
    }
}

func duplication_files(file_path string, template_path string) {
    template, err := ioutil.ReadFile(template_path)
    check(err)
    err = ioutil.WriteFile(file_path, template, 0644)
    check(err)
}

func get_namelist(prefix string, postfix string, file_path string) []string {
    var arr_words []string
    r, _ := regexp.Compile("^" + prefix + "[A-z]*" + postfix + "$")

    file, err := os.Open(file_path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if r.FindString(scanner.Text()) != "" {

            str := r.FindString(scanner.Text())
            teststr := strings.Replace(str,     prefix,  "", -len(prefix) )
            teststr  = strings.Replace(teststr, postfix, "", -len(postfix))

            arr_words = append(arr_words, teststr)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    return arr_words
}

func text_addition(file_path string, postfix string) {
    editable_file, err := os.OpenFile(file_path, os.O_APPEND|os.O_WRONLY, 0600)
    check(err)

    defer editable_file.Close()

    if _, err = editable_file.WriteString(postfix); err != nil {
        panic(err)
    }
}

func reading_file(file_path string) string {
    text_file := ""
    file, err := os.Open(file_path)
    check(err)
    defer file.Close()

    data := make([]byte, 64)

    for{
        n, err := file.Read(data)
        if err == io.EOF{
            break
        }
        text_file = text_file + string(data[:n])
    }
    return text_file
}
