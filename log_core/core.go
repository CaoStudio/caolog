package main

//import (
//	"io"
//	"os"
//)
//
//type Logger struct{}
//
//// Debug starts a new message with debug level.
////
//// You must call Msg on the returned event in order to send the event.
//func (l *Logger) Debug() *Event {
//	return l.newEvent(DebugLevel, nil)
//}
//
//// Info starts a new message with info level.
////
//// You must call Msg on the returned event in order to send the event.
//func (l *Logger) Info() *Event {
//	return l.newEvent(InfoLevel, nil)
//}
//
//// Warn starts a new message with warn level.
////
//// You must call Msg on the returned event in order to send the event.
//func (l *Logger) Warn() *Event {
//	return l.newEvent(WarnLevel, nil)
//}
//
//// Error starts a new message with error level.
////
//// You must call Msg on the returned event in order to send the event.
//func (l *Logger) Error() *Event {
//	return l.newEvent(ErrorLevel, nil)
//}
//
//// Err starts a new message with error level with err as a field if not nil or
//// with info level if err is nil.
////
//// You must call Msg on the returned event in order to send the event.
//func (l *Logger) Err(err error) *Event {
//	if err != nil {
//		return l.Error().Err(err)
//	}
//
//	return l.Info()
//}
//
//// Fatal starts a new message with fatal level. The os.Exit(1) function
//// is called by the Msg method, which terminates the program immediately.
////
//// You must call Msg on the returned event in order to send the event.
//func (l *Logger) Fatal() *Event {
//	return l.newEvent(FatalLevel, func(msg string) {
//		if closer, ok := l.w.(io.Closer); ok {
//			// Close the writer to flush any buffered message. Otherwise the message
//			// will be lost as os.Exit() terminates the program immediately.
//			closer.Close()
//		}
//		os.Exit(1)
//	})
//}
//
//// Panic starts a new message with panic level. The panic() function
//// is called by the Msg method, which stops the ordinary flow of a goroutine.
////
//// You must call Msg on the returned event in order to send the event.
//func (l *Logger) Panic() *Event {
//	return l.newEvent(PanicLevel, func(msg string) { panic(msg) })
//}
