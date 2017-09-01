// Package configo provides functions to handle configuration items.
//
// Example:
//  type Config struct {
//    SomeThing *SomeThing
//    Mysql  MysqlConfig `toml:"mysql"`
//    Port  *Int `default:"8080"`
//  }
//
//  type MysqlConfig struct {
//    Dsn *string
//  }
//
//  mydb := "/mydb"
//
//  config := Config{
//    Mysql: MysqlConfig{
//      Dsn: &mydb,
//    },
//  }
//
//  err := configo.UnmarshalFile("/path/to/config.toml", &config)
//  // error check
//
//  fmt.Printf("Listening on port %d", *config.Port)
package configo

// UnmarshalFile decodes the contents of the file `f` in TOML format into a pointer `v`. If `v` contains data, that data will be used as "defaults".
//
// A field's value will be determined based on the following order:
//
// 1. If an "env" tag exists for a field and an environment variable matching the tag's value exists, the environment variable's value will be used, subject to type casting.
// 2. If the field exists in the file, its value will be used. The `toml` tag may be used to map TOML keys to fields that don't match the key name exactly.
// 3. If `v` already contains a value for the field, it will be used.
// 4. If a "default" tag exists for a field, its value will be used, subject to type casting.
// 5. The field will be initialized to its zero value (i.e., "" for string, 0 for int, etc).
func UnmarshalFile(f string, v interface{}) error {
	var err error

	cc := NewConfigoChain(
		NewDefaultsConfigo(),
		NewTomlConfigo(f),
		NewEnvConfigo(),
	)

	err = cc.Load(v)
	if err != nil {
		return err
	}

	return nil
}
