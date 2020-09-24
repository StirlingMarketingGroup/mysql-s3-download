package main

// #include <string.h>
// #include <stdbool.h>
// #include <mysql.h>
// #cgo CFLAGS: -O3 -I/usr/include/mysql -fno-omit-frame-pointer
import "C"
import (
	"log"
	"os"
	"unsafe"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

// main function is needed even for generating shared object files
func main() {}

var l = log.New(os.Stderr, "s3-download: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

//export s3_download_init
func s3_download_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 3 {
		msg(message, "`s3_download` requires 3 parameters: the region, the bucket, and the key")
		return C.bool(true)
	}

	argsTypes := (*[2]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	initid.maybe_null = 1

	return C.bool(false)
}

//export s3_download
func s3_download(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 3
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	*length = 0
	*isNull = 1
	if argsArgs[0] == nil ||
		argsArgs[1] == nil ||
		argsArgs[2] == nil {
		return nil
	}

	a := make([]string, c, c)
	for i, argsArg := range argsArgs {
		a[i] = C.GoStringN(argsArg, C.int(argsLengths[i]))
	}

	sess, err := session.NewSession(&aws.Config{Region: &a[0]})
	if err != nil {
		l.Println(errors.Wrapf(err, "failed to create AWS session"))
		return nil
	}
	buff := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: &a[1],
			Key:    &a[2],
		})
	if err != nil {
		l.Println(errors.Wrapf(err, "failed to download file from S3"))
		return nil
	}

	b := buff.Bytes()
	*length = uint64(len(b))
	*isNull = 0
	return C.CString(string(b))
}
