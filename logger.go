/*
 * Project: Application Utility Library
 * Filename: /logger.go
 * Created Date: Sunday September 3rd 2023 18:22:42 +0800
 * Author: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * Company: BerryPay (M) Sdn. Bhd.
 * --------------------------------------
 * Last Modified: Tuesday June 25th 2024 16:52:38 +0800
 * Modified By: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * --------------------------------------
 * Copyright (c) 2023 BerryPay (M) Sdn. Bhd.
 */

package apputil

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var SLogger *zap.SugaredLogger
var ZapLevel zapcore.Level
var AppZapConfig zapcore.EncoderConfig
var AppZapCore zapcore.Core
var AppZapFileEncoder zapcore.Encoder
var AppZapConsoleEncoder zapcore.Encoder
var AppZapWriter zapcore.WriteSyncer

func init() {
	logBaseDir, present := os.LookupEnv("LOG_BASE_DIR")
	if present {
		if logBaseDir != "" {
			InitializeLogger(128, 3, 28, true, logBaseDir)
			return
		}
	}
	InitializeLogger(128, 3, 28, true, "")
}

func InitializeLogger(maxSize int, maxBackups int, maxAge int, compress bool, logBaseDir string) {
	AppZapConfig = zap.NewProductionEncoderConfig()
	AppZapConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	AppZapFileEncoder = zapcore.NewJSONEncoder(AppZapConfig)
	AppZapConsoleEncoder = zapcore.NewConsoleEncoder(AppZapConfig)

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeName := filepath.Base(exePath)
	logFilename := exeName + ".log"

	logDir := filepath.Join(filepath.Dir(exePath), "logs")
	if logBaseDir != "" {
		logDir = filepath.Join(filepath.Dir(exePath), "logs")
	}

	logPath := filepath.Join(logDir, logFilename)

	// Check if directory exists
	_, err = os.Stat(filepath.Dir(logPath))
	if os.IsNotExist(err) {
		// Create directory if it does not exist
		err := os.MkdirAll(filepath.Dir(logPath), 0755)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	// Check if file exists
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		currentTime := time.Now()
		formattedTime := currentTime.Format("20060102150405")
		if err := os.Rename(logPath, filepath.Join(logDir, exeName+"-"+formattedTime+".log")); err != nil {
			fmt.Printf("Error renaming pre-existing log file: %v\n", err)
		} else {
			fmt.Println("Pre-existing log file renamed successfully.")
		}
	}

	AppZapWriter = zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize,    // size in megabytes
		MaxBackups: maxBackups, // # of backups
		MaxAge:     maxAge,     // # of days
		Compress:   compress,
	})

	envLogLevel := os.Getenv("LOG_LEVEL")
	switch envLogLevel {
	case "debug":
		ZapLevel = zapcore.DebugLevel
	case "info":
		ZapLevel = zapcore.InfoLevel
	case "warn":
		ZapLevel = zapcore.WarnLevel
	case "error":
		ZapLevel = zapcore.ErrorLevel
	case "panic":
		ZapLevel = zapcore.PanicLevel
	case "fatal":
		ZapLevel = zapcore.FatalLevel
	case "":
		fmt.Println("No log level specified... defaulting to info.")
		ZapLevel = zapcore.InfoLevel
	default:
		fmt.Println("Unrecognized log level specified... defaulting to info.")
		ZapLevel = zapcore.InfoLevel
	}

	// Default log level for debug mode is "debug"
	envAppEnv := os.Getenv("APP_ENV")
	if envAppEnv == "debug" {
		fmt.Println("Debug app environment detected... forcing log level to debug.")
		ZapLevel = zapcore.DebugLevel
	}

	AppZapCore = zapcore.NewTee(
		zapcore.NewCore(AppZapFileEncoder, AppZapWriter, ZapLevel),
		zapcore.NewCore(AppZapConsoleEncoder, zapcore.AddSync(os.Stdout), ZapLevel),
	)
	Logger = zap.New(AppZapCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	SLogger = Logger.Sugar()
}

func ChangeLogLevel(level string) {
	if ZapLevel.String() == level {
		return
	}

	switch level {
	case "debug":
		ZapLevel = zapcore.DebugLevel
	case "info":
		ZapLevel = zapcore.InfoLevel
	case "warn":
		ZapLevel = zapcore.WarnLevel
	case "error":
		ZapLevel = zapcore.ErrorLevel
	case "panic":
		ZapLevel = zapcore.PanicLevel
	case "fatal":
		ZapLevel = zapcore.FatalLevel
	default:
		Logger.Warn("Unrecognized log level specified... defaulting to info.")
		ZapLevel = zapcore.InfoLevel
	}

	SLogger.Infof("Changing log level to %v", ZapLevel.String())

	AppZapCore = zapcore.NewTee(
		zapcore.NewCore(AppZapFileEncoder, AppZapWriter, ZapLevel),
		zapcore.NewCore(AppZapConsoleEncoder, zapcore.AddSync(os.Stdout), ZapLevel),
	)
	Logger = zap.New(AppZapCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	SLogger = Logger.Sugar()
}
