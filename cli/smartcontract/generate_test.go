package smartcontract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/manifest"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestGenerate(t *testing.T) {
	m := manifest.NewManifest("MyContract")
	m.ABI.Methods = append(m.ABI.Methods,
		manifest.Method{
			Name:       manifest.MethodDeploy,
			ReturnType: smartcontract.VoidType,
		},
		manifest.Method{
			Name: "sum",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("first", smartcontract.IntegerType),
				manifest.NewParameter("second", smartcontract.IntegerType),
			},
			ReturnType: smartcontract.IntegerType,
		},
		manifest.Method{
			Name: "sum", // overloaded method
			Parameters: []manifest.Parameter{
				manifest.NewParameter("first", smartcontract.IntegerType),
				manifest.NewParameter("second", smartcontract.IntegerType),
				manifest.NewParameter("third", smartcontract.IntegerType),
			},
			ReturnType: smartcontract.IntegerType,
		},
		manifest.Method{
			Name:       "sum3",
			Parameters: []manifest.Parameter{},
			ReturnType: smartcontract.IntegerType,
			Safe:       true,
		},
		manifest.Method{
			Name: "zum",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("type", smartcontract.IntegerType),
				manifest.NewParameter("typev", smartcontract.IntegerType),
				manifest.NewParameter("func", smartcontract.IntegerType),
			},
			ReturnType: smartcontract.IntegerType,
		},
		manifest.Method{
			Name: "justExecute",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("arr", smartcontract.ArrayType),
			},
			ReturnType: smartcontract.VoidType,
		},
		manifest.Method{
			Name:       "getPublicKey",
			Parameters: nil,
			ReturnType: smartcontract.PublicKeyType,
		},
		manifest.Method{
			Name: "otherTypes",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("ctr", smartcontract.Hash160Type),
				manifest.NewParameter("tx", smartcontract.Hash256Type),
				manifest.NewParameter("sig", smartcontract.SignatureType),
				manifest.NewParameter("data", smartcontract.AnyType),
			},
			ReturnType: smartcontract.BoolType,
		},
		manifest.Method{
			Name: "searchStorage",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("ctx", smartcontract.InteropInterfaceType),
			},
			ReturnType: smartcontract.InteropInterfaceType,
		},
		manifest.Method{
			Name: "getFromMap",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("intMap", smartcontract.MapType),
				manifest.NewParameter("indices", smartcontract.ArrayType),
			},
			ReturnType: smartcontract.ArrayType,
		},
		manifest.Method{
			Name: "doSomething",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("bytes", smartcontract.ByteArrayType),
				manifest.NewParameter("str", smartcontract.StringType),
			},
			ReturnType: smartcontract.InteropInterfaceType,
		},
		manifest.Method{
			Name:       "getBlockWrapper",
			Parameters: []manifest.Parameter{},
			ReturnType: smartcontract.InteropInterfaceType,
		},
		manifest.Method{
			Name: "myFunc",
			Parameters: []manifest.Parameter{
				manifest.NewParameter("in", smartcontract.MapType),
			},
			ReturnType: smartcontract.ArrayType,
		})

	manifestFile := filepath.Join(t.TempDir(), "manifest.json")
	outFile := filepath.Join(t.TempDir(), "out.go")

	rawManifest, err := json.Marshal(m)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestFile, rawManifest, os.ModePerm))

	h := util.Uint160{
		0x04, 0x08, 0x15, 0x16, 0x23, 0x42, 0x43, 0x44, 0x00, 0x01,
		0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF, 0x03, 0x04,
	}
	app := cli.NewApp()
	app.Commands = []cli.Command{generateWrapperCmd}

	rawCfg := `package: wrapper
hash: ` + h.StringLE() + `
overrides:
    searchStorage.ctx: storage.Context
    searchStorage: iterator.Iterator
    getFromMap.intMap: "map[string]int"
    getFromMap.indices: "[]string"
    getFromMap: "[]int"
    getBlockWrapper: ledger.Block
    myFunc.in: "map[int]github.com/heyitsme/mycontract.Input"
    myFunc: "[]github.com/heyitsme/mycontract.Output"
callflags:
    doSomething: ReadStates
`
	cfgPath := filepath.Join(t.TempDir(), "binding.yml")
	require.NoError(t, os.WriteFile(cfgPath, []byte(rawCfg), os.ModePerm))

	require.NoError(t, app.Run([]string{"", "generate-wrapper",
		"--manifest", manifestFile,
		"--config", cfgPath,
		"--out", outFile,
		"--hash", h.StringLE(),
	}))

	const expected = `// Package wrapper contains wrappers for MyContract contract.
package wrapper

import (
	"github.com/heyitsme/mycontract"
	"github.com/nspcc-dev/neo-go/pkg/interop"
	"github.com/nspcc-dev/neo-go/pkg/interop/contract"
	"github.com/nspcc-dev/neo-go/pkg/interop/iterator"
	"github.com/nspcc-dev/neo-go/pkg/interop/native/ledger"
	"github.com/nspcc-dev/neo-go/pkg/interop/neogointernal"
	"github.com/nspcc-dev/neo-go/pkg/interop/storage"
)

// Hash contains contract hash in big-endian form.
const Hash = "\x04\x08\x15\x16\x23\x42\x43\x44\x00\x01\xca\xfe\xba\xbe\xde\xad\xbe\xef\x03\x04"

// Sum invokes ` + "`sum`" + ` method of contract.
func Sum(first int, second int) int {
	return neogointernal.CallWithToken(Hash, "sum", int(contract.All), first, second).(int)
}

// Sum_3 invokes ` + "`sum`" + ` method of contract.
func Sum_3(first int, second int, third int) int {
	return neogointernal.CallWithToken(Hash, "sum", int(contract.All), first, second, third).(int)
}

// Sum3 invokes ` + "`sum3`" + ` method of contract.
func Sum3() int {
	return neogointernal.CallWithToken(Hash, "sum3", int(contract.ReadOnly)).(int)
}

// Zum invokes ` + "`zum`" + ` method of contract.
func Zum(typev int, typev_ int, funcv int) int {
	return neogointernal.CallWithToken(Hash, "zum", int(contract.All), typev, typev_, funcv).(int)
}

// JustExecute invokes ` + "`justExecute`" + ` method of contract.
func JustExecute(arr []any) {
	neogointernal.CallWithTokenNoRet(Hash, "justExecute", int(contract.All), arr)
}

// GetPublicKey invokes ` + "`getPublicKey`" + ` method of contract.
func GetPublicKey() interop.PublicKey {
	return neogointernal.CallWithToken(Hash, "getPublicKey", int(contract.All)).(interop.PublicKey)
}

// OtherTypes invokes ` + "`otherTypes`" + ` method of contract.
func OtherTypes(ctr interop.Hash160, tx interop.Hash256, sig interop.Signature, data any) bool {
	return neogointernal.CallWithToken(Hash, "otherTypes", int(contract.All), ctr, tx, sig, data).(bool)
}

// SearchStorage invokes ` + "`searchStorage`" + ` method of contract.
func SearchStorage(ctx storage.Context) iterator.Iterator {
	return neogointernal.CallWithToken(Hash, "searchStorage", int(contract.All), ctx).(iterator.Iterator)
}

// GetFromMap invokes ` + "`getFromMap`" + ` method of contract.
func GetFromMap(intMap map[string]int, indices []string) []int {
	return neogointernal.CallWithToken(Hash, "getFromMap", int(contract.All), intMap, indices).([]int)
}

// DoSomething invokes ` + "`doSomething`" + ` method of contract.
func DoSomething(bytes []byte, str string) any {
	return neogointernal.CallWithToken(Hash, "doSomething", int(contract.ReadStates), bytes, str).(any)
}

// GetBlockWrapper invokes ` + "`getBlockWrapper`" + ` method of contract.
func GetBlockWrapper() ledger.Block {
	return neogointernal.CallWithToken(Hash, "getBlockWrapper", int(contract.All)).(ledger.Block)
}

// MyFunc invokes ` + "`myFunc`" + ` method of contract.
func MyFunc(in map[int]mycontract.Input) []mycontract.Output {
	return neogointernal.CallWithToken(Hash, "myFunc", int(contract.All), in).([]mycontract.Output)
}
`

	data, err := os.ReadFile(outFile)
	require.NoError(t, err)
	require.Equal(t, expected, string(data))
}

