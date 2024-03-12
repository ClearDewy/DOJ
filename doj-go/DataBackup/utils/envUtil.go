/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package utils

import "os"

func GetEnvDefault(env string, value string) string {
	env = os.Getenv(env)
	if env == "" {
		return value
	}
	return env
}
