package util

import (
    "fmt"

    "github.com/olebedev/config"
)

// ConfigReader represents a configuration with convenient access methods.
// It is based on Config and offers additional reading functions.
type ConfigReader struct {
    *config.Config
}

// MustGet returns a nested config according to a dotted path or panics on error.
func (reader *ConfigReader) MustGet(path string) *ConfigReader {
    val, err := reader.Get(path)
    if err != nil {
        panic(err)
    }
    return &ConfigReader{Config: val}
}

// MustString returns a string according to a dotted path or panics on error.
func (reader *ConfigReader) MustString(path string) string {
    val, err := reader.String(path)
    if err != nil {
        panic(err)
    }
    return val
}

// MustInt returns an int according to a dotted path or panics on error.
func (reader *ConfigReader) MustInt(path string) int {
    val, err := reader.Int(path)
    if err != nil {
        panic(err)
    }
    return val
}

// MustBool returns a bool according to a dotted path or panics on error.
func (reader *ConfigReader) MustBool(path string) bool {
    val, err := reader.Bool(path)
    if err != nil {
        panic(err)
    }
    return val
}

// MustFloat64 returns a float64 according to a dotted path or panics on error.
func (reader *ConfigReader) MustFloat64(path string) float64 {
    val, err := reader.Float64(path)
    if err != nil {
        panic(err)
    }
    return val
}

// MustMap returns a map according to a dotted path or panics on error.
func (reader *ConfigReader) MustMap(path string) map[string]interface{} {
    val , err := reader.Map(path)
    if err != nil {
        panic(err)
    }
    return val
}

// MustStringList returns a []string according to a dotted path or panics on error.
func (reader *ConfigReader) MustStringList(path string) []string {
    val, err := reader.List(path)
    if err != nil {
        panic(err)
    }
    list := make([]string, len(val))
    for index, rawValue := range val {
       stringValue, wasString := rawValue.(string)
       if !wasString {
           panic(fmt.Errorf("expected a string value, got %v", rawValue))
       }
       list[index] = stringValue
    }
    return list
}
