package service

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	pb "ucode/ucode_go_sms_service/genproto/sms_service"
)

var (
	client pb.SmsServiceClient
)

func TestSmsService_Send(t *testing.T) {
	conn, err := grpc.Dial("localhost:2323", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect %v", err)
	}

	client = pb.NewSmsServiceClient(conn)

	_, err = uuid.NewRandom()
	if err != nil {
		t.Error("Error while generating UUID")
	}

	_, err = client.Send(context.Background(), &pb.Sms{
		Id:        "cffce8d8-d87f-4486-b9a9-534ca59c32c5",
		Text:      "Hello from service ",
		Recipient: "+998908082596",
	})
	if err != nil {
		t.Error("Error while creating client: ", err)
	}
}
