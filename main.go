package main

import (
	"bufio"
	"embed"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
)

//go:embed templates/*.yml templates/configs/*.conf
var templatesFS embed.FS

// struct to store key templates for docker-compose.yml file
type TemplateFile struct {
	Name string							`yaml:"name"`
	Services map[string]interface{}		`yaml:"services"`
	Volumes map[string]interface{}		`yaml:"volumes"`
}


// struct to store options for user to select from docker-compose.yml file
type ServiceOptions struct {
	Key string							`yaml:"key"`
	Label string						`yaml:"label"`
	Services map[string]interface{}		`yaml:"services"`
	Volumes map[string]interface{}		`yaml:"volumes"`
}



func scaffoldExtraFiles(selected []ServiceOptions){

	reader := bufio.NewReader(os.Stdin)
	for _, s := range selected {
		if s.Key == "nginx"{
			targetConfig := "configs/nginx.conf"

			
			if _, err := os.Stat(targetConfig); err == nil{
				
				fmt.Print("Nginx config file already exists (" + targetConfig + "). Overwrite? (y/n) [default: n]: ")

				confirm , err := reader.ReadString('\n')

				if err != nil {
					fmt.Println("Error reading input: " , err)
					continue
				}

				confirm = strings.ToLower(strings.TrimSpace(confirm))

				if confirm != "y" && confirm != "yes" {
					fmt.Println("Keeping existing " + targetConfig)
					continue
				}

			}

			confData , err := templatesFS.ReadFile("templates/configs/nginx.conf")
			if err != nil {
				fmt.Println("Error reading template file: " , err)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(targetConfig) , 0755); err != nil {
				fmt.Println("Error creating directory: " , err)
				continue
			}

			if err := os.WriteFile(targetConfig , confData , 0644); err != nil {
				fmt.Println("Error writing config file: " , err)
				continue
			}

			fmt.Println("Generated nginx.conf in " + targetConfig)

			if err := os.MkdirAll("html" , 0755); err != nil {
				fmt.Println("Error creating directory: " , err)
				continue
			}

			
		}

		
	}
}

// this function will load all template from templates folder
func loadTemplate()([]ServiceOptions , error){
	entries , err := templatesFS.ReadDir("templates")

	if err != nil {
		return nil , err
	}


	var options []ServiceOptions

	for _, entry := range entries {
		if entry.IsDir(){
			continue
		}

		data , err := templatesFS.ReadFile("templates/" + entry.Name())
		if err != nil {
			continue
		}

		var tpl TemplateFile

		// unmarshal template file into tpl struct
		if err := yaml.Unmarshal(data , &tpl); err != nil {
			fmt.Printf("Warning: failed to parse template %s: %v\n", entry.Name(), err)
			continue
		}

		// extract key and label
		ext := filepath.Ext(entry.Name())
		key := strings.TrimSuffix(entry.Name() , ext)
		label := tpl.Name

		
		if label == "" {
			label = key
		}

		// load options by mapping to ServiceOptions struct
		options = append(options , ServiceOptions{
			Key : key,
			Label : label,
			Services : tpl.Services,
			Volumes : tpl.Volumes,
		})



		
	}

	return options, nil
}


func main(){

	// load templates
	templates , err := loadTemplate()

	if err != nil || len(templates) == 0 {
		println("No templates found")
		os.Exit(1)
	}

	// init read user input for selecting templates
	reader := bufio.NewReader(os.Stdin)

	// print available templates
	fmt.Println("Available templates:")
	for i , tpl := range templates {
		fmt.Printf("%d. %s\n" , i + 1 , tpl.Label)
	}

	fmt.Print("\nEnter choices: ")


	// read user input 
	input , _ := reader.ReadString('\n')


	input = strings.TrimSpace(input)

	if input == "" {
		println("No choices selected")
		os.Exit(1)
	}

	var selected []ServiceOptions

	// if user input is all , select all templates
	if strings.ToLower(input ) == "all"{
		selected = templates
	} else {
		// if user input is comma separated , split and select templates
		indices := strings.Split(input , ",")
		for _, indexStr := range indices{
			idx , err := strconv.Atoi(strings.TrimSpace(indexStr))
			if err != nil || idx < 1 || idx > len(templates) {
				println("Invalid choice: " + indexStr)
				os.Exit(1)
			}

			selected = append(selected , templates[idx - 1])
		}
	}

	if len(selected) == 0 {
		println("No valid choices selected")
		os.Exit(1)
	}

	// init docker compose file
	targetFile := "docker-compose.yml"

	if _ , err := os.Stat(targetFile); err == nil {
		println("docker-compose.yml already exists. Overwrite? (y/n): ")
		confirm , _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
			println("Aborted.")
			os.Exit(0)
		}
	}

	// init compose map with services
	composeMap := map[string]interface{}{
		"services": make(map[string]interface{}),
	}
	
	// init volume map
	volumeMap := make(map[string]interface{})

	// loop through selected templates and add them to compose map
	for _, s := range selected{
		composeMap["services"].(map[string]interface{})[s.Key] = s.Services

		if s.Volumes != nil {
			for volKey, volVal := range s.Volumes {
				volumeMap[volKey] = volVal
			}
		}
	}

	if len(volumeMap) > 0 {
		composeMap["volumes"] = volumeMap
	}
	
	// call the function if nginx is chosen by the user
	scaffoldExtraFiles(selected)

	// marshal compose map to yaml format
	outData , err := yaml.Marshal(composeMap)

	if err != nil {
		println("Failed to marshal docker-compose.yml: " + err.Error())
		os.Exit(1)
	}

	// write docker compose file with chosen services
	if err := os.WriteFile(targetFile , outData , 0644); err != nil {
		println("Failed to write docker-compose.yml: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("docker-compose.yml generated successfully in " + targetFile )


	
}


