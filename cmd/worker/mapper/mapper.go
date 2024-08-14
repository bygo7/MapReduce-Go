package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"unicode"
	gateway "workshop/common/gateways"
	pb "workshop/protos"

	"github.com/google/uuid"

	glog "github.com/golang/glog"
)

const (
	extraRangeToFindSpace = 30
)

type MapOutput struct {
	OutputFilePaths []string `json:"output_file_paths"`
	OutputFileBytes []int    `json:"output_file_bytes"`
}

func handleError(err error) {
	if err != nil {
		glog.Error(err.Error())
	}
}

func Map(azureGateway *gateway.AzureGateway, ctx context.Context, request *pb.MapRequest) (*pb.MapReply, error) {
	// Parse shards
	mapTaskNum := request.GetMapTaskNum()
	shards := request.GetShardInfos()
	reduceTasks := request.GetReduceTasks()
	containerNameGlobal := shards[0].GetContainerName()
	downloadedBlobs := make([]string, 0)
	keyValueInfos := make([]*pb.BlobInfo, 0)

	for _, shard := range shards {
		containerName, blobName := shard.GetContainerName(), shard.GetBlobName()
		shardRangeStart, shardRangeEnd := shard.GetRangeStart(), shard.GetRangeEnd()
		// Get a full word at the beginning if cut off and truncate the last word if cut off
		rangeStart, rangeEnd, err := GetRangeForFullWords(azureGateway, containerName, blobName, shardRangeStart, shardRangeEnd)
		handleError(err)

		// Download blob
		glog.Infof("Downloading shard of %s/%s for map %d", containerName, blobName, mapTaskNum)
		shardBlobFile, err := azureGateway.DownloadFile(containerName, blobName, rangeStart, rangeEnd)
		handleError(err)
		downloadedBlobs = append(downloadedBlobs, shardBlobFile.Name())
		defer shardBlobFile.Close()
	}

	// Perform map
	glog.Infof("Performing map task %d", mapTaskNum)
	mappedFilePaths, mappedFileBytes, err := MapBlob(downloadedBlobs, reduceTasks)
	handleError(err)

	// Create new blobs according to reduce task
	uuid := uuid.New().String()
	mapperContainerName := containerNameGlobal + "-mapped"
	for i := range mappedFilePaths {
		mappedFilePath := mappedFilePaths[i]
		mappedFileByteCount := mappedFileBytes[i]
		keyValueBlobName := "m" + fmt.Sprint(mapTaskNum) + "_r" + fmt.Sprint(i) + "_" + fmt.Sprint(uuid) + ".txt"

		glog.Infof("Uploading shard of %s for map %d, reduce %d", containerNameGlobal, mapTaskNum, i)
		uploadErr := azureGateway.UploadFile(mapperContainerName, keyValueBlobName, mappedFilePath)
		handleError(uploadErr)
		// Add to key value info
		keyValueInfo := &pb.BlobInfo{
			ContainerName: mapperContainerName,
			BlobName:      keyValueBlobName,
			RangeStart:    0,
			RangeEnd:      int64(mappedFileByteCount),
		}

		keyValueInfos = append(keyValueInfos, keyValueInfo)
	}

	return &pb.MapReply{
		Success:       true,
		MapTaskNum:    mapTaskNum,
		KeyValueInfos: keyValueInfos,
	}, nil
}

func MapBlob(inputFilePaths []string, reduceTasks int32) ([]string, []int, error) {
	/* Python API:
	Inputs:
		input_file_paths: str[]
		reduce_tasks: int
	Outputs:
		output_file_paths: str[]
		output_file_bytes: int[]
	*/
	inputFilePathsJson, err := json.Marshal(inputFilePaths)
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.Command("python3", "map.py", string(inputFilePathsJson), fmt.Sprint(reduceTasks))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil, err
	}
	glog.Infof("Map output: %s", string(out))

	var mapOutput MapOutput
	err = json.Unmarshal(out, &mapOutput)
	if err != nil {
		return nil, nil, err
	}

	return mapOutput.OutputFilePaths, mapOutput.OutputFileBytes, nil
}

// Get the full word from the beginning if cut off and truncate the last word if cut off
func GetRangeForFullWords(azureGateway *gateway.AzureGateway, containerName string, blobName string, rangeStart int64, rangeEnd int64) (int64, int64, error) {
	// Get a full word at the beginning by looking extra range before rangeStart
	wordRangeStart, wordRangeEnd := rangeStart, rangeEnd
	if rangeStart != 0 {
		// Find the nearest start of the word within rangeStart - extraRangeToFindSpace ~ rangeStart
		extraRangeBlob, err := azureGateway.DownloadBuffer(containerName, blobName, int64(math.Max(0, float64(rangeStart-extraRangeToFindSpace))), rangeStart)
		handleError(err)
		rangeForFullWord := 0
		for i := len(extraRangeBlob) - 1; i >= 0; i-- {
			if unicode.IsSpace(rune(extraRangeBlob[i])) {
				rangeForFullWord = len(extraRangeBlob) - i - 1
				break
			}
		}
		wordRangeStart = rangeStart - int64(rangeForFullWord)
	}

	// Truncate the full word at the end by trimming characters after the last space
	extraRangeTruncateBlob, err := azureGateway.DownloadBuffer(containerName, blobName, int64(math.Max(float64(rangeStart), float64(rangeEnd-extraRangeToFindSpace))), rangeEnd+1)
	handleError(err)
	// If the last character is a space, then we don't need to truncate
	if unicode.IsSpace(rune(extraRangeTruncateBlob[len(extraRangeTruncateBlob)-1])) {
		return wordRangeStart, wordRangeEnd, nil
	}
	// Iterate from end to beginning to find the nearest space
	for i := len(extraRangeTruncateBlob) - 1; i >= 0; i-- {
		if unicode.IsSpace(rune(extraRangeTruncateBlob[i])) {
			wordRangeEnd = rangeEnd - int64(len(extraRangeTruncateBlob)-i-1)
			break
		}
	}

	return wordRangeStart, wordRangeEnd, nil
}
