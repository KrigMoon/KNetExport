package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"github.com/KrigMoon/KNet/KNet"
	"github.com/KrigMoon/KNet/KNetCore"
	knetpb "github.com/KrigMoon/KNet/Proto/gen"
)

var (
	KNUnit KNet.KNet
)

//region Lifecircle

//export RunKNet
func RunKNet(settings *C.char) *C.char {
	setsStr := C.GoString(settings)

	var sets KNSettings
	if err := fromJson(setsStr, &sets); err != nil {
		return C.CString(errorResponse(err.Error()))
	}

	extCmds := make(map[string]struct{}, len(sets.ExternalCommandsMap))
	for _, cmd := range sets.ExternalCommandsMap {
		extCmds[cmd] = struct{}{}
	}

	myNode := &knetpb.Node{}
	if err := fromProtoJson(string(sets.MyNode), myNode); err != nil {
		return C.CString(errorResponse(err.Error()))
	}

	setsKN := KNet.KNSettings{
		Journal:             &KJUnit,
		ExternalCommandsMap: extCmds,
		AutoRun:             sets.AutoRun,
		FileStorageDir:      sets.FileStorageDir,
		RemoteControl:       sets.RemoteControl,
		Core: KNetCore.KNCSettings{
			Journal:  &KJUnit,
			MyNode:   myNode,
			Password: sets.Password,
		},
	}

	if err := KNUnit.Run(setsKN); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse("Successful RunKNet"))
}

//export EndworkKNet
func EndworkKNet() *C.char {
	if err := KNUnit.Endwork(); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse("Successful EndworkKNet"))
}

//endregion

//region Metrics

//export KNetMetrics
func KNetMetrics() *C.char {
	str, err := toJson(KNUnit.GetMetrics())
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//endregion

//region Commands

//export KNExecute
func KNExecute(command *C.char) *C.char {
	if err := KNUnit.Execute(C.GoString(command)); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(""))
}

//export KNResponse
func KNResponse() *C.char {
	resp := KNUnit.Response()
	if resp == nil {
		return C.CString(successResponse(""))
	}

	str, err := toJson(responseFromKP(resp))
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//export KNWriteResponse
func KNWriteResponse(response *C.char) *C.char {
	var raw Response
	if err := fromJson(C.GoString(response), &raw); err != nil {
		return C.CString(errorResponse(err.Error()))
	}

	if err := KNUnit.WriteResponse(responseToKP(raw)); err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(""))
}

//export KNGetCursor
func KNGetCursor() *C.char {
	str, err := toJson(KNUnit.GetCursor())
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//export KNExternalPrograms
func KNExternalPrograms() *C.char {
	str, err := toJson(KNUnit.ExternalPrograms())
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//endregion

//region Peers

//export KNPeerList
func KNPeerList() *C.char {
	str, err := toJson(KNUnit.GetPeerList())
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//endregion

//region Packages

//export KNAddReaderPackages
func KNAddReaderPackages(id *C.char, buff C.int) *C.char {
	var idPtr *string
	if id != nil {
		s := C.GoString(id)
		idPtr = &s
	}

	ID := KNUnit.AddReaderPackages(idPtr, int(buff), nil)
	return C.CString(successResponse(ID))
}

//export KNDelReaderPackages
func KNDelReaderPackages(id *C.char) *C.char {
	KNUnit.DelReaderPackages(C.GoString(id))
	return C.CString(successResponse(""))
}

//export KNPackage
func KNPackage(id *C.char) *C.char {
	pkg := KNUnit.Package(C.GoString(id))
	if pkg == nil {
		return C.CString(successResponse(""))
	}

	str, err := toProtoJson(pkg)
	if err != nil {
		return C.CString(errorResponse(err.Error()))
	}
	return C.CString(successResponse(str))
}

//endregion
