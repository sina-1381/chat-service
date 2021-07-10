package controllers

import (
	"fmt"
	"ginGorm/validations"
	"github.com/NaySoftware/go-fcm"
	"github.com/gin-gonic/gin"
)
const (
	serverKey = "AAAA8OUx-wE:APA91bEW-Hw70ccL2dn9plBjuzEgvij4hnEt8yhZj8LCkE-zdJqBNT4rajJZwQC7OU9Ab7EzqpKtCx4YatI2V-h68_fX4NP9o3eIwEn067DVdLc46EgIhw2p-jHl2lOTOAS1dRcYQq0g"
)
func Fcm(c *gin.Context)  {
	var input validations.FcmModel
	validations.CheckValidate(&input, c)
	notif:= fcm.NotificationPayload{
		Title:            input.Title,
		Body:             input.Msg,
	}
	s := fcm.NewFcmClient(serverKey)
	s.SetNotificationPayload(&notif)
	s.NewFcmRegIdsMsg(input.UserIDs, notif)
	status, err := s.Send()
	if err == nil {
		status.PrintResults()
	} else {
		fmt.Println(err)
	}
}
