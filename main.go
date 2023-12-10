package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
)

const (
	tmpDir          = "/tmp"
	namespace       = "http://schemas.android.com/apk/res/android"
	versionCodeAttr = "versionCode"
	versionNameAttr = "versionName"
)

type Config struct {
	versionCode int32
	versionName string
	packageName string
	minSdkVersion int
}

func main() {
	versionCode := flag.Uint("versionCode", 0, "The versionCode to set")
	versionName := flag.String("versionName", "", "The versionName to set")
	packageName := flag.String("package", "", "The package to set")
	minSdkVersion := flag.Int("minSdkVersion", 0, "The minSdkVersion to set")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(flag.CommandLine.Output(), "Error: File path is required.")
		flag.Usage()
		os.Exit(2)
	}
	config := &Config{
		versionCode: int32(*versionCode),
		versionName: *versionName,
		packageName: *packageName,
		minSdkVersion: *minSdkVersion,
	}

	path := flag.Arg(0)

	if strings.HasSuffix(path, ".apk") {
		updateApk(path, config)
	} else if strings.HasSuffix(path, ".aab") {
		updateAab(path, config)
	} else {
		updateManifest(path, config)
	}
}

func updateApk(path string, config *Config) {
	file, err := ioutil.TempFile(tmpDir, "*.aar")
	if err != nil {
		log.Fatalln("Failed creating temp file:", err)
	}
	defer os.Remove(file.Name())

	out, err := exec.Command("aapt2", "convert", "-o", file.Name(), "--output-format", "proto", path).CombinedOutput()
	if err != nil {
		log.Fatalln("Failed executing aapt2:", err, string(out))
	}

	updateManifestPbInZip(file.Name(), "AndroidManifest.xml", config)

	out, err = exec.Command("aapt2", "convert", "-o", path, "--output-format", "binary", file.Name()).CombinedOutput()
	if err != nil {
		log.Fatalln("Failed executing aapt2:", err, string(out))
	}
}

func updateAab(path string, config *Config) {
	updateManifestPbInZip(path, "base/manifest/AndroidManifest.xml", config)
}

func updateManifestPbInZip(path string, manifestPath string, config *Config) {
	manifest, err := ioutil.TempFile(tmpDir, "AndroidManifest.*.xml")
	if err != nil {
		log.Fatalln("Failed creating temp file:", err)
	}
	defer os.Remove(manifest.Name())

	extractFromZip(path, manifestPath, manifest)
	updateManifest(manifest.Name(), config)
	addToZip(path, manifestPath, manifest)
}

func addToZip(zipPath string, name string, source *os.File) {
	manifestDir, err := ioutil.TempDir(tmpDir, "*")
	if err != nil {
		log.Fatalln("Failed creating temp dir:", err)
	}
	defer os.RemoveAll(manifestDir)

	tmpPath := path.Join(manifestDir, name)
	os.MkdirAll(path.Dir(tmpPath), 0700)
	f, err := os.Create(tmpPath)
	if err != nil {
		log.Fatalln("Failed opening file:", err)
	}
	defer f.Close()
	source.Seek(0, 0)
	io.Copy(f, source)

	absZipPath, err := filepath.Abs(zipPath)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("zip", absZipPath, name)
	cmd.Dir = manifestDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("Failed executing zip:", err, string(out))
	}
}

func extractFromZip(path string, name string, target *os.File) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	f := findFile(r, name)
	if f == nil {
		log.Fatalln(errors.New("file is missing"))
	}

	innerFile, err := f.Open()
	if err != nil {
		log.Fatalln("Failed opening zip file's AndroidManifest.xml:", err)
	}
	defer innerFile.Close()
	_, err = io.Copy(target, innerFile)
	if err != nil {
		log.Fatal(err)
	}
}

func findFile(r *zip.ReadCloser, name string) *zip.File {
	for _, f := range r.File {
		if f.Name != name {
			continue
		}
		return f
	}
	return nil
}

func updateManifest(path string, config *Config) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	xmlNode := &XmlNode{}
	if err := proto.Unmarshal(in, xmlNode); err != nil {
		log.Fatalln("Failed to parse manifest:", err)
	}
	for _, node := range xmlNode.GetElement().GetChild() {
		if elem, ok := node.GetNode().(*XmlNode_Element); ok {
			element := elem.Element
			if element.GetName() == "uses-sdk" {
				for _, attr := range element.GetAttribute() {
					if attr.GetNamespaceUri() == "http://schemas.android.com/apk/res/android" {
						switch attr.GetName() {
						case "minSdkVersion":
							if config.minSdkVersion > 0 {
								fmt.Println("Changing minSdkVersion from", attr.Value, "to", config.minSdkVersion)
								attr.Value = strconv.Itoa(config.minSdkVersion)
							}
						}
					}
				}
			}
		}
	}
	for _, attr := range xmlNode.GetElement().GetAttribute() {
		if attr.GetNamespaceUri() == "" && attr.GetName() == "package" {
			if config.packageName != "" {
				fmt.Println("Changing packageName from", attr.Value, "to", config.packageName)
				attr.Value = config.packageName
			}
		}
		if attr.GetNamespaceUri() != namespace {
			continue
		}
		switch attr.GetName() {
		case versionCodeAttr:
			if config.versionCode > 0 {
				prim := attr.GetCompiledItem().GetPrim()
				if x, ok := prim.GetOneofValue().(*Primitive_IntDecimalValue); ok {
					fmt.Println("Changing versionCode from", x.IntDecimalValue, "to", config.versionCode)
					x.IntDecimalValue = int32(config.versionCode)
				}
				// In AABs the value exists, but when using aapt2 to convert the binary manifest the value is gone
				if attr.Value != "" {
					attr.Value = fmt.Sprint(config.versionCode)
				}
			}
		case versionNameAttr:
			if config.versionName != "" {
				fmt.Println("Changing versionName from", attr.Value, "to", config.versionName)
				attr.Value = config.versionName
			}
		}
	}

	// We use MarshalVT because it keeps the correct field ordering.
	// With the standard Marshal function, Android Studio can't read the resulting proto file inside aab files. :-/
	out, err := xmlNode.MarshalVT()
	if err != nil {
		log.Fatalln("Error marshalling XML:", err)
	}
	if err := ioutil.WriteFile(path, out, 0600); err != nil {
		log.Fatalln("Error writing file:", err)
	}
}
