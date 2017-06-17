package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"

	"golang.org/x/net/context"

	pb "say-grpc/backend/api"

	"os"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("start Listening...%d\n", *port)

	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}

type server struct{}

func (server) Say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("error temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close %s: %v", f.Name(), err)
	}

	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read tmp file: %v", err)
	}
	if err := os.Remove(f.Name()); err != nil {
		return nil, fmt.Errorf("could not remove tmp file: %v", err)
	}

	return &pb.Speech{Audio: data}, nil
}

/*	cmd := exec.Command("flite", "-t", os.Args[1], "-o", "output.wav")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}*/
