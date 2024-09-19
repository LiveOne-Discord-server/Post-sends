package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

const (
	owner       = "BANSAFAn"
	repo        = "https://github.com/LiveOne-Discord-server/Post-sends"
	folderPath  = "./Posts-publics"
	githubToken = "J3C4N-3CC323-C32C2M-C32CAZ"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for {
		err := processFolder(ctx, client)
		if err != nil {
			log.Printf("Error processing folder: %v", err)
		}

		time.Sleep(5 * time.Minute) // Check every 5 minutes
	}
}

func processFolder(ctx context.Context, client *github.Client) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(folderPath, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file %s: %v", file.Name(), err)
			continue
		}

		err = createOrUpdateFile(ctx, client, file.Name(), content)
		if err != nil {
			log.Printf("Error uploading file %s: %v", file.Name(), err)
			continue
		}

		log.Printf("Successfully uploaded %s", file.Name())

		// Move the file to a 'processed' folder
		processedPath := filepath.Join(folderPath, "processed", file.Name())
		err = os.Rename(filePath, processedPath)
		if err != nil {
			log.Printf("Error moving file %s to processed folder: %v", file.Name(), err)
		}
	}

	return nil
}

func createOrUpdateFile(ctx context.Context, client *github.Client, fileName string, content []byte) error {
	fileContent := base64.StdEncoding.EncodeToString(content)
	message := "Update " + fileName

	_, _, err := client.Repositories.CreateFile(ctx, owner, repo, fileName, &github.RepositoryContentFileOptions{
		Message: &message,
		Content: []byte(fileContent),
	})

	if err != nil {
		// If the file already exists, try updating it
		file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, fileName, nil)
		if err != nil {
			return fmt.Errorf("error getting existing file: %v", err)
		}

		_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, fileName, &github.RepositoryContentFileOptions{
			Message: &message,
			Content: []byte(fileContent),
			SHA:     file.SHA,
		})

		if err != nil {
			return fmt.Errorf("error updating file: %v", err)
		}
	}

	return nil
}