package builders

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/minio/blake2b-simd"

	"github.com/chenjianmei111/test-vectors/schema"
)

type Metadata = schema.Metadata

// Generator is a batch generator and organizer of test vectors.
//
// Test vector scripts are simple programs (main function). Test vector scripts
// can delegate to the Generator to handle the execution, reporting and capture
// of emitted test vectors into files.
//
// Generator supports the following CLI flags:
//
//  -o <directory>
//		directory where test vector JSON files will be saved; if omitted,
//		vectors will be written to stdout.
//
//  -u
//		update any existing test vector files in the output directory IF their
//		content has changed. Note `_meta` is ignored when checking equality.
//
//  -f
//		force regeneration and overwrite any existing vectors in the output
//		directory.
//
//  -i <include regex>
//		regex inclusion filter to select a subset of vectors to execute; matched
//		against the vector's ID.
//
//  TODO
//  -v <protocol versions, comma-separated>
//      protocol version variants to generate; if not provided, all supported
//      protocol versions as declared by the vector will be attempted
//
// Scripts can bundle test vectors into "groups". The generator will execute
// each group in parallel, and will write each vector in a file:
// <output_dir>/<group>--<vector_id>.json
type Generator struct {
	OutputPath    string
	Mode          OverwriteMode
	IncludeFilter *regexp.Regexp

	wg sync.WaitGroup
}

const brokenVectorPrefix = "x--"

// OverwriteMode is the mode used when overwriting existing test vector files.
type OverwriteMode int

const (
	// OverwriteNone will not overwrite existing test vector files.
	OverwriteNone OverwriteMode = iota
	// OverwriteUpdate will update test vector files if they're different.
	OverwriteUpdate
	// OverwriteForce will force overwrite the vector files.
	OverwriteForce
)

var GenscriptCommit = "dirty"

// genData is the generation data to stamp into vectors.
var genData = []schema.GenerationData{
	{
		Source:  "genscript",
		Version: GenscriptCommit,
	},
}

func init() {
	genData = append(genData, getBuildInfo()...)
}

func getBuildInfo() []schema.GenerationData {
	deps := []string{"github.com/chenjianmei111/lotus", "github.com/chenjianmei111/specs-actors"}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("cant read build info")
	}

	var result []schema.GenerationData

	for _, v := range bi.Deps {
		for _, dep := range deps {
			if strings.HasPrefix(v.Path, dep) {
				result = append(result, schema.GenerationData{Source: v.Path, Version: v.Version})
			}
		}
	}

	return result
}

type VectorDef struct {
	Metadata *schema.Metadata
	Selector schema.Selector

	// SupportedVersions enumerates the versions this vector is supported
	// against. If nil or empty, this vector is valid for all known versions
	// (as per KnownProtocolVersions).
	SupportedVersions []ProtocolVersion

	// Hints are arbitrary flags that convey information to the driver.
	// Use hints to express facts like this vector is knowingly incorrect
	// (e.g. when the reference implementation is broken), or that drivers
	// should negate the postconditions (i.e. test that they are NOT the ones
	// expressed in the vector), etc.
	//
	// Refer to the schema.Hint* constants for common hints.
	Hints []string

	// Mode tunes certain elements of how the generation and assertion of
	// a test vector will be conducted, such as being lenient to assertion
	// failures when a vector is knowingly incorrect. Refer to the Mode*
	// constants for further information.
	Mode Mode

	// MessageFunc if non-nil, declares this vector as a message-class vector,
	// generated by the specified builder.
	MessageFunc func(v *MessageVectorBuilder)

	// TipsetFunc if non-nil, declares this vector as a tipset-class vector,
	// generated by the specified builder.
	TipsetFunc func(v *TipsetVectorBuilder)
}

func NewGenerator() *Generator {
	// Consume CLI parameters.
	var outputDir string
	const outputDirUsage = "directory where test vector JSON files will be saved; if omitted, vectors will be written to stdout."
	flag.StringVar(&outputDir, "o", "", outputDirUsage)
	flag.StringVar(&outputDir, "out", "", outputDirUsage)

	var update bool
	const updateUsage = "update any existing test vector files in the output directory IF their content has changed. Note `_meta` is ignored when checking equality."
	flag.BoolVar(&update, "u", false, updateUsage)
	flag.BoolVar(&update, "update", false, updateUsage)

	var force bool
	const forceUsage = "force regeneration and overwrite any existing vectors in the output directory."
	flag.BoolVar(&force, "f", false, forceUsage)
	flag.BoolVar(&force, "force", false, forceUsage)

	var includeFilter string
	const includeFilterUsage = "regex inclusion filter to select a subset of vectors to execute; matched against the vector's ID or the vector group name."
	flag.StringVar(&includeFilter, "i", "", includeFilterUsage)
	flag.StringVar(&includeFilter, "include", "", includeFilterUsage)

	flag.Parse()

	var mode OverwriteMode
	switch {
	case force:
		mode = OverwriteForce
	case update:
		mode = OverwriteUpdate
	default:
		mode = OverwriteNone
	}

	gen := Generator{Mode: mode}

	// If output directory is provided, we ensure it exists, or create it.
	// Else, we'll output to stdout.
	if outputDir != "" {
		err := ensureDirectory(outputDir)
		if err != nil {
			log.Fatal(err)
		}
		gen.OutputPath = outputDir
	}

	// If a filter has been provided, compile it into a regex.
	if includeFilter != "" {
		exp, err := regexp.Compile(includeFilter)
		if err != nil {
			log.Fatalf("supplied inclusion filter regex %s is invalid: %s", includeFilter, err)
		}
		gen.IncludeFilter = exp
	}

	return &gen
}

