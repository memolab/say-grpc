package main

import (
	"flag"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"fmt"
	"io/ioutil"
	"os"
	pb "say-grpc/backend/api"
)

func main() {
	backend := flag.String("b", "localhost:8080", "backend addr")
	output := flag.String("o", "output.wav", "wav file")

	flag.Parse()

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speak\"\n", os.Args[0])
		os.Exit(1)
	}

	text := &pb.Text{Text: flag.Arg(0)}
	rsp, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(*output, rsp.Audio, 0666); err != nil {
		log.Fatal(err)
	}
}
