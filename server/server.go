package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sRRRs-7/gRPC_Blog_API/server/blog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {}

type blogItem struct {
	ID 			int32
	AuthorID 	int32
	Title 		string
	Content 	string
	CreatedAt 	string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog service started...")
	listen, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Couldn't listen on %s: %v", listen, err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blog.RegisterBlogApiServer(s, &server{})

	go func() {
		fmt.Println("Starting server...")
		if err = s.Serve(listen); err != nil {
			log.Fatalf("Error serving %s: %v", listen, err)
		}
	}()
	// wait ctrl+C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-ch
	fmt.Println("Stopping server")
	s.Stop()
	fmt.Println("Closing listener")
	listen.Close()
	fmt.Println("EOP")
}

func ( *server) CreateBlog(ctx context.Context, req *blog.CreateBlogReq) (*blog.CreateBlogRes, error) {
	fmt.Println("Create Blog request")

	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/blogapi?sslmode=disable")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully created connection to database")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS blogapi (id bigserial PRIMARY KEY, author_id bigint NOT NULL, title varchar NOT NULL, content varchar NOT NULL, created_at timestamptz NOT NULL DEFAULT now());")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Insert error: %v", err),
		)
	}

	args := req.GetBlog()
	data := blogItem{
		AuthorID: args.AuthorId,
		Title: args.Title,
		Content: args.Content,
	}

	_, err = db.Exec("INSERT INTO blogapi (author_id, title, content) VALUES ($1, $2, $3) RETURNING *;", data.AuthorID, data.Title, data.Content)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Insert error: %v", err),
		)
	}

	row, err := db.Query("SELECT id FROM blogapi ORDER BY id DESC LIMIT 1")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Insert error: %v", err),
		)
	}
	var id int32
	row.Next()
	row.Scan(&id)

	return &blog.CreateBlogRes{
		Result: &blog.Blog{
			Id: id,
			AuthorId: args.AuthorId,
			Title: args.Title,
			Content: args.Content,
		},
	}, nil
}

func ( *server) FindBlog(req *blog.FindBlogReq, stream blog.BlogApi_FindBlogServer) error {
	fmt.Println("Get Blog data from postgres")

	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/blogapi?sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Successfully created connection to database")

	rows, err := db.Query("select * from blogapi;")
	if err != nil {
		return status.Errorf(codes.Internal, "Internal error: %v", err)
	}
	defer rows.Close()

	data := blogItem{}
	for rows.Next() {
		err = rows.Scan(&data.ID, &data.AuthorID, &data.Title, &data.Content, &data.CreatedAt)
		if err != nil {
			return status.Errorf(codes.Internal, "could not scan rows: %v", err)
		}

		jst := time.FixedZone("Asia/Tokyo", 9*60*60)
		date, err := time.Parse(time.RFC3339, data.CreatedAt)
		if err != nil {
			log.Fatalf("could not parse date: %v", err)
		}
		jstDate := date.In(jst)
		strDate := jstDate.Format("2006-01-02 15:04:05")

		stream.Send(&blog.FindBlogRes{
			Result: &blog.Blog{
				Id: data.ID,
				AuthorId: data.AuthorID,
				Title: data.Title,
				Content: data.Content,
				CreatedAt: strDate,
			},
		})
	}
	if err := rows.Err(); err != nil {
		return status.Errorf(codes.Internal, "Rows error: %v", err)
	}
	return nil
}