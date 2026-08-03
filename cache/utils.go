package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-glide/go/v2/constants"
	"github.com/valkey-io/valkey-glide/go/v2/options"
)

// Reusuable function to reduce boilerplate dealing with JSON data.

// Make sure to use generics using your own struct,
// so that it's used in json.unmarshal,
// and you can use it directly without manual conversion
func GetJSON[T any](ctx context.Context, key string, returned_struct T) (T, error) {
	value, err := Client.Get(context.TODO(), key)
	if err != nil || value.IsNil() {
		// Make sure to return custom error, so that if err != nil but value is nil,
		// It'll fail your outer check for err != nil.
		return returned_struct, fmt.Errorf("Couldn't get value")
	}
	json_str := value.Value()

	// Unmarshal the JSON string into the map
	// Note: jsonStr must be converted to a byte slice []byte
	err = json.Unmarshal([]byte(json_str), &returned_struct)
	if err != nil {
		return returned_struct, err
	}

	return returned_struct, nil
}

// Used to Set JSON data as a string, as I don't need atomic updates.
func SetJSON(ctx context.Context, key string, value any) error {
	// Transform map into JSON byte slice
	json_bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	json_str := string(json_bytes) // Convert byte slice to string
	_, err = Client.SetWithOptions(
		ctx,
		key,
		json_str,
		options.SetOptions{Expiry: &options.Expiry{Type: constants.Seconds, Duration: 60 * 15}},
	)
	if err != nil {
		return err
	}

	return nil
}

func DelKey(ctx context.Context, key string) error {
	_, err := Client.Del(ctx, []string{key})
	if err != nil {
		return err
	}
	return nil
}

func FormatKeyByID(prefix string, id uuid.UUID) string {
	return prefix + ":" + id.String()
}
