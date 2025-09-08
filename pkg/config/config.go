/*
 *  Copyright (c) 2022 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: DingoCli
 * Created Date: 2022-05-09
 * Author: chengyi (Cyber-SiKu)
 */
package config

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var (
	MDSApiV1 bool // is MDS API v2
)

const (
	FORMAT = "format"
	// global
	SHOWERROR                      = "showerror"
	VIPER_GLOBALE_SHOWERROR        = "global.showError"
	HTTPTIMEOUT                    = "httptimeout"
	VIPER_GLOBALE_HTTPTIMEOUT      = "global.httpTimeout"
	DEFAULT_HTTPTIMEOUT            = 100 * time.Millisecond
	RPCTIMEOUT                     = "rpctimeout"
	VIPER_GLOBALE_RPCTIMEOUT       = "global.rpcTimeout"
	DEFAULT_RPCTIMEOUT             = 30000 * time.Millisecond
	RPCRETRYTIMES                  = "rpcretrytimes"
	VIPER_GLOBALE_RPCRETRYTIMES    = "global.rpcRetryTimes"
	DEFAULT_RPCRETRYTIMES          = int32(5)
	RPCRETRYDElAY                  = "rpcretrydelay"
	VIPER_GLOBALE_RPCRETRYDELAY    = "global.rpcRetryDelay"
	DEFAULT_RPCRETRYDELAY          = 200 * time.Millisecond
	VERBOSE                        = "verbose"
	VIPER_GLOBALE_VERBOSE          = "global.verbose"
	DEFAULT_VERBOSE                = false
	VIPER_GLOBALE_MAX_CHANNEL_SIZE = "global.maxChannelSize"
	DEFAULT_MAX_CHANNEL_SIZE       = int32(4)
	VIPER_GLOBALE_MDS_API_VERSION  = "global.mds_api_version"
	LOGFILE                        = "logfile"
	VIPER_GLOBALE_LOGFILE          = "global.logfile"
	DEFAULT_LOGFILE                = "dingo.log"
	LOGLEVEL                       = "loglevel"
	VIPER_GLOBALE_LOGLEVEL         = "global.loglevel"
	DEFAULT_LOGLEVEL               = "info"
	LOGFMT                         = "logfmt"
	VIPER_GLOBALE_LOGFMT           = "global.logfmt"
	DEFAULT_LOGFMT                 = "text"
)

const (
	STATUS_SUBURI  = "/vars/dingofs_mds_status"
	VERSION_SUBURI = "/vars/dingo_version"
	ROOTINODEID    = uint64(1)
)

var (
	FLAFG_GLOBAL = []string{
		SHOWERROR, HTTPTIMEOUT, RPCTIMEOUT, RPCRETRYTIMES, VERBOSE,
	}
)

func InitConfig(confFile string) {
	// configure file priority
	// command line (--conf dingo.yaml) > environment variables(CONF=/opt/dingo.yaml) > default (~/.dingo/dingo.yaml)
	if confFile == "" {
		confFile = os.Getenv("CONF") //check environment variable
	}
	if confFile != "" {
		viper.SetConfigFile(confFile)
	} else {
		// using home directory and /etc/dingo as default configuration file path
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home + "/.dingo")
		viper.AddConfigPath("/etc/dingo")
		viper.SetConfigType("yaml")
		viper.SetConfigName("dingo")
	}

	// viper.SetDefault("format", "plain")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("config file name: %v", viper.ConfigFileUsed())
			cobra.CheckErr(err)
		}
	}

	// check mds api version, env MDS_API_VERSION priority > config file
	MDSApiV1 = false
	mdsApiVersion := os.Getenv("MDS_API_VERSION")
	if mdsApiVersion == "1" {
		MDSApiV1 = true
		return
	}
	if len(mdsApiVersion) == 0 { // env not set, check config file
		mdsApiVersion = viper.GetString(VIPER_GLOBALE_MDS_API_VERSION)
		if mdsApiVersion == "1" {
			MDSApiV1 = true
		}
	}
}

// global
// format
const (
	FORMAT_JSON  = "json"
	FORMAT_PLAIN = "plain"
	FORMAT_NOOUT = "noout"
)

func AddFormatFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("format", "", FORMAT_PLAIN, "output format (json|plain)")
	err := viper.BindPFlag("format", cmd.Flags().Lookup("format"))
	if err != nil {
		cobra.CheckErr(err)
	}
}

