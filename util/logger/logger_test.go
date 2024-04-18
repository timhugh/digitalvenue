package logger_test

import (
	"context"
	"encoding/json"
	"github.com/matryer/is"
	"github.com/timhugh/digitalvenue/util/logger"
	"os"
	"strings"
	"testing"
)

func initTest(t *testing.T) (*is.I, strings.Builder, map[string]interface{}) {
	is := is.New(t)
	out := strings.Builder{}
	msg := make(map[string]interface{})

	err := os.Setenv("LOG_LEVEL", "debug")
	if err != nil {
		t.Fatalf("Error setting environment variable: %s", err)
	}

	return is, out, msg
}

func TestContextLogger_BasicLogging(t *testing.T) {
	is, out, msg := initTest(t)

	log := logger.New(&out)
	clearOutputs(&out, &msg)

	log.Debug("This is a %s level message", "debug")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "debug")
	is.Equal(msg["message"], "This is a debug level message")

	clearOutputs(&out, &msg)

	log.Info("This is an %s level message", "info")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "info")
	is.Equal(msg["message"], "This is an info level message")

	clearOutputs(&out, &msg)

	log.Warn("This is a %s level message", "warn")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "warn")
	is.Equal(msg["message"], "This is a warn level message")

	clearOutputs(&out, &msg)

	log.Error("This is an %s level message", "error")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "error")
	is.Equal(msg["message"], "This is an error level message")

	clearOutputs(&out, &msg)

	log.Fatal("This is a %s level message", "fatal")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "fatal")
	is.Equal(msg["message"], "This is a fatal level message")
}

func TestContextLogger_Params(t *testing.T) {
	is, out, msg := initTest(t)

	log := logger.New(&out)
	clearOutputs(&out, &msg)

	log.AddParam("key", "value")
	log.AddParams(map[string]interface{}{"key2": 123.0, "key3": true})

	log.Debug("Message")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "debug")
	is.Equal(msg["message"], "Message")
	is.Equal(msg["key"], "value")
	is.Equal(msg["key2"], 123.0)
	is.Equal(msg["key3"], true)
}

func TestContextLogger_NewContext(t *testing.T) {
	is, out, msg := initTest(t)

	ctx := logger.NewContext(&out)
	_, log := logger.FromContext(ctx)
	clearOutputs(&out, &msg)

	log.Debug("Message")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "debug")
	is.Equal(msg["message"], "Message")
}

func TestContextLogger_ExistingContext(t *testing.T) {
	is, out, msg := initTest(t)

	logIn := logger.New(&out)
	clearOutputs(&out, &msg)
	logIn.AddParam("key", "value")

	ctx := logger.Attach(context.Background(), logIn)

	_, logOut := logger.FromContext(ctx)
	is.Equal(logIn, logOut)

	logOut.Debug("Message")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["key"], "value")
}

func TestContextLogger_SubLogger(t *testing.T) {
	is, out, msg := initTest(t)

	log := logger.New(&out)
	clearOutputs(&out, &msg)
	log.AddParam("key", "value")

	sub := log.Sub()
	sub.AddParam("key2", "value2")

	sub.Debug("Message")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["key"], "value")
	is.Equal(msg["key2"], "value2")

	clearOutputs(&out, &msg)

	log.Debug("Message")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["key"], "value")
	is.Equal(msg["key2"], nil)
}

func TestContextLogger_Chaining(t *testing.T) {
	is, out, msg := initTest(t)

	log := logger.New(&out).AddParam("key", "value").AddParams(map[string]interface{}{"key2": 123.0, "key3": true})
	clearOutputs(&out, &msg)

	log.Debug("Message")
	is.NoErr(json.Unmarshal([]byte(out.String()), &msg))
	is.Equal(msg["level"], "debug")
	is.Equal(msg["message"], "Message")
	is.Equal(msg["key"], "value")
	is.Equal(msg["key2"], 123.0)
	is.Equal(msg["key3"], true)
}

func clearOutputs(out *strings.Builder, msg *map[string]interface{}) {
	out.Reset()
	*msg = make(map[string]interface{})
}
