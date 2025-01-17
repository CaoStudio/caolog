package caolog_test

import (
	"bytes"
	"context"
	caolog "github.com/CaoStudio/cao-log"
	"github.com/CaoStudio/cao-log/plugin"
	"github.com/goccy/go-json"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

var (
	fl64         = 23426339.234
	fl32 float32 = 23426339.234
	i            = 23426339
	ui   uint    = 23426339
	i8   int8    = 123
	ui8  uint8   = 234
	i16  int16   = 12345
	ui16 uint16  = 23456
	i32  int32   = 23426339
	ui32 uint32  = 23426339
	i64  int64   = 23426339
	ui64 uint64  = 23426339
	s            = "hello, world; 你好，世界"
	by           = []byte("hello, world; 你好，世界")
	i32l         = []int32{45627194, 2404854, 7809539, 26405593, 35493885, 40958275, 3574030, 83653665, 52735216, 74129751}
)

func BenchmarkLogFaster(b *testing.B) {
	caolog.InitLogger(caolog.InfoLevel)
	tracer := plugin.NewTrace()
	caolog.With(tracer.Option)
	ctx := context.Background()
	//_ = make([]byte, 1<<34)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			caolog.Info(ctx, fl64, fl32, i, ui, i8, ui8, i16, ui16, i32, ui32, i64, ui64, s, by, i32l) //, i32l
			//logFaster.Info(ctx, fl64, fl64, fl64, fl64)
		}
	})
}

func BenchmarkAll(b *testing.B) {
	//_ = make([]byte, 1<<32)
	runtime.GOMAXPROCS(2)
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("BenchmarkLogFaster", BenchmarkLogFaster)
}

func BenchmarkIntList2String(b *testing.B) {
	b.Run("BenchmarkIntList2String", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			byL, _ := json.Marshal(i32l)
			_ = *(*string)(unsafe.Pointer(&byL))
		}
	})
	buffer := bytes.Buffer{}
	b.Run("BenchmarkIntList2StringByteBuffer", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buffer.WriteString("[")
			for no, item := range i32l {
				buffer.WriteString(strconv.FormatInt(int64(item), 10))
				if no < len(i32l)-1 {
					buffer.WriteString(",")
				}
			}
			buffer.WriteString("[")
			byL := buffer.Bytes()
			_ = *(*string)(unsafe.Pointer(&byL))
		}
	})
}

func BenchmarkWriteAll(b *testing.B) {
	_ = make([]byte, 1<<36)
	b.ReportAllocs()
	b.Run("BenchmarkWrite", func(b *testing.B) {
		bu := bytes.Buffer{}
		for item := 0; item < b.N; item++ {
			bu.Write(by)
		}
	})
	b.Run("BenchmarkWriteString", func(b *testing.B) {
		bu := bytes.Buffer{}
		for item := 0; item < b.N; item++ {
			bu.WriteString(s)
		}
	})
}

func ListToString(il []int32) string {
	bu := bytes.Buffer{}
	bu.WriteString("[")
	for no, item := range il {
		bu.WriteString(strconv.FormatInt(int64(item), 10))
		if no < len(il)-1 {
			bu.WriteString(",")
		}
	}
	bu.WriteString("]")
	st := bu.Bytes()
	return *(*string)(unsafe.Pointer(&st))
}

func ListToStringByte(il []int32) string {
	//bu.WriteString("[")
	var formatByte []byte
	formatByte = append(formatByte, '[')
	for no, item := range il {
		formatByte = append(formatByte, []byte(strconv.FormatInt(int64(item), 10))...)
		if no < len(il)-1 {
			formatByte = append(formatByte, ',')
		}
	}
	formatByte = append(formatByte, ']')
	return *(*string)(unsafe.Pointer(&formatByte))
}

func ListToStringBuilder(il []int32) string {
	bu := strings.Builder{}
	bu.WriteString("[")
	for no, item := range il {
		bu.WriteString(strconv.FormatInt(int64(item), 10))
		if no < len(il)-1 {
			bu.WriteString(",")
		}
	}
	bu.WriteString("]")
	return bu.String()
}

func BenchmarkListToString(b *testing.B) {
	var list []int32
	for tem := 0; tem < 100; tem++ {
		list = append(list, int32(tem))
	}
	b.ReportAllocs()
	var st string
	b.Run("json", func(b *testing.B) {
		b.ReportAllocs()
		for item := 0; item < b.N; item++ {
			byt, _ := json.Marshal(list)
			st = *(*string)(unsafe.Pointer(&byt))
		}
	})
	b.Log(st)
	b.Run("Buffer", func(b *testing.B) {
		b.ReportAllocs()
		for item := 0; item < b.N; item++ {
			st = ListToString(list)
		}
	})
	b.Log(st)
	b.Run("BytesSum", func(b *testing.B) {
		b.ReportAllocs()
		for item := 0; item < b.N; item++ {
			st = ListToStringByte(list)
		}
	})
	b.Log(st)
	b.Run("Builder", func(b *testing.B) {
		b.ReportAllocs()
		for item := 0; item < b.N; item++ {
			st = ListToStringBuilder(list)
		}
	})
	b.Log(st)
}

