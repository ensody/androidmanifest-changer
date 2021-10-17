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
}

func main() {
	versionCode := flag.Uint("versionCode", 0, "The versionCode to set")
	versionName := flag.String("versionName", "", "The versionName to set")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(flag.CommandLine.Output(), "Error: File path is required.")
		flag.Usage()
		os.Exit(2)
	}
	config := &Config{
		versionCode: int32(*versionCode),
		versionName: *versionName,
	}

	path := flag.Arg(0)

	if strings.HasSuffix(path, ".apk") {
		updateApk(path, config)
	} else if strings.HasSuffix(path, ".aar") {
		updateAar(path, config)
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

	updateAar(file.Name(), config)

	out, err = exec.Command("aapt2", "convert", "-o", path, "--output-format", "binary", file.Name()).CombinedOutput()
	if err != nil {
		log.Fatalln("Failed executing aapt2:", err, string(out))
	}
}

func updateAar(path string, config *Config) {
	manifest, err := ioutil.TempFile(tmpDir, "AndroidManifest.*.xml")
	if err != nil {
		log.Fatalln("Failed creating temp file:", err)
	}
	defer os.Remove(manifest.Name())

	extractFromZip(path, "AndroidManifest.xml", manifest)
	updateManifest(manifest.Name(), config)
	addToZip(path, "AndroidManifest.xml", manifest)
}

func addToZip(zipPath string, name string, source *os.File) {
	manifestDir, err := ioutil.TempDir(tmpDir, "*")
	if err != nil {
		log.Fatalln("Failed creating temp dir:", err)
	}
	defer os.RemoveAll(manifestDir)

	f, err := os.Create(path.Join(manifestDir, "AndroidManifest.xml"))
	if err != nil {
		log.Fatalln("Failed opening file:", err)
	}
	defer f.Close()
	source.Seek(0, 0)
	io.Copy(f, source)

	out, err := exec.Command("zip", "-j", zipPath, f.Name()).CombinedOutput()
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
	for _, attr := range xmlNode.GetElement().GetAttribute() {
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
			}
		case versionNameAttr:
			if config.versionName != "" {
				fmt.Println("Changing versionName from", attr.Value, "to", config.versionName)
				attr.Value = config.versionName
			}
		}
	}

	out, err := proto.Marshal(xmlNode)
	if err != nil {
		log.Fatalln("Error marshalling XML:", err)
	}
	if err := ioutil.WriteFile(path, out, 0644); err != nil {
		log.Fatalln("Error writing file:", err)
	}
}
