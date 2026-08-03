package cache

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/valkey-io/valkey-glide/go/v2/constants"
	"github.com/valkey-io/valkey-glide/go/v2/options"
)

// Reusuable function to reduce boilerplate dealing with JSON data.

func GetJSON(ctx context.Context, key string) (map[string]any, error) {
	value, err := Client.Get(context.TODO(), key)
	if err != nil || value.IsNil() {
		return nil, err
	}
	var data map[string]any
	json_str := value.Value()

	// Unmarshal the JSON string into the map
	// Note: jsonStr must be converted to a byte slice []byte
	err = json.Unmarshal([]byte(json_str), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
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