func (g *Generator) Close() {
	g.wg.Wait()
}

func (g *Generator) Group(group string, vectors ...*VectorDef) {
	// validate and filter vectors.
	var generate []*VectorDef
	for _, v := range vectors {
		if v.MessageFunc != nil && v.TipsetFunc != nil {
			panic(fmt.Sprintf("vector with id %s had more than one function", v.Metadata.ID))
		}
		if v.MessageFunc == nil && v.TipsetFunc == nil {
			panic(fmt.Sprintf("vector with id %s had no functions", v.Metadata.ID))
		}
		if id := v.Metadata.ID; g.IncludeFilter != nil && !g.IncludeFilter.MatchString(id) && !g.IncludeFilter.MatchString(group) {
			log.Printf("skipping %s: does not match inclusion filter", id)
			continue
		}
		generate = append(generate, v)
	}

	if len(generate) == 0 {
		log.Printf("no vectors to generate for group %s", group)
		return
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		var tmpDir string
		if g.OutputPath != "" {
			dir, err := ioutil.TempDir("", group)
			if err != nil {
				log.Printf("failed to create temp output directory: %s", err)
				return
			}
			defer func() {
				if err := os.RemoveAll(dir); err != nil {
					log.Printf("failed to remove temp output directory: %s", err)
				}
			}()
			tmpDir = dir
		}

		var wg sync.WaitGroup
		for _, item := range generate {
			wg.Add(1)
			go func(item VectorDef) {
				defer wg.Done()

				// generate variants.
				variants := g.generateVariants(item)

				// print to stdout.
				if g.OutputPath == "" {
					for _, v := range variants {
						fmt.Println(string(v.MustMarshalJSON()))
					}
					return
				}

				for _, v := range variants {
					var (
						tmp      = vectorPath(tmpDir, group, &item, v)
						existing = vectorPath(g.OutputPath, group, &item, v)
					)

					out, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
					if err != nil {
						log.Printf("failed to open file for writing %s: %s", tmp, err)
						continue
					}

					enc := json.NewEncoder(out)
					enc.SetIndent("", "\t")
					if err := enc.Encode(v); err != nil {
						log.Printf("failed to write json into file %s: %s", tmp, err)
						continue
					}

					_ = out.Close()

					switch _, err := os.Stat(existing); {
					case err == nil:
						// file exists.
						if g.Mode == OverwriteForce {
							// skip straight to the writing.
							break
						}
						eql, err := g.vectorsEqual(tmp, existing)
						if err != nil {
							log.Printf("failed to check new vs existing vector equality: %s", err)
							continue
						}
						if eql {
							log.Printf("not writing %s: no changes", existing)
							continue
						}
						if g.Mode == OverwriteNone {
							// no overwrite requested, warn that the vector has changed but we're refusing to overwrite.
							log.Printf("⚠️ WARNING: not writing %s: vector changed, use -u or -f to overwrite", existing)
							continue
						}
						if g.Mode == OverwriteUpdate {
							// acknowledge that the vector has changed and we will overwrite
							log.Printf("test vector exists, is not equal, and update was requested, overwriting: %s", existing)
						}
					case os.IsNotExist(err):
						// file doesn't exist, write it.
					default:
						log.Printf("failed unexpectedly while checking if file exists: %s; err: %s", existing, err)
						continue
					}

					// Move vector from tmp dir to final location
					if err := os.Rename(tmp, existing); err != nil {
						log.Printf("failed to move generated test vector: %s", err)
					}

					// If this vector was broken and became fixed, then remove the broken
					// vector file (and vice versa).
					if err := removePrevious(g.OutputPath, group, &item, v); err != nil {
						log.Printf("failed to remove previously broken vector: %s", err)
					}
					log.Printf("wrote test vector: %s", existing)
				}

			}(*item)
		}

		wg.Wait()
	}()
}

// vectorPath returns the filepath for the supplied vector, in the supplied
// group, under the supplied directory. It prefixes files with `x--` if the
// vector is known to be broken (i.e. carrying the schema.HintIncorrect hint).
func vectorPath(dir string, group string, item *VectorDef, vector *schema.TestVector) string {
	path := filepath.Join(dir, vectorFilename(group, item, vector))
	return path
}