func TestGenerateValidPackageName(t *testing.T) {
	m := manifest.NewManifest("My space\tcontract")
	m.ABI.Methods = append(m.ABI.Methods,
		manifest.Method{
			Name:       "get",
			Parameters: []manifest.Parameter{},
			ReturnType: smartcontract.IntegerType,
			Safe:       true,
		},
	)

	manifestFile := filepath.Join(t.TempDir(), "manifest.json")
	outFile := filepath.Join(t.TempDir(), "out.go")

	rawManifest, err := json.Marshal(m)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestFile, rawManifest, os.ModePerm))

	h := util.Uint160{
		0x04, 0x08, 0x15, 0x16, 0x23, 0x42, 0x43, 0x44, 0x00, 0x01,
		0xCA, 0xFE, 0xBA, 0xBE, 0xDE, 0xAD, 0xBE, 0xEF, 0x03, 0x04,
	}
	app := cli.NewApp()
	app.Commands = []cli.Command{generateWrapperCmd, generateRPCWrapperCmd}
	require.NoError(t, app.Run([]string{"", "generate-wrapper",
		"--manifest", manifestFile,
		"--out", outFile,
		"--hash", "0x" + h.StringLE(),
	}))

	data, err := os.ReadFile(outFile)
	require.NoError(t, err)
	require.Equal(t, `// Package myspacecontract contains wrappers for My space	contract contract.
package myspacecontract

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/contract"
	"github.com/nspcc-dev/neo-go/pkg/interop/neogointernal"
)

// Hash contains contract hash in big-endian form.
const Hash = "\x04\x08\x15\x16\x23\x42\x43\x44\x00\x01\xca\xfe\xba\xbe\xde\xad\xbe\xef\x03\x04"

// Get invokes `+"`get`"+` method of contract.
func Get() int {
	return neogointernal.CallWithToken(Hash, "get", int(contract.ReadOnly)).(int)
}
`, string(data))
	require.NoError(t, app.Run([]string{"", "generate-rpcwrapper",
		"--manifest", manifestFile,
		"--out", outFile,
		"--hash", "0x" + h.StringLE(),
	}))

	data, err = os.ReadFile(outFile)
	require.NoError(t, err)
	require.Equal(t, `// Package myspacecontract contains RPC wrappers for My space	contract contract.
package myspacecontract

import (
	"github.com/nspcc-dev/neo-go/pkg/neorpc/result"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/unwrap"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"math/big"
)

// Hash contains contract hash.
var Hash = util.Uint160{0x4, 0x8, 0x15, 0x16, 0x23, 0x42, 0x43, 0x44, 0x0, 0x1, 0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef, 0x3, 0x4}

// Invoker is used by ContractReader to call various safe methods.
type Invoker interface {
	Call(contract util.Uint160, operation string, params ...any) (*result.Invoke, error)
}

// ContractReader implements safe contract methods.
type ContractReader struct {
	invoker Invoker
	hash util.Uint160
}

// NewReader creates an instance of ContractReader using Hash and the given Invoker.
func NewReader(invoker Invoker) *ContractReader {
	var hash = Hash
	return &ContractReader{invoker, hash}
}

// Get invokes `+"`get`"+` method of contract.
func (c *ContractReader) Get() (*big.Int, error) {
	return unwrap.BigInt(c.invoker.Call(c.hash, "get"))
}
`, string(data))
}