func BenchmarkMakeBuilder(b *testing.B) {
	var list []int32
	for tem := 0; tem < 100; tem++ {
		list = append(list, int32(tem))
	}
	b.ReportAllocs()
	b.Run("makeClean", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bu := &strings.Builder{}
				bu.Reset()
			}
		})
	})

	//pool := sync.Pool{New: func() interface{} {
	//	return &strings.Builder{}
	//}}
	//b.Run("getPoolClean", func(b *testing.B) {
	//	b.ReportAllocs()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			bu := pool.Get().(*strings.Builder)
	//			bu.Reset()
	//		}
	//	})
	//})

	//b.Run("makeGrow", func(b *testing.B) {
	//	b.ReportAllocs()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			bu := &strings.Builder{}
	//			bu.Grow(1024)
	//		}
	//	})
	//})

	//poolGrow := sync.Pool{New: func() interface{} {
	//	bu := &strings.Builder{}
	//	bu.Grow(1024)
	//	return bu
	//}}
	//b.Run("poolGrow", func(b *testing.B) {
	//	b.ReportAllocs()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			bu := poolGrow.Get().(*strings.Builder)
	//			bu.Reset()
	//		}
	//	})
	//})
	//
	//b.Run("makeGrowWrite", func(b *testing.B) {
	//	b.ReportAllocs()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			bu := &strings.Builder{}
	//			bu.Grow(1024)
	//			bu.WriteString("hello, world")
	//		}
	//	})
	//})

	//poolGrowWrite := sync.Pool{New: func() interface{} {
	//	bu := &strings.Builder{}
	//	bu.Grow(1024)
	//	return bu
	//}}
	//b.Run("poolGrowWrite", func(b *testing.B) {
	//	b.ReportAllocs()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			bu := poolGrowWrite.Get().(*strings.Builder)
	//			bu.Reset()
	//			bu.WriteString("hello, world")
	//		}
	//	})
	//})
	//
	//b.Run("makeNewWrite", func(b *testing.B) {
	//	b.ReportAllocs()
	//	b.RunParallel(func(pb *testing.PB) {
	//		for pb.Next() {
	//			bu := strings.Builder{}
	//			bu.WriteString("hello, world")
	//		}
	//	})
	//})

	b.Run("builderResetWrite", func(b *testing.B) {
		b.ReportAllocs()
		for item := 0; item < b.N; item++ {
			bu := strings.Builder{}
			bu.WriteString("hello, world")
		}
	})

	bur := strings.Builder{}
	b.Run("makeNewWrite", func(b *testing.B) {
		b.ReportAllocs()
		for item := 0; item < b.N; item++ {
			bur.WriteString("hello, world")
			bur.Reset()
		}
	})

}

var (
	fileName = "greeter/greeter.go"
	fileLine = "123"
)

func BenchmarkBuildFileLineBuilder(b *testing.B) {
	b.ReportAllocs()
	var out string
	for item := 0; item < b.N; item++ {
		bu := strings.Builder{}
		bu.WriteString(fileName)
		bu.WriteString(":")
		bu.WriteString(fileLine)
		out = bu.String()
	}
	_ = out
	b.Log(out)
}

func BenchmarkBuildFileLineBytes(b *testing.B) {
	b.ReportAllocs()
	var out string
	for item := 0; item < b.N; item++ {
		msgBytes := make([]byte, len(fileName)+len(fileLine)+1)
		copy(msgBytes, fileName)
		msgBytes[len(fileName)] = ':'
		//copy(msgBytes[len(fileName):], ":")
		copy(msgBytes[len(fileName)+1:], fileLine)
		out = *(*string)(unsafe.Pointer(&msgBytes))
	}
	_ = out
	b.Log(out)
}

func BenchmarkBuildFileLineBytesNoCopy(b *testing.B) {
	b.ReportAllocs()
	var out string
	for item := 0; item < b.N; item++ {
		msgBytes := make([]byte, 0, len(fileName)+len(fileLine)+1)
		msgBytes = append(msgBytes, fileName...)
		msgBytes = append(msgBytes, ':')
		msgBytes = append(msgBytes, fileLine...)
		out = *(*string)(unsafe.Pointer(&msgBytes))
	}
	_ = out
	b.Log(out)
}

func BenchmarkBuildFileLineString(b *testing.B) {
	b.ReportAllocs()
	var out string
	for item := 0; item < b.N; item++ {
		out = fileName + ":" + fileLine
	}
	_ = out
	b.Log(out)
}

func BenchmarkBuildFileLineAll(b *testing.B) {
	//_ = make([]byte, 1<<32)
	//runtime.GOMAXPROCS(2)
	b.ReportAllocs()
	b.ResetTimer()
	b.Run("BenchmarkBuildFileLineBuilder", BenchmarkBuildFileLineBuilder)
	b.Run("BenchmarkBuildFileLineBytes", BenchmarkBuildFileLineBytes)
	b.Run("BenchmarkBuildFileLineBytesNoCopy", BenchmarkBuildFileLineBytesNoCopy)
	b.Run("BenchmarkBuildFileLineString", BenchmarkBuildFileLineString)
}

func BenchmarkBuildMsgBuilder(b *testing.B) {

}

func TestName(t *testing.T) {
	a := int32(-1)
	t.Log(uint32(a))
}