// vectorPath returns the file name for the supplied vector, in the supplied
// group. It prefixes with `x--` if the vector is known to be broken (i.e.
// carrying the schema.HintIncorrect hint).
func vectorFilename(group string, item *VectorDef, vector *schema.TestVector) string {
	filename := fmt.Sprintf("%s--%s--%s.json", group, item.Metadata.ID, vector.Pre.Variants[0].ID)

	// Prefix the file with "x--" if the vector is known to be broken.
	var broken = map[string]struct{}{schema.HintIncorrect: {}}
	for _, hint := range item.Hints {
		if _, ok := broken[hint]; ok {
			filename = brokenVectorPrefix + filename
			break
		}
	}

	return filename
}

// parseVectorFile unnmarshals a JSON serialized test vector stored at the
// given file path and returns it.
func (g *Generator) parseVectorFile(p string) (*schema.TestVector, error) {
	raw, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("reading test vector file: %w", err)
	}
	var vector schema.TestVector
	err = json.Unmarshal(raw, &vector)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling test vector: %w", err)
	}
	return &vector, nil
}

// vectorBytesNoMeta parses the vector at the given file path and returns the
// serialized bytes for the vector after stripping the metadata.
func (g *Generator) vectorBytesNoMeta(p string) ([]byte, error) {
	v, err := g.parseVectorFile(p)
	if err != nil {
		return nil, err
	}
	v.Meta = nil
	return json.Marshal(v)
}

// vectorsEqual determines if two vectors are "equal". They are considered
// equal if they serialize to the same bytes without a `_meta` property.
func (g *Generator) vectorsEqual(apath, bpath string) (bool, error) {
	abytes, err := g.vectorBytesNoMeta(apath)
	if err != nil {
		return false, err
	}
	bbytes, err := g.vectorBytesNoMeta(bpath)
	if err != nil {
		return false, err
	}
	return bytes.Equal(abytes, bbytes), nil
}

func (g *Generator) generateVariants(b VectorDef) []*schema.TestVector {
	if len(b.SupportedVersions) == 0 {
		b.SupportedVersions = KnownProtocolVersions
	}

	// stamp with our generation data.
	b.Metadata.Gen = genData

	var result []*schema.TestVector
	for _, version := range b.SupportedVersions {
		log.Printf("generating vector [%s] ~~>> pv: [%s]", b.Metadata.ID, version.ID)

		var vector Builder
		// TODO: currently if an assertion fails, we call os.Exit(1), which
		//  aborts all ongoing vector generations. The Asserter should
		//  call runtime.Goexit() instead so only that goroutine is
		//  cancelled. The assertion error must bubble up somehow.
		switch {
		case b.MessageFunc != nil:
			v := MessageVector(b.Metadata, b.Selector, b.Mode, b.Hints, version)
			b.MessageFunc(v)
			vector = v
		case b.TipsetFunc != nil:
			v := TipsetVector(b.Metadata, b.Selector, b.Mode, b.Hints, version)
			b.TipsetFunc(v)
			vector = v
		default:
			panic("no generation function provided")
		}

		// Finish the vector.
		v := vector.Finish()
		result = append(result, v)
	}

	uniq := make(map[[32]byte]*schema.TestVector)
	for _, v := range result {
		variants := v.Pre.Variants                  // stash the variants
		v.Pre.Variants = nil                        // compare without variants
		hash := blake2b.Sum256(v.MustMarshalJSON()) // hash the serialized form
		if merged, ok := uniq[hash]; ok {
			// dedup.
			merged.Pre.Variants = append(merged.Pre.Variants, variants...)
			continue
		}
		v.Pre.Variants = variants // restore the variant
		uniq[hash] = v
	}

	var merged [][]string
	for _, vector := range uniq {
		var ids []string
		for _, v := range vector.Pre.Variants {
			ids = append(ids, v.ID)
		}
		merged = append(merged, ids)
	}

	log.Printf("merged equivalent variants for vector %s; deduped groups: %v", b.Metadata.ID, merged)

	var ret []*schema.TestVector
	for _, v := range uniq {
		ret = append(ret, v)
	}

	return ret
}

// ensureDirectory checks if the provided path is a directory. If yes, it
// returns nil. If the path doesn't exist, it creates the directory and
// returns nil. If the path is not a directory, or another error occurs, an
// error is returned.
func ensureDirectory(path string) error {
	switch stat, err := os.Stat(path); {
	case os.IsNotExist(err):
		// create directory.
		log.Printf("creating directory %s", path)
		err := os.MkdirAll(path, 0700)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %s", path, err)
		}

	case err == nil && !stat.IsDir():
		return fmt.Errorf("path %s exists, but it's not a directory", path)

	case err != nil:
		return fmt.Errorf("failed to stat directory %s: %w", path, err)
	}
	return nil
}

// removePrevious will remove a previously broken vector file if it
// became fixed & removes a previously working vector file if it became broken.
func removePrevious(dir string, group string, item *VectorDef, vector *schema.TestVector) error {
	filename := vectorFilename(group, item, vector)
	var filepath string

	if strings.HasPrefix(filename, brokenVectorPrefix) {
		filepath = path.Join(dir, strings.Replace(filename, brokenVectorPrefix, "", 1))
	} else {
		filepath = path.Join(dir, brokenVectorPrefix+filename)
	}

	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
