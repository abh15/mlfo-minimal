package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

//struct to make sense of response received from db
type response struct {
	id           string
	uri          string
	accessType   string
	trainingTime string
	accuracy     string
	resourceReq  string
}

func main() {
	arg := os.Args[1]

	m := make(map[string]interface{})
	accessType := ""
	trainingTime := ""
	resourceReq := ""
	yamlFile, err := ioutil.ReadFile(arg) //Read yaml file
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//iterate over the yaml structure and decide model requirements based on use case
	for k, v := range m {
		if k == "model" {
			vnew, _ := v.(map[interface{}]interface{})
			if vnew["usecase"] == "edge" {
				accessType = "private"
				trainingTime = "low"
				resourceReq = "low"
				fmt.Println("Selecting model for Edge use case...")
				fmt.Println("\nModel selection requirements: accessType = " + accessType + "\t trainingTime = " + trainingTime + "\t resourceRequirements = " + resourceReq)

			} else if vnew["usecase"] == "cloud" {
				accessType = "public"
				trainingTime = "high"
				resourceReq = "high"
				fmt.Println("Selecting model for Cloud use case...")
				fmt.Println("\nModel selection requirements: accessType = " + accessType + "\t trainingTime = " + trainingTime + "\t resourceRequirements = " + resourceReq)

			} else {
				log.Printf("Invalid use case. Now exiting")
				return
			}

		}
	}

	//start connection to the db
	db, err := sql.Open("mysql", "root:mlfo1234@tcp(127.0.0.1:3306)/modelrepo")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished
	defer db.Close()

	fmt.Println("\nQuerying model repository for model...")

	// perform a db Query
	query := "SELECT * FROM models WHERE accessType='" + accessType + "' AND trainingTime='" + trainingTime + "' AND resourceReq='" + resourceReq + "';"
	result, err := db.Query(query)

	// if there was error while querying, handle it
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var r response
		err = result.Scan(&r.id, &r.uri, &r.accessType, &r.trainingTime, &r.accuracy, &r.resourceReq)
		if err != nil {
			panic(err.Error())
		}
		// print received model information
		fmt.Println("\nReceived model...")
		fmt.Printf("\nModel id: %s\n", r.id)
		fmt.Printf("Model URI: %s\n", r.uri)
		fmt.Printf("Model accessType: %s\n", r.accessType)
		fmt.Printf("Model trainingTime: %s\n", r.trainingTime)
		fmt.Printf("Model accuracy: %s\n", r.accuracy)
		fmt.Printf("Model resourceRequirements: %s\n", r.resourceReq)

	}

	defer result.Close()

}
