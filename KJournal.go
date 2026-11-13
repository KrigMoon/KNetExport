package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"github.com/KrigMoon/KJournal"
)

var (
	KJUnit KJournal.KJournal
)

//region Lifecircle

//export RunKJournal
func RunKJournal(settings *C.char) *C.char {
	setsStr := C.GoString(settings)

	var sets KJSettings
	if err := fromJson(setsStr, &sets); err != nil {
		return C.CString(errorResponse(err.Error()))
	}

	setsKJ := KJournal.KJSettings{
		DefaultReaderBuf: sets.DefaultReaderBuf,
		DefaultReaderFilter: func(r *KJournal.Record) bool {
			return true
		},
		MaxInGroup:     sets.MaxInGroup,
		FileDumperDir:  sets.FileDumperDir,
		AlarmRWGroup:   fromDTToDuration(sets.AlarmRWGroup),
		AlarmWrite:     fromDTToDuration(sets.AlarmWrite),
		CrashReportDir: sets.CrashReportDir,
	}

	if err := KJUnit.Run(setsKJ); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse("Successful RunKJournal"))
}

//export EndworkKJournal
func EndworkKJournal() *C.char {
	if err := KJUnit.Endwork(); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse("Successful EndworkKJournal"))
}

//endregion

//region Metrics

//export KJournalMetrics
func KJournalMetrics() *C.char {
	str, err := toJson(KJUnit.GetMetrics())
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//endregion

//region Writer

//export KJNewWriter
func KJNewWriter() *C.char {
	if err := KJUnit.NewWriter(); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(""))
}

//export KJDelWriter
func KJDelWriter() *C.char {
	if err := KJUnit.DelWriter(); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(""))
}

//export KJWrite
func KJWrite(record *C.char) *C.char {
	var rRaw Record
	if err := fromJson(C.GoString(record), &rRaw); err != nil {
		return C.CString(errorResponse(err.Error()))
	}

	KJUnit.Write(KJournal.NewRecord(
		rRaw.Author,
		rRaw.Level,
		rRaw.Code,
		rRaw.Message,
	))
	return C.CString(successResponse(""))
}

//endregion

//region Reader

//export KJNewReader
func KJNewReader(id *C.char, buff C.int) *C.char {
	var idPtr *string
	if id != nil {
		s := C.GoString(id)
		idPtr = &s
	}

	ID, err := KJUnit.NewReader(idPtr, int(buff), nil)
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(ID))
}

//export KJDelReader
func KJDelReader(id *C.char) *C.char {
	if err := KJUnit.DelReader(C.GoString(id)); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(""))
}

//export KJRead
func KJRead(id *C.char) *C.char {
	record, err := KJUnit.Read(C.GoString(id))
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	if record == nil {
		return C.CString(successResponse(""))
	}

	str, err := toJson(recordFromKJ(record))
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//export KJReadGroup
func KJReadGroup(id *C.char) *C.char {
	records, err := KJUnit.ReadGroup(C.GoString(id))
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}

	out := make([]Record, 0, len(records))
	for _, r := range records {
		if r == nil {
			continue
		}
		out = append(out, recordFromKJ(r))
	}

	str, err := toJson(out)
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//endregion
