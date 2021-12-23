package statelint

import (
	"flag"
	"fmt"
	"log"
	"os"
	config2 "statelint/config"
	"statelint/j2119"
	"statelint/localization"
	"statelint/reader"
)

var (
	languageUsage = fmt.Sprintf("Sets the language of the error output, value should be "+
		"name of file in folder \"%s\" without extension", localization.DefaultLocalizationFolder)
	language = flag.String("localization", "en", languageUsage)

	minioFileUsage = "Filepath to minio validation file, should be string in this " +
		"format:\"bucket_name/.../validation_file.json\""
	minioFilePath = flag.String("minio_file", "", minioFileUsage)

	localFileUsage = "path to local json file. Can be reference or absolute."
	localFilePath  = flag.String("local_file", "", localFileUsage)

	help = flag.Bool("help", false, "print this help")
)

func init() {
	flag.StringVar(language, "l", "en", languageUsage)
	flag.StringVar(minioFilePath, "mf", "", minioFileUsage)
	flag.StringVar(localFilePath, "lf", "", localFileUsage)
}

func main() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()

		return
	}

	problems := j2119.NewProblems()
	defer func() {
		if problems.Len() != 0 {
			problemsCountStr := fmt.Sprintf(localization.GetLocalizerOrPanic().
				GetString("ProblemsCount"),
				problems.Len())
			fmt.Fprintln(os.Stdout, problemsCountStr)
		}

		for _, p := range problems.GetProblems() {
			_, err := fmt.Fprintln(os.Stdout, p)
			if err != nil {
				log.Fatalf("can not write to stdout %s", err.Error())
			}
		}
	}()

	stateLint, hasError := setup(problems)
	if hasError {
		return
	}

	switch {
	case *minioFilePath != "":
		config, err := config2.ReadConfig()
		if err != nil {
			problems.Append(err.Error())

			return
		}

		json, err := reader.GetJSONFromMinio(config.MinioConfig, *minioFilePath)
		if err != nil {
			problems.Append(err.Error())

			return
		}

		problems = stateLint.ValidateJSONStruct(json)
	case *localFilePath != "":
		json, err := reader.GetJSONFromLocalFile(*localFilePath)
		if err != nil {
			problems.Append(err.Error())

			return
		}

		problems = stateLint.ValidateJSONStruct(json)
	}
}

func setup(problems *j2119.Problems) (*j2119.StateLinter, bool) {
	localizer, err := localization.GetLocalizer()
	if err != nil {
		problems.Append(err.Error())

		return nil, true
	}

	err = localizer.SetLocalization(*language)
	if err != nil {
		problems.Append(err.Error())

		return nil, true
	}

	stateLint, err := j2119.NewStateLinter()
	if err != nil {
		problems.Append(err.Error())

		return nil, true
	}

	return stateLint, false
}
