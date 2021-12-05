// Package generator ...
package generator

//go:generate mockgen -destination=mocks/kafka/kafka.go -source=libs/kafka/kafka.go
//go:generate mockgen -destination=mocks/db/db.go -source=libs/db/provider.go
