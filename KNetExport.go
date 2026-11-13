package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"time"
	"unsafe"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

//region Helps

//region^ DataTime struct

type DataTime struct {
	Year        int `json:"Year"`
	Month       int `json:"Month"`
	Day         int `json:"Day"`
	Hour        int `json:"Hour"`
	Minute      int `json:"Minute"`
	Second      int `json:"Second"`
	Miliseconds int `json:"Miliseconds"`
}

func toDTFromTime(t time.Time) DataTime {
	return DataTime{
		Year:        t.Year(),
		Month:       int(t.Month()),
		Day:         t.Day(),
		Hour:        t.Hour(),
		Minute:      t.Minute(),
		Second:      t.Second(),
		Miliseconds: t.Nanosecond() / 1_000_000,
	}
}

func fromDTToTime(dt DataTime) time.Time {
	return time.Date(
		dt.Year,
		time.Month(dt.Month),
		dt.Day,
		dt.Hour,
		dt.Minute,
		dt.Second,
		dt.Miliseconds*1_000_000,
		time.Local,
	)
}

func fromDTToDuration(dt DataTime) time.Duration {
	return time.Duration(dt.Year)*365*24*time.Hour +
		time.Duration(dt.Month)*30*24*time.Hour +
		time.Duration(dt.Day)*24*time.Hour +
		time.Duration(dt.Hour)*time.Hour +
		time.Duration(dt.Minute)*time.Minute +
		time.Duration(dt.Second)*time.Second +
		time.Duration(dt.Miliseconds)*time.Millisecond
}

//endregion^

//region^ FuncResponse

type FuncResponse struct {
	Level    string `json:"Level"`
	Response string `json:"Response"`
}

func errorResponse(response string) string {
	str, _ := toJson(FuncResponse{Level: "Error", Response: response})
	return str
}

func successResponse(response string) string {
	str, _ := toJson(FuncResponse{Level: "Success", Response: response})
	return str
}

//endregion^

//region^ Struct`s Convert

func toJson(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func fromJson(j string, v any) error {
	return json.Unmarshal([]byte(j), v)
}

func toProtoJson(v proto.Message) (string, error) {
	b, err := protojson.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func fromProtoJson(j string, v proto.Message) error {
	return protojson.Unmarshal([]byte(j), v)
}

//endregion^

//region^ FreeString

//export FreeString
func FreeString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

//endregion^

//endregion

func main() {}
