package log;

// Fairly shallow logging wrapper for whatever logging library we go with.
// The only additional functionality is Panic() and PanicE().
// These will behave the same as Erorr() and ErrorE(), except they will panic after logging.

import (
   "github.com/Sirupsen/logrus"
)

// If you need to pass around a logger, you can use this to get ahold of one.
type Logger struct {}

func SetDebug(debug bool) {
   if (debug) {
      logrus.SetLevel(logrus.DebugLevel);
   } else {
      logrus.SetLevel(logrus.InfoLevel);
   }
}

func Debug(msg string) {
   logrus.Debug(msg);
}

func Info(msg string) {
   logrus.Info(msg);
}

func Warn(msg string) {
   WarnE(msg, nil)
}

func WarnE(msg string, err error) {
   if (err == nil) {
      logrus.Warn(msg);
   } else {
      logrus.Warnf("%s [%v]", msg, err);
   }
}

func Error(msg string) {
   ErrorE(msg, nil)
}

func ErrorE(msg string, err error) {
   if (err == nil) {
      logrus.Error(msg);
   } else {
      logrus.Errorf("%s [%v]", msg, err);
   }
}

func Panic(msg string) {
   PanicE(msg, nil)
}

func PanicE(msg string, err error) {
   if (err == nil) {
      logrus.Error(msg);
   } else {
      logrus.Errorf("%s [%v]", msg, err);
   }
   panic(msg);
}

func Fatal(msg string) {
   FatalE(msg, nil)
}

func FatalE(msg string, err error) {
   if (err == nil) {
      logrus.Fatal(msg);
   } else {
      logrus.Fatalf("%s [%v]", msg, err);
   }
}

// Attach each of the logging methods to Logger.
func (log Logger) Debug(msg string) { Debug(msg); }
func (log Logger) Info(msg string) { Info(msg); }
func (log Logger) Warn(msg string) { Warn(msg); }
func (log Logger) WarnE(msg string, err error) { WarnE(msg, err); }
func (log Logger) Error(msg string) { Error(msg); }
func (log Logger) ErrorE(msg string, err error) { ErrorE(msg, err); }
func (log Logger) Panic(msg string) { Panic(msg); }
func (log Logger) PanicE(msg string, err error) { PanicE(msg, err); }
func (log Logger) Fatal(msg string) { Fatal(msg); }
func (log Logger) FatalE(msg string, err error) { FatalE(msg, err); }
