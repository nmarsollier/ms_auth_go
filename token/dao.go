package token

import (
	"context"
	"fmt"
	"strings"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/nmarsollier/authgo/tools/db"
	"github.com/nmarsollier/authgo/tools/errors"
)

func collection() (*mongo.Collection, error) {
	db, err := db.Get()
	if err != nil {
		return nil, err
	}

	collection := db.Collection("tokens")

	_, err = collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.String("userId", ""),
			),
			Options: bson.NewDocument(),
		},
	)
	if err != nil {
		fmt.Print(err.Error())
	}

	return db.Collection("tokens"), nil
}

// Save agrega un token a la base de datos
func save(token Token) (Token, error) {
	if err := validateSchema(token); err != nil {
		return token, err
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return token, err
	}

	if len(token.ID()) > 0 {
		_id, _ := objectid.FromHex(token.ID())

		_, err := collection.UpdateOne(context.Background(),
			bson.NewDocument(bson.EC.ObjectID("_id", _id)),
			bson.NewDocument(
				bson.EC.SubDocumentFromElements("$set",
					bson.EC.Boolean("enabled", token.Enabled),
				),
			))

		if err != nil {
			db.HandleError(err)
			return token, err
		}
	} else {
		res, err := collection.InsertOne(context.Background(), token)
		if err != nil {
			db.HandleError(err)
			return token, err
		}

		token.SetID(res.InsertedID.(objectid.ObjectID))
	}

	return token, nil
}

func validateSchema(token Token) error {
	token.UserID = strings.TrimSpace(token.UserID)

	result := make(validator.ValidationErrors)

	if len(token.ID()) > 0 {
		if _, err := objectid.FromHex(token.ID()); err != nil {
			result["id"] = &validator.FieldError{
				Field: "id",
				Tag:   "Invalid",
			}
		}
	}
	if len(token.UserID) == 0 {
		result["userId"] = &validator.FieldError{
			Field: "userId",
			Tag:   "Requerido",
		}
	} else {
		if _, err := objectid.FromHex(token.UserID); err != nil {
			result["userId"] = &validator.FieldError{
				Field: "userId",
				Tag:   "Invalid",
			}
		}
	}

	if len(result) > 0 {
		return result
	}

	return nil
}

func findByID(tokenID string) (*Token, error) {
	_id, err := getID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	result := bson.NewDocument()
	filter := bson.NewDocument(bson.EC.ObjectID("_id", *_id))
	err = collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		} else {
			return nil, err
		}
	}

	token := newTokenFromBson(*result)

	return &token, nil
}

func findByUserID(tokenID string) (*Token, error) {
	_id, err := getID(tokenID)
	if err != nil {
		return nil, errors.Unauthorized
	}

	collection, err := collection()
	if err != nil {
		db.HandleError(err)
		return nil, err
	}

	result := bson.NewDocument()

	filter := bson.NewDocument(
		bson.EC.String("userId", _id.Hex()),
		bson.EC.Boolean("enabled", true),
	)
	err = collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		db.HandleError(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.Unauthorized
		} else {
			return nil, err
		}
	}

	token := newTokenFromBson(*result)

	return &token, nil
}

func delete(tokenID string) error {
	token, err := findByID(tokenID)
	if err != nil {
		db.HandleError(err)
		return err
	}

	token.Enabled = false
	_, err = save(*token)

	if err != nil {
		db.HandleError(err)
		return err
	}

	return nil
}

func getID(ID string) (*objectid.ObjectID, error) {
	_id, err := objectid.FromHex(ID)
	if err != nil {
		return nil, ErrID
	}
	return &_id, nil
}