// rewriteExpectedOutputs denotes whether expected output files should be rewritten
// for TestGenerateRPCBindings and TestAssistedRPCBindings.
const rewriteExpectedOutputs = false

func TestGenerateRPCBindings(t *testing.T) {
	tmpDir := t.TempDir()
	app := cli.NewApp()
	app.Commands = []cli.Command{generateWrapperCmd, generateRPCWrapperCmd}

	var checkBinding = func(manifest string, hash string, good string) {
		t.Run(manifest, func(t *testing.T) {
			outFile := filepath.Join(tmpDir, "out.go")
			require.NoError(t, app.Run([]string{"", "generate-rpcwrapper",
				"--manifest", manifest,
				"--out", outFile,
				"--hash", hash,
			}))

			data, err := os.ReadFile(outFile)
			require.NoError(t, err)
			data = bytes.ReplaceAll(data, []byte("\r"), []byte{}) // Windows.
			if rewriteExpectedOutputs {
				require.NoError(t, os.WriteFile(good, data, os.ModePerm))
			} else {
				expected, err := os.ReadFile(good)
				require.NoError(t, err)
				expected = bytes.ReplaceAll(expected, []byte("\r"), []byte{}) // Windows.
				require.Equal(t, string(expected), string(data))
			}
		})
	}

	checkBinding(filepath.Join("testdata", "nex", "nex.manifest.json"),
		"0xa2a67f09e8cf22c6bfd5cea24adc0f4bf0a11aa8",
		filepath.Join("testdata", "nex", "nex.go"))
	checkBinding(filepath.Join("testdata", "nameservice", "nns.manifest.json"),
		"0x50ac1c37690cc2cfc594472833cf57505d5f46de",
		filepath.Join("testdata", "nameservice", "nns.go"))
	checkBinding(filepath.Join("testdata", "gas", "gas.manifest.json"),
		"0xd2a4cff31913016155e38e474a2c06d08be276cf",
		filepath.Join("testdata", "gas", "gas.go"))
	checkBinding(filepath.Join("testdata", "verifyrpc", "verify.manifest.json"),
		"0x00112233445566778899aabbccddeeff00112233",
		filepath.Join("testdata", "verifyrpc", "verify.go"))
	checkBinding(filepath.Join("testdata", "nonepiter", "iter.manifest.json"),
		"0x00112233445566778899aabbccddeeff00112233",
		filepath.Join("testdata", "nonepiter", "iter.go"))

	require.False(t, rewriteExpectedOutputs)
}

