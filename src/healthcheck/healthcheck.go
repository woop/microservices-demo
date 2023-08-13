package main

import (
	"context"
	"log"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:3550", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewProductCatalogServiceClient(conn)

	listProductsResponse, err := client.ListProducts(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Failed to list products: %v", err)
	}

	for _, product := range listProductsResponse.Products {
		log.Printf("Product ID: %s, Name: %s, Price: $%d", product.Id, product.Name, product.PriceUsd.Units)
	}
}
