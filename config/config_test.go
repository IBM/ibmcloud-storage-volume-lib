/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2017, 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package config

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

type testConfig struct {
	Header sectionTestConfig
}

type sectionTestConfig struct {
	ID      int
	Name    string
	YesOrNo bool
	Pi      float64
	List    string
}

var testConf = testConfig{
	Header: sectionTestConfig{
		ID:      1,
		Name:    "test",
		YesOrNo: true,
		Pi:      3.14,
		List:    "1, 2",
	},
}

func getContextLogger() (*zap.Logger, zap.AtomicLevel) {
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	traceLevel := zap.NewAtomicLevel()
	traceLevel.SetLevel(zap.InfoLevel)
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleDebugging, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (lvl >= traceLevel.Level()) && (lvl < zapcore.ErrorLevel)
		})),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleErrors, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, traceLevel
}

var testLogger, _ = getContextLogger()

func TestReadConfig(t *testing.T) {
	t.Log("Testing ReadConfig")

	configPath := "test.toml"
	expectedConf, _ := ReadConfig(configPath, testLogger)

	assert.NotNil(t, expectedConf)
}

func TestReadConfigEmptyPath(t *testing.T) {
	t.Log("Testing ReadConfig")

	configPath := ""
	expectedConf, _ := ReadConfig(configPath, testLogger)

	assert.NotNil(t, expectedConf)
}

func TestParseConfig(t *testing.T) {
	t.Log("Testing config parsing")
	var testParseConf testConfig

	configPath := "test.toml"
	ParseConfig(configPath, &testParseConf, testLogger)

	expected := testConf
	assert.Exactly(t, expected, testParseConf)
}

func TestParseConfigNoMatch(t *testing.T) {
	t.Log("Testing config parsing false positive")
	var testParseConf testConfig

	configPath := "test.toml"
	ParseConfig(configPath, &testParseConf, testLogger)

	expected := testConfig{
		Header: sectionTestConfig{
			ID:      1,
			Name:    "testnomatch",
			YesOrNo: true,
			Pi:      3.14,
			List:    "1, 2",
		}}

	assert.NotEqual(t, expected, testParseConf)

}

func TestParseConfigNoMatchTwo(t *testing.T) {
	t.Log("Testing config parsing false positive")
	var testParseConf testConfig

	configPath := "test1.toml"
	ParseConfig(configPath, &testParseConf, testLogger)

	expected := testConfig{
		Header: sectionTestConfig{
			ID:      1,
			Name:    "testnomatch",
			YesOrNo: true,
			Pi:      3.14,
			List:    "1, 2",
		}}

	assert.NotEqual(t, expected, testParseConf)

}

func TestGetGoPath(t *testing.T) {
	t.Log("Testing getting GOPATH")
	goPath := "/tmp"
	os.Setenv("GOPATH", goPath)

	path := GetGoPath()

	assert.Equal(t, goPath, path)
}

func TestGetEnv(t *testing.T) {
	t.Log("Testing getting ENV")
	goPath := "/tmp"
	os.Setenv("ENVTEST", goPath)

	path := getEnv("ENVTEST")

	assert.Equal(t, goPath, path)
}

func TestGetGoPathNullPath(t *testing.T) {
	t.Log("Testing getting GOPATH NULL Path")
	goPath := ""
	os.Setenv("GOPATH", goPath)

	path := GetGoPath()

	assert.Equal(t, goPath, path)
}

func TestGetEtcPath(t *testing.T) {
	t.Log("Testing GetEtcPath")
	expectedEtcPath := "src/github.com/IBM/ibmcloud-storage-volume-lib/etc"

	etcPath := GetEtcPath()

	assert.Equal(t, expectedEtcPath, etcPath)
}

func TestGetConfPath(t *testing.T) {
	t.Log("Testing GetEtcPath")
	expectedEtcPath := "src/github.com/IBM/ibmcloud-storage-volume-lib/etc/libconfig.toml"

	defaultEtcPath := GetConfPath()

	assert.Equal(t, expectedEtcPath, defaultEtcPath)
}

func TestGetConfPathWithEnv(t *testing.T) {
	t.Log("Testing GetEtcPath")
	os.Setenv("SECRET_CONFIG_PATH", "src/github.com/IBM/ibmcloud-storage-volume-lib/etc")
	expectedEtcPath := "src/github.com/IBM/ibmcloud-storage-volume-lib/etc/libconfig.toml"

	defaultEtcPath := GetConfPath()

	assert.Equal(t, expectedEtcPath, defaultEtcPath)
}

func TestGetDefaultConfPath(t *testing.T) {
	t.Log("Testing GetEtcPath")
	expectedEtcPath := "src/github.com/IBM/ibmcloud-storage-volume-lib/etc/libconfig.toml"

	defaultEtcPath := GetDefaultConfPath()

	assert.Equal(t, expectedEtcPath, defaultEtcPath)
}

func TestGetConfPathDir(t *testing.T) {
	t.Log("Testing GetConfPathDir")

	configPath := "test.toml"
	conf, _ := ReadConfig(configPath, testLogger)

	maxTimeout, _, _ := conf.VPC.GetTimeOutParameters()
	assert.Equal(t, maxTimeout, 120)
}

func TestGetTimeOutParameters(t *testing.T) {
	t.Log("Testing GetTimeOutParameters")

	configPath := "test.toml"
	conf, _ := ReadConfig(configPath, testLogger)

	maxTimeout, _, _ := conf.VPC.GetTimeOutParameters()
	assert.Equal(t, maxTimeout, 120)
}
