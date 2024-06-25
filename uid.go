/*
 * Project: Application Utility Library
 * Filename: /uid.go
 * Created Date: Sunday September 3rd 2023 18:25:10 +0800
 * Author: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * Company: BerryPay (M) Sdn. Bhd.
 * --------------------------------------
 * Last Modified: Tuesday June 25th 2024 15:09:22 +0800
 * Modified By: Sallehuddin Abdul Latif (sallehuddin@berrypay.com)
 * --------------------------------------
 * Copyright (c) 2023 BerryPay (M) Sdn. Bhd.
 */

package apputil

import (
	"time"

	"github.com/muyo/sno"
	"github.com/vishal-bihani/go-tsid"
)

const (
	MSG_CREATE_FACTORY_ERROR = "Error fetching customized TsidFactory, fallback to default"
	MSG_TSID_GENERATE_ERROR  = "Error generating tsid, fallback to default generator"
)

var meta byte
var nodeBits int32
var nodeId int32
var epoch int64
var prefix string
var suffix string

type UidConfig struct {
	Meta     byte
	NodeId   int32
	NodeBits int32
	Epoch    int64
	Prefix   string
	Suffix   string
}

func init() {
	meta = 88
	nodeBits = 8      // 8 bits for node bits -> max 254 nodes
	nodeId = 88       // Default node identifier: 88
	epoch = 259345800 // default epoch: 1978-03-22
	prefix = ""
	suffix = ""
}

func InitWithConfig(config *UidConfig) {
	SetIDMeta(config.Meta)

	if config.NodeBits > 0 && config.NodeBits <= 10 {
		SetNodeBits(config.NodeBits)
	}

	if config.NodeId >= 0 && config.NodeId <= 254 {
		SetNodeId(config.NodeId)
	}

	if config.Epoch < 0 {
		SetEpoch(config.Epoch)
	}

	if config.Prefix != "" {
		SetPrefix(config.Prefix)
	}

	if config.Suffix != "" {
		SetSuffix(config.Suffix)
	}
}

func SetIDMeta(b byte) {
	meta = b
}

func SetNodeBits(b int32) {
	// our minimum number of nodes are set at 254 (which requires 8 bits representation)
	// and maximum supported by the tsid library is 10
	if b <= 0 {
		nodeBits = 8
	} else if b >= 10 {
		nodeBits = 10
	} else {
		nodeBits = b
	}
}

func GetNodeBits() int32 {
	return nodeBits
}

func SetNodeId(b int32) {
	if b < 0 {
		nodeId = 0
	} else if b > 254 {
		nodeId = 254
	} else {
		nodeId = b
	}
}

func GetNodeId() int32 {
	return nodeId
}

func SetEpoch(b int64) {
	if b < 0 {
		epoch = 0
	} else {
		epoch = b
	}
}

func GetEpoch() int64 {
	return epoch
}

func SetPrefix(s string) {
	prefix = s
}

func GetPrefix() string {
	return prefix
}

func SetSuffix(s string) {
	suffix = s
}

func GetSuffix() string {
	return suffix
}

func GetIDMeta() byte {
	return meta
}

func GenerateRequestIDByte() []byte {
	return sno.NewWithTime(meta, time.Now()).Bytes()
}

func GenerateRequestID() sno.ID {
	return sno.NewWithTime(meta, time.Now())
}

func GenerateRequestIDString() string {
	return sno.NewWithTime(meta, time.Now()).String()
}

func getExtendedTsid(newIdStr string, options ...bool) string {
	usePrefix := false
	useSuffix := false

	if len(options) > 0 && options[0] {
		usePrefix = true
	}

	if len(options) > 1 && options[1] {
		useSuffix = true
	}

	if usePrefix && useSuffix {
		return GetPrefix() + newIdStr + GetSuffix()
	} else if usePrefix && !useSuffix {
		return GetPrefix() + newIdStr
	} else if !usePrefix && useSuffix {
		return newIdStr + GetSuffix()
	} else {
		return newIdStr
	}
}

/**
 * GenerateTsidString
 * @params options - [usePrefix, useSuffix]
 */
func GenerateTsidString(options ...bool) string {
	tsidFactory, err := tsid.TsidFactoryBuilder().
		WithNodeBits(nodeBits).
		WithNode(nodeId).
		WithCustomEpoch(epoch).
		Build()
	if err != nil {
		Logger.Warn(MSG_CREATE_FACTORY_ERROR)

		return getExtendedTsid(tsid.Fast().ToString(), options...)
	}

	id, err := tsidFactory.Generate()
	if err != nil {
		Logger.Warn(MSG_TSID_GENERATE_ERROR)
		id = tsid.Fast()
	}

	return getExtendedTsid(id.ToString(), options...)
}

func GenerateTsidNumber() int64 {
	tsidFactory, err := tsid.TsidFactoryBuilder().
		WithNodeBits(nodeBits).
		WithNode(nodeId).
		WithCustomEpoch(epoch).
		Build()
	if err != nil {
		Logger.Warn(MSG_CREATE_FACTORY_ERROR)
		return tsid.Fast().ToNumber()
	}

	id, err := tsidFactory.Generate()
	if err != nil {
		Logger.Warn(MSG_TSID_GENERATE_ERROR)
		id = tsid.Fast()
	}

	return id.ToNumber()
}

func GenerateTsidBytes(options ...bool) []byte {
	tsidFactory, err := tsid.TsidFactoryBuilder().
		WithNodeBits(nodeBits).
		WithNode(nodeId).
		WithCustomEpoch(epoch).
		Build()
	if err != nil {
		Logger.Warn(MSG_CREATE_FACTORY_ERROR)
		return tsid.Fast().ToBytes()
	}

	id, err := tsidFactory.Generate()
	if err != nil {
		Logger.Warn(MSG_TSID_GENERATE_ERROR)
		id = tsid.Fast()
	}

	return id.ToBytes()
}
