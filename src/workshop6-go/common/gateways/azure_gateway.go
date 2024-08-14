package gateway

import (
	"context"
	"fmt"
	"os"
	"time"

	glog "github.com/golang/glog"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureGateway struct {
	AccountName string
	AccountKey  string
	Credential  *azblob.SharedKeyCredential
	Client      *azblob.Client
}

func (gateway *AzureGateway) Init() (*AzureGateway, error) {
	credential, err := azblob.NewSharedKeyCredential(gateway.AccountName, gateway.AccountKey)
	if err != nil {
		return nil, err
	}
	gateway.Credential = credential
	return gateway, nil
}

func (gateway *AzureGateway) InitBlobClient() (*azblob.Client, error) {
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/", gateway.AccountName)
	client, err := azblob.NewClientWithSharedKeyCredential(blobURL, gateway.Credential, nil)
	if err != nil {
		return nil, err
	}
	gateway.Client = client
	return client, nil
}

func (gateway *AzureGateway) GetBlobClient(options *azblob.ClientOptions) (*azblob.Client, error) {
	if gateway.Credential == nil {
		_, err := gateway.Init()
		if err != nil {
			return nil, err
		}
	}

	return gateway.InitBlobClient()
}

func (gateway *AzureGateway) DownloadBuffer(containerName string, blobName string, rangeStart int64, rangeEnd int64) ([]byte, error) {
	client, clientErr := gateway.GetBlobClient(nil)
	if clientErr != nil {
		return nil, clientErr
	}

	downloadBuffer := make([]byte, rangeEnd-rangeStart)
	downloadBufferOptions := azblob.DownloadBufferOptions{
		Range: azblob.HTTPRange{
			Offset: rangeStart,
			Count:  rangeEnd - rangeStart,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, downloadErr := client.DownloadBuffer(ctx, containerName, blobName, downloadBuffer, &downloadBufferOptions)
	if downloadErr != nil {
		return nil, downloadErr
	}
	return downloadBuffer, nil
}

func (gateway *AzureGateway) UploadBuffer(containerName string, blobName string, data []byte) error {
	client, clientErr := gateway.GetBlobClient(nil)
	if clientErr != nil {
		return clientErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, uploadErr := client.UploadBuffer(ctx, containerName, blobName, data, nil)
	if uploadErr != nil {
		return uploadErr
	}
	return nil
}

// DownloadFile downloads a blob file from Azure Blob Storage and returns a pointer to the file. Remember to close the file after!
func (gateway *AzureGateway) DownloadFile(containerName string, blobName string, rangeStart int64, rangeEnd int64) (*os.File, error) {
	client, clientErr := gateway.GetBlobClient(nil)
	if clientErr != nil {
		return nil, clientErr
	}

	file, err := os.Create(blobName + "_" + fmt.Sprint(rangeStart) + "_" + fmt.Sprint(rangeEnd) + ".txt")
	if err != nil {
		return nil, err
	}
	counter := 0
	downloadFileOptions := azblob.DownloadFileOptions{
		Range: azblob.HTTPRange{
			Offset: rangeStart,
			Count:  rangeEnd - rangeStart,
		},
		BlockSize: int64(1024 * 1024),
		Progress: func(bytesTransferred int64) {
			counter++
			if counter%20 == 0 {
				glog.Infof("Downloaded %d of %d bytes for %s/%s\n", bytesTransferred, rangeEnd-rangeStart, containerName, blobName)
			}
		},
		Concurrency: uint16(20),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, downloadErr := client.DownloadFile(ctx, containerName, blobName, file, &downloadFileOptions)
	if downloadErr != nil {
		return nil, downloadErr
	}
	return file, nil
}

// UploadFile uploads a file to Azure Blob Storage
func (gateway *AzureGateway) UploadFile(containerName string, blobName string, filePath string) error {
	client, clientErr := gateway.GetBlobClient(nil)
	if clientErr != nil {
		return clientErr
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	_, uploadErr := client.UploadFile(ctx, containerName, blobName, file, nil)
	if uploadErr != nil {
		return uploadErr
	}
	return nil
}
