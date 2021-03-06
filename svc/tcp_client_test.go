package svc_test

import (
	"fmt"
	"github.com/hetianyi/godfs/common"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/gpip"
	"github.com/hetianyi/gox/logger"
	json "github.com/json-iterator/go"
	"io"
	"net"
	"strconv"
	"testing"
)

func TestSendMsg(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(6577))
	if err != nil {
		logger.Fatal("error start client:", err)
	}
	pip := &gpip.Pip{
		Conn: conn,
	}
	err = pip.Send(&common.Header{
		Operation:  common.OPERATION_CONNECT,
		Attributes: map[string]string{"secret": "123456"},
	}, nil, 0)
	if err != nil {
		logger.Fatal("error send data:", err)
	}
	err = pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
		header := _header.(*common.Header)
		bs, _ := json.Marshal(header)
		logger.Info("client got message:", string(bs))
		return nil
	})
	if err != nil {
		logger.Error("error:", err)
	}
}

func TestUpload(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(3456))
	if err != nil {
		logger.Fatal("error start client:", err)
	}
	pip := &gpip.Pip{
		Conn: conn,
	}

	// validate
	err = pip.Send(&common.Header{
		Operation:  common.OPERATION_CONNECT,
		Attributes: map[string]string{"secret": "kasd3123"},
	}, nil, 0)
	if err != nil {
		logger.Fatal("error send data:", err)
	}
	err = pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
		header := _header.(*common.Header)
		printResult(header)
		if header.Result != common.SUCCESS {
			logger.Fatal("error status: ", header.Msg)
		}
		return nil
	})
	if err != nil {
		logger.Error("error:", err)
	}

	// upload
	srcFile, _ := file.GetFile("E:\\TEMP\\9.jpg")
	fi, _ := srcFile.Stat()
	err = pip.Send(&common.Header{
		Operation: common.OPERATION_UPLOAD,
	}, srcFile, fi.Size())
	if err != nil {
		logger.Fatal("error upload file:", err)
	}
	err = pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
		header := _header.(*common.Header)
		printResult(header)
		return nil
	})
	if err != nil {
		logger.Error("error:", err)
	}
}

var pip *gpip.Pip

func pre() {
	logger.Init(nil)
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(6578))
	if err != nil {
		logger.Fatal("error start client:", err)
	}
	pip = &gpip.Pip{
		Conn: conn,
	}

	// validate
	err = pip.Send(&common.Header{
		Operation:  common.OPERATION_CONNECT,
		Attributes: map[string]string{"secret": "123456"},
	}, nil, 0)
	if err != nil {
		logger.Fatal("error send data:", err)
	}
	err = pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
		header := _header.(*common.Header)
		printResult(header)
		if header.Result != common.SUCCESS {
			logger.Fatal("error status: ", header.Msg)
		}
		return nil
	})
	if err != nil {
		logger.Error("error:", err)
	}
}

func TestRegisterSync(t *testing.T) {
	pre()
	instance := &common.Instance{
		Server: common.Server{
			Host:       "192.168.0.101",
			Port:       1234,
			Secret:     "123123",
			InstanceId: "xxx111",
		},
		Role: common.ROLE_STORAGE,
	}
	bs, _ := json.Marshal(instance)

	// register------
	err := pip.Send(&common.Header{
		Operation: common.OPERATION_CONNECT,
		Attributes: map[string]string{
			"instance": string(bs),
		},
	}, nil, 0)
	if err != nil {
		logger.Fatal("error register instance:", err)
	}
	err = pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
		header := _header.(*common.Header)
		printResult(header)
		return nil
	})
	if err != nil {
		logger.Error("error:", err)
	}

	// sync-------
	err = pip.Send(&common.Header{
		Operation: common.OPERATION_SYNC_INSTANCES,
	}, nil, 0)
	if err != nil {
		logger.Fatal("error register instance:", err)
	}
	err = pip.Receive(&common.Header{}, func(_header interface{}, bodyReader io.Reader, bodyLength int64) error {
		header := _header.(*common.Header)
		printResult(header)
		return nil
	})
	if err != nil {
		logger.Error("error:", err)
	}
}

func printResult(o interface{}) {
	bs, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println(string(bs))
}
