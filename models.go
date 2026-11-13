package main

import (
	"encoding/json"

	"github.com/KrigMoon/KCommandParser"
	"github.com/KrigMoon/KJournal"
)

//region Settings

type KJSettings struct {
	DefaultReaderBuf int      `json:"DefaultReaderBuf"`
	MaxInGroup       int      `json:"MaxInGroup"`
	FileDumperDir    string   `json:"FileDumperDir"`
	AlarmRWGroup     DataTime `json:"AlarmRWGroup"`
	AlarmWrite       DataTime `json:"AlarmWrite"`
	CrashReportDir   string   `json:"CrashReport"`
}

type KNSettings struct {
	ExternalCommandsMap []string        `json:"ExternalCommandsMap"`
	AutoRun             []string        `json:"AutoRun"`
	FileStorageDir      string          `json:"FileStorageDir"`
	RemoteControl       bool            `json:"RemoteControl"`
	MyNode              json.RawMessage `json:"MyNode"`
	Password            string          `json:"Password"`
}

//endregion

//region Record

type Record struct {
	Author  string   `json:"Author"`
	Level   string   `json:"Level"`
	Code    string   `json:"Code"`
	Time    DataTime `json:"Time"`
	Message string   `json:"Message"`
}

func recordFromKJ(r *KJournal.Record) Record {
	return Record{
		Author:  r.Author,
		Level:   r.Level,
		Code:    r.Code,
		Time:    toDTFromTime(r.Time),
		Message: r.Message,
	}
}

//endregion

//region Response

type Response struct {
	CommandID       string   `json:"CommandID"`
	Level           string   `json:"Level"`
	Time            DataTime `json:"Time"`
	ResponseMessage string   `json:"ResponseMessage"`
}

func responseFromKP(r *KCommandParser.Response) Response {
	return Response{
		CommandID:       r.CommandID,
		Level:           r.Level,
		Time:            toDTFromTime(r.Time),
		ResponseMessage: r.ResponseMessage,
	}
}

func responseToKP(r Response) *KCommandParser.Response {
	return &KCommandParser.Response{
		CommandID:       r.CommandID,
		Level:           r.Level,
		Time:            fromDTToTime(r.Time),
		ResponseMessage: r.ResponseMessage,
	}
}

//endregion
