package worker

import (
	"birthdayreminder/config"
	"birthdayreminder/db"
	"birthdayreminder/models"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func RunDailyTask() {
	// Define the IST location
	loc, err := time.LoadLocation("Asia/Kolkata") // IST timezone
	if err != nil {
		fmt.Println("Error loading IST location:", err)
		return
	}

	// Get current time in IST
	now := time.Now().In(loc)
	day := now.Day()
	month := int(now.Month()) // Converts the month to int (1-12)
	fmt.Printf("Running task for day %d of month %d in IST\n", day, month)
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error in loading config in worker", err)
		return
	}
	dbSerVice, err := db.NewDBService(config)
	if err != nil {
		fmt.Println("Error in init DBService in worker", err)
		return
	}
	birthdays, err := dbSerVice.GetAllBirthDaysByDayMonth(context.Background(), day+1, month)
	if err != nil {
		fmt.Println("Error in getting birthdays", err)
		return
	}
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(config.AWSRegion),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     config.AWSAccessKey, // Assume these are specified in your local config
				SecretAccessKey: config.AWSSecretKey,
			},
		}),
	)
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}
	sesClient := ses.NewFromConfig(awsCfg)
	for _, b := range birthdays {
		go sendBirthDayEmail(&b, sesClient, config)
	}
	// Implement your task logic here
}

func sendBirthDayEmail(b *models.BirthDay, sesClient *ses.Client, conf *config.Config) {
	fmt.Printf("Sending for %s\n", b.Name)
	destination := &types.Destination{
		ToAddresses: []string{conf.TargetEmail},
	}
	emailMsg := fmt.Sprintf("Hello Bta,\n%s birthday a hnai ta.Hei a birthday a nia %d/%d\n\nKa lawm e", b.Name, b.Day, b.Month)
	subject := fmt.Sprintf("%s Birthday hriattirna", b.Name)
	message := &types.Message{
		Body: &types.Body{
			Text: &types.Content{
				Data: &emailMsg,
			},
		},
		Subject: &types.Content{
			Data: &subject,
		},
	}
	source := conf.SourceEmail
	_, err := sesClient.SendEmail(context.Background(), &ses.SendEmailInput{
		Destination: destination,
		Message:     message,
		Source:      &source,
	})
	if err != nil {
		fmt.Printf("Sending Failed for %s\n", b.Name)
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Sent succesfully for %s\n", b.Name)
}