func TestAssistedRPCBindings(t *testing.T) {
	tmpDir := t.TempDir()
	app := cli.NewApp()
	app.Commands = NewCommands()

	var checkBinding = func(source string, hasDefinedHash bool, guessEventTypes bool, suffix ...string) {
		testName := source
		if len(suffix) != 0 {
			testName += suffix[0]
		}
		testName += fmt.Sprintf(", predefined hash: %t", hasDefinedHash)
		t.Run(testName, func(t *testing.T) {
			configFile := filepath.Join(source, "config.yml")
			expectedFile := filepath.Join(source, "rpcbindings.out")
			if len(suffix) != 0 {
				configFile = filepath.Join(source, "config"+suffix[0]+".yml")
				expectedFile = filepath.Join(source, "rpcbindings"+suffix[0]+".out")
			} else if !hasDefinedHash {
				expectedFile = filepath.Join(source, "rpcbindings_dynamic_hash.out")
			}
			manifestF := filepath.Join(tmpDir, "manifest.json")
			bindingF := filepath.Join(tmpDir, "binding.yml")
			nefF := filepath.Join(tmpDir, "out.nef")
			cmd := []string{"", "contract", "compile",
				"--in", source,
				"--config", configFile,
				"--manifest", manifestF,
				"--bindings", bindingF,
				"--out", nefF,
			}
			if guessEventTypes {
				cmd = append(cmd, "--guess-eventtypes")
			}
			require.NoError(t, app.Run(cmd))

			cmds := []string{"", "contract", "generate-rpcwrapper",
				"--config", bindingF,
				"--manifest", manifestF,
				"--out", expectedFile,
			}
			if hasDefinedHash {
				cmds = append(cmds, "--hash", "0x00112233445566778899aabbccddeeff00112233")
			}
			require.NoError(t, app.Run(cmds))

			data, err := os.ReadFile(expectedFile)
			require.NoError(t, err)
			data = bytes.ReplaceAll(data, []byte("\r"), []byte{}) // Windows.
			if rewriteExpectedOutputs {
				require.NoError(t, os.WriteFile(expectedFile, data, os.ModePerm))
			} else {
				expected, err := os.ReadFile(expectedFile)
				require.NoError(t, err)
				expected = bytes.ReplaceAll(expected, []byte("\r"), []byte{}) // Windows.
				require.Equal(t, string(expected), string(data))
			}
		})
	}

	for _, hasDefinedHash := range []bool{true, false} {
		checkBinding(filepath.Join("testdata", "types"), hasDefinedHash, false)
		checkBinding(filepath.Join("testdata", "structs"), hasDefinedHash, false)
	}
	checkBinding(filepath.Join("testdata", "notifications"), true, false)
	checkBinding(filepath.Join("testdata", "notifications"), true, false, "_extended")
	checkBinding(filepath.Join("testdata", "notifications"), true, true, "_guessed")

	require.False(t, rewriteExpectedOutputs)
}