// http timeout
func AddHttpTimeoutFlag(cmd *cobra.Command) {
	cmd.Flags().Duration(HTTPTIMEOUT, DEFAULT_HTTPTIMEOUT, "http timeout")
	err := viper.BindPFlag(VIPER_GLOBALE_HTTPTIMEOUT, cmd.Flags().Lookup(HTTPTIMEOUT))
	if err != nil {
		cobra.CheckErr(err)
	}
}

// rpc time out [option]
func AddRpcTimeoutFlag(cmd *cobra.Command) {
	AddDurationOptionFlag(cmd, RPCTIMEOUT, "rpc timeout")
}

// rpc retry times
func AddRpcRetryTimesFlag(cmd *cobra.Command) {
	AddInt32OptionFlag(cmd, RPCRETRYTIMES, "rpc retry times")
}

// rpc retry delay
func AddRpcRetryDelayFlag(cmd *cobra.Command) {
	AddDurationOptionFlag(cmd, RPCRETRYDElAY, "rpc retry delay")
}

// channel size
func MaxChannelSize() int {
	return viper.GetInt(VIPER_GLOBALE_MAX_CHANNEL_SIZE)
}

// show errors
func AddShowErrorPFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool(SHOWERROR, false, "display all errors in command")
	err := viper.BindPFlag(VIPER_GLOBALE_SHOWERROR, cmd.PersistentFlags().Lookup(SHOWERROR))
	if err != nil {
		cobra.CheckErr(err)
	}
}

func AddLogfilePFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(LOGFILE, DEFAULT_LOGFILE, "logfile name")
	err := viper.BindPFlag(VIPER_GLOBALE_LOGFILE, cmd.PersistentFlags().Lookup(LOGFILE))
	if err != nil {
		cobra.CheckErr(err)
	}
}

func AddLoglevelPFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(LOGLEVEL, DEFAULT_LOGLEVEL, "loglevel (info|warn|error|debug|fatal|panic)")
	err := viper.BindPFlag(VIPER_GLOBALE_LOGLEVEL, cmd.PersistentFlags().Lookup(LOGLEVEL))
	if err != nil {
		cobra.CheckErr(err)
	}
}

func AddLogfmtPFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().String(LOGFMT, DEFAULT_LOGFMT, "log format (text|json)")
	err := viper.BindPFlag(VIPER_GLOBALE_LOGFMT, cmd.PersistentFlags().Lookup(LOGFMT))
	if err != nil {
		cobra.CheckErr(err)
	}
}

// Align the flag (changed) in the caller with the callee
func AlignFlagsValue(caller *cobra.Command, callee *cobra.Command, flagNames []string) {
	callee.Flags().VisitAll(func(flag *pflag.Flag) {
		index := slices.IndexFunc(flagNames, func(i string) bool {
			return flag.Name == i
		})
		if index == -1 {
			return
		}
		callerFlag := caller.Flag(flag.Name)
		if callerFlag != nil && callerFlag.Changed {
			if flag.Value.Type() == callerFlag.Value.Type() {
				flag.Value = callerFlag.Value
				flag.Changed = callerFlag.Changed
			} else {
				flag.Value.Set(callerFlag.Value.String())
				flag.Changed = callerFlag.Changed
			}
		}
	})
	// golobal flag
	for _, flagName := range FLAFG_GLOBAL {
		callerFlag := caller.Flag(flagName)
		if callerFlag != nil {
			if callee.Flag(flagName) != nil {
				callee.Flag(flagName).Value = callerFlag.Value
				callee.Flag(flagName).Changed = callerFlag.Changed
			} else {
				callee.Flags().AddFlag(callerFlag)
			}
		}
	}
}

func GetFlagChanged(cmd *cobra.Command, flagName string) bool {
	flag := cmd.Flag(flagName)
	if flag != nil {
		return flag.Changed
	}
	return false
}

const (
	IP_PORT_REGEX = "((\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5]):([0-9]|[1-9]\\d{1,3}|[1-5]\\d{4}|6[0-4]\\d{4}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5]))|(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])"
)

func IsValidAddr(addr string) bool {
	matched, err := regexp.MatchString(IP_PORT_REGEX, addr)
	if err != nil || !matched {
		return false
	}
	return true
}
