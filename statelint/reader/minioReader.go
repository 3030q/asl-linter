package reader

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	path2 "path"
	"statelint/config"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrBadMinioPath = errors.New("bad minio path")

// minimum path is myBucket/file
const (
	minimumMinioPathArgs = 2
	separator            = "/"
)

func GetJSONFromMinio(minioConfig config.Minio, minioPath string) (interface{}, error) {
	minioPath = path2.Clean(minioPath)
	minioPath = strings.ReplaceAll(minioPath, "\\", "/")
	parts := strings.Split(minioPath, separator)

	if len(parts) < minimumMinioPathArgs {
		return nil, ErrBadMinioPath
	}

	minioBucketName := parts[0]
	minioPath = strings.Join(parts[1:], separator)

	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.Username, minioConfig.Password, ""),
		Secure: minioConfig.UseSsl,
	})
	if err != nil {
		return nil, fmt.Errorf("can not initialize minio client: %w", err)
	}

	ctx := context.Background()
	object, err := client.GetObject(
		ctx,
		minioBucketName,
		minioPath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("can not get minio object: %w", err)
	}

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("can not read minio object: %w", err)
	}

	var j interface{}
	err = json.Unmarshal(data, &j)

	if err != nil {
		return nil, fmt.Errorf("can not unmarshal minio object: %w", err)
	}

	return j, nil
}
