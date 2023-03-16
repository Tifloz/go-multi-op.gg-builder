package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Parsing des arguments de la ligne de commande
	filename := flag.String("file", "", "nom du fichier csv")
	flag.Parse()

	// Vérification que le nom du fichier a été fourni
	if *filename == "" {
		fmt.Println("Veuillez spécifier le nom du fichier csv en utilisant l'argument -file.")
		return
	}

	// Ouverture du fichier csv
	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Création d'un lecteur csv pour le fichier
	reader := csv.NewReader(file)

	// On ignore la première ligne
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	// Création d'une liste de noms d'utilisateurs
	var usernames []string

	// Itération sur chaque ligne du fichier csv
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}

		// Extraction des données de la ligne
		opGG := line[2]

		// Si la troisième colonne est vide, on ignore cette ligne
		if opGG == "" {
			continue
		}

		// Extraction du nom d'utilisateur à partir de l'URL
		u, err := url.Parse(opGG)
		if err != nil {
			panic(err)
		}

		var username string
		if strings.Contains(opGG, "euw.op.gg/summoners") {
			// Si l'URL est de la forme https://www.op.gg/summoner/userName=Kus%C3%B4%20Gaki
			// on extrait le nom d'utilisateur à partir de la chaîne de requête
			username = u.Query().Get("userName")
		} else {
			// Si l'URL est de la forme https://www.op.gg/summoners/euw/nathelis
			// on extrait le nom d'utilisateur à partir du dernier segment du chemin
			path := u.Path
			parts := strings.Split(path, "/")
			username = parts[len(parts)-1]
		}

		// Ajout du nom d'utilisateur à la liste, en échappant les caractères réservés dans l'URL
		usernames = append(usernames, url.QueryEscape(username))
	}

	// Création du lien op.gg avec la liste de noms d'utilisateurs
	query := "summoners=" + strings.Join(usernames, ",")
	cleanedLink := strings.Replace(query, "userName%3D", "", -1)
	fmt.Println(query)
	fmt.Println(cleanedLink)

	// Encodage de l'URL pour qu'elle soit valide
	link := "https://www.op.gg/multisearch/euw?" + url.PathEscape(cleanedLink)
	// Aff
	fmt.Println(link)
}
