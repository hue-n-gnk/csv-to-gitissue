package env

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/joho/godotenv"
)

type EnvMapper struct {
	PostgrestDatabase string `config:"SECRET_POSTGRES_DATABASE,required"`
	PostgrestHostname string `config:"SECRET_POSTGRES_HOSTNAME,required"`
	PostgrestPort     string `config:"SECRET_POSTGRES_PORT,required"`
	PostgrestPassword string `config:"SECRET_POSTGRES_PASSWORD,required"`
	PostgrestUser     string `config:"SECRET_POSTGRES_USER,required"`

	BucketId  string `config:"BUCKET_ID,required"`
	GHToken   string `config:"GH_TOKEN,required"`
	GHOwner   string `config:"GH_OWNER,required"`
	GHRepo    string `config:"GH_REPO,required"`
	AccessKey string `config:"AWS_ACCESS_KEY_ID"`
	SecretKey string `config:"AWS_SECRET_ACCESS_KEY"`
	Region    string `config:"AWS_DEFAULT_REGION"`
}

var (
	envData = EnvMapper{}
)

func Load() EnvMapper {
	loadEnv(os.Getenv("ENV"))
	return envData
}

func loadEnv(target string) {
	// get filepath of this file
	_, base, _, _ := runtime.Caller(1)
	base = filepath.Dir(base)
	f := filepath.Join(base, "..", "..", fmt.Sprintf(".env.%s", target))

	if f == filepath.Join("..", "..", ".env.") {
		fmt.Println("NOTE: did you set 'ENV' environment variable? For testing, set 'ENV=test'.")
	}

	if err := godotenv.Load(f); err != nil {
		// When dotenv files are not found, print warning and continue the context.
		fmt.Println("NOTE:", err)
		fmt.Println("NOTE: falling back to environment variable checks..")
	}

	loader := confita.NewLoader(
		env.NewBackend(),
	)

	if err := loader.Load(context.Background(), &envData); err != nil {
		panic(err)
	}
}

func Clear() {
	envData = EnvMapper{}
}

func GetSession(typ ...string) (*session.Session, error) {
	if IsAWSRuntime() {
		// Refer to IAM Role in AWS runtime
		return session.NewSession()
	}

	// Docker must specify access/secret keys explicitly
	return session.NewSession(&aws.Config{
		Region: aws.String(envData.Region),
		Credentials: credentials.NewStaticCredentials(
			envData.AccessKey,
			envData.SecretKey,
			"",
		),
	})
}

func IsAWSRuntime() bool {
	if envData == (EnvMapper{}) || envData.AccessKey == "" || envData.SecretKey == "" {
		return true
	}
	return false
}
