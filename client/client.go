package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sRRRs-7/gRPC_Blog_API/server/blog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	create := flag.Bool("create", false, "Create a new blog")
	find := flag.Bool("find", false, "Find blog")
	flag.Parse()

	fmt.Println("Blog client...")
	conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	switch {
	case *create:
		client := blog.NewBlogApiClient(conn)
		CreateBlog(client)
	case *find:
		client := blog.NewBlogApiClient(conn)
		FindBlog(client)
	}
}

func CreateBlog(c blog.BlogApiClient) {
	fmt.Println("Create new blog...")

	fmt.Print("Input AuthorId: ")
	authorId, err := input(os.Stdin, flag.Args()...)
	if err != nil {
		log.Fatalf("Error author id input: %v", err)
	}
	intAuthorId, err := strconv.Atoi(authorId)
	if err != nil {
		log.Fatalf("Error convert author id: %v", err)
	}

	fmt.Print("Input titled: ")
	title, err := input(os.Stdin, flag.Args()...)
	if err != nil {
		log.Fatalf("Error title input: %v", err)
	}

	fmt.Print("Input content: ")
	content, err := input(os.Stdin, flag.Args()...)
	if err != nil {
		log.Fatalf("Error content input: %v", err)
	}

	req := &blog.CreateBlogReq{
		Blog: &blog.Blog{
			AuthorId: int32(intAuthorId),
			Title: title,
			Content: content,
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to calling Create blog RPC: %v", err)
	}

	log.Printf("Server Response: %v", res)
}

func FindBlog(c blog.BlogApiClient) {
	fmt.Println("Find blog")

	stream, err := c.FindBlog(context.Background(), &blog.FindBlogReq{})
	if err != nil {
		log.Fatalf("Failed to calling find blog RPC")
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			log.Fatalf("Failed to recieve server response: %v", err)
		}
		log.Printf("Server Response: %v", res)
	}
}

func input(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	text := scanner.Text()
	if len(text) == 0 {
		return "", nil
	}
	return text, nil
}