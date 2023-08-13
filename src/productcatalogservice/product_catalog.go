package main

import (
	"context"
	"fmt"
	"strings"

	"database/sql"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type productCatalog struct {
	db      *sql.DB
	catalog pb.ListProductsResponse
}

func (p *productCatalog) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (p *productCatalog) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (p *productCatalog) ListProducts(context.Context, *pb.Empty) (*pb.ListProductsResponse, error) {
	rows, err := p.db.Query(`SELECT id, name, description, picture, price_currency, price_units, price_nanos, categories FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var id string
		var name string
		var description string
		var picture string
		var price_currency string
		var price_units int64
		var price_nanos int32
		var categories string

		if err2 := rows.Scan(&id, &name, &description, &picture, &price_currency, &price_units, &price_nanos, &categories); err2 != nil {
			fmt.Printf("Error scanning row: %v", err2)
			return nil, err2
		}

		var product pb.Product
		product.Categories = strings.Split(categories, ",")
		product.Id = id
		product.Name = name
		product.Description = description
		product.Picture = picture
		product.PriceUsd = new(pb.Money)
		product.PriceUsd.CurrencyCode = price_currency
		product.PriceUsd.Units = price_units
		product.PriceUsd.Nanos = price_nanos
		products = append(products, &product)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.ListProductsResponse{Products: products}, nil
}

func (p *productCatalog) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	row := p.db.QueryRow(`SELECT id, name, description, picture, price_currency, price_units, price_nanos, categories FROM products WHERE id = $1`, req.Id)

	var id string
	var name string
	var description string
	var picture string
	var price_currency string
	var price_units int64
	var price_nanos int32
	var categories string

	if err := row.Scan(&id, &name, &description, &picture, &price_currency, &price_units, &price_nanos, &categories); err != nil {
		return nil, err
	}

	var product pb.Product
	product.Categories = strings.Split(categories, ",")
	product.Id = id
	product.Name = name
	product.Description = description
	product.Picture = picture
	product.PriceUsd = new(pb.Money)
	product.PriceUsd.CurrencyCode = price_currency
	product.PriceUsd.Units = price_units
	product.PriceUsd.Nanos = price_nanos

	return &product, nil
}

func (p *productCatalog) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	rows, err := p.db.Query(`SELECT id, name, description, picture, price_currency, price_units, price_nanos, categories FROM products WHERE LOWER(name) LIKE $1 OR LOWER(description) LIKE $1`, "%"+strings.ToLower(req.Query)+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*pb.Product
	for rows.Next() {
		var id string
		var name string
		var description string
		var picture string
		var price_currency string
		var price_units int64
		var price_nanos int32
		var categories string

		if err2 := rows.Scan(&id, &name, &description, &picture, &price_currency, &price_units, &price_nanos, &categories); err2 != nil {
			fmt.Printf("Error scanning row: %v", err2)
			return nil, err2
		}

		var product pb.Product
		product.Categories = strings.Split(categories, ",")
		product.Id = id
		product.Name = name
		product.Description = description
		product.Picture = picture
		product.PriceUsd = new(pb.Money)
		product.PriceUsd.CurrencyCode = price_currency
		product.PriceUsd.Units = price_units
		product.PriceUsd.Nanos = price_nanos
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &pb.SearchProductsResponse{Results: products}, nil
}
