// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package mongo_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/juju/clock"
	"github.com/juju/clock/testclock"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/network"
	"github.com/juju/juju/internal/mongo"
	"github.com/juju/juju/internal/mongo/mongotest"
	"github.com/juju/juju/internal/packaging"
	"github.com/juju/juju/internal/service/common"
	"github.com/juju/juju/internal/service/snap"
	coretesting "github.com/juju/juju/testing"
)

type MongoSuite struct {
	coretesting.BaseSuite

	clock            clock.Clock
	mongodConfigPath string

	mongoSnapService *mongotest.MockMongoSnapService
}

var _ = gc.Suite(&MongoSuite{})

var testInfo = struct {
	StatePort    int
	Cert         string
	PrivateKey   string
	SharedSecret string
}{
	StatePort:    25252,
	Cert:         "foobar-cert",
	PrivateKey:   "foobar-privkey",
	SharedSecret: "foobar-sharedsecret",
}

func makeEnsureServerParams(dataDir, configDir string) mongo.EnsureServerParams {
	return mongo.EnsureServerParams{
		StatePort:    testInfo.StatePort,
		Cert:         testInfo.Cert,
		PrivateKey:   testInfo.PrivateKey,
		SharedSecret: testInfo.SharedSecret,

		DataDir:           dataDir,
		ConfigDir:         configDir,
		JujuDBSnapChannel: "latest",

		OplogSize: 1,
	}
}

func (s *MongoSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	testing.PatchExecutable(c, s, "juju-db.mongod", "#!/bin/bash\n\nprintf %s 'db version v6.6.6'\n")
	jujuMongodPath, err := exec.LookPath("juju-db.mongod")
	c.Assert(err, jc.ErrorIsNil)
	s.PatchValue(&mongo.JujuDbSnapMongodPath, jujuMongodPath)

	// Patch "df" such that it always reports there's 1MB free.
	s.PatchValue(mongo.AvailSpace, func(dir string) (float64, error) {
		info, err := os.Stat(dir)
		if err != nil {
			return 0, err
		}
		if info.IsDir() {
			return 1, nil

		}
		return 0, fmt.Errorf("not a directory")
	})
	s.PatchValue(mongo.SmallOplogSizeMB, 1)

	s.clock = testclock.NewClock(time.Now())
}

func (s *MongoSuite) setupMocks(c *gc.C) *gomock.Controller {
	ctrl := gomock.NewController(c)

	s.mongoSnapService = mongotest.NewMockMongoSnapService(ctrl)

	return ctrl
}

func (s *MongoSuite) expectInstallMongoSnap() {
	mExp := s.mongoSnapService.EXPECT()
	mExp.Name().Return("not-juju-db")
	mExp.Install().Return(nil)
	mExp.ConfigOverride().Return(nil)
	mExp.Start().Return(nil).AnyTimes()
	mExp.Running().Return(true, nil).AnyTimes()

	s.PatchValue(mongo.NewSnapService, func(mainSnap, serviceName string, conf common.Conf, snapPath, configDir, channel string, confinementPolicy snap.ConfinementPolicy, backgroundServices []snap.BackgroundService, prerequisites []snap.Installable) (mongo.MongoSnapService, error) {
		return s.mongoSnapService, nil
	})
}

func (s *MongoSuite) assertTLSKeyFile(c *gc.C, dataDir string) {
	contents, err := os.ReadFile(mongo.SSLKeyPath(dataDir))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(contents), gc.Equals, testInfo.Cert+"\n"+testInfo.PrivateKey)
}

func (s *MongoSuite) assertSharedSecretFile(c *gc.C, dataDir string) {
	contents, err := os.ReadFile(mongo.SharedSecretPath(dataDir))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(contents), gc.Equals, testInfo.SharedSecret)
}

func (s *MongoSuite) assertMongoConfigFile(c *gc.C, dataDir string, ipV6 bool) {
	contents, err := os.ReadFile(s.mongodConfigPath)
	c.Assert(err, jc.ErrorIsNil)
	part1 := fmt.Sprintf(`
# WARNING
# autogenerated by juju on .*
# manual changes to this file are likely to be overwritten
auth = true
bind_ip_all = true
dbpath = %s/db`[1:], dataDir)
	if ipV6 {
		part1 += "\nipv6 = true"
	}

	part2 := fmt.Sprintf(`
journal = true
keyFile = %s/shared-secret
logpath = %s/logs/mongodb.log
oplogSize = 1
port = 25252
quiet = true
replSet = juju
slowms = 1000
storageEngine = wiredTiger
tlsCertificateKeyFile = %s/server.pem
tlsCertificateKeyFilePassword=ignored
tlsMode = requireTLS`, dataDir, dataDir, dataDir)

	c.Assert(string(contents), gc.Matches, part1+part2)
}

