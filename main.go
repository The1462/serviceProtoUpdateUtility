package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

/*type ServiceItems struct {
	Services []ServiceItem
}*/

type ServiceItem struct {
	Name string 	`json:"name"`
	Folder string 	`json:"folder"`
}

func main()  {

	jsonFileMainPaths, err := os.Open("mainServiceFiles.json")

	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("--Successfully Opened mainServiceFiles.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFileMainPaths.Close()


	byteValue, _ := ioutil.ReadAll(jsonFileMainPaths)

	// we initialize our Users array
	var serviceList = []ServiceItem{}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &serviceList)


	/*for i := 0; i < len(serviceList); i++ {
		fmt.Println("User Name: " + serviceList[i].Name)
		fmt.Println("User Folder: " + serviceList[i].Folder)
	}*/


	for i := 0; i < len(serviceList); i++ {
		// sourceProto = "./" + serviceList[i].Folder + "/" + serviceList[i].Name + "/" + serviceList[i].Name + ".proto"
		var sourceProto = serviceList[i].Folder + "/services/" + serviceList[i].Name + "/" + serviceList[i].Name + ".proto"
		var sourceIndex = serviceList[i].Folder + "/services/" + serviceList[i].Name + "/index.js"
		//fmt.Println(sourceProto)

		/*fmt.Println("User Name: " + serviceList[i].Name)
		fmt.Println("User Folder: " + serviceList[i].Folder)*/

		//Tum Dosyalari Listele
		projectFolders, err := ioutil.ReadDir("./")
		if err != nil {
			log.Fatal(err)
		}

		//Tum Dosyalari Tara
		for _, projectFolderItem := range projectFolders {
			//Dosya Bir Klasor Mu eger klasor ise muhtemelen proje klasoru
			if projectFolderItem.IsDir() && projectFolderItem.Name()!=serviceList[i].Folder {
				//Icinde Services Klasoru var mi
				isThereServiceFolder := false
				//Proje klasoru icindeki tum dosyalari cek
				projectInternalFolders, err2 := ioutil.ReadDir("./" + projectFolderItem.Name())
				if err2 != nil {
					log.Fatal(err2)
				}

				for _, projectInternalFolderItem := range projectInternalFolders {
					if projectInternalFolderItem.IsDir() && projectInternalFolderItem.Name()=="services" {
						isThereServiceFolder = true
					}
				}

				if isThereServiceFolder {
					currentProjectServices, err3 := ioutil.ReadDir("./" + projectFolderItem.Name() + "/services")
					if err3 != nil {
						log.Fatal(err3)
					}

					for _, currentProjectServiceItem := range currentProjectServices {
						if currentProjectServiceItem.IsDir() && currentProjectServiceItem.Name()==serviceList[i].Name {

							var targetFolder = "./" + projectFolderItem.Name() + "/services/" + currentProjectServiceItem.Name() + "/"
							var targetProto = targetFolder + currentProjectServiceItem.Name() + ".proto"

							ReplaceProtoFiles(sourceProto, sourceIndex, targetFolder, targetProto)
						}
					}
				}
			}
		}
	}
}

func ReplaceProtoFiles(sourceProto, sourceIndexPath, targetFolder, targetProto string) (bool, error) {
	os.Remove(targetProto)

	sourceFileStat, err2 := os.Stat(sourceProto)
	if err2 != nil {
		log.Fatal(err2)
	}

	//Checking File Mode
	if !sourceFileStat.Mode().IsRegular() {
		fmt.Errorf("%s is not a regular file", sourceProto)

		return false, nil
	} else {
		//Reading Proto File
		source, err3 := os.Open(sourceProto)
		if err3 != nil {
			log.Fatal(err3)
		}
		defer source.Close() //Close File When Process Ends


		//Proto Destination Create
		destination, err4 := os.Create(targetProto)
		if err4 != nil {
			log.Fatal(err4)
		}

		defer destination.Close() //Close Destination File When Process Ends
		_, err5 := io.Copy(destination, source) //Copy Proto File

		if err5 != nil {
			return false, err5
		} else {

			//Proto replaced

			//Check Index File
			sourceIndexFileStat, siSErr := os.Stat(sourceIndexPath)
			if siSErr != nil {
				log.Fatal(siSErr)
			}

			if !sourceIndexFileStat.Mode().IsRegular() {
				fmt.Errorf("%s is not a regular file", sourceIndexPath)

				return false, nil
			} else {

				//Read Source Index.js file
				sourceIndexF, siFErr := os.Open(sourceIndexPath)
				if siFErr != nil {
					log.Fatal(siFErr)
				}
				defer sourceIndexF.Close()


				//Index.js Destination Create
				destinationIndex, destinationIndexErr := os.Create(targetFolder + "index.js")
				if destinationIndexErr != nil {
					log.Fatal(destinationIndexErr)
				}
				defer destinationIndex.Close() //Close Destination File When Process Ends


				_, err6 := io.Copy(destinationIndex, sourceIndexF) //Copy Proto File

				if err6 != nil {
					return false, err6
				} else {
					//All Done
					return true, nil
				}


			}
		}
	}
}