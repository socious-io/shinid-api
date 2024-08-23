package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

var (
	queryCache = make(map[string]string)
	cacheLock  sync.RWMutex
)

// Get retrieves multiple records from the database with pagination
func Fetch(dest interface{}, ids ...any) error {
	_, err := cb.Execute(func() (interface{}, error) {
		isSlice := false
		var queryName string
		switch v := dest.(type) {
		case Model:
			queryName = v.FetchQuery()
		case []Model:
			queryName = v[0].FetchQuery()
			isSlice = true
		default:
			log.Fatal("invalid type to cast")
		}
		q, err := LoadQuery(queryName)
		if err != nil {
			log.Fatal(err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}

		query, args, err := sqlx.In(q, ids)
		if err != nil {
			return nil, err
		}

		query = db.Rebind(query)
		if isSlice {
			if err := db.Select(dest, query, args...); err != nil {
				return nil, err
			}
		} else {
			if err := db.Get(dest, query, args...); err != nil {
				return nil, err
			}
		}
		if err := UnmarshalJSONTextFields(dest); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func Get(dest interface{}, queryName string, args ...interface{}) error {
	_, err := cb.Execute(func() (interface{}, error) {
		q, err := LoadQuery(queryName)
		if err != nil {
			log.Fatal(err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}

		var jsonBytes []byte
		if err := db.QueryRowx(q, args...).Scan(&jsonBytes); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(jsonBytes, dest); err != nil {
			return nil, err
		}

		if err := db.Get(dest, q, args...); err != nil {
			return nil, err
		}

		return nil, nil
	})
	return err
}

// LoadQuery reads the SQL query from the file and caches it.
func LoadQuery(queryName string) (string, error) {
	cacheLock.RLock()
	if query, found := queryCache[queryName]; found {
		cacheLock.RUnlock()
		return query, nil
	}
	cacheLock.RUnlock()
	// Load query from file
	filePath := filepath.Join(sqlDir, queryName+".sql")
	queryBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	query := string(queryBytes)

	// Cache the query
	cacheLock.Lock()
	queryCache[queryName] = query
	cacheLock.Unlock()

	return query, nil
}

// ExecuteQuery is a general function to execute a write operation (INSERT, UPDATE, DELETE) with a circuit breaker
func ExecuteQuery(queryName string, data interface{}) (sql.Result, error) {
	var result sql.Result
	_, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()
		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}

		result, err = db.NamedExec(query, data)
		if err != nil {
			return nil, err
		}

		return result, nil
	})

	return result, err
}

// TxExecuteQuery is a general function to execute a write operation (INSERT, UPDATE, DELETE) with a circuit breaker
func TxExecuteQuery(tx *sqlx.Tx, queryName string, data interface{}) (sql.Result, error) {
	var result sql.Result
	_, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		result, err = tx.NamedExec(query, data)
		if err != nil {
			return nil, err
		}

		return result, nil
	})

	return result, err
}

// Query is a general function to execute a read operation that returns multiple rows with a circuit breaker
func Query(ctx context.Context, queryName string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}
		return db.QueryxContext(ctx, query, args...)
	})

	if err != nil {
		return nil, err
	}

	return rows.(*sqlx.Rows), nil
}

// Query is a general function to execute a read operation that returns multiple rows with a circuit breaker
func Queryx(queryName string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}
		return db.Queryx(query, args...)
	})

	if err != nil {
		return nil, err
	}

	return rows.(*sqlx.Rows), nil
}

func TxQuery(ctx context.Context, tx *sqlx.Tx, queryName string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		return tx.QueryxContext(ctx, query, args...)
	})

	if err != nil {
		return nil, err
	}

	return rows.(*sqlx.Rows), nil
}

// QuerSelect is a general function to execute a read operation that returns multiple rows with a circuit breaker
func QuerySelect(queryName string, dest interface{}, args ...interface{}) error {
	_, err := cb.Execute(func() (interface{}, error) {
		query, err := LoadQuery(queryName)
		if err != nil {
			return nil, fmt.Errorf("could not load query: %v", err)
		}

		db := GetDB()

		if db == nil {
			return nil, fmt.Errorf("database not connected")
		}
		return nil, db.Select(dest, query, args...)
	})
	return err
}

// UnmarshalJSONTextFields processes a single struct or a slice of structs to unmarshal
// fields of type types.JSONText into corresponding Go struct fields based on matching `db` and `json` tags.
func UnmarshalJSONTextFields(input interface{}) error {
	// Handle slice input
	v := reflect.ValueOf(input)

	if v.Kind() == reflect.Slice {
		// Iterate over each element in the slice and call UnmarshalJSONTextFields recursively
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i).Addr().Interface()
			if err := UnmarshalJSONTextFields(element); err != nil {
				return err
			}
		}
		return nil
	}

	// Handle single struct input
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("input must be a non-nil pointer to a struct or a slice")
	}

	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct or a slice of structs")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		// Look for fields of type types.JSONText
		if fieldType.Type == reflect.TypeOf(types.JSONText("")) {
			// Get the field name from the `db` tag
			dbTag := fieldType.Tag.Get("db")
			if dbTag == "" || dbTag == "-" {
				continue
			}

			// Look for a corresponding field with the same json tag
			for j := 0; j < v.NumField(); j++ {
				targetField := v.Field(j)
				targetFieldType := t.Field(j)
				if targetFieldType.Tag.Get("json") != dbTag {
					continue
				}
				// Ensure that the target field is a struct pointer and not the JSONText itself
				if targetField.Kind() == reflect.Ptr && targetFieldType.Type.Kind() == reflect.Ptr {
					targetField.Set(reflect.New(targetFieldType.Type.Elem())) // Initialize the struct pointer

					// Unmarshal into the corresponding field
					err := json.Unmarshal([]byte(field.Interface().(types.JSONText)), targetField.Interface())
					if err != nil {
						log.Println("parse json foreign Key : ", err)
						continue
					}
				}

				if targetField.Kind() == reflect.Slice && targetFieldType.Type.Kind() == reflect.Slice {
					sliceType := targetFieldType.Type.Elem() // Get the type of the slice element

					// Create a new slice with the appropriate type
					slicePtr := reflect.New(reflect.SliceOf(sliceType)).Interface()

					// Unmarshal into the slice
					err := json.Unmarshal([]byte(field.Interface().(types.JSONText)), slicePtr)
					if err != nil {
						log.Println("parse Json array foreign Key : ", err)
						continue
					}

					// Set the unmarshaled slice to the target field
					targetField.Set(reflect.ValueOf(slicePtr).Elem())
				}
			}
		}
	}

	return nil
}

// extractFields extracts the fields from a struct and returns them as a slice of interface{}
func extractFields(model interface{}) []interface{} {
	v := reflect.ValueOf(model).Elem()
	numFields := v.NumField()
	fields := make([]interface{}, numFields)

	for i := 0; i < numFields; i++ {
		fields[i] = v.Field(i).Interface()
	}

	return fields
}