func (s *MongoSuite) TestEnsureServerInstalled(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectInstallMongoSnap()

	dataDir := s.assertEnsureServerIPv6(c, true)

	s.assertTLSKeyFile(c, dataDir)
	s.assertSharedSecretFile(c, dataDir)
	s.assertMongoConfigFile(c, dataDir, true)

	// make sure that we log the version of mongodb as we get ready to
	// start it
	tlog := c.GetTestLog()
	anyExp := `(.|\n)*`
	start := "^" + anyExp
	tail := anyExp + "$"
	c.Assert(tlog, gc.Matches, start+`using mongod: .*mongod --version:\sdb version v\d\.\d\.\d`+tail)
}

func (s *MongoSuite) TestEnsureServerInstalledNoIPv6(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectInstallMongoSnap()

	dataDir := s.assertEnsureServerIPv6(c, false)

	s.assertTLSKeyFile(c, dataDir)
	s.assertSharedSecretFile(c, dataDir)
	s.assertMongoConfigFile(c, dataDir, false)
}

func (s *MongoSuite) TestEnsureServerInstalledSetsSysctlValues(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectInstallMongoSnap()

	dataDir := c.MkDir()
	dataFilePath := filepath.Join(dataDir, "mongoKernelTweaks")
	dataFile, err := os.Create(dataFilePath)
	c.Assert(err, jc.ErrorIsNil)
	_, err = dataFile.WriteString("original value")
	c.Assert(err, jc.ErrorIsNil)
	_ = dataFile.Close()

	testing.PatchExecutableAsEchoArgs(c, s, "snap")

	contents, err := os.ReadFile(dataFilePath)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(contents), gc.Equals, "original value")

	configDir := c.MkDir()
	err = mongo.SysctlEditableEnsureServer(
		context.Background(),
		makeEnsureServerParams(dataDir, configDir),
		map[string]string{dataFilePath: "new value"},
	)
	c.Assert(err, jc.ErrorIsNil)

	contents, err = os.ReadFile(dataFilePath)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(contents), gc.Equals, "new value")
}

func (s *MongoSuite) TestEnsureServerInstalledError(c *gc.C) {
	defer s.setupMocks(c).Finish()

	dataDir := c.MkDir()
	configDir := c.MkDir()

	testing.PatchExecutableAsEchoArgs(c, s, "snap")

	failure := errors.New("boom")
	s.PatchValue(mongo.InstallMongo, func(dep packaging.Dependency, series string) error {
		return failure
	})

	err := mongo.EnsureServerInstalled(context.Background(), makeEnsureServerParams(dataDir, configDir))
	c.Assert(errors.Cause(err), gc.Equals, failure)
}

func (s *MongoSuite) assertEnsureServerIPv6(c *gc.C, ipv6 bool) string {
	dataDir := c.MkDir()
	configDir := c.MkDir()
	s.mongodConfigPath = filepath.Join(dataDir, "juju-db.config")

	testing.PatchExecutableAsEchoArgs(c, s, "snap")

	s.PatchValue(mongo.SupportsIPv6, func() bool {
		return ipv6
	})
	testParams := makeEnsureServerParams(dataDir, configDir)
	err := mongo.EnsureServerInstalled(context.Background(), testParams)
	c.Assert(err, jc.ErrorIsNil)
	return dataDir
}

func (s *MongoSuite) TestNoMongoDir(c *gc.C) {
	defer s.setupMocks(c).Finish()
	s.expectInstallMongoSnap()

	// Make a non-existent directory that can nonetheless be
	// created.
	testing.PatchExecutableAsEchoArgs(c, s, "snap")

	dataDir := filepath.Join(c.MkDir(), "dir", "data")
	configDir := c.MkDir()
	err := mongo.EnsureServerInstalled(context.Background(), makeEnsureServerParams(dataDir, configDir))
	c.Check(err, jc.ErrorIsNil)

	_, err = os.Stat(filepath.Join(dataDir, "db"))
	c.Assert(err, jc.ErrorIsNil)
}

func (s *MongoSuite) TestSelectPeerAddress(c *gc.C) {
	addresses := network.ProviderAddresses{
		network.NewMachineAddress("126.0.0.1", network.WithScope(network.ScopeMachineLocal)).AsProviderAddress(),
		network.NewMachineAddress("10.0.0.1", network.WithScope(network.ScopeCloudLocal)).AsProviderAddress(),
		network.NewMachineAddress("8.8.8.8", network.WithScope(network.ScopePublic)).AsProviderAddress(),
	}

	address := mongo.SelectPeerAddress(addresses)
	c.Assert(address, gc.Equals, "10.0.0.1")
}

func (s *MongoSuite) TestGenerateSharedSecret(c *gc.C) {
	secret, err := mongo.GenerateSharedSecret()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(secret, gc.HasLen, 1024)
	_, err = base64.StdEncoding.DecodeString(secret)
	c.Assert(err, jc.ErrorIsNil)
}
