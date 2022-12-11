package email

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hack-31/point-app-backend/config"
)

const (
	// The character encoding for the email.
	CharSet = "UTF-8"
)

// メールを送信する
//
// @params recipient 送信者のメールアドレス
//
// @params subject 件名
//
// @params textBody テキストボディ
func SendMail(recipient string, subject string, textBody string) (*ses.SendEmailOutput, error) {
	// 環境変数の読み込み
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}
	// セッションの作成
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AWSRegion),
		Endpoint:    aws.String(cfg.AWSEndpoint),
		Credentials: credentials.NewStaticCredentials(cfg.AWSId, cfg.AWSSecret, "")},
	)
	if err != nil {
		return nil, err
	}
	svc := ses.New(sess)

	// 送信元メールアドレスの検証
	_, err = svc.VerifyEmailIdentity(&ses.VerifyEmailIdentityInput{EmailAddress: aws.String(cfg.SenderMailAddress)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return nil, fmt.Errorf("message rejected: %w", aerr)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return nil, fmt.Errorf("mail from domain not verified: %w", aerr)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return nil, fmt.Errorf("configuration set does not exist: %w", aerr)
			default:
				return nil, aerr
			}
		}
		return nil, err
	}

	// 送信先メールアドレスの検証
	_, err = svc.VerifyEmailIdentity(&ses.VerifyEmailIdentityInput{EmailAddress: aws.String(recipient)})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return nil, fmt.Errorf("message rejected: %w", aerr)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return nil, fmt.Errorf("mail from domain not verified: %w", aerr)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return nil, fmt.Errorf("configuration set does not exist: %w", aerr)
			default:
				return nil, aerr
			}
		}
		return nil, err
	}

	// 送信メールの作成
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(cfg.SenderMailAddress),
	}

	// メールの送信
	result, err := svc.SendEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				return nil, fmt.Errorf("message rejected: %w", aerr)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				return nil, fmt.Errorf("mail from domain not verified: %w", aerr)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				return nil, fmt.Errorf("configuration set does not exist: %w", aerr)
			default:
				return nil, aerr
			}
		}
		return nil, err
	}
	return result, nil
}
