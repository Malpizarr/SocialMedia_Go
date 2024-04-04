package utils

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func UploadFileToBlobStorage(containerName, filename string, fileBytes []byte) (string, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", fmt.Errorf("error al crear credenciales: %w", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, err := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	if err != nil {
		return "", fmt.Errorf("error al parsear la URL: %w", err)
	}
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	blobURL := containerURL.NewBlockBlobURL(filename)
	ctx := context.Background()

	_, err = azblob.UploadBufferToBlockBlob(ctx, fileBytes, blobURL, azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "application/octet-stream",
		},
	})
	if err != nil {
		return "", fmt.Errorf("error al subir el archivo: %w", err)
	}

	blobURLString := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, filename)
	return blobURLString, nil
}
