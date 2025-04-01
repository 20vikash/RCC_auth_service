package env

import "os"

func GetDBUserName() string {
	return os.Getenv("DB_USERNAME")
}

func GetDBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func GetDBName() string {
	return os.Getenv("DB_NAME")
}

func GetGmailAppPassword() string {
	return os.Getenv("GMAIL_APP_PASSWORD")
}