func TestGenerate_Errors(t *testing.T) {
	app := cli.NewApp()
	app.Commands = []cli.Command{generateWrapperCmd}
	app.ExitErrHandler = func(*cli.Context, error) {}

	checkError := func(t *testing.T, msg string, args ...string) {
		// cli.ExitError doesn't implement wraping properly, so we check for an error message.
		err := app.Run(append([]string{"", "generate-wrapper"}, args...))
		require.True(t, strings.Contains(err.Error(), msg), "got: %v", err)
	}
	t.Run("invalid hash", func(t *testing.T) {
		checkError(t, "invalid contract hash", "--hash", "xxx", "--manifest", "yyy", "--out", "zzz")
	})
	t.Run("missing manifest argument", func(t *testing.T) {
		checkError(t, "Required flag \"manifest\" not set", "--hash", util.Uint160{}.StringLE(), "--out", "zzz")
	})
	t.Run("missing manifest file", func(t *testing.T) {
		checkError(t, "can't read contract manifest", "--manifest", "notexists", "--hash", util.Uint160{}.StringLE(), "--out", "zzz")
	})
	t.Run("empty manifest", func(t *testing.T) {
		manifestFile := filepath.Join(t.TempDir(), "invalid.json")
		require.NoError(t, os.WriteFile(manifestFile, []byte("[]"), os.ModePerm))
		checkError(t, "json: cannot unmarshal array into Go value of type manifest.Manifest", "--manifest", manifestFile, "--hash", util.Uint160{}.StringLE(), "--out", "zzz")
	})
	t.Run("invalid manifest", func(t *testing.T) {
		manifestFile := filepath.Join(t.TempDir(), "invalid.json")
		m := manifest.NewManifest("MyContract") // no methods
		rawManifest, err := json.Marshal(m)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(manifestFile, rawManifest, os.ModePerm))
		checkError(t, "ABI: no methods", "--manifest", manifestFile, "--hash", util.Uint160{}.StringLE(), "--out", "zzz")
	})

	manifestFile := filepath.Join(t.TempDir(), "manifest.json")
	m := manifest.NewManifest("MyContract")
	m.ABI.Methods = append(m.ABI.Methods, manifest.Method{
		Name:       "method0",
		Offset:     0,
		ReturnType: smartcontract.AnyType,
		Safe:       true,
	})
	rawManifest, err := json.Marshal(m)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestFile, rawManifest, os.ModePerm))

	t.Run("missing config", func(t *testing.T) {
		checkError(t, "can't read config file",
			"--manifest", manifestFile, "--hash", util.Uint160{}.StringLE(),
			"--config", filepath.Join(t.TempDir(), "not.exists.yml"), "--out", "zzz")
	})
	t.Run("invalid config", func(t *testing.T) {
		rawCfg := `package: wrapper
callflags:
    someFunc: ReadSometimes 
`
		cfgPath := filepath.Join(t.TempDir(), "binding.yml")
		require.NoError(t, os.WriteFile(cfgPath, []byte(rawCfg), os.ModePerm))

		checkError(t, "can't parse config file",
			"--manifest", manifestFile, "--hash", util.Uint160{}.StringLE(),
			"--config", cfgPath, "--out", "zzz")
	})
}

func TestCompile_GuessEventTypes(t *testing.T) {
	app := cli.NewApp()
	app.Commands = NewCommands()
	app.ExitErrHandler = func(*cli.Context, error) {}

	checkError := func(t *testing.T, msg string, args ...string) {
		// cli.ExitError doesn't implement wraping properly, so we check for an error message.
		err := app.Run(args)
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), msg), "got: %v", err)
	}
	check := func(t *testing.T, source string, expectedErrText string) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(source, "invalid.yml")
		manifestF := filepath.Join(tmpDir, "invalid.manifest.json")
		bindingF := filepath.Join(tmpDir, "invalid.binding.yml")
		nefF := filepath.Join(tmpDir, "invalid.out.nef")
		cmd := []string{"", "contract", "compile",
			"--in", source,
			"--config", configFile,
			"--manifest", manifestF,
			"--bindings", bindingF,
			"--out", nefF,
			"--guess-eventtypes",
		}
		checkError(t, expectedErrText, cmd...)
	}

	t.Run("not declared in manifest", func(t *testing.T) {
		check(t, filepath.Join("testdata", "invalid5"), "inconsistent usages of event `Non declared event`: not declared in the contract config")
	})
	t.Run("invalid number of params", func(t *testing.T) {
		check(t, filepath.Join("testdata", "invalid6"), "inconsistent usages of event `SomeEvent` against config: number of params mismatch: 2 vs 1")
	})
	/*
		// TODO: this on is a controversial one. If event information is provided in the config file, then conversion code
		// will be emitted by the compiler according to the parameter type provided via config. Thus, we can be sure that
		// either event parameter has the type specified in the config file or the execution of the contract will fail.
		// Thus, this testcase is always failing (no compilation error occures).
		// Question: do we want to compare `RealType` of the emitted parameter with the one expected in the manifest?
		t.Run("SC parameter type mismatch", func(t *testing.T) {
			check(t, filepath.Join("testdata", "invalid7"), "inconsistent usages of event `SomeEvent` against config: number of params mismatch: 2 vs 1")
		})
	*/
	t.Run("extended types mismatch", func(t *testing.T) {
		check(t, filepath.Join("testdata", "invalid8"), "inconsistent usages of event `SomeEvent`: extended type of param #0 mismatch")
	})
	t.Run("named types redeclare", func(t *testing.T) {
		check(t, filepath.Join("testdata", "invalid9"), "configured declared named type intersects with the contract's one: `invalid9.NamedStruct`")
	})
}
