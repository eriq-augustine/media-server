package log;

// Fairly shallow logging wrapper for whatever logging library we go with.
// The only additional functionality is Panic() and PanicE().
// These will behave the same as Erorr() and ErrorE(), except they will panic after logging.

import (
   "github.com/Sirupsen/logrus"
)

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
