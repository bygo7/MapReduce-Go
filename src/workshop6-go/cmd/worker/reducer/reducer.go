package reducer

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	gateway "workshop/common/gateways"
	pb "workshop/protos"

	glog "github.com/golang/glog"
)

type ReduceOutput struct {
	OutputFilePath  string `json:"output_file_path"`
	OutputFileBytes int    `json:"output_file_bytes"`
}

func handleError(err error) {
	if err != nil {
		glog.Error(err.Error())
	}
}

func Reduce(azureGateway *gateway.AzureGateway, ctx context.Context, request *pb.ReduceRequest) (*pb.ReduceReply, error) {
	reduceTaskNum := request.GetReduceTaskNum()
	keyValueInfos := request.GetKeyValueInfos()
	filePaths := make([]string, len(keyValueInfos))
	// Parse key value infos
	for i := range keyValueInfos {
		keyValueInfo := keyValueInfos[i]
		containerName, blobName := keyValueInfo.GetContainerName(), keyValueInfo.GetBlobName()
		rangeStart, rangeEnd := int64(keyValueInfo.GetRangeStart()), int64(keyValueInfo.GetRangeEnd())
		// Download blob
		keyValueFile, err := azureGateway.DownloadFile(containerName, blobName, rangeStart, rangeEnd)
		handleError(err)
		defer keyValueFile.Close()

		filePaths[i] = keyValueFile.Name()
	}
	// Perform reduce
	glog.Infof("Performing reduce for %s/%s", keyValueInfos[0].GetContainerName(), keyValueInfos[0].GetBlobName())
	resultFilePath, resultFileByteNum, err := ReduceBlobs(filePaths)
	handleError(err)
	// Create new blob
	containerName := keyValueInfos[0].GetContainerName()
	blobName := keyValueInfos[0].GetBlobName()
	resultBlobName := "R" + fmt.Sprint(reduceTaskNum) + "_" + blobName + ".txt"
	reducerContainerName := containerName + "-reduced"
	glog.Infof("Uploading result to %s/%s", reducerContainerName, resultBlobName)
	uploadErr := azureGateway.UploadFile(reducerContainerName, resultBlobName, resultFilePath)
	handleError(uploadErr)
	// Add to result info
	resultInfo := &pb.BlobInfo{
		ContainerName: reducerContainerName,
		BlobName:      resultBlobName,
		RangeStart:    0,
		RangeEnd:      int64(resultFileByteNum),
	}

	return &pb.ReduceReply{
		Success:       true,
		ReduceTaskNum: reduceTaskNum,
		ResultInfo:    resultInfo,
	}, nil
}

func ReduceBlobs(inputFilePaths []string) (string, int, error) {
	/* Python API:
	Inputs:
		input_file_paths: str[]
	Outputs:
		output_file_path: str
		output_file_bytes: int
	*/
	inputFilePathsJson, err := json.Marshal(inputFilePaths)
	if err != nil {
		return "", 0, err
	}
	cmd := exec.Command("python3", "reduce.py", string(inputFilePathsJson))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", 0, err
	}
	glog.Infof("Reduce output: %s", output)

	var reduceOutput ReduceOutput
	err = json.Unmarshal(output, &reduceOutput)
	if err != nil {
		return "", 0, err
	}

	return reduceOutput.OutputFilePath, reduceOutput.OutputFileBytes, nil
}
