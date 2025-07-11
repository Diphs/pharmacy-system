package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"fmt"
	"pharmacy/graphql/graph/model"
	"pharmacy/graphql/internal/db"
	"pharmacy/graphql/internal/rabbitmq"
)

// CreateTransaction is the resolver for the createTransaction field.
func (r *mutationResolver) CreateTransaction(ctx context.Context, input model.TransactionInput) (*model.Transaction, error) {
	// Validation
	if input.TransactionID == "" || input.MedicineName == "" {
		return nil, fmt.Errorf("transactionId and medicineName cannot be empty")
	}
	if input.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}
	if input.Price <= 0 {
		return nil, fmt.Errorf("price must be positive")
	}

	tx := db.Transaction{
		TransactionID: input.TransactionID,
		MedicineName:  input.MedicineName,
		Quantity:      int(input.Quantity),
		Price:         input.Price,
	}

	// Save to database
	if err := r.DB.SaveTransaction(tx); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// Publish to RabbitMQ
	if err := r.Publisher.Publish(ctx, rabbitmq.Transaction{
		TransactionID: input.TransactionID,
		MedicineName:  input.MedicineName,
		Quantity:      int(input.Quantity),
		Price:         input.Price,
	}); err != nil {
		return nil, fmt.Errorf("failed to publish transaction: %w", err)
	}

	return &model.Transaction{
		ID:            "", // Will be set by database
		TransactionID: tx.TransactionID,
		MedicineName:  tx.MedicineName,
		Quantity:      int32(tx.Quantity),
		Price:         tx.Price,
		CreatedAt:     "", // Will be set by database
	}, nil
}

// Transactions is the resolver for the transactions field.
func (r *queryResolver) Transactions(ctx context.Context) ([]*model.Transaction, error) {
	dbTransactions, err := r.DB.GetAllTransactions()
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	var transactions []*model.Transaction
	for _, tx := range dbTransactions {
		transactions = append(transactions, &model.Transaction{
			ID:            fmt.Sprintf("%d", tx.ID),
			TransactionID: tx.TransactionID,
			MedicineName:  tx.MedicineName,
			Quantity:      int32(tx.Quantity),
			Price:         tx.Price,
			CreatedAt:     tx.CreatedAt,
		})
	}

	return transactions, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}
*/
